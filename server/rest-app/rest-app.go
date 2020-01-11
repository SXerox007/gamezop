package main

import (
	"context"
	"crypto/x509"
	"flag"
	"gamezop/protos"
	"log"
	"net/http"
	"path"

	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type myService struct{}

var (
	authpoint    = flag.String("auth_end_points", "localhost:50051", "expose gamezop end point ")
	demoCertPool *x509.CertPool
)

func newServer() *myService {
	return new(myService)
}

func ExposePoint(address string, opts ...runtime.ServeMuxOption) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	err := gamezoppb.RegisterGamezopServiceHandlerFromEndpoint(ctx, mux, *authpoint, dialOpts)
	if err != nil {
		return err
	}
	grpcMux := http.NewServeMux()
	grpcMux.Handle("/", mux)
	// serveSwagger(grpcMux)
	grpcMux.HandleFunc("/swagger/", serveSwagger)
	log.Println("Starting Endpoint Exposed Server: localhost:5051")
	http.ListenAndServe(address, grpcMux)
	return nil
}

func main() {
	Init()
}

// initilization
func Init() {
	if err := ExposePoint(":5051"); err != nil {
		log.Fatal("Error in serve", err)
	}
}

// serve swagger
func serveSwagger(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join("swagger/third_party/swagger-new/swagger-ui/", p)
	http.ServeFile(w, r, p)
}
