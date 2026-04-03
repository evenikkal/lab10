package main

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "task_v1_grpc/pb"
)

const bufSize = 1024 * 1024

// newTestServer spins up an in-process gRPC server using bufconn.
func newTestServer(t *testing.T) pb.GreeterClient {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, NewGreeterServer())

	go func() {
		if err := s.Serve(lis); err != nil {
			// ignore server-closed errors during test cleanup
		}
	}()

	t.Cleanup(func() { s.Stop() })

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return pb.NewGreeterClient(conn)
}

func TestSayHello_WithName(t *testing.T) {
	client := newTestServer(t)

	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Evgenia"})

	require.NoError(t, err)
	assert.Equal(t, "Hello, Evgenia!", resp.GetMessage())
}

func TestSayHello_EmptyName(t *testing.T) {
	client := newTestServer(t)

	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: ""})

	require.NoError(t, err)
	assert.Equal(t, "Hello, stranger!", resp.GetMessage())
}

func TestSayHello_MultipleRequests(t *testing.T) {
	client := newTestServer(t)

	names := []string{"Alice", "Bob", "Charlie"}
	for _, name := range names {
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		require.NoError(t, err)
		assert.Contains(t, resp.GetMessage(), name)
	}
}

// TestGreeterServer_Unit tests the handler directly, without any gRPC transport.
func TestGreeterServer_Unit(t *testing.T) {
	srv := NewGreeterServer()

	reply, err := srv.SayHello(context.Background(), &pb.HelloRequest{Name: "World"})
	require.NoError(t, err)
	assert.Equal(t, "Hello, World!", reply.GetMessage())
}
