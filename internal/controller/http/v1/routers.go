package v1

import (
	"go.uber.org/zap"
	"id-maker/internal/controller/http/router"
	"id-maker/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

type V1Router struct {
	s  usecase.Segment
	lg *zap.Logger
}

func NewV1Router(s usecase.Segment, lg *zap.Logger) *V1Router {
	return &V1Router{s, lg}
}

func RegisterRouteSrv(s usecase.Segment, lg *zap.Logger) {
	log.Println("init user router")
	router.Register(NewV1Router(s, lg)) // V1Router是实现了Router接口的新的路由
}

func (v *V1Router) Route(r *gin.Engine) {
	// 此处申请路由
	h := NewSegmentService(v.s, v.lg) // HandlerUser实现了每一个业务模块对应的处理函数
	base := r.Group("/v1")
	{
		base.GET("/ping", h.pong)
		base.GET("/id", h.GetId)
		base.GET("/snowid", h.GetSnowId)
		base.POST("/tag", h.CreateTag)
	}
}
