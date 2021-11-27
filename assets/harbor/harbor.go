package harbor

type Project struct {
	ProjectID         int64  `json:"project_id"`
	OwnerID           int64  `json:"owner_id"`
	CreationTime      string `json:"creation_time"`
	UpdateTime        string `json:"update_time"`
	Delete            bool   `json:"delete"`
	OwnerName         string `json:"owner_name"`
	CurrentUserRoleID int64  `json:"current_user_role_id"`
	Name              string `json:"name"`
	RepoCount         int64  `json:"repo_count"`
}

type Registry struct {
	CreationTime string `json:"creation_time,omitempty"`
	Description  string `json:"description,omitempty"`
	ID           int64  `json:"id,omitempty"`
	Insecure     bool   `json:"insecure,omitempty"`
	Name         string `json:"name,omitempty"`
	Status       string `json:"status,omitempty"`
	Type         string `json:"type,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
	URL          string `json:"url,omitempty"`
}

// 镜像列表字段
type Repository struct {
	ProjectID      int64  `json:"project_id"`
	ProjectName    string `json:"project_name"`
	Project_public string `json:"project___public"`
	PullCount      int64  `json:"pull_count"`
	RepositoryName string `json:"repository_name"`
	TagCount       int64  `json:"tag_count"`
}
