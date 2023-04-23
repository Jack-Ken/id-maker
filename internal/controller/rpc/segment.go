package rpc

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"id-maker/internal/controller/rpc/proto"
	"id-maker/internal/entity"
	"id-maker/internal/usecase"
	"id-maker/pkg/grpcserver"
	"net/http"
)

type SegmentRpc struct {
	proto.UnimplementedGidServer
	t  usecase.Segment
	lg *zap.Logger
}

func NewSegment(s usecase.Segment, l *zap.Logger) {
	proto.RegisterGidServer(grpcserver.RpcServer, &SegmentRpc{t: s, lg: l})
}

func (s *SegmentRpc) Ping(ctx context.Context, g *empty.Empty) (out *proto.PingReply, err error) {
	out = &proto.PingReply{
		Status: &proto.Status{
			Code: http.StatusOK,
		},
		Data: "pong",
	}
	s.lg.Info("grpc - Ping")
	return
}

func (s *SegmentRpc) GetId(ctx context.Context, in *proto.IdRequest) (out *proto.IdReply, err error) {
	var id int64

	out = &proto.IdReply{
		Status: &proto.Status{
			Code: http.StatusOK,
		},
	}

	tag := in.GetTag()
	if tag == "" {
		out.Status.Code = http.StatusBadRequest
		out.Status.Msg = "tag cannot empty"
		err = errors.New("tag cannot empty")
		s.lg.Error("grpc - GetId\n", zap.Error(err))
		return
	}

	if id, err = s.t.GetId(tag); err != nil {
		out.Status.Code = http.StatusInternalServerError
		out.Status.Msg = err.Error()
		s.lg.Error("grpc - GetId\n", zap.Error(err))
		return
	}
	out.Id = id
	return
}

func (s *SegmentRpc) GetSnowId(ctx context.Context, g *empty.Empty) (out *proto.SnowIdReply, err error) {
	out = &proto.SnowIdReply{
		Status: &proto.Status{
			Code: http.StatusOK,
		},
		Id: s.t.SnowFlakeGetId(),
	}
	s.lg.Info("grpc - GetSnowId")
	return

}

func (s *SegmentRpc) CreateTag(ctx context.Context, in *proto.CreateTagRequest) (out *proto.CreateTagReply, err error) {
	out = &proto.CreateTagReply{
		Status: &proto.Status{
			Code: http.StatusOK,
		},
	}
	if in.GetTag() == "" || in.GetStep() == 0 {
		out.Status.Code = http.StatusInternalServerError
		out.Status.Msg = "param error"
		err = errors.New("param error")
		s.lg.Error("grpc - CreateTag\n", zap.Error(err))
		return
	}

	if err = s.t.CreateTag(&entity.Segments{
		BizTag: in.GetTag(),
		MaxId:  in.GetMaxId(),
		Step:   in.GetStep(),
		Remark: in.GetRemark(),
	}); err != nil {
		out.Status.Code = http.StatusInternalServerError
		out.Status.Msg = err.Error()
		s.lg.Error("grpc - CreateTag\n", zap.Error(err))
		return
	}
	return
}
