# Disclaimer
Notwithstanding anything that may be contained to the contrary in your agreement(s) with Sysdig, Sysdig provides no support, no updates, and no warranty or guarantee of any kind with respect to these script(s), including as to their functionality or their ability to work in your environment(s). Sysdig disclaims all liability and responsibility with respect to any use of these scripts.

# sysdig-audit-vm-runtime #

## Introduction ##
A simple GO application aimed at validating your running workloads to ensure they are also showing in the runtime vuln view.  This is an inital POC so there may still be work to do here.

## Compilation / Running ##
Example of how to compile app
```
GOOS=linux GOARCH=amd64 go build sysdig-audit-vm-runtime.go 
```
I wrote it against Go 1.23.3 darwin/arm64

## Usage ##
```
Sysdig-Audit-VM-Runtime 0.2

Usage of sysdig-audit-vm-runtime:
  SECURE_API_TOKEN=xxx sysdig-audit-vm-runtime [OPTIONS]

Options:
  --help		Display Help
  --api			Specify Sysdig API URL
  --cluster		Cluster to process (Default is all)
  --debug		Log extra debug information
  --ignorejobs		Ignore Jobs or CronJobs
```

## Example Output ##

```
SECURE_API_TOKEN=xxxx go run sysdig-audit-vm-runtime.go                   
Sysdig-Audit-VM-Runtime 0.1


Workloads found that are missing from VM Runtime Scanning:
Cluster / Namespace / Workload

Total workloads detected from Kubernetes Live: 44
Total workloads detected in VM Runtime Scanning: 61
Total workloads missing in VM Runtime Scanning: 0
```
## Example with cluster specified ##
```
SECURE_API_TOKEN=xxxx go run sysdig-audit-vm-runtime.go --cluster kubernetes        
Sysdig-Audit-VM-Runtime 0.2


Workloads found that are missing from VM Runtime Scanning:
Cluster / Namespace / Workload

Total workloads detected from Kubernetes Live: 44
Total workloads detected in VM Runtime Scanning: 61
Total workloads missing in VM Runtime Scanning: 0

```

