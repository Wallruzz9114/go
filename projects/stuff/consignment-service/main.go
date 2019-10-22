package main

// Import the generated protobuf code
import (
	"context"
	"log"
	"net"
	"sync"

	// Import the generated protobuf code
	pb "github.com/Wallruzz9114/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

// Create a new consignment
func (repository *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repository.mu.Lock()
	updated := append(repository.consignments, consignment)
	repository.consignments = updated
	repository.mu.Unlock()

	return consignment, nil
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repository repository
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (service *service) CreateConsignment(context context.Context, request *pb.Consignment) (*pb.Response, error) {
	// Save our consignment
	consignment, err := service.repository.Create(request)

	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repository := &Repository{}

	// Set-up our gRPC server.
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterShippingServiceServer(server, &service{repository})

	// Register reflection service on gRPC server.
	reflection.Register(server)

	log.Println("Running on port:", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
