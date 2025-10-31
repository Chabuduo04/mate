package services

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Chabuduo04/mate/back_end/dao"
	"github.com/Chabuduo04/mate/back_end/db"
	"github.com/Chabuduo04/mate/back_end/models"
)

type RoleService struct {
	roles map[string]*models.Role
}

func NewRoleService(path string) (*RoleService, error) {
	rs := &RoleService{
		roles: make(map[string]*models.Role),
	}
	if err := rs.LoadRolesFromJSON(path); err != nil {
		return nil, err
	}
	return rs, nil
}

func (s *RoleService) LoadRolesFromJSON(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var roles []*models.Role
	if err := json.Unmarshal(b, &roles); err != nil {
		return err
	}

	for _, r := range roles {
		s.roles[r.ID] = r
	}

	dao.InsertRolesToDB(db.GetConn(), roles)
	return nil
}

func (s *RoleService) ListRoles() []*models.Role {
	out := make([]*models.Role, 0, len(s.roles))
	for _, v := range s.roles {
		out = append(out, v)
	}
	return out
}

func (s *RoleService) GetRole(id string) (*models.Role, bool) {
	r, ok := s.roles[id]
	return r, ok
}
