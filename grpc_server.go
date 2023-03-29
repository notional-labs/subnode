package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
)

func StartGrpcServer() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		outCtx := metadata.NewOutgoingContext(ctx, md.Copy()) // Copy the inbound metadata explicitly.
		if ok {
			prunedNode := SelectPrunedNodeGrpc()
			selectedHost := prunedNode.Backend.Grpc // default to pruned node

			// Decide on which backend to dial
			m0 := md["x-cosmos-block-height"]
			if len(m0) > 0 {
				xCosmosBlockHeight := m0[0]
				fmt.Printf("xCosmosBlockHeight=%s\n", xCosmosBlockHeight)

				height, err := strconv.ParseInt(xCosmosBlockHeight, 10, 64)
				if err != nil {
					return nil, nil, status.Errorf(codes.InvalidArgument, "Invalid x-cosmos-block-height")
				}

				node, err := SelectMatchedNodeGrpc(height)
				if err != nil {
					return nil, nil, status.Errorf(codes.InvalidArgument, "Invalid x-cosmos-block-height")
				}

				selectedHost = node.Backend.Grpc
			} else {
				selectedHost = prunedNode.Backend.Grpc
			}

			// no SSL: grpc.WithTransportCredentials(insecure.NewCredentials())
			conn, err := grpc.DialContext(ctx, selectedHost, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
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
