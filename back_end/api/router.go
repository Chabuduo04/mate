package api

import (
	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, svc *services.Services) {
	// user routing group
	u := NewUserApi(*svc.UserService)
	{
		userRouterGroup := r.Group("/user")
		userRouterGroup.POST("/register", u.UserRegister)
		userRouterGroup.POST("/login", u.UserLogin)
	}
	role := NewRoleApi(*svc.RoleService)
	{
		roleRouterGroup := r.Group("/role")
		roleRouterGroup.GET("/list", role.ListRoles)
	}
	v := NewVoiceApi(svc.TTS)
	{
		voiceRouterGroup := r.Group("/voice")
		voiceRouterGroup.GET("/list", v.VoiceList)
	}
	c := NewChatApi(*svc)
	api := r.Group("/api")
	{
		// protected routes - require JWT auth
		protected := api.Group("")
		protected.Use(AuthRequired())
		{
			protected.POST("/llm", c.ChatWithText)
			protected.POST("/voice-chat", c.ChatWithAudio)
		}
	}
}
