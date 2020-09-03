package main

import (
	pb "go-grpc/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type StreamService struct{}

const PORT = "9002"

func main() {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen error:%v", err)
	}
	server.Serve(lis) //nolint
}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for i := 0; i < 6; i++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.GetPt().GetName(),
				Value: r.GetPt().GetValue() + int32(i),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{
				Pt: &pb.StreamPoint{
					Name:  "gRPC Stream Server: Record",
					Value: 1,
				},
			})
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.GetPt().GetName(), r.GetPt().GetValue())
	}
	return nil
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gPRC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
	return nil
}
