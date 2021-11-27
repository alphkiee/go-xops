package docker

type Image struct {
	Repository []string `json:"repository"`
	ImageID    string   `json:"image_id"`
	Created    int64    `json:"created"`
	Size       int64    `json:"size"`
	Tag        []string `json:"tag"`
}
