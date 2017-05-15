package job

import (
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	utils "smartsc/utils"
	"strconv"
	"k8s.io/client-go/pkg/util"
)

const (
	ContainerNamePrefix = "training-job"
	BaseImage           = "reg.skycloud.com:5000/sky-firmament/training:v1"
	TotalContainers     = 0
	Threshold           = 1
)

type JobController struct {
	k8sClient *utils.K8sClient
}

func NewJobController() *JobController {
	return &JobController{}
}

// maybe we need volumes config later
func componentDeployment(container v1.Container) *v1beta1.Deployment {
	return &v1beta1.Deployment{
		TypeMeta: unversioned.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "Deployment",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      container.Name,
			Namespace: "sky-firmament",
			Labels:    map[string]string{"component": container.Name, "tier": "training-job"},
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: util.Int32Ptr(1),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{container},
				},
			},
		},
	}
}

func componentResources(cpu string) v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceName(v1.ResourceCPU): resource.MustParse(cpu),
		},
	}
}

// use template to start training job in kubernetes cluster,
// maybe one job corresponding to one node(or to many nodes)
func (jc *JobController) StartTrainingJob(node string) {
	dep := componentDeployment(v1.Container{
		Name:      ContainerNamePrefix + strconv.Itoa(TotalContainers),
		Image:     BaseImage,
		Command:   []string{"/training"},
		Args:      []string{"--hosts=" + node},
		Resources: componentResources("500m")})

	jc.k8sClient.CreateDeployment(dep)
}

// not sure what to do.
func (jc *JobController) UpdateTrainingJob(node string) {

}

// delete the corresponding deployments
func (jc *JobController) DeleteTrainingJob(node string) {

}
