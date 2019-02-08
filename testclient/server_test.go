package server

import (
	"context"
	fmt "fmt"
	"io"
	"os"
	"strings"
	"time"

	"./configstoreExample"

	"google.golang.org/grpc"

	"github.com/rs/xid"

	"testing"

	"gotest.tools/assert"
)

var ctx context.Context
var client configstoreExample.UserServiceClient

func TestMain(m *testing.M) {
	conn, err := grpc.Dial("127.0.0.1:13389", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println()
		return
	}
	defer conn.Close()

	ctx = context.Background()
	client = configstoreExample.NewUserServiceClient(conn)
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	resp, err := client.Create(ctx, &configstoreExample.CreateUserRequest{
		Entity: &configstoreExample.User{
			Id:           "",
			EmailAddress: "hello@example.com",
			PasswordHash: "what",
		},
	})
	assert.NilError(t, err)
	assert.Assert(t, resp.Entity.Id != "")
	assert.Equal(t, resp.Entity.EmailAddress, "hello@example.com")
	assert.Equal(t, resp.Entity.PasswordHash, "what")
}

func TestList(t *testing.T) {
	_, err := client.List(ctx, &configstoreExample.ListUserRequest{
		Limit: 10,
	})
	assert.NilError(t, err)
}

func TestCreateThenGet(t *testing.T) {
	resp, err := client.Create(ctx, &configstoreExample.CreateUserRequest{
		Entity: &configstoreExample.User{
			Id:           "",
			EmailAddress: "hello@example.com",
			PasswordHash: "what",
		},
	})
	assert.NilError(t, err)
	assert.Assert(t, resp.Entity.Id != "")
	assert.Equal(t, resp.Entity.EmailAddress, "hello@example.com")
	assert.Equal(t, resp.Entity.PasswordHash, "what")

	resp2, err := client.Get(ctx, &configstoreExample.GetUserRequest{
		Id: resp.Entity.Id,
	})
	assert.NilError(t, err)
	assert.Equal(t, resp2.Entity.Id, resp.Entity.Id)
	assert.Equal(t, resp2.Entity.EmailAddress, "hello@example.com")
	assert.Equal(t, resp2.Entity.PasswordHash, "what")
}

func TestWatchThenCreate(t *testing.T) {
	watcher, err := client.Watch(ctx, &configstoreExample.WatchUserRequest{})
	assert.NilError(t, err)

	mutex := make(chan bool, 1)
	timeout := make(chan bool, 1)

	testID := xid.New()

	var watchError error
	go func() {
		for {
			change, err := watcher.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				watchError = err
			}
			if change.Type == configstoreExample.WatchEventType_Created &&
				change.Entity.PasswordHash == testID.String() {
				mutex <- true
			}
		}
	}()

	resp, err := client.Create(ctx, &configstoreExample.CreateUserRequest{
		Entity: &configstoreExample.User{
			Id:           "",
			EmailAddress: "hello@example.com",
			PasswordHash: testID.String(),
		},
	})
	assert.NilError(t, err)
	assert.Assert(t, resp.Entity.Id != "")
	assert.Equal(t, resp.Entity.EmailAddress, "hello@example.com")
	assert.Equal(t, resp.Entity.PasswordHash, testID.String())

	go func() {
		time.Sleep(20 * time.Second)
		timeout <- true
	}()

	select {
	case <-mutex:
		assert.NilError(t, watchError)
	case <-timeout:
		assert.Assert(t, false, "timed out waiting for watch event")
	}
}

func TestCreateThenUpdateThenGet(t *testing.T) {
	resp, err := client.Create(ctx, &configstoreExample.CreateUserRequest{
		Entity: &configstoreExample.User{
			Id:           "",
			EmailAddress: "hello@example.com",
			PasswordHash: "what",
		},
	})
	assert.NilError(t, err)
	assert.Assert(t, resp.Entity.Id != "")
	assert.Equal(t, resp.Entity.EmailAddress, "hello@example.com")
	assert.Equal(t, resp.Entity.PasswordHash, "what")

	resp.Entity.EmailAddress = "update@example.com"

	resp2, err := client.Update(ctx, &configstoreExample.UpdateUserRequest{
		Entity: resp.Entity,
	})
	assert.NilError(t, err)
	assert.Equal(t, resp2.Entity.Id, resp.Entity.Id)
	assert.Equal(t, resp2.Entity.EmailAddress, "update@example.com")
	assert.Equal(t, resp2.Entity.PasswordHash, "what")

	resp3, err := client.Get(ctx, &configstoreExample.GetUserRequest{
		Id: resp.Entity.Id,
	})
	assert.NilError(t, err)
	assert.Equal(t, resp3.Entity.Id, resp2.Entity.Id)
	assert.Equal(t, resp3.Entity.EmailAddress, "update@example.com")
	assert.Equal(t, resp3.Entity.PasswordHash, "what")
}

func TestCreateThenDeleteThenGet(t *testing.T) {
	resp, err := client.Create(ctx, &configstoreExample.CreateUserRequest{
		Entity: &configstoreExample.User{
			Id:           "",
			EmailAddress: "hello@example.com",
			PasswordHash: "what",
		},
	})
	assert.NilError(t, err)
	assert.Assert(t, resp.Entity.Id != "")
	assert.Equal(t, resp.Entity.EmailAddress, "hello@example.com")
	assert.Equal(t, resp.Entity.PasswordHash, "what")

	resp2, err := client.Delete(ctx, &configstoreExample.DeleteUserRequest{
		Id: resp.Entity.Id,
	})
	assert.NilError(t, err)
	assert.Equal(t, resp2.Entity.Id, resp.Entity.Id)
	assert.Equal(t, resp2.Entity.EmailAddress, "hello@example.com")
	assert.Equal(t, resp2.Entity.PasswordHash, "what")

	_, err = client.Get(ctx, &configstoreExample.GetUserRequest{
		Id: resp.Entity.Id,
	})
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(fmt.Sprintf("%v", err), "code = NotFound"))
}
