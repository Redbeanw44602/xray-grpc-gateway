package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	gwLogger "xray-grpc-gateway/gen/app/log/command"
	gwObservatory "xray-grpc-gateway/gen/app/observatory/command"
	gwHandler "xray-grpc-gateway/gen/app/proxyman/command"
	gwRouting "xray-grpc-gateway/gen/app/router/command"
	gwStats "xray-grpc-gateway/gen/app/stats/command"
)

type registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

var grpcGatewayRegistry = []registerFunc{
	gwLogger.RegisterLoggerServiceHandlerFromEndpoint,
	gwObservatory.RegisterObservatoryServiceHandlerFromEndpoint,
	gwHandler.RegisterHandlerServiceHandlerFromEndpoint,
	gwRouting.RegisterRoutingServiceHandlerFromEndpoint,
	gwStats.RegisterStatsServiceHandlerFromEndpoint,
}

var (
	grpcServerEndpoint = flag.String("endpoint", "localhost:10010", "Xray grpc server endpoint")
	listenPort         = flag.String("listen", ":10020", "Proxy server listening port")
)

func RegisterAllServices(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	for _, register := range grpcGatewayRegistry {
		if err := register(ctx, mux, endpoint, opts); err != nil {
			return err
		}
	}
	return nil
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := RegisterAllServices(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(*listenPort, mux)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
