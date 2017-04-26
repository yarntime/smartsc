package job

type JobController struct {
}

func NewJobController() *JobController {
	return &JobController{}
}

// use template to start training job in kubernetes cluster,
// maybe one job corresponding to one node(or to many nodes)
func (jc *JobController) StartTrainingJob(node string) {

}

func (jc *JobController) UpdateTrainingJob(node string) {

}

func (jc *JobController) DeleteTrainingJob(node string) {

}
