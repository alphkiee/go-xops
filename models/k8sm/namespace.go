package k8sm

type NameSpace struct {
	Name   string `json:"name"`
	UID    string `json:"uid"`
	Status string `json:"status"`
}
