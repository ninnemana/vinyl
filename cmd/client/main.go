package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ninnemana/vinyl/pkg/vinyl"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.DialContext(context.Background(), ":8080", grpc.WithInsecure())
	// conn, err := grpc.Dial(":8080")
	if err != nil {
		log.Fatalf("failed to dial service: %v", err)
	}
	defer conn.Close()

	svc := vinyl.NewVinylClient(conn)
	client, err := svc.Search(context.Background(), &vinyl.SearchParams{
		Artist: "Kendrick Lamar",
	})
	if err != nil {
		log.Fatalf("failed to create search scanner: %v", err)
	}

	for {
		msg, err := client.Recv()
		if err != nil {
			log.Printf("failed to receive release: %v\n", err)
			break
		}

		fmt.Println(msg.GetYear(), msg.GetArtist(), msg.GetTitle(), msg.GetType())
	}
}
