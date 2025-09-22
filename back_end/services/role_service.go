package services

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/Chabuduo04/mate/back_end/models"
)

type RoleService struct {
	roles map[string]*models.Role
	mu    sync.RWMutex
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*models.Role, 0, len(s.roles))
	for _, v := range s.roles {
		out = append(out, v)
	}
	return out
}

func (s *RoleService) GetRole(id string) (*models.Role, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.roles[id]
	return r, ok
}
