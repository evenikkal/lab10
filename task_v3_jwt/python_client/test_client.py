"""
Tests for the Python JWT client (task V3).
HTTP calls are mocked; JWT verification uses the real PyJWT library.
"""

import time
import unittest
from unittest.mock import patch, MagicMock

import jwt
import requests as req

from client import (
    login, call_protected, call_profile,
    verify_token_locally,
    JWT_SECRET, JWT_ALGORITHM,
)


def make_real_token(username: str = "alice", expired: bool = False) -> str:
    """Create a real HS256 JWT using the same secret as the Go service."""
    now = int(time.time())
    payload = {
        "username": username,
        "iss": "lab10-go-service",
        "iat": now,
        "exp": now - 10 if expired else now + 3600,
    }
    return jwt.encode(payload, JWT_SECRET, algorithm=JWT_ALGORITHM)


class TestVerifyTokenLocally(unittest.TestCase):
    """Test PyJWT verification against tokens the Go service would produce."""

    def test_verify_valid_token(self):
        token = make_real_token("alice")
        claims = verify_token_locally(token)
        self.assertEqual(claims["username"], "alice")
        self.assertEqual(claims["iss"], "lab10-go-service")

    def test_verify_expired_token_raises(self):
        token = make_real_token("alice", expired=True)
        with self.assertRaises(jwt.ExpiredSignatureError):
            verify_token_locally(token)

    def test_verify_wrong_secret_raises(self):
        bad_token = jwt.encode({"username": "alice"}, "wrong-secret", algorithm="HS256")
        with self.assertRaises(jwt.InvalidSignatureError):
            verify_token_locally(bad_token)

    def test_verify_garbage_raises(self):
        with self.assertRaises(jwt.DecodeError):
            verify_token_locally("this.is.garbage")

    def test_verify_contains_username(self):
        token = make_real_token("bob")
        claims = verify_token_locally(token)
        self.assertEqual(claims["username"], "bob")


class TestLogin(unittest.TestCase):

    @patch("client.requests.post")
    def test_login_success(self, mock_post):
        token = make_real_token("alice")
        mock_response = MagicMock()
        mock_response.json.return_value = {"token": token, "expires_in": 3600}
        mock_response.raise_for_status = MagicMock()
        mock_post.return_value = mock_response

        result = login("alice", "password123")
        self.assertEqual(result, token)

    @patch("client.requests.post")
    def test_login_wrong_credentials_raises(self, mock_post):
        mock_response = MagicMock()
        mock_response.raise_for_status.side_effect = req.HTTPError("401 Unauthorized")
        mock_post.return_value = mock_response

        with self.assertRaises(req.HTTPError):
            login("alice", "wrong-password")

    @patch("client.requests.post")
    def test_login_sends_correct_body(self, mock_post):
        token = make_real_token()
        mock_response = MagicMock()
        mock_response.json.return_value = {"token": token, "expires_in": 3600}
        mock_response.raise_for_status = MagicMock()
        mock_post.return_value = mock_response

        login("alice", "password123")

        call_kwargs = mock_post.call_args.kwargs
        self.assertEqual(call_kwargs["json"]["username"], "alice")
        self.assertEqual(call_kwargs["json"]["password"], "password123")


class TestCallProtected(unittest.TestCase):

    @patch("client.requests.get")
    def test_call_protected_success(self, mock_get):
        mock_response = MagicMock()
        mock_response.json.return_value = {
            "message": "welcome to the protected zone",
            "username": "alice",
        }
        mock_response.raise_for_status = MagicMock()
        mock_get.return_value = mock_response

        result = call_protected("some-token")
        self.assertEqual(result["username"], "alice")

    @patch("client.requests.get")
    def test_call_protected_sends_bearer_header(self, mock_get):
        mock_response = MagicMock()
        mock_response.json.return_value = {}
        mock_response.raise_for_status = MagicMock()
        mock_get.return_value = mock_response

        call_protected("my-jwt-token")

        call_kwargs = mock_get.call_args.kwargs
        auth_header = call_kwargs["headers"]["Authorization"]
        self.assertTrue(auth_header.startswith("Bearer "))
        self.assertIn("my-jwt-token", auth_header)

    @patch("client.requests.get")
    def test_call_protected_no_token_raises(self, mock_get):
        mock_response = MagicMock()
        mock_response.raise_for_status.side_effect = req.HTTPError("401 Unauthorized")
        mock_get.return_value = mock_response

        with self.assertRaises(req.HTTPError):
            call_protected("")


class TestCallProfile(unittest.TestCase):

    @patch("client.requests.get")
    def test_call_profile_success(self, mock_get):
        mock_response = MagicMock()
        mock_response.json.return_value = {
            "username": "alice",
            "role": "user",
            "email": "alice@example.com",
        }
        mock_response.raise_for_status = MagicMock()
        mock_get.return_value = mock_response

        result = call_profile("token")
        self.assertEqual(result["email"], "alice@example.com")


if __name__ == "__main__":
    unittest.main()
