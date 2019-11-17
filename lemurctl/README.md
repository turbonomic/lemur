<p align="center">
    <img src="https://user-images.githubusercontent.com/10012486/68960926-2af33d80-079f-11ea-928f-ccfd21d56982.png">
</p>

<!--
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2019 Turbonomic

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

## Overview

**Lemurctl** is a command line utility to access and control Lemur powered by Turbonomic. It offers the capability to view real time entities discovered by Lemur, and the resource consumer and provider relationships between these entities (i.e., the supply chain). Using Lemur, user can easily sort entities of certain type (e.g., applications, or containers) based on their resource consumption. User can also take a top-down approach to examine the metric and resource usage of a specific application and all the other entities along the supply chain of that application to quickly determine the resource bottle neck that may affect the performance of the application.

## Getting Started 
Leverage the [wiki pages](https://github.com/turbonomic/lemur/wiki/Lemurctl) for updated details
### Prerequisites
* Go 1.13 (Go version prior to 1.13 requires module-aware mode, e.g., [GO111MODULE](https://golang.org/cmd/go/#hdr-Module_support)=on)
* Lemur server is installed and running

### Lemurctl Installation
* Clone the repository
```
$ git clone https://github.com/turbonomic/lemur.git
```
* Build the binary
```
$ cd lemur/lemurctl
$ make
```
* Move the binary into your PATH
```
$ sudo mv ./lemurctl /usr/local/bin/lemurctl
```
* Test to ensure the version you installed is up-to-date
```
$ lemurctl -v
```

## Use Cases
Get first familiar with the concept and taxonomy of the Lemur/Turbonomic Supply Chain [here](https://github.com/turbonomic/lemur/wiki/Lemur-Use-Cases)
* Always use the `-h` or `--help` option to get a description of the command syntax and options
```
$ lemurctl -h

Usage: lemurctl [OPTIONS] COMMAND [arg...]

lemurctl controls Lemur, powered by Turbonomic.

...

Options:
  --influxdb value             specify the endpoint of the InfluxDB server (default: "192.168.1.96:8086") [$INFLUXDB_SERVER]
  --debug, -d                  enable debug mode [$DEBUG]
  --log-level value, -l value  specify log level (debug, info, warn, error, fatal, panic) (default: "info") [$LOG_LEVEL]
  --help, -h                   show help
  --version, -v                print the version
  
Commands:
  get, g  Display one or many entities or groups of entities
  help    Shows a list of commands or help for one command
  
Run 'lemurctl COMMAND --help' for more information on a command.

```
* Get a list of kubernetes clusters discovered by Lemur
```
$ lemurctl get cluster
ID                                      TYPE                     
6eb24347-92fe-11e9-96bd-005056b8136b    vm                       
d8eb33cf-f4e5-11e7-9653-005056802f41    vm  
```
You can match the ID to a specific kubernetes cluster by running the following command:
```
$ kubectl -n default get svc kubernetes -o=jsonpath='{.metadata.uid}'
6eb24347-92fe-11e9-96bd-005056b8136b
```
All the following use cases require a cluster ID to be set either through the `LEMUR_CLUSTER` environment variable, or through the `--cluster` command line option such that entities retrieved only belong to one cluster scope. For example:
```
$ export LEMUR_CLUSTER=6eb24347-92fe-11e9-96bd-005056b8136b
```
* Get a list of applications that belong to a kubernetes cluster sorted by VCPU
```
$ lemurctl get app
NAME                                                                                                        VCPU (MHz)VMEM (GB) 
App-istio-system/istio-galley-664b9468d6-kmx6n/galley                                                       52.86     0.02      
App-istio-system/prometheus-776fdf7479-4c6mw/prometheus                                                     51.55     0.58      
App-kube-system/calico-node-85txg/calico-node                                                               45.00     0.15         
App-kube-system/kube-controller-manager-ip-10-0-0-38.ca-central-1.compute.internal/kube-controller-manager  32.71     0.05      
App-istio-system/istio-sidecar-injector-7c5d49854d-5nfqv/sidecar-injector-webhook                           30.31     0.01      
App-turbonomic/kafka-77689b4464-mvkv7/kafka                                                                 25.33     0.74      
App-kube-system/kube-proxy-mhtsh/kube-proxy                                                                 12.59     0.01      
App-turbonomic/topology-processor-7747fb5b97-gc2xx/topology-processor                                       11.64     0.85      
App-kube-system/kube-proxy-hhkfl/kube-proxy                                                                 10.31     0.02      
App-istio-system/istio-pilot-66795bbf77-njdvg/discovery                                                     10.19     0.02      
App-cluster-api-system/cluster-api-controller-manager-0/manager                                             9.35      0.01      
App-turbonomic/auth-7c86754889-qbdmj/auth                                                                   8.72      0.39      
App-turbonomic/consul-7d8dfc9dc9-vpqzn/consul                                                               8.26      0.01      
App-turbonomic/prometheus-node-exporter-jxr4n/prometheus-node-exporter                                      7.84      0.01      
App-turbonomic/ml-datastore-c75fbff7b-jrqgz/ml-datastore                                                    7.79      0.28      
App-turbonomic/repository-768fbcc754-kcl8b/repository                                                       7.61      0.41      
App-turbonomic/zookeeper-55769465c8-l496v/zookeeper                                                         7.56      0.07      
App-kube-system/kube-proxy-27558/kube-proxy                                                                 7.51      0.02      
App-turbonomic/api-7bf4c9dfb8-284z4/api                                                                     7.41      0.44      
App-turbonomic/group-74cc4db8d6-wf7x9/group                                                                 7.37      0.39      
App-istio-system/istio-ingressgateway-cb48c6ffc-lprsc/istio-proxy                                           7.01      0.02      
App-istio-system/istio-policy-85469bdfbb-jmw77/istio-proxy                                                  6.92      0.02      
```
* Get a list of VMs that belong to a kubernetes cluster sorted by VMEM
```
$ lemurctl get vm --sort VMEM
NAME                                        VCPU (MHz)          VMEM (GB)    
meng-xl-machinedeployment-78d5dd8785-2gws6  971.53 [21.12%]     5.77 [70.68%]
meng-xl-machinedeployment-78d5dd8785-8v4ss  364.28 [7.92%]      5.54 [67.85%]
meng-xl-machinedeployment-78d5dd8785-cvzpr  375.34 [8.16%]      3.84 [47.01%]
meng-xl-machinedeployment-78d5dd8785-bjgtl  216.55 [4.71%]      3.65 [44.72%]
meng-xl-machinedeployment-78d5dd8785-2vltq  254.30 [5.53%]      3.12 [38.16%]
meng-xl-machinedeployment-78d5dd8785-x8brm  304.19 [6.61%]      3.02 [36.95%]
controlplane-0                              447.38 [9.32%]      2.02 [24.76%]
```
* View the supply chain summary originating from applications
```
$ lemurctl get app --supply-chain
TYPE                     COUNT                    PROVIDERS                     CONSUMERS                     
APPLICATION              45                       CONTAINER                                                   
CONTAINER                45                       CONTAINER_POD                 APPLICATION                   
CONTAINER_POD            42                       VIRTUAL_MACHINE               CONTAINER                     
VIRTUAL_MACHINE          7                        COMPUTE_TIER,STORAGE_TIER     CONTAINER_POD,APPLICATION     
COMPUTE_TIER             2                                                      VIRTUAL_MACHINE               
STORAGE_TIER             1                                                      VIRTUAL_MACHINE               
```
* View the supply chain starting from an application
```
$ lemurctl get app App-istio-system/prometheus-776fdf7479-4c6mw/prometheus --supply-chain
TYPE                     COUNT                    PROVIDERS                     CONSUMERS                     
APPLICATION              1                        CONTAINER                                                   
NAME                                              VCPU (MHz)                    VMEM (GB)                     
*stem/prometheus-776fdf7479-4c6mw/prometheus      69.38                         0.57                          

TYPE                     COUNT                    PROVIDERS                     CONSUMERS                     
CONTAINER                1                        CONTAINER_POD                 APPLICATION                   
NAME                                              VCPU (MHz)                    VMEM (GB)                     
*stem/prometheus-776fdf7479-4c6mw/prometheus      69.38 [1.51%]                 0.57 [6.99%]                  

TYPE                     COUNT                    PROVIDERS                     CONSUMERS                     
CONTAINER_POD            1                        VIRTUAL_MACHINE               CONTAINER                     
NAME                                              VCPU (MHz)                    VMEM (GB)                     
istio-system/prometheus-776fdf7479-4c6mw          69.38 [1.51%]                 0.57 [6.99%]                  

TYPE                     COUNT                    PROVIDERS                     CONSUMERS                     
VIRTUAL_MACHINE          1                        COMPUTE_TIER,STORAGE_TIER     APPLICATION,CONTAINER_POD     
NAME                                              VCPU (MHz)                    VMEM (GB)                     
meng-xl-machinedeployment-78d5dd8785-2gws6        1152.54 [25.05%]              5.26 [64.44%]                 

TYPE                     COUNT                    PROVIDERS                     CONSUMERS                     
COMPUTE_TIER             1                                                      VIRTUAL_MACHINE               
NAME                                              
m4.large                                          

TYPE                     COUNT                    PROVIDERS                     CONSUMERS                     
STORAGE_TIER             1                                                      VIRTUAL_MACHINE               
NAME                                              
GP2                                              

```
## Coming Soon
* The ability to auto-discover influxdb service without user having to manually provide the information
* The ability to manage (list/add/remove) targets
* The ability to execute actions
* The ability to show a meaningful display name of a cluster
