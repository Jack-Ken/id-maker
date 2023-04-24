package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"id-maker/internal/entity"
	"id-maker/internal/usecase"
	"log"
	"net/http"
)

type tagRequest struct {
	BizTag string `json:"biz_tag"`
	MaxId  int64  `json:"max_id"`
	Step   int64  `json:"step"`
	Remark string `json:"remark"`
}

type SegmentService struct {
	s  usecase.Segment
	lg *zap.Logger
}

func NewSegmentService(s usecase.Segment, l *zap.Logger) *SegmentService {
	return &SegmentService{s, l}
}

func (r *SegmentService) pong(c *gin.Context) {
	r.lg.Info("http - v1 - Ping")
	successResponse(c, http.StatusOK, "pong")
	//c.JSON(http.StatusOK, "pong")
}

func (r *SegmentService) GetId(c *gin.Context) {
	var (
		tag string
		id  int64
		err error
	)

	//tag = c.Param("tag")
	tag = c.Query("tag")
	if tag == "" {
		errorResponse(c, http.StatusBadRequest, "tag cannot empty")

		return
	}
	log.Print(tag)

	if id, err = r.s.GetId(tag); err != nil {
		r.lg.Error("http - v1 - GetId\n", zap.Error(err))
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}
	successIdResponse(c, http.StatusOK, id)
	//c.JSON(http.StatusOK, id)
}

func (r *SegmentService) GetSnowId(c *gin.Context) {
	id := r.s.SnowFlakeGetId()
	r.lg.Info("http - v1 - GetSnowId")
	successIdResponse(c, http.StatusOK, id)
	//c.JSON(http.StatusOK, id)
}

func (r *SegmentService) CreateTag(c *gin.Context) {
	var request tagRequest
	if err := c.ShouldBind(&request); err != nil {
		r.lg.Error("http - v1 - CreateTag\n", zap.Error(err))
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}
	r.lg.Info("http - v1 - CreateTag " + request.BizTag)

	tagSegment := entity.Segments{
		BizTag: request.BizTag,
		MaxId:  request.MaxId,
		Step:   request.Step,
		Remark: request.Remark,
	}
	if err := r.s.CreateTag(&tagSegment); err != nil {
		r.lg.Error("http - v1 - CreateTag\n", zap.Error(err))
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}
	successResponse(c, http.StatusOK, "success")
	//c.JSON(http.StatusOK, nil)
}
