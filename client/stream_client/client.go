package main

import (
	"context"
	pb "go-grpc/proto"
	"io"
	"log"

	"google.golang.org/grpc"
)

const PORT = "9002"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial error :%v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)
	//err = printLists(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "lists", Value: 111}})
	//if err != nil {
	//	log.Fatalf("printLists error: %v", err)
	//}
	//err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "printRecord", Value: 222}})
	//if err != nil {
	//	log.Fatalf("printLists error: %v", err)
	//}
	err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "printRoute", Value: 333}})
	if err != nil {
		log.Fatalf("printLists error: %v", err)
	}
}

func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pj.name: %s, pt.value: %d", resp.GetPt().GetName(), resp.GetPt().GetValue())
	}

	return nil
}

func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < 6; i++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Fatalf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	return nil
}

func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n < 6; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}

		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pj.name: %s, pt.value: %d", response.Pt.Name, response.Pt.Value)
	}
	stream.CloseSend()

	return nil
}
