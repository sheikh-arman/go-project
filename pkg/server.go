package pkg

import (
	"context"
	"fmt"
	pb "github.com/sheikh-arman/go-project/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	ctx      context.Context
	req      *pb.RequestHelloWorld
	resp     chan *pb.ResponseHelloWorld
	err      chan error
	workerID int
}

type server struct {
	pb.UnimplementedHelloWorldServiceServer
	tasks chan Task
	wg    sync.WaitGroup
}

func (s *server) startWorker(workerCount int) {
	for i := 0; i < workerCount; i++ {
		s.wg.Add(workerCount)
		go func(workerId int) {
			defer s.wg.Done()
			log.Printf("worker %d starting process", workerId)
			for task := range s.tasks {
				task.workerID = workerId
				s.processTask(task)
			}
			log.Printf("worker %d stopping process", workerId)
		}(i + 1)
	}
}

func (*server) processTask(task Task) {
	log.Println(time.Now(), " processing request", task.req.Hello)
	resp := &pb.ResponseHelloWorld{
		World: fmt.Sprintln("I've received your request: " + task.req.GetHello() + " worker id " + strconv.Itoa(task.workerID)),
	}
	select {
	case <-task.ctx.Done():
		task.err <- task.ctx.Err()
	default:
		task.resp <- resp
	}
}

func (s *server) SayHello(ctx context.Context, req *pb.RequestHelloWorld) (*pb.ResponseHelloWorld, error) {
	task := Task{
		ctx:  ctx,
		req:  req,
		resp: make(chan *pb.ResponseHelloWorld, 1),
		err:  make(chan error, 1),
	}
	select {
	case s.tasks <- task:
	case <-task.ctx.Done():
		return nil, task.ctx.Err()
	}

	select {
	case resp := <-task.resp:
		return resp, nil
	case err := <-task.err:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

}

func Server() {
	server := &server{
		tasks: make(chan Task, 100),
	}
	server.startWorker(100)
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloWorldServiceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
