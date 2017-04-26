package controller

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/kubernetes/pkg/api/v1"
	clientv1 "k8s.io/client-go/pkg/api/v1"
	model1 "k8s.io/kubernetes/plugin/pkg/scheduler/api/v1"
	api "smartsc/internal/api"
	"smartsc/internal/cache"
	"smartsc/internal/job"
	"smartsc/internal/predicates"
	"smartsc/internal/priorities"
	utils "smartsc/utils"
)

type Controller struct {
	k8sClient     *utils.K8sClient
	fitFactory    *predicates.PredicateFactory
	priFactory    *priorities.PriorityFactory
	jobController *job.JobController
}

func NewController() *Controller {
	cache.Cache = api.SchedulerInfoCache{
		Infos: make(map[string]*api.SchedulerInfo),
	}

	c := &Controller{
		k8sClient:     utils.NewK8sClint(),
		fitFactory:    predicates.NewPriorityFactory(),
		priFactory:    priorities.NewPriorityFactory(),
		jobController: job.NewJobController(),
	}

	go c.watchNodes()

	return c
}

// Register registers this to the provided container
func (rr *Controller) Register(container *restful.Container) {

	ws := new(restful.WebService)
	ws.Path("/api/v1/").
		Doc("smart scheduler").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// GET /api/v1/fliter
	ws.Route(ws.GET("/filter").To(rr.filterNode).
		Doc("filter nodes").
		Operation("filterNodes").
		Writes(model1.ExtenderFilterResult{}))
	container.Add(ws)

	// GET /api/v1/prioritize
	ws.Route(ws.GET("/prioritize").To(rr.prioritizeNode).
		Doc("prioritize nodes").
		Operation("prioritizeNodes").
		Writes(model1.HostPriorityList{}))
	container.Add(ws)

	// GET /api/v1/feedback
	ws.Route(ws.GET("/feedback").To(rr.feedBack).
		Doc("feedbacks from training jobs").
		Operation("feedBack"))

	container.Add(ws)
}

func (rr *Controller) watchNodes() {
	listOptions := clientv1.ListOptions{}

	nw, err := rr.k8sClient.WatchNodes(listOptions)
	if err != nil {
		fmt.Printf("failed to watch nodes\n")
		return
	}
	defer nw.Stop()

	for {
		select {
		case t := <-nw.ResultChan():
			node := t.Object.(*clientv1.Node)
			if t.Type == watch.Added {
				rr.jobController.StartTrainingJob(node.Name)
			}
			if t.Type == watch.Modified {
				rr.jobController.UpdateTrainingJob(node.Name)
			}
			if t.Type == watch.Deleted {
				rr.jobController.DeleteTrainingJob(node.Name)
			}
		}
	}
}

func (rr *Controller) filterNode(req *restful.Request, res *restful.Response) {
	args := model1.ExtenderArgs{}
	if err := req.ReadEntity(&args); err != nil {
		errorResponse(res, errFailToReadResponse)
		return
	}

	failedNodes := make(map[string]string)

	var filtered []v1.Node
	var error string

	pod := args.Pod
	for i := 0; i < len(args.Nodes.Items); i++ {
		nodeName := args.Nodes.Items[i].Name
		for _, predicate := range rr.fitFactory.FitMap {
			fit, reasons, err := predicate(&pod, nodeName)
			if err != nil {
				error = error + err.Error()
			}
			if !fit {
				failedNodes[nodeName] = string(reasons[0])
			} else {
				filtered = append(filtered, args.Nodes.Items[i])
			}
		}
	}

	result := model1.ExtenderFilterResult{
		Nodes: v1.NodeList{
			Items: filtered,
		},
		FailedNodes: failedNodes,
		Error:       error,
	}

	if err := res.WriteEntity(result); err != nil {
		errorResponse(res, errFailToWriteResponse)
	}
}

func (rr *Controller) prioritizeNode(req *restful.Request, res *restful.Response) {
	args := model1.ExtenderArgs{}
	if err := req.ReadEntity(&args); err != nil {
		errorResponse(res, errFailToReadResponse)
		return
	}

	result := model1.HostPriorityList{}

	pod := args.Pod
	for i := 0; i < len(args.Nodes.Items); i++ {
		nodeName := args.Nodes.Items[i].Name
		score := 0
		for k, priority := range rr.priFactory.PriMap {
			pri, err := priority(&pod, nodeName)
			if err != nil {
				fmt.Printf("failed to calculate score for priority " + k)
				continue
			}
			score += pri.Score
		}
		result = append(result, model1.HostPriority{
			Host:  nodeName,
			Score: score,
		})
	}

	if err := res.WriteEntity(result); err != nil {
		errorResponse(res, errFailToWriteResponse)
	}
}

func (rr *Controller) feedBack(req *restful.Request, res *restful.Response) {
	feedbackList := api.FeedBackList{}
	if err := req.ReadEntity(&feedbackList); err != nil {
		errorResponse(res, errFailToReadResponse)
		return
	}

	for i := 0; i < len(feedbackList.FeedBacks); i++ {
		node := feedbackList.FeedBacks[i].Node
		cacheInfo := cache.Cache.Infos[node]
		cacheInfo.Mutex.Lock()
		defer cacheInfo.Mutex.Unlock()
		cacheInfo.Predicates = feedbackList.FeedBacks[i].Predicates
		cacheInfo.Priorities = feedbackList.FeedBacks[i].Priorities
	}
}
