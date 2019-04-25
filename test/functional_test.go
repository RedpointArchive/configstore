package test

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"google.golang.org/grpc"

	"testing"

	"gotest.tools/assert"
)

func TestFunctional(t *testing.T) {
	// start server
	fmt.Printf("starting server...\n")
	server := exec.Command("powershell", "../server.ps1")
	server.Stdout = os.Stdout
	server.Stderr = os.Stderr
	err := server.Start()
	assert.NilError(t, err)
	defer server.Process.Kill()
	serverExited := false
	go func() {
		server.Wait()
		serverExited = true
		fmt.Printf("server exited.\n")
	}()

	// try to connect
	fmt.Printf("waiting for connection...\n")
	var connected bool
	for !connected && !serverExited {
		_, err := grpc.Dial("127.0.0.1:13389", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*1))
		if err != nil {
			continue
		}
		fmt.Printf("connected to server!\n")
		connected = true
		break
	}
	if serverExited {
		assert.Assert(t, false, "server exited before client could start")
	}

	// start client
	fmt.Printf("starting client...\n")
	client := exec.Command("powershell", "../client.ps1", "-SkipDeps")
	client.Stdout = os.Stdout
	client.Stderr = os.Stderr
	err = client.Run()
	assert.NilError(t, err)
}
