package rpc

import (
	pb "CDcoding2333/scaffold/apps/rpc/proto"
	"context"
	"fmt"
)

// Handler ...
type Handler struct {
}

// Ping ...
func (h *Handler) Ping(ctx context.Context, in *pb.PingReq) (*pb.PingResp, error) {
	return &pb.PingResp{
		Pong: fmt.Sprintf("pong:%s", in.Msg),
	}, nil
}
