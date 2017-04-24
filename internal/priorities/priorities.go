package priorities

import (
	"k8s.io/kubernetes/pkg/api/v1"
	schedulerapi "k8s.io/kubernetes/plugin/pkg/scheduler/api"
	"smartsc/internal/cache"
)

var (
	Total_Share   = 10 * 4
	Cpu_Share     = 5
	Memory_Share  = 3
	DiskIO_Share  = 1
	Network_Share = 1
)

type PriorityMapFunction func(pod *v1.Pod, node string) (schedulerapi.HostPriority, error)

type PriorityFactory struct {
	PriMap map[string]PriorityMapFunction
}

func NewPriorityFactory() *PriorityFactory {
	f := PriorityFactory{
		PriMap: map[string]PriorityMapFunction{
			"cpuloadpriority":    PriorityMapFunction(LeastCpuLoadPriority),
			"memoryloadpriority": PriorityMapFunction(LeastMemoryLoadPriority),
			"diskiopriority":     PriorityMapFunction(LeastDiskIOPriority),
			"networkiopriority":  PriorityMapFunction(LeastNetworkIOPriority),
		},
	}
	return &f
}

func LeastCpuLoadPriority(pod *v1.Pod, node string) (schedulerapi.HostPriority, error) {
	cacheInfo := cache.Cache.Infos[node]
	return schedulerapi.HostPriority{
		Host:  node,
		Score: int((100 - cacheInfo.Priorities.CpuLoad) / Total_Share * Cpu_Share), //TODO: set each resources share Dynamically
	}, nil
}

func LeastMemoryLoadPriority(pod *v1.Pod, node string) (schedulerapi.HostPriority, error) {
	cacheInfo := cache.Cache.Infos[node]
	return schedulerapi.HostPriority{
		Host:  node,
		Score: int((100 - cacheInfo.Priorities.MemoryLoad) / Total_Share * Memory_Share),
	}, nil
}

func LeastDiskIOPriority(pod *v1.Pod, node string) (schedulerapi.HostPriority, error) {
	cacheInfo := cache.Cache.Infos[node]
	return schedulerapi.HostPriority{
		Host:  node,
		Score: int((100 - cacheInfo.Priorities.DiskIO) / Total_Share * DiskIO_Share),
	}, nil
}

func LeastNetworkIOPriority(pod *v1.Pod, node string) (schedulerapi.HostPriority, error) {
	cacheInfo := cache.Cache.Infos[node]
	return schedulerapi.HostPriority{
		Host:  node,
		Score: int((100 - cacheInfo.Priorities.NetworkIO) / Total_Share * Network_Share),
	}, nil
}
