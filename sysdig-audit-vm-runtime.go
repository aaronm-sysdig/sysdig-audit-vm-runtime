package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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
	// Set out custom -h/--help usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of sysdig-audit-vm-runtime:\n")
		fmt.Fprintf(os.Stderr, "  SECURE_API_TOKEN=xxx sysdig-audit-vm-runtime [OPTIONS]\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -h/--help\tDisplay Help\n")
		fmt.Fprintf(os.Stderr, "  -a/--api\tSpecify Sysdig API URL\n")
		fmt.Fprintf(os.Stderr, "  -c/--cluster\tCluster to process (Default is all)\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
	strAPIKey := getOSEnvString("SECURE_API_TOKEN", false)
	if strAPIKey == "" {
		dlog.Fatalf("main:: Please set SECURE_API_TOKEN variable.  Exiting...")
	}
	flag.BoolVar(&DebugLog.Debug, "debug", false, "Enable debug logging")
	clusterName := flag.String("cluster", "", "Name of the Kubernetes cluster")
	flag.Parse()

	fmt.Println("Sysdig-Audit-VM-Runtime 0.1")
	fmt.Print("\n\n")

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
	dtFromMicro := dtToMicro - 86400000000

	request.Requests[0].Time["to"] = dtToMicro
	dlog.Printf("main:: From epoch: %d", dtFromMicro)
	dlog.Printf("main:: To epoch: %d", dtToMicro)

	request.Requests[0].Time["from"] = dtFromMicro
	if clusterName != nil {
		// set the cluster name explicitly
		request.Requests[0].Scope = fmt.Sprintf("kubernetes.cluster.name = \"%s\"", *clusterName)
	}

	configHTTPK8sLive := sysdighttp.DefaultSysdigRequestConfig()
	configHTTPK8sLive.Method = "POST"
	configHTTPK8sLive.URL = "https://app.au1.sysdig.com/api/data/batch?emptyValuesAsNull=true&dynamicSampling=true"
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
	intK8sCount := 0
	arrK8sLiveWorkloads := make(map[string]types.WorkloadStruct)
	for index, item := range jsonK8sLiveResponse.Responses[0].Data {
		entry := types.WorkloadStruct{
			ClusterName:   item["k2"],
			NamespaceName: item["k1"],
			WorkloadName:  item["k0"],
		}
		key := fmt.Sprintf("%s / %s / %s", entry.ClusterName, entry.NamespaceName, entry.WorkloadName)
		arrK8sLiveWorkloads[key] = entry
		dlog.Printf("Index: %d, arrK8sLiveWorkloads Workload: %s / %s / %s", index, item["k2"], item["k1"], item["k0"])
		intK8sCount += 1
	}

	// Now get the scanning data we need
	configScanning := sysdighttp.DefaultSysdigRequestConfig()
	configScanning.URL = "https://app.au1.sysdig.com/api/scanning/runtime/v2/workflows/results"
	configScanning.Headers = map[string]string{
		"Authorization": "bearer " + strAPIKey,
	}
	configScanning.Params = map[string]interface{}{
		"limit": 9999,
	}

	objScanningResponse, err := sysdighttp.SysdigRequest(configScanning)
	if err != nil {
		dlog.Fatalf("main:: Could not query K8Live API.  Exiting %v...", err)

	}
	var jsonScanningResponse types.ScanningResponse
	err = sysdighttp.ResponseBodyToJson(objScanningResponse, &jsonScanningResponse)
	if err != nil {
		dlog.Println("main:: failed to convert body to JSON: %w", err)
	}

	//Build the runtime map to use
	intRuntimeCount := 0
	arrRuntimeWorkloads := make(map[string]types.WorkloadStruct)
	for index, item := range jsonScanningResponse.Data {
		if item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["asset.type"].(string) == "workload" {
			entry := types.WorkloadStruct{
				ClusterName:   item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["kubernetes.cluster.name"].(string),
				NamespaceName: item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["kubernetes.namespace.name"].(string),
				WorkloadName:  item["recordDetails"].(map[string]interface{})["labels"].(map[string]interface{})["kubernetes.workload.name"].(string),
			}
			key := fmt.Sprintf("%s / %s / %s", entry.ClusterName, entry.NamespaceName, entry.WorkloadName)
			arrRuntimeWorkloads[key] = entry
			dlog.Printf("main:: Index %d arrRuntimeWorkloads Workload: %s", index, key)
			intRuntimeCount += 1
		}
	}
	dlog.Println("main:: Finished builing lookup arrary 'arrRuntimeWorkloads'")
	dlog.Println("\nmain:: Being logging results...")

	fmt.Println("Workloads found that are missing from VM Runtime Scanning:")
	fmt.Println("Cluster / Namespace / Workload")
	count := 0
	for key := range arrK8sLiveWorkloads {
		if _, exists := arrRuntimeWorkloads[key]; !exists {
			// The workload from arrK8sLiveWorkloads is not in arrRuntimeWorkloads
			fmt.Println(key)
			count += 1
		}
	}
	fmt.Printf("\nTotal workloads detected from Kubernetes Live: %d", intK8sCount)
	fmt.Printf("\nTotal workloads detected in VM Runtime Scanning: %d", intRuntimeCount)
	fmt.Printf("\nTotal workloads missing in VM Runtime Scanning: %d\n", count)
	dlog.Println("Finished...")
}
