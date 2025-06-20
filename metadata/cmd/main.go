package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller/metadata"
	grpcHandler "movieexample.com/metadata/internal/handler/grpc"
	"movieexample.com/metadata/internal/repository/memory"
)

// Below is the Implementation of HTTP with JSON communication

// const serviceName = "metadata"

// func main() {
// 	var port int
// 	flag.IntVar(&port, "port", 8081, "API handler port")
// 	flag.Parse()
// 	log.Printf("Starting the metadata service on port %d", port)

// 	registry, err := consul.NewRegistry("localhost:8500")
// 	if err != nil {
// 		panic(err)
// 	}

// 	ctx := context.Background()
// 	instanceID := discovery.GenerateInstanceID(serviceName)

// 	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
// 		panic(err)
// 	}

// 	go func() {
// 		for {
// 			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
// 				log.Println("Failed to report healthy state: " + err.Error())
// 			}
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()

// 	defer registry.Deregister(ctx, instanceID, serviceName)

// 	repo := memory.New()
// 	ctrl := metadata.New(repo)
// 	h := httphandler.New(ctrl)

// 	http.Handle("/metdata", http.HandlerFunc(h.GetMetadata))
// 	if err := http.ListenAndServe(":8081", nil); err != nil {
// 		panic(err)
// 	}
// }

// Below is the implementation of gRPC with Protobuf communicaiton b/w services
func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	srv.Serve(lis)

}
