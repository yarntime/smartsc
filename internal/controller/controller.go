package controller

import (
      "github.com/emicklei/go-restful"
	utils "smartsc/utils"
	model "k8s.io/kubernetes/plugin/pkg/scheduler/api/v1"
)

type Controller struct {
	k8sClient *utils.K8sClient
}

func NewController() *Controller {
	return &Controller{
		k8sClient: utils.NewK8sClint(),
	}
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
	Writes(model.ExtenderFilterResult{}))
	container.Add(ws)

	// GET /api/v1/prioritize
	ws.Route(ws.GET("/prioritize").To(rr.prioritizeNode).
	Doc("prioritize nodes").
	Operation("prioritizeNodes").
	Writes(model.HostPriorityList{}))
	container.Add(ws)
}

func (rr *Controller) filterNode(req *restful.Request, res *restful.Response) {

}

func (rr *Controller) prioritizeNode(req *restful.Request, res *restful.Response) {

}