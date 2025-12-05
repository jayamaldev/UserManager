package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"user-manager/database"
	"user-manager/dto"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateUser(ctx context.Context, user dto.User, q database.Querier) (*database.User, string, int) {
	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		var errs validator.ValidationErrors
		var validationFiledErr validator.FieldError
		errors.As(err, &errs)
		for _, validationError := range errs {
			validationFiledErr = validationError
		}
		msg := "Validation Failed on: " + formatError(validationFiledErr)
		fmt.Println(msg)
		return nil, msg, http.StatusBadRequest
	}

	dbUser, err := q.CreateUser(ctx, database.CreateUserParams{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Phone:     pgtype.Text{String: user.Phone, Valid: true},
		Age:       pgtype.Int4{Int32: user.Age, Valid: true},
		UserStatus: database.NullUserstatus{
			Userstatus: database.Userstatus(user.Status),
			Valid:      true,
		},
	})

	if err != nil {
		return nil, "Internal Server Error", http.StatusInternalServerError
	}
	return &dbUser, "", http.StatusCreated
}

func UpdateUser(ctx context.Context, id int, user dto.User, q database.Querier) (string, int) {
	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		var errs validator.ValidationErrors
		var validationFiledErr validator.FieldError
		errors.As(err, &errs)
		for _, validationError := range errs {
			validationFiledErr = validationError
		}
		msg := "Validation Failed on: " + formatError(validationFiledErr)
		fmt.Println(msg)
		return msg, http.StatusBadRequest
	}

	updateErr := q.UpdateUser(ctx, database.UpdateUserParams{
		Userid:    int32(id),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Phone:     pgtype.Text{String: user.Phone, Valid: true},
		Age:       pgtype.Int4{Int32: user.Age, Valid: true},
		UserStatus: database.NullUserstatus{
			Userstatus: database.Userstatus(user.Status),
			Valid:      true,
		},
	})

	if updateErr != nil {
		return "Internal Server Error", http.StatusInternalServerError
	}
	return "", http.StatusOK
}

func formatError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is a required field", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s long", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s long", err.Field(), err.Param())
	case "gt":
		return fmt.Sprintf("%s must be positive value", err.Field())
	case "e164":
		return fmt.Sprintf("%s must be a valid phone number", err.Field())
	default:
		return fmt.Sprintf("%s failed validation with tag %s", err.Field(), err.Tag())
	}
}
