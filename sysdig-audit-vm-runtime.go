package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	DebugLog "sysdig-audit-vm-runtime/debuglog"
	"sysdig-audit-vm-runtime/payloads"
	"sysdig-audit-vm-runtime/sysdighttp"
	"sysdig-audit-vm-runtime/types"
	"time"
)

var dlog = DebugLog.DebugLogger{}

func getOSEnvString(environmentVariable string, optional bool) string {
	env := os.Getenv(environmentVariable)
	if env == "" {
		if !optional {
			dlog.Fatalf("Fatal Error: Could not find %s environment variable, exiting status code (1)", environmentVariable)
			os.Exit(1)
		} else {
			dlog.Printf("Warning: Could not find %s environment variable, continuing anyway...", environmentVariable)
		}
	} else {
		dlog.Printf("Found %s Variable, continuing ...", environmentVariable)
	}
	return env
}

func main() {
	fmt.Println("Sysdig-Audit-VM-Runtime 0.3")
	fmt.Print("\n")

	// Set out custom -h/--help usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of sysdig-audit-vm-runtime:\n")
		fmt.Fprintf(os.Stderr, "  SECURE_API_TOKEN=xxx sysdig-audit-vm-runtime [OPTIONS]\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  --help\t\tDisplay Help\n")
		fmt.Fprintf(os.Stderr, "  --api\t\t\tSpecify Sysdig API URL\n")
		fmt.Fprintf(os.Stderr, "  --cluster\t\tCluster to process (Default is all)\n")
		fmt.Fprintf(os.Stderr, "  --debug\t\tLog extra debug information\n")
		fmt.Fprintf(os.Stderr, "  --ignorejobs\t\tIgnore Jobs or CronJobs\n")
		fmt.Fprintf(os.Stderr, "  --last\t\tLast xx minutes of data to compare (default: 60 minutes)\n")
		fmt.Fprintf(os.Stderr, "\n")
	}
	strAPIKey := getOSEnvString("SECURE_API_TOKEN", false)
	if strAPIKey == "" {
		dlog.Fatalf("main:: Please set SECURE_API_TOKEN variable.  Exiting...")
	}
	var boolIgnoreJobs bool
	flag.BoolVar(&boolIgnoreJobs, "ignorejobs", false, "Ignore Jobs / Cronjobs")
	flag.BoolVar(&DebugLog.Debug, "debug", false, "Enable debug logging")
	intLast := flag.Int("last", 60, "Last xx minutes of data to compare (default: 60 minutes)")
	clusterName := flag.String("cluster", "", "Name of the Kubernetes cluster")
	ApiURL := flag.String("api", "", "Specify Sysdig API URL")
	flag.Parse()
	if strings.HasSuffix(*ApiURL, "/") {
		*ApiURL = strings.TrimSuffix(*ApiURL, "/")
	}

	if *ApiURL == "" {
		dlog.Fatalf("main:: Please specify a sysdig --api URL")
	} else {
		dlog.Printf("main:: Using URL: %s", *ApiURL)
	}
	requestStr := payloads.K8sLiveJson
	request := types.K8sLiveRequestWrapper{}
	if err := json.Unmarshal([]byte(requestStr), &request); err != nil {
		dlog.Println("main:: Error unmarshaling JSON:", err)
		return
	}

	dtNow := time.Now()

	// Round down to the nearest 10 minutes, in microseconds then the 24 hours prior
	dtNow = dtNow.
		Add(-time.Minute * time.Duration(dtNow.Minute()%10)).
		Add(-time.Second * time.Duration(dtNow.Second())).
		Add(-time.Nanosecond * time.Duration(dtNow.Nanosecond()))

	dtToMicro := dtNow.UnixNano() / 1000
	dtFromMicro := dtToMicro - (int64(*intLast) * 60000000)

	// Convert epoch time to time.Time
	dtTo := time.Unix(dtToMicro/1e6, 0)
	dtFrom := time.Unix(dtFromMicro/1e6, 0)

	request.Requests[0].Time["to"] = dtToMicro
	dlog.Printf("main:: From epoch: %d", dtFromMicro)
	dlog.Printf("main:: To epoch: %d", dtToMicro)
	dlog.Printf("main:: From Date/Time: %s", dtFrom)
	dlog.Printf("main:: To Date/Time: %s", dtTo)

	request.Requests[0].Time["from"] = dtFromMicro
	if *clusterName != "" {
		// set the cluster name explicitly
		request.Requests[0].Scope = fmt.Sprintf("kubernetes.cluster.name = \"%s\"", *clusterName)
	}
	dlog.Printf("main:: Processing Kubernetes Live data, please wait.")
	dlog.Printf("main:: While you wait.")
	dlog.Printf("main:: What do you call a mislabeled orange juice container?")

	configHTTPK8sLive := sysdighttp.DefaultSysdigRequestConfig()
	configHTTPK8sLive.Method = "POST"
	configHTTPK8sLive.URL = *ApiURL + "/api/data/batch?emptyValuesAsNull=true&dynamicSampling=true"
	configHTTPK8sLive.Headers = map[string]string{
		"Authorization":     "bearer " + strAPIKey,
		"emptyValuesAsNull": "true",
		"dynamicSampling":   "true",
	}
	configHTTPK8sLive.JSON = request
	objResponse, err := sysdighttp.SysdigRequest(configHTTPK8sLive)
	if err != nil {
		dlog.Fatalf("main:: Could not query K8Live API.  Exiting %v...", err)

	}
	var jsonK8sLiveResponse types.K8sLiveResponse
	err = sysdighttp.ResponseBodyToJson(objResponse, &jsonK8sLiveResponse)
	if err != nil {
		dlog.Fatalf("main:: failed to convert body to JSON: %v", err)
	}

	//build the map of K8s live data to use later
	arrK8sLiveWorkloads := make(map[string]types.WorkloadStruct)
	arrsortedK8sLiveWorkloads := []string{}

	for _, item := range jsonK8sLiveResponse.Responses[0].Data {
		entry := types.WorkloadStruct{
			ClusterName:   item["k2"],
			NamespaceName: item["k1"],
			WorkloadName:  item["k0"],
			JobName:       item["k3"],
			WorkloadType:  item["k4"],
		}
		if entry.JobName != "" && boolIgnoreJobs {
			dlog.Printf("Ignoring 'Job/Cronjob' Workload %s (JobName:%s)", entry.WorkloadName, entry.JobName)
		} else {
			key := fmt.Sprintf("%s / %s / %s / %s", entry.ClusterName, entry.NamespaceName, entry.WorkloadName, entry.WorkloadType)
			arrK8sLiveWorkloads[key] = entry
		}
	}

	for key := range arrK8sLiveWorkloads {
		arrsortedK8sLiveWorkloads = append(arrsortedK8sLiveWorkloads, key)
	}

	sort.Strings(arrsortedK8sLiveWorkloads)
	for index, item := range arrsortedK8sLiveWorkloads {
		dlog.Printf("Index: %d, arrK8sLiveWorkloads Workload: %s", index, item)
	}
	intK8sCount := len(arrsortedK8sLiveWorkloads)

	// Now get the scanning data we need
	configScanning := sysdighttp.DefaultSysdigRequestConfig()
	configScanning.URL = *ApiURL + "/api/scanning/runtime/v2/workflows/results"
	configScanning.Headers = map[string]string{
		"Authorization": "bearer " + strAPIKey,
	}
	configScanning.Params = map[string]interface{}{}
	// Max is 1000, so lets return 1000 then iterate through the pages until we get no more cursors
	configScanning.Params["limit"] = 1000
	if *clusterName != "" {
		configScanning.Params["filter"] = fmt.Sprintf("kubernetes.cluster.name = \"%s\"", *clusterName)
		dlog.Printf("main:: filtering on kubernetes.cluster.name = \"%s\"", *clusterName)
	}

	objScanningResponse, err := sysdighttp.SysdigRequest(configScanning)
	if err != nil {
		dlog.Fatalf("main:: Could not query Scanning API.  Exiting %v...", err)

	}

	var jsonScanningResponse types.ScanningResponse

	err = sysdighttp.ResponseBodyToJson(objScanningResponse, &jsonScanningResponse)
	if err != nil {
		dlog.Println("main:: failed to convert body to JSON: %w", err)
	}
	dlog.Printf("main:: Pulp Fiction")
	dlog.Printf("main:: *cue laughter... or eye rolling*")
	//Build the runtime map to use
	arrRuntimeWorkloads := make(map[string]types.WorkloadStruct)
	arrSortedRuntimeWorkloads := []string{}

	for {
		dlog.Printf("main:: Processing Runtime Recordcount: %d/%d", int(jsonScanningResponse.Page["returned"].(float64)), int(jsonScanningResponse.Page["matched"].(float64)))
		for _, item := range jsonScanningResponse.Data {
			if item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["asset.type"].(string) == "workload" {
				entry := types.WorkloadStruct{
					ClusterName:   item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["kubernetes.cluster.name"].(string),
					NamespaceName: item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["kubernetes.namespace.name"].(string),
					WorkloadName:  item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["kubernetes.workload.name"].(string),
					WorkloadType:  item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["kubernetes.workload.type"].(string),
				}
				key := fmt.Sprintf("%s / %s / %s / %s", entry.ClusterName, entry.NamespaceName, entry.WorkloadName, entry.WorkloadType)
				arrRuntimeWorkloads[key] = entry
			}
		}
		// Now check if we have another page to process.. if so.. then do
		if jsonScanningResponse.Page["next"] != nil {
			dlog.Printf("main:: Found cursor, querying next set of data.  Cursor = %s", jsonScanningResponse.Page["next"])
			configScanning.Params["cursor"] = jsonScanningResponse.Page["next"]
			objScanningResponse, err = sysdighttp.SysdigRequest(configScanning)
			if err != nil {
				dlog.Fatalf("main:: Could not query Scanning API.  Exiting %v...", err)
			}
			err = sysdighttp.ResponseBodyToJson(objScanningResponse, &jsonScanningResponse)
			if err != nil {
				dlog.Println("main:: failed to convert body to JSON: %w", err)
			}
		} else {
			break
		}
	}

	for key := range arrRuntimeWorkloads {
		arrSortedRuntimeWorkloads = append(arrSortedRuntimeWorkloads, key)
	}

	sort.Strings(arrSortedRuntimeWorkloads)
	for index, item := range arrSortedRuntimeWorkloads {
		dlog.Printf("main:: Index %d arrRuntimeWorkloads Workload: %s", index, item)
	}
	intRuntimeCount := len(arrRuntimeWorkloads)

	dlog.Println("main:: Finished builing lookup arrary 'arrRuntimeWorkloads'")
	dlog.Println("\nmain:: Being logging results...")

	fmt.Println("Workloads found that are missing from VM Runtime Scanning:")
	fmt.Println("Cluster / Namespace / Workload Name / Workload Type")
	count := 0

	var arrSortedResult []string

	for key := range arrK8sLiveWorkloads {
		if _, exists := arrRuntimeWorkloads[key]; !exists {
			arrSortedResult = append(arrSortedResult, key)
			count += 1
		}
	}
	sort.Strings(arrSortedResult)
	for _, item := range arrSortedResult {
		fmt.Println(item)
	}
	fmt.Printf("\nTotal workloads detected from Kubernetes Live: %d", intK8sCount)
	fmt.Printf("\nTotal workloads detected in VM Runtime Scanning: %d", intRuntimeCount)
	fmt.Printf("\nTotal workloads missing in VM Runtime Scanning: %d\n", count)
	dlog.Println("Finished...")
}
