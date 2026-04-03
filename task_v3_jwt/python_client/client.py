"""
Python client for the Go JWT service (task V3).

Demonstrates:
  1. Calling /login to obtain a JWT from the Go service.
  2. Using the token to call a protected endpoint.
  3. Independently verifying the token with PyJWT (shared secret).

Usage:
    python3 client.py

The Go service must be running:
    cd ../go_service && go run main.go
"""

import sys
import requests
import jwt  # PyJWT

GO_SERVICE_URL = "http://localhost:8083"
JWT_SECRET = "super-secret-lab10-key"
JWT_ALGORITHM = "HS256"


def login(username: str, password: str) -> str:
    """POST /login and return the JWT token string."""
    response = requests.post(
        f"{GO_SERVICE_URL}/login",
        json={"username": username, "password": password},
        timeout=10,
    )
    response.raise_for_status()
    data = response.json()
    return data["token"]


def call_protected(token: str) -> dict:
    """GET /protected with the Bearer token."""
    response = requests.get(
        f"{GO_SERVICE_URL}/protected",
        headers={"Authorization": f"Bearer {token}"},
        timeout=10,
    )
    response.raise_for_status()
    return response.json()


def call_profile(token: str) -> dict:
    """GET /profile with the Bearer token."""
    response = requests.get(
        f"{GO_SERVICE_URL}/profile",
        headers={"Authorization": f"Bearer {token}"},
        timeout=10,
    )
    response.raise_for_status()
    return response.json()


def verify_token_locally(token: str) -> dict:
    """
    Verify the JWT locally using the shared secret.
    This mirrors what the Go service does on every protected request.
    """
    payload = jwt.decode(
        token,
        JWT_SECRET,
        algorithms=[JWT_ALGORITHM],
        options={"verify_exp": True},
    )
    return payload


if __name__ == "__main__":
    username = sys.argv[1] if len(sys.argv) > 1 else "alice"
    password = sys.argv[2] if len(sys.argv) > 2 else "password123"

    print(f"Logging in as '{username}'...")
    try:
        token = login(username, password)
        print(f"Token received: {token[:40]}...")

        print("\nVerifying token locally with PyJWT...")
        claims = verify_token_locally(token)
        print(f"Claims: {claims}")

        print("\nCalling /protected endpoint...")
        protected_data = call_protected(token)
        print(f"Response: {protected_data}")

        print("\nCalling /profile endpoint...")
        profile_data = call_profile(token)
        print(f"Profile: {profile_data}")

    except requests.HTTPError as e:
        print(f"HTTP error: {e.response.status_code} – {e.response.text}")
    except jwt.InvalidTokenError as e:
        print(f"JWT verification failed: {e}")
