package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://127.0.0.1:8080"

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestRegisterUser(t *testing.T) {
	userData := UserData{
		Username: "user_test",
		Password: "password228",
	}
	body, _ := json.Marshal(userData)

	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLoginUser(t *testing.T) {
	userData := UserData{
		Username: "user_test",
		Password: "password228",
	}
	body, _ := json.Marshal(userData)

	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
}

func TestCreateTask(t *testing.T) {
	// Login to get auth token
	userData := UserData{
		Username: "user_test",
		Password: "password228",
	}
	body, _ := json.Marshal(userData)

	loginResp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	var loginResponse map[string]string
	err = json.NewDecoder(loginResp.Body).Decode(&loginResponse)
	assert.NoError(t, err)
	token := loginResponse["token"]

	// Create a task
	req, _ := http.NewRequest("POST", baseURL+"/task", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "value")
}
