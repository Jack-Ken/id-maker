package router

import (
	"github.com/gin-gonic/gin"
	"id-maker/internal/initialize"
)

type Router interface {
	Route(r *gin.Engine)
}

var routers []Router

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(initialize.GinLogger(), initialize.GinRecovery(true))
	for _, router := range routers {
		router.Route(r)
	}
	return r
}

// 路由注册

func Register(ro ...Router) {
	routers = append(routers, ro...)
}
