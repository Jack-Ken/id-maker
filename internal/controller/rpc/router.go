package rpc

import (
	"go.uber.org/zap"
	"id-maker/internal/usecase"
)

func NewRouter(s usecase.Segment, l *zap.Logger) {
	NewSegment(s, l)
}
