package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- helpers ----

func getValidToken(t *testing.T, router interface{ ServeHTTP(http.ResponseWriter, *http.Request) }) string {
	t.Helper()
	body, _ := json.Marshal(LoginRequest{Username: "alice", Password: "password123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	token, ok := resp["token"].(string)
	require.True(t, ok, "token must be a string")
	return token
}

// ---- GenerateToken / ParseToken unit tests ----

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken("alice")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, "alice", claims.Username)
	assert.Equal(t, "lab10-go-service", claims.Issuer)
}

func TestParseToken_InvalidSignature(t *testing.T) {
	_, err := ParseToken("this.is.garbage")
	assert.Error(t, err)
}

func TestParseToken_WrongSecret(t *testing.T) {
	// Manually assembled token whose signature was made with a different key
	badToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
		"eyJ1c2VybmFtZSI6ImFsaWNlIn0." +
		"aaaaabbbbbcccccdddddeeeeefffff00000111112"
	_, err := ParseToken(badToken)
	assert.Error(t, err)
}

// ---- /login endpoint ----

func TestLogin_ValidCredentials(t *testing.T) {
	router := SetupRouter()
	body, _ := json.Marshal(LoginRequest{Username: "alice", Password: "password123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["token"])
}

func TestLogin_WrongPassword(t *testing.T) {
	router := SetupRouter()
	body, _ := json.Marshal(LoginRequest{Username: "alice", Password: "wrong"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_UnknownUser(t *testing.T) {
	router := SetupRouter()
	body, _ := json.Marshal(LoginRequest{Username: "unknown", Password: "abc"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_MissingFields(t *testing.T) {
	router := SetupRouter()
	body, _ := json.Marshal(map[string]string{"username": "alice"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- /protected endpoint ----

func TestProtected_WithValidToken(t *testing.T) {
	router := SetupRouter()
	token := getValidToken(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "welcome")
}

func TestProtected_NoToken(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtected_InvalidToken(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not.a.real.token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProfile_WithValidToken(t *testing.T) {
	router := SetupRouter()
	token := getValidToken(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "alice", resp["username"])
}
