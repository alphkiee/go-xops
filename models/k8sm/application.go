package k8sm

type Application struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Format     string      `json:"format"`
	Containers []Container `json:"containers"`
}
