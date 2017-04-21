package api

type SchedulerInfoCache struct {
	cache []*SchedulerInfo
}

type FeedBackList struct {
	FeedBacks []*SchedulerInfo `json:"feedbacks"`
}

type SchedulerInfo struct {
	Node string `json:"node"`
	Predicates *PredicatesInfo `json:"predicates"`
	Priorities *PrioritiesInfo `json:"priorities"`
}

type PredicatesInfo struct {
	NodeStates bool `json:"nodeStatus"`
	NodeAlarms int `json:"nodeAlarms"`
}

type PrioritiesInfo struct {
	CpuLoad float32 `json:"cpuLoad"`
	MemoryLoad float32 `json:"memoryLoad"`
	DiskIO  float32  `json:"diskIO"`
	NetworkIO float32 `json:"networkIO"`
}
