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

	client := configstoreExample.NewUserServiceClient(conn)
	resp, err := client.Get(ctx, &configstoreExample.GetUserRequest{
		Id: "lBpJE5piqLi7mk75cKwF",
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%s", resp)
	if resp != nil && resp.Entity != nil {
		fmt.Printf("%s", resp.Entity.EmailAddress)
	}
}
