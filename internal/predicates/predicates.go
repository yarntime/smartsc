package predicates

import (
	"k8s.io/kubernetes/pkg/api/v1"
)

type PredicateFailureReason string

type FitPredicate func(pod *v1.Pod, node string) (bool, []PredicateFailureReason, error)

type PredicateFactory struct {
	FitMap map[string]FitPredicate
}

func NewPriorityFactory() *PredicateFactory {
	f := PredicateFactory{
		FitMap: map[string]FitPredicate{
			"nodestatuspredicate": FitPredicate(NodeStatusPredicate),
			"nodealarmpredicate":  FitPredicate(NodeAlarmPredicate),
		},
	}
	return &f
}

func NodeStatusPredicate(pod *v1.Pod, node string) (bool, []PredicateFailureReason, error) {
	return true, nil, nil
}

func NodeAlarmPredicate(pod *v1.Pod, node string) (bool, []PredicateFailureReason, error) {
	return true, nil, nil
}
