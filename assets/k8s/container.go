package k8s

type Container struct {
	Name       string     `json:"name"`
	Image      string     `json:"image"`
	Version    string     `json:"version"`
	Deployment Deployment `json:"deployment"`
}
