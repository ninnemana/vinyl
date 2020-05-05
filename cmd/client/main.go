package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ninnemana/vinyl/pkg/vinyl"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	conn, err := grpc.DialContext(context.Background(), ":8081", grpc.WithInsecure())
	// conn, err := grpc.Dial(":8080")
	if err != nil {
		log.Fatalf("failed to dial service: %v", err)
	}
	defer conn.Close()

	svc := vinyl.NewVinylClient(conn)
	client, err := svc.Search(context.Background(), &vinyl.SearchParams{
		Artist: os.Args[1],
	})
	if err != nil {
		log.Fatalf("failed to create search scanner: %v", err)
	}

	for {
		msg, err := client.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("failed to receive release: %v\n", err)
			break
		}

		fmt.Printf("Message: %+v\n", msg)
	}
}

func usage() {
	fmt.Printf("Vinyltap.io CLI\n----------\n\n")
	fmt.Printf("\t Sub Commands\n")
	fmt.Printf("\t\t search\n")
}
