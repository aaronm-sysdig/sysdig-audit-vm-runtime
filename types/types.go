package types

type K8sLiveRequestWrapper struct {
	Requests []K8sLiveRequest `json:"requests"`
}

type K8sLiveRequest struct {
	Format  map[string]interface{} `json:"format"`
	Time    map[string]interface{} `json:"time"`
	Metrics map[string]string      `json:"metrics"`
	Group   map[string]interface{} `json:"group"`
	Paging  map[string]interface{} `json:"paging"`
	Sort    []map[string]string    `json:"sort"`
	Scope   string                 `json:"scope"`
}

type K8sLiveDataStruct struct {
	Requests []K8sLiveDataStruct `json:"requests"`
}

type K8sLiveResponse struct {
	Responses []struct {
		Format  map[string]interface{} `json:"format"`
		Time    map[string]interface{} `json:"time"`
		Group   map[string]interface{} `json:"group"`
		Sort    []map[string]string    `json:"sort"`
		Paging  map[string]interface{} `json:"paging"`
		Data    []map[string]string    `json:"data"`
		Metrics map[string]string      `json:"metrics"`
	} `json:"responses"`
}

type ScanningResponse struct {
	Page map[string]interface{}   `json:"page"`
	Data []map[string]interface{} `json:"data"`
}

type WorkloadStruct struct {
	ClusterName   string
	NamespaceName string
	WorkloadName  string
	JobName       string
	CronJobName   string
}

type RunningResponse struct {
	Data []map[string]interface{} `json:"data"`
}
