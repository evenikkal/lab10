"""
Tests for the Python gRPC client (task V1).
All gRPC calls are mocked — no running Go server required.
grpc-dependent tests are skipped if grpcio is not installed.
"""

import unittest
from unittest.mock import MagicMock, patch

import hello_pb2

try:
    import grpc as _grpc
    import hello_pb2_grpc
    from client import run, DEFAULT_ADDRESS
    GRPC_AVAILABLE = True
except ImportError:
    GRPC_AVAILABLE = False


def make_stub_mock(reply_message: str = "Hello, Evgenia!"):
    stub = MagicMock()
    stub.SayHello.return_value = hello_pb2.HelloReply(message=reply_message)
    return stub


class TestHelloPb2(unittest.TestCase):
    """Unit tests for the generated protobuf message classes (no grpc needed)."""

    def test_hello_request_serialization(self):
        req = hello_pb2.HelloRequest(name="Alice")
        data = req.SerializeToString()
        restored = hello_pb2.HelloRequest()
        restored.ParseFromString(data)
        self.assertEqual(restored.name, "Alice")

    def test_hello_reply_serialization(self):
        reply = hello_pb2.HelloReply(message="Hello, Alice!")
        data = reply.SerializeToString()
        restored = hello_pb2.HelloReply()
        restored.ParseFromString(data)
        self.assertEqual(restored.message, "Hello, Alice!")

    def test_hello_request_empty_name(self):
        req = hello_pb2.HelloRequest()
        self.assertEqual(req.name, "")

    def test_hello_reply_empty_message(self):
        reply = hello_pb2.HelloReply()
        self.assertEqual(reply.message, "")

    def test_round_trip_unicode(self):
        req = hello_pb2.HelloRequest(name="Евгения")
        restored = hello_pb2.HelloRequest()
        restored.ParseFromString(req.SerializeToString())
        self.assertEqual(restored.name, "Евгения")

    def test_multiple_requests_independent(self):
        r1 = hello_pb2.HelloRequest(name="Alice")
        r2 = hello_pb2.HelloRequest(name="Bob")
        self.assertNotEqual(r1.name, r2.name)


@unittest.skipUnless(GRPC_AVAILABLE, "grpcio not installed — run: pip install grpcio")
class TestGreeterClient(unittest.TestCase):
    """Integration tests for client.py using a mocked gRPC stub."""

    @patch("client.hello_pb2_grpc.GreeterStub")
    @patch("client.grpc.insecure_channel")
    def test_run_calls_say_hello(self, mock_channel, mock_stub_cls):
        mock_stub = make_stub_mock()
        mock_stub_cls.return_value = mock_stub
        mock_channel.return_value.__enter__ = MagicMock(return_value=MagicMock())
        mock_channel.return_value.__exit__ = MagicMock(return_value=False)
        run(DEFAULT_ADDRESS)
        self.assertTrue(mock_stub.SayHello.called)

    @patch("client.hello_pb2_grpc.GreeterStub")
    @patch("client.grpc.insecure_channel")
    def test_run_sends_evgenia(self, mock_channel, mock_stub_cls):
        mock_stub = make_stub_mock()
        mock_stub_cls.return_value = mock_stub
        mock_channel.return_value.__enter__ = MagicMock(return_value=MagicMock())
        mock_channel.return_value.__exit__ = MagicMock(return_value=False)
        run(DEFAULT_ADDRESS)
        sent_names = [c.args[0].name for c in mock_stub.SayHello.call_args_list]
        self.assertIn("Evgenia", sent_names)

    @patch("client.hello_pb2_grpc.GreeterStub")
    @patch("client.grpc.insecure_channel")
    def test_run_handles_rpc_error_gracefully(self, mock_channel, mock_stub_cls):
        rpc_error = _grpc.RpcError()
        rpc_error.code = lambda: _grpc.StatusCode.UNAVAILABLE
        rpc_error.details = lambda: "connection refused"
        mock_stub = MagicMock()
        mock_stub.SayHello.side_effect = rpc_error
        mock_stub_cls.return_value = mock_stub
        mock_channel.return_value.__enter__ = MagicMock(return_value=MagicMock())
        mock_channel.return_value.__exit__ = MagicMock(return_value=False)
        try:
            run(DEFAULT_ADDRESS)
        except Exception as e:
            self.fail(f"run() raised unexpectedly: {e}")


if __name__ == "__main__":
    unittest.main()
