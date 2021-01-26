package service

import (
	"go-xops/models/system"
	"go-xops/pkg/common"
)

// 获取符合条件的casbin规则, 按角色
func (s *MysqlService) GetRoleCasbins(c system.SysRoleCasbin) []system.SysRoleCasbin {
	e := common.Casbin
	policies := e.GetFilteredPolicy(0, c.Keyword, c.Path, c.Method)
	cs := make([]system.SysRoleCasbin, 0)
	for _, policy := range policies {
		cs = append(cs, system.SysRoleCasbin{
			Keyword: policy[0],
			Path:    policy[1],
			Method:  policy[2],
		})
	}
	return cs
}

// 创建一条casbin规则, 按角色
func (s *MysqlService) CreateRoleCasbin(c system.SysRoleCasbin) (bool, error) {
	e := common.Casbin
	return e.AddPolicy(c.Keyword, c.Path, c.Method)
}

// 批量创建多条casbin规则, 按角色
func (s *MysqlService) BatchCreateRoleCasbins(cs []system.SysRoleCasbin) (bool, error) {
	e := common.Casbin
	// 按角色构建
	rules := make([][]string, 0)
	for _, c := range cs {
		rules = append(rules, []string{
			c.Keyword,
			c.Path,
			c.Method,
		})
	}
	return e.AddPolicies(rules)
}

// 删除一条casbin规则, 按角色
func (s *MysqlService) DeleteRoleCasbin(c system.SysRoleCasbin) (bool, error) {
	e := common.Casbin
	return e.RemovePolicy(c.Keyword, c.Path, c.Method)
}

// 批量删除多条casbin规则, 按角色
func (s *MysqlService) BatchDeleteRoleCasbins(cs []system.SysRoleCasbin) (bool, error) {
	e := common.Casbin
	// 按角色构建
	rules := make([][]string, 0)
	for _, c := range cs {
		rules = append(rules, []string{
			c.Keyword,
			c.Path,
			c.Method,
		})
	}
	return e.RemovePolicies(rules)
}

// 根据权限编号读取casbin规则
func (s *MysqlService) GetCasbinListByRoleId(roleId uint) ([]system.SysCasbin, error) {
	casbins := make([]system.SysCasbin, 0)
	var role system.SysRole
	err := s.db.Where("id = ?", roleId).First(&role).Error
	if err != nil {
		return casbins, err
	}
	e := common.Casbin
	// 查询符合字段v0=role.Keyword所有casbin规则
	list := e.GetFilteredPolicy(0, role.Keyword)
	for _, v := range list {
		casbins = append(casbins, system.SysCasbin{
			PType: "p",
			V0:    v[0],
			V1:    v[1],
			V2:    v[2],
		})
	}
	return casbins, nil
}
