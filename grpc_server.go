package main

import (
	"context"
	"crypto/tls"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

func StartGrpcServer() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		outCtx := metadata.NewOutgoingContext(ctx, md.Copy()) // Copy the inbound metadata explicitly.
		if ok {
			// Decide on which backend to dial

			// no SSL: grpc.WithTransportCredentials(insecure.NewCredentials())
			conn, err := grpc.DialContext(ctx, "grpc-osmosis-ia.cosmosia.notional.ventures:443", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
			return outCtx, conn, err
		}
		return nil, nil, status.Errorf(codes.Unimplemented, "Unknown method")
	}

	srv := grpc.NewServer(grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}

	srv.Serve(lis)
}
