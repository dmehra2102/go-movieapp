package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"movieexample.com/gen"
	"movieexample.com/rating/internal/controller/rating"
	grpcHandler "movieexample.com/rating/internal/handler/grpc"
	"movieexample.com/rating/internal/repository/memory"
)

// var serviceName = "rating"

// func main() {
// 	var port int
// 	flag.IntVar(&port, "port", 8082, "API handler port")
// 	flag.Parse()
// 	log.Printf("Starting the rating service on port %d", port)
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
// 	ctrl := rating.New(repo)
// 	h := httphandler.New(ctrl)
// 	http.Handle("/rating", http.HandlerFunc(h.Handle))
// 	if err := http.ListenAndServe(":8082", nil); err != nil {
// 		panic(err)
// 	}
// }

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	svc := rating.New(repo)
	g := grpcHandler.New(svc)
	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, g)
	srv.Serve(lis)
}
