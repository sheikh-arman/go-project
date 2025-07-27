package pkg

import (
	"context"
	pb "github.com/sheikh-arman/go-project/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

func Client() {
	conn, err := grpc.DialContext(context.TODO(), "localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	client := pb.NewHelloWorldServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	i := 0
	for {
		resp, err := client.SayHello(ctx, &pb.RequestHelloWorld{Hello: "Hello World " + strconv.Itoa(i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Println(resp.GetWorld())
		i++
	}

}
