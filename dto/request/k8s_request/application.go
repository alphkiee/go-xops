package k8s_request

type ApplicationReq struct {
	ID        string `json:"id"`
	NameSpace string `json:"namespace"`
	Name      string `json:"name"`
	Format    string `json:"format"`
}

type DeploymentReq struct {
	DeployName    string            `json:"deployment"`
	Replicas      int32             `json:"replicas"`
	MatchLabels   map[string]string `json:"matchlabels"`
	Labels        map[string]string `json:"labels"`
	ImageName     string            `json:"imagename"`
	Image         string            `json:"Image"`
	NameSpace     string            `json:"namespace"`
	ContainerPort uint              `json:"containerport"`
}

type UpDeploymentReq struct {
	Replicas       int32  `json:"replicas"`
	NameSpace      string `json:"namespace"`
	DeploymentName string `json:"deploymentname"`
	Image          string `json:"image"`
	ContainerPort  uint   `json:"containerport"`
}
