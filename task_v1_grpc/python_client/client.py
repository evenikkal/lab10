"""
Python gRPC client for the Go Greeter service (task V1).

Usage:
    python3 client.py [host:port]

The Go gRPC server must be running on the given address (default localhost:50051).
Start it with:  cd ../go_server && go run main.go
"""

import sys
import grpc

import hello_pb2
import hello_pb2_grpc


DEFAULT_ADDRESS = "localhost:50051"


def run(address: str = DEFAULT_ADDRESS) -> None:
    """Connect to the Greeter service and call SayHello."""
    with grpc.insecure_channel(address) as channel:
        stub = hello_pb2_grpc.GreeterStub(channel)

        names = ["Evgenia", "World", ""]
        for name in names:
            request = hello_pb2.HelloRequest(name=name)
            try:
                response = stub.SayHello(request, timeout=5)
                print(f"[SayHello] name={name!r:12} -> {response.message}")
            except grpc.RpcError as e:
                print(f"[SayHello] RPC error: {e.code()} – {e.details()}")


if __name__ == "__main__":
    addr = sys.argv[1] if len(sys.argv) > 1 else DEFAULT_ADDRESS
    run(addr)
