package api

import (
	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/gin-gonic/gin"
)

type RoleApi struct {
	Service services.RoleService
}

func NewRoleApi(service services.RoleService) RoleApi {
	return RoleApi{Service: service}
}

func (r *RoleApi) ListRoles(c *gin.Context) {
	roles := r.Service.ListRoles()
	c.JSON(200, gin.H{"roles": roles})
}