## Example Debug Output ##
```
Sysdig-Audit-VM-Runtime 0.2

2023/10/15 23:36:11 main:: Using URL: https://app.au1.sysdig.com
2023/10/15 23:36:11 main:: From epoch: 1697286600000000
2023/10/15 23:36:11 main:: To epoch: 1697373000000000
2023/10/15 23:36:11 Index: 0, arrK8sLiveWorkloads Workload: kubernetes / calico-apiserver / calico-apiserver
2023/10/15 23:36:11 Index: 1, arrK8sLiveWorkloads Workload: kubernetes / calico-system / calico-kube-controllers
2023/10/15 23:36:11 Index: 2, arrK8sLiveWorkloads Workload: kubernetes / calico-system / calico-node
2023/10/15 23:36:11 Index: 3, arrK8sLiveWorkloads Workload: kubernetes / calico-system / calico-typha
2023/10/15 23:36:11 Index: 4, arrK8sLiveWorkloads Workload: kubernetes / calico-system / csi-node-driver
2023/10/15 23:36:11 Index: 5, arrK8sLiveWorkloads Workload: kubernetes / example-voting-app / worker
2023/10/15 23:36:11 Index: 6, arrK8sLiveWorkloads Workload: kubernetes / harbor / harbor-core
2023/10/15 23:36:11 Index: 7, arrK8sLiveWorkloads Workload: kubernetes / harbor / harbor-database
2023/10/15 23:36:11 Index: 8, arrK8sLiveWorkloads Workload: kubernetes / harbor / harbor-jobservice
2023/10/15 23:36:11 Index: 9, arrK8sLiveWorkloads Workload: kubernetes / harbor / harbor-nginx
2023/10/15 23:36:11 Index: 10, arrK8sLiveWorkloads Workload: kubernetes / harbor / harbor-portal
2023/10/15 23:36:11 Index: 11, arrK8sLiveWorkloads Workload: kubernetes / harbor / harbor-redis
2023/10/15 23:36:11 Index: 12, arrK8sLiveWorkloads Workload: kubernetes / harbor / harbor-registry
2023/10/15 23:36:11 Index: 13, arrK8sLiveWorkloads Workload: kubernetes / kube-system / coredns
2023/10/15 23:36:11 Index: 14, arrK8sLiveWorkloads Workload: kubernetes / kube-system / etcd-lab-master
2023/10/15 23:36:11 Index: 15, arrK8sLiveWorkloads Workload: kubernetes / kube-system / kube-apiserver-lab-master
2023/10/15 23:36:11 Index: 16, arrK8sLiveWorkloads Workload: kubernetes / kube-system / kube-controller-manager-lab-master
2023/10/15 23:36:11 Index: 17, arrK8sLiveWorkloads Workload: kubernetes / kube-system / kube-proxy
2023/10/15 23:36:11 Index: 18, arrK8sLiveWorkloads Workload: kubernetes / kube-system / kube-scheduler-lab-master
2023/10/15 23:36:11 Index: 19, arrK8sLiveWorkloads Workload: kubernetes / local-path-storage / local-path-provisioner
2023/10/15 23:36:11 Index: 20, arrK8sLiveWorkloads Workload: kubernetes / registry / forever-sleep
2023/10/15 23:36:11 Index: 21, arrK8sLiveWorkloads Workload: kubernetes / registry / regscanjob
2023/10/15 23:36:11 Index: 22, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / carts
2023/10/15 23:36:11 Index: 23, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / carts-db
2023/10/15 23:36:11 Index: 24, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / catalogue
2023/10/15 23:36:11 Index: 25, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / catalogue-db
2023/10/15 23:36:11 Index: 26, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / front-end
2023/10/15 23:36:11 Index: 27, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / load-gen
2023/10/15 23:36:11 Index: 28, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / orders
2023/10/15 23:36:11 Index: 29, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / orders-db
2023/10/15 23:36:11 Index: 30, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / payment
2023/10/15 23:36:11 Index: 31, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / queue-master
2023/10/15 23:36:11 Index: 32, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / rabbitmq
2023/10/15 23:36:11 Index: 33, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / redis-exporter-sock-shop-session-db-deploy
2023/10/15 23:36:11 Index: 34, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / session-db
2023/10/15 23:36:11 Index: 35, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / shipping
2023/10/15 23:36:11 Index: 36, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / user
2023/10/15 23:36:11 Index: 37, arrK8sLiveWorkloads Workload: kubernetes / sock-shop / user-db
2023/10/15 23:36:11 Index: 38, arrK8sLiveWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent
2023/10/15 23:36:11 Index: 39, arrK8sLiveWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-admissioncontroller-scanner
2023/10/15 23:36:11 Index: 40, arrK8sLiveWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-admissioncontroller-webhook
2023/10/15 23:36:11 Index: 41, arrK8sLiveWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-kspmcollector
2023/10/15 23:36:11 Index: 42, arrK8sLiveWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-node-analyzer
2023/10/15 23:36:11 Index: 43, arrK8sLiveWorkloads Workload: kubernetes / temp / nginx
2023/10/15 23:36:11 Index: 44, arrK8sLiveWorkloads Workload: kubernetes / temp / test
2023/10/15 23:36:11 Index: 45, arrK8sLiveWorkloads Workload: kubernetes / tigera-operator / tigera-operator
2023/10/15 23:36:11 main:: filtering on kubernetes.cluster.name = "kubernetes"
2023/10/15 23:36:11 main:: Index 0 arrRuntimeWorkloads Workload: kubernetes / calico-apiserver / calico-apiserver
2023/10/15 23:36:11 main:: Index 1 arrRuntimeWorkloads Workload: kubernetes / calico-system / calico-kube-controllers
2023/10/15 23:36:11 main:: Index 2 arrRuntimeWorkloads Workload: kubernetes / calico-system / calico-node
2023/10/15 23:36:11 main:: Index 3 arrRuntimeWorkloads Workload: kubernetes / calico-system / calico-typha
2023/10/15 23:36:11 main:: Index 4 arrRuntimeWorkloads Workload: kubernetes / calico-system / csi-node-driver
2023/10/15 23:36:11 main:: Index 5 arrRuntimeWorkloads Workload: kubernetes / calico-system / csi-node-driver
2023/10/15 23:36:11 main:: Index 6 arrRuntimeWorkloads Workload: kubernetes / example-voting-app / db
2023/10/15 23:36:11 main:: Index 7 arrRuntimeWorkloads Workload: kubernetes / example-voting-app / redis
2023/10/15 23:36:11 main:: Index 8 arrRuntimeWorkloads Workload: kubernetes / example-voting-app / result
2023/10/15 23:36:11 main:: Index 9 arrRuntimeWorkloads Workload: kubernetes / example-voting-app / vote
2023/10/15 23:36:11 main:: Index 10 arrRuntimeWorkloads Workload: kubernetes / example-voting-app / worker
2023/10/15 23:36:11 main:: Index 11 arrRuntimeWorkloads Workload: kubernetes / harbor / harbor-core
2023/10/15 23:36:11 main:: Index 12 arrRuntimeWorkloads Workload: kubernetes / harbor / harbor-database
2023/10/15 23:36:11 main:: Index 13 arrRuntimeWorkloads Workload: kubernetes / harbor / harbor-jobservice
2023/10/15 23:36:11 main:: Index 14 arrRuntimeWorkloads Workload: kubernetes / harbor / harbor-nginx
2023/10/15 23:36:11 main:: Index 15 arrRuntimeWorkloads Workload: kubernetes / harbor / harbor-portal
2023/10/15 23:36:11 main:: Index 16 arrRuntimeWorkloads Workload: kubernetes / harbor / harbor-redis
2023/10/15 23:36:11 main:: Index 17 arrRuntimeWorkloads Workload: kubernetes / harbor / harbor-registry
2023/10/15 23:36:11 main:: Index 18 arrRuntimeWorkloads Workload: kubernetes / harbor / harbor-registry
2023/10/15 23:36:11 main:: Index 19 arrRuntimeWorkloads Workload: kubernetes / kube-system / coredns
2023/10/15 23:36:11 main:: Index 20 arrRuntimeWorkloads Workload: kubernetes / kube-system / etcd-lab-master
2023/10/15 23:36:11 main:: Index 21 arrRuntimeWorkloads Workload: kubernetes / kube-system / kube-apiserver-lab-master
2023/10/15 23:36:11 main:: Index 22 arrRuntimeWorkloads Workload: kubernetes / kube-system / kube-controller-manager-lab-master
2023/10/15 23:36:11 main:: Index 23 arrRuntimeWorkloads Workload: kubernetes / kube-system / kube-proxy
2023/10/15 23:36:11 main:: Index 24 arrRuntimeWorkloads Workload: kubernetes / kube-system / kube-scheduler-lab-master
2023/10/15 23:36:11 main:: Index 25 arrRuntimeWorkloads Workload: kubernetes / local-path-storage / local-path-provisioner
2023/10/15 23:36:11 main:: Index 26 arrRuntimeWorkloads Workload: kubernetes / registry / forever-sleep-cronjob
2023/10/15 23:36:11 main:: Index 27 arrRuntimeWorkloads Workload: kubernetes / security-playground / security-playground
2023/10/15 23:36:11 main:: Index 28 arrRuntimeWorkloads Workload: kubernetes / sock-shop / carts
2023/10/15 23:36:11 main:: Index 29 arrRuntimeWorkloads Workload: kubernetes / sock-shop / carts-db
2023/10/15 23:36:11 main:: Index 30 arrRuntimeWorkloads Workload: kubernetes / sock-shop / catalogue
2023/10/15 23:36:11 main:: Index 31 arrRuntimeWorkloads Workload: kubernetes / sock-shop / catalogue-db
2023/10/15 23:36:11 main:: Index 32 arrRuntimeWorkloads Workload: kubernetes / sock-shop / front-end
2023/10/15 23:36:11 main:: Index 33 arrRuntimeWorkloads Workload: kubernetes / sock-shop / load-gen
2023/10/15 23:36:11 main:: Index 34 arrRuntimeWorkloads Workload: kubernetes / sock-shop / orders
2023/10/15 23:36:11 main:: Index 35 arrRuntimeWorkloads Workload: kubernetes / sock-shop / orders-db
2023/10/15 23:36:11 main:: Index 36 arrRuntimeWorkloads Workload: kubernetes / sock-shop / payment
2023/10/15 23:36:11 main:: Index 37 arrRuntimeWorkloads Workload: kubernetes / sock-shop / queue-master
2023/10/15 23:36:11 main:: Index 38 arrRuntimeWorkloads Workload: kubernetes / sock-shop / rabbitmq
2023/10/15 23:36:11 main:: Index 39 arrRuntimeWorkloads Workload: kubernetes / sock-shop / redis-exporter-sock-shop-session-db-deploy
2023/10/15 23:36:11 main:: Index 40 arrRuntimeWorkloads Workload: kubernetes / sock-shop / session-db
2023/10/15 23:36:11 main:: Index 41 arrRuntimeWorkloads Workload: kubernetes / sock-shop / shipping
2023/10/15 23:36:11 main:: Index 42 arrRuntimeWorkloads Workload: kubernetes / sock-shop / user
2023/10/15 23:36:11 main:: Index 43 arrRuntimeWorkloads Workload: kubernetes / sock-shop / user-db
2023/10/15 23:36:11 main:: Index 44 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent
2023/10/15 23:36:11 main:: Index 45 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-admissioncontroller-scanner
2023/10/15 23:36:11 main:: Index 46 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-admissioncontroller-webhook
2023/10/15 23:36:11 main:: Index 47 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-kspmcollector
2023/10/15 23:36:11 main:: Index 48 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-node-analyzer
2023/10/15 23:36:11 main:: Index 49 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-node-analyzer
2023/10/15 23:36:11 main:: Index 50 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-node-analyzer
2023/10/15 23:36:11 main:: Index 51 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-node-analyzer
2023/10/15 23:36:11 main:: Index 52 arrRuntimeWorkloads Workload: kubernetes / sysdig-agent / sysdig-agent-node-analyzer
2023/10/15 23:36:11 main:: Index 53 arrRuntimeWorkloads Workload: kubernetes / temp / nginx
2023/10/15 23:36:11 main:: Index 54 arrRuntimeWorkloads Workload: kubernetes / temp / test
2023/10/15 23:36:11 main:: Index 55 arrRuntimeWorkloads Workload: kubernetes / text4shell-deployment / text4shell-patched-v1
2023/10/15 23:36:11 main:: Index 56 arrRuntimeWorkloads Workload: kubernetes / text4shell-deployment / text4shell-patched-v2
2023/10/15 23:36:11 main:: Index 57 arrRuntimeWorkloads Workload: kubernetes / text4shell-deployment / text4shell-vuln
2023/10/15 23:36:11 main:: Index 58 arrRuntimeWorkloads Workload: kubernetes / text4shell-patched / text4shell-patched-v1
2023/10/15 23:36:11 main:: Index 59 arrRuntimeWorkloads Workload: kubernetes / text4shell-patched / text4shell-patched-v2
2023/10/15 23:36:11 main:: Index 60 arrRuntimeWorkloads Workload: kubernetes / text4shell-vuln / text4shell-vuln
2023/10/15 23:36:11 main:: Index 61 arrRuntimeWorkloads Workload: kubernetes / tigera-operator / tigera-operator
2023/10/15 23:36:11 main:: Finished builing lookup arrary 'arrRuntimeWorkloads'
2023/10/15 23:36:11 
main:: Being logging results...
Workloads found that are missing from VM Runtime Scanning:
Cluster / Namespace / Workload
kubernetes / registry / forever-sleep
kubernetes / registry / regscanjob

Total workloads detected from Kubernetes Live: 46
Total workloads detected in VM Runtime Scanning: 62
Total workloads missing in VM Runtime Scanning: 2
2023/10/15 23:36:11 Finished...
```
