package services

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Chabuduo04/mate/back_end/models"
)

type RoleService struct {
	roles map[string]*models.Role
}

func NewRoleService(path string) (*RoleService, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var roles []*models.Role
	if err := json.Unmarshal(b, &roles); err != nil {
		return nil, err
	}
	rs := &RoleService{
		roles: make(map[string]*models.Role),
	}
	for _, r := range roles {
		rs.roles[r.ID] = r
	}
	return rs, nil
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
