package middleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	res, err := handler(ctx, req)
	if err != nil {
		log.Println(err)

		if st, ok := status.FromError(err); ok {
			return nil, status.Errorf(st.Code(), st.Message())
		}
		
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return res, err
}
