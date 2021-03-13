package k8s_request

type ApplicationReq struct {
	ID        string `json:"id"`
	NameSpace string `json:"namespace"`
	Name      string `json:"name"`
	Format    string `json:"format"`
}
