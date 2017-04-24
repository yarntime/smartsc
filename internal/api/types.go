package api

import "sync"

type SchedulerInfoCache struct {
	Infos map[string]*SchedulerInfo
}

type FeedBackList struct {
	FeedBacks []*SchedulerInfo `json:"feedbacks"`
}

type SchedulerInfo struct {
	sync.Mutex
	Node       string          `json:"node"`
	Predicates *PredicatesInfo `json:"predicates"`
	Priorities *PrioritiesInfo `json:"priorities"`
}

type PredicatesInfo struct {
	NodeStates bool `json:"nodeStatus"`
	NodeAlarms int  `json:"nodeAlarms"`
}

type PrioritiesInfo struct {
	CpuLoad    int `json:"cpuLoad"`
	MemoryLoad int `json:"memoryLoad"`
	DiskIO     int `json:"diskIO"`
	NetworkIO  int `json:"networkIO"`
}
