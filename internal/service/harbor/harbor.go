package harbor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/mittwald/goharbor-client/v4/apiv1/model"
	"go-xops/assets/harbor"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"io/ioutil"
	"net/http"
)

type Projects struct {
	P []harbor.Project
}

type RP struct {
	project    []harbor.Project
	repository []harbor.Repository
}

type ProjectsListReq struct {
	Name         string          `json:"name"`
	CreationTime strfmt.DateTime `json:"creation_time"`
	RepoCount    int64           `json:"repo_count"`
	OwnerName    string          `json:"owner_name"`
}

type ProjectDeleteReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProjectCreateReq struct {
	Name         string `json:"name"`
	CountLimit   int    `json:"count_limit"`
	StorageLimit int    `json:"storage_limit"`
}

type ProjectUpdate struct {
	ProjectID         int64  `json:"project_id"`
	OwnerID           int64  `json:"owner_id"`
	CreationTime      string `json:"creation_time"`
	UpdateTime        string `json:"update_time"`
	Delete            bool   `json:"delete"`
	OwnerName         string `json:"owner_name"`
	CurrentUserRoleID int64  `json:"current_user_role_id"`
	Name              string `json:"name"`
	RepoCount         int64  `json:"repo_count"`
	CountLimit        int    `json:"count_limit"`
	StorageLimit      int    `json:"storage_limit"`
}

// http 方式获取project
func GetProjectsList() (Projects, error) {
	var b Projects
	response, err := http.Get(common.Conf.HarborApi.Url + "projects")
	if err != nil || response.StatusCode != http.StatusOK {
		return b, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &b.P)
	if err != nil {
		return b, err
	}
	return b, nil
}

// 使用第三方插件获取project
func GetProjectList() (Projects, error) {
	var b Projects
	var c harbor.Project
	project, err := common.HarborClient.ListProjects(context.TODO(), "")
	if err != nil {
		return b, err
	}
	for _, v := range project {
		c.Name = v.Name
		c.Delete = v.Deleted
		c.ProjectID = int64(v.ProjectID)
		c.Delete = v.Deleted
		c.OwnerName = v.OwnerName
		c.OwnerID = int64(v.OwnerID)
		c.RepoCount = v.RepoCount
		c.CurrentUserRoleID = v.CurrentUserRoleID
		c.CreationTime = v.CreationTime
		c.UpdateTime = v.UpdateTime
		b.P = append(b.P, c)
	}
	return b, nil
}

// 新增project
func CreateProject(name string, countLimit int, storageLimit int) (harbor.Project, error) {

	var p harbor.Project

	// harbor project 不允许重复
	ps, err := GetProjectList()
	if err != nil {
		panic(err)
	}

	//name是否存在ps.P里面
	t := utils.ContainsStr(ps.P, name)
	if t {
		p.Name = name
		return p, errors.New("name is exist!")
	} else {
		project, err := common.HarborClient.NewProject(context.TODO(), name, countLimit, storageLimit)
		if err != nil {
			return p, err
		}
		p.Name = project.Name
		p.Delete = project.Deleted
		p.ProjectID = int64(project.ProjectID)
		p.Delete = project.Deleted
		p.OwnerName = project.OwnerName
		p.OwnerID = int64(project.OwnerID)
		p.RepoCount = project.RepoCount
		p.CurrentUserRoleID = project.CurrentUserRoleID
		p.CreationTime = project.CreationTime
		p.UpdateTime = project.UpdateTime
	}
	return p, nil
}

// 删除project
func DeleteProject(p harbor.Project) error {
	m := model.Project{}
	m.Name = p.Name
	// m.Deleted = p.Delete
	m.ProjectID = int32(p.ProjectID)
	m.OwnerID = int32(p.OwnerID)
	m.CurrentUserRoleID = p.CurrentUserRoleID
	m.RepoCount = p.RepoCount
	m.OwnerName = p.OwnerName
	err := common.HarborClient.DeleteProject(context.TODO(), &m)
	if err != nil {
		return err
	}
	return nil
}

// 修改project
func UpdateProject(p harbor.Project, countLimit, storageLimit int) error {
	m := model.Project{}
	m.Name = p.Name
	m.Deleted = p.Delete
	m.ProjectID = int32(p.ProjectID)
	m.OwnerID = int32(p.OwnerID)
	m.CurrentUserRoleID = p.CurrentUserRoleID
	m.RepoCount = p.RepoCount
	m.OwnerName = p.OwnerName
	err := common.HarborClient.UpdateProject(context.TODO(), &m, countLimit, storageLimit)
	if err != nil {
		return err
	}
	return nil
}

func GetRegistry(name string) (harbor.Registry, error) {
	var R harbor.Registry
	r, err := common.HarborClient.GetRegistry(context.TODO(), name)
	if err != nil {
		return R, err
	}
	R.Status = r.Status
	R.Name = r.Name
	R.UpdateTime = r.UpdateTime
	R.ID = r.ID
	R.Type = r.Type
	R.URL = r.URL
	R.Insecure = r.Insecure
	R.Description = r.Description
	R.CreationTime = r.CreationTime

	return R, nil
}

// 获取project里面的所有镜像列表
func GetSearch(name string) (map[string]interface{}, error) {

	var m map[string]interface{}
	response, err := http.Get("http://192.168.212.176/api/search?q=" + name)
	if err != nil || response.StatusCode != http.StatusOK {
		return m, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}
