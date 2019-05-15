# kubectl-plugin

build :
```
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/kubectl-view-node-resource cmd/view-node-resource/main.go

$ GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/kubectl-view-node-taints cmd/view-node-taints/main.go
```

start :
```
$ mv bin/kubectl-view-node-resource /usr/bin/  

$ mv bin/kubectl-view-node-taints /usr/bin/  
```

demo:
```
$ kubectl view node resource
 Name            PodCount  CPURequests  MemoryRequests  CPULimits     MemoryLimits
 192.168.1.110   4         0 (0.00%)    6.4 (41.26%)    8 (100.00%)   16.0 (103.14%)


$ kubectl view node taints
 Name            Status                       Age   Version                         Taints
 192.168.1.110   Ready,SchedulingDisabled     49d   v1.8.1-35+9406f9d9909c61-dirty  enabledDiskSchedule=true:NoSchedule
```


