package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"user-manager/api"
	"user-manager/database"
	"user-manager/dto"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	var err error
	var cleanup func()
	server, cleanup, err := connectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer cleanup()

	r.Route("/users", server.UserRouter)

	fmt.Println("Test Server is Running")
	ts = httptest.NewServer(r)

	code := m.Run()

	os.Exit(code)
}

func TestUserLifeCycle(t *testing.T) {
	t.Run("Create", CreateUserTest)
	t.Run("Get All", GetUsersTest)
	t.Run("Get Single", GetUserTest)
	t.Run("Update", UpdateUserTest)
	t.Run("Delete", DeleteUserTest)
}

func GetUsersTest(t *testing.T) {
	resp, err := ts.Client().Get(ts.URL + "/users")
	if err != nil {
		log.Fatal("Can not call users endpoint")
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for Get Users List. Received %d", resp.StatusCode)
	}
}

func CreateUserTest(t *testing.T) {
	// test create user
	user := dto.User{
		Firstname: "jay",
		Lastname:  "vas",
		Email:     "jay@gmail.com",
		Phone:     "+0722134567",
		Age:       30,
		Status:    string(database.UserstatusActive),
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Can not create request by parsing json")
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/users", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Can not create request")
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := ts.Client().Do(req)
	if err != nil {
		log.Fatal("Can not call create user endpoint")
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201 for Create User. Received %d", resp.StatusCode)
	}
}

func GetUserTest(t *testing.T) {
	// test /users/id GET endpoint
	resp, err := ts.Client().Get(ts.URL + "/users/1")
	if err != nil {
		log.Fatal("Can not call users/id endpoint")
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for Get a User. Received %d", resp.StatusCode)
	}
}

func UpdateUserTest(t *testing.T) {
	// test update user
	user := dto.User{
		Firstname: "jay",
		Lastname:  "vas",
		Email:     "jay@gmail.com",
		Phone:     "+0722134567",
		Age:       35,
		Status:    string(database.UserstatusActive),
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Can not update request by parsing json")
	}

	req, err := http.NewRequest(http.MethodPatch, ts.URL+"/users/1", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Can not create update request")
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := ts.Client().Do(req)
	if err != nil {
		log.Fatal("Can not call update user endpoint")
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for Update User. Received %d", resp.StatusCode)
	}
}

func DeleteUserTest(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/users/1", nil)
	if err != nil {
		log.Fatal("Can not create delete request")
	}

	resp, err := ts.Client().Do(req)
	if err != nil {
		log.Fatal("Can not call delete user endpoint")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for Delete a User. Received %d", resp.StatusCode)
	}

}

func connectDatabase() (*api.Server, func(), error) {
	ctx := context.Background()

	fmt.Println("starting test containers")

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "usermanager_testdb",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, err
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://postgres:password@%s:%s/usermanager_testdb?sslmode=disable", host, port.Port())
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Fatal("Could not connect to test DB")
	}
	fmt.Println("connected to test db")

	queries := database.New(pool)
	server := api.NewServer(queries, pool)

	schema, err := os.ReadFile("./test_schema.sql")
	if err != nil {
		log.Fatal("Could not read the Schema file", err)
	}

	_, err = pool.Exec(ctx, string(schema))
	if err != nil {
		log.Fatal("Could not execute DB Schema", err)
	}

	cleanup := func() {
		pool.Close()
		err := container.Terminate(ctx)
		if err != nil {
			log.Fatal("error on terminating the test container")
		}
	}

	return server, cleanup, nil
}
