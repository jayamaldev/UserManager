package services

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"user-manager/database"
	"user-manager/dto"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateUserInvalidFirstName(t *testing.T) {
	user := dto.User{
		Firstname: "a",
		Lastname:  "abd",
		Email:     "jay@gmail.com",
		Phone:     "+0722134567",
		Age:       30,
		Status:    string(database.UserstatusActive),
	}

	_, msg, status := CreateUser(t.Context(), user, nil)
	fmt.Println("error message: ", msg, " status: ", status)

	if status != http.StatusBadRequest {
		t.Errorf("Test Failure! Incorrect status")
	}
}

func TestCreateUserInvalidLastName(t *testing.T) {
	user := dto.User{
		Firstname: "abc",
		Lastname:  "a",
		Email:     "jay@gmail.com",
		Phone:     "+0722134567",
		Age:       30,
		Status:    string(database.UserstatusActive),
	}

	_, msg, status := CreateUser(t.Context(), user, nil)
	fmt.Println("error message: ", msg, " status: ", status)

	if status != http.StatusBadRequest {
		t.Errorf("Test Failure! Incorrect status")
	}
}

func TestCreateUserInvalidEmail(t *testing.T) {
	user := dto.User{
		Firstname: "abc",
		Lastname:  "ajuuoi",
		Email:     "jaygmail.com",
		Phone:     "+0722134567",
		Age:       30,
		Status:    string(database.UserstatusActive),
	}

	_, msg, status := CreateUser(t.Context(), user, nil)
	fmt.Println("error message: ", msg, " status: ", status)

	if status != http.StatusBadRequest {
		t.Errorf("Test Failure! Incorrect status")
	}
}

func TestCreateUserInvalidPhone(t *testing.T) {
	user := dto.User{
		Firstname: "abc",
		Lastname:  "anhgd",
		Email:     "jay@gmail.com",
		Phone:     "722134567",
		Age:       30,
		Status:    string(database.UserstatusActive),
	}

	_, msg, status := CreateUser(t.Context(), user, nil)
	fmt.Println("error message: ", msg, " status: ", status)

	if status != http.StatusBadRequest {
		t.Errorf("Test Failure! Incorrect status")
	}
}

func TestUpdateUserInvalidAge(t *testing.T) {
	user := dto.User{
		Firstname: "abc",
		Lastname:  "abdsd",
		Email:     "jay@gmail.com",
		Phone:     "+0722134567",
		Age:       -30,
		Status:    string(database.UserstatusActive),
	}

	msg, status := UpdateUser(t.Context(), 2, user, nil)
	fmt.Println("error message: ", msg, " status: ", status)

	if status != http.StatusBadRequest {
		t.Errorf("Test Failure! Incorrect status")
	}
}

func TestCreateUserSuccess(t *testing.T) {
	user := dto.User{
		Firstname: "Jay",
		Lastname:  "Vas",
		Email:     "jay@gmail.com",
		Phone:     "+0722134567",
		Age:       30,
		Status:    string(database.UserstatusActive),
	}

	mockDb := &MockDb{}

	_, msg, status := CreateUser(t.Context(), user, mockDb)
	fmt.Println("error message: ", msg, " status: ", status)

	if status != http.StatusCreated {
		t.Errorf("Test Failure! Incorrect status")
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	user := dto.User{
		Firstname: "Jay",
		Lastname:  "Vas",
		Email:     "jay@gmail.com",
		Phone:     "+0722134567",
		Age:       30,
		Status:    string(database.UserstatusActive),
	}

	mockDb := &MockDb{}

	msg, status := UpdateUser(t.Context(), 1, user, mockDb)
	fmt.Println("error message: ", msg, " status: ", status)

	if status != http.StatusOK {
		t.Errorf("Test Failure! Incorrect status")
	}
}

type MockDb struct {
}

func (m *MockDb) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	dbUser := database.User{
		Userid:    1,
		Firstname: "Jay",
		Lastname:  "Vas",
		Email:     "jay@gmail.com",
		Phone:     pgtype.Text{String: "+0722134567", Valid: true},
		Age:       pgtype.Int4{Int32: int32(30), Valid: true},
		UserStatus: database.NullUserstatus{
			Userstatus: database.UserstatusActive,
			Valid:      true,
		},
	}

	return dbUser, nil
}

func (m *MockDb) UpdateUser(ctx context.Context, arg database.UpdateUserParams) error {
	return nil
}
