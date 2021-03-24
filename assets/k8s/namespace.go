package k8s

type NameSpace struct {
	Name   string `json:"name"`
	UID    string `json:"uid"`
	Status string `json:"status"`
}
