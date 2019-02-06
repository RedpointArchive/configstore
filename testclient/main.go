package main

import (
	"context"
	fmt "fmt"

	"./configstoreExample"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:13389", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer conn.Close()

	ctx := context.Background()

	client := configstoreExample.NewProjectServiceClient(conn)
	resp, err := client.List(ctx, &configstoreExample.ListProjectRequest{})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%s", resp)
}
