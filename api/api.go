package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"user-manager/database"
	"user-manager/dto"
	services "user-manager/internal"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Queries *database.Queries
	Pool    *pgxpool.Pool
}

func NewServer(queries *database.Queries, pool *pgxpool.Pool) *Server {
	return &Server{
		Queries: queries,
		Pool:    pool,
	}
}

func (s *Server) UserRouter(r chi.Router) {
	r.Get("/", s.getUsers)
	r.Post("/", s.createUser)
	r.Get("/{id}", s.getUser)
	r.Patch("/{id}", s.updateUser)
	r.Delete("/{id}", s.deleteUser)
}

// @Summary Get all users
// @Description Retrieve a list of all users
// @Produce json
// @Success 200 {array} dto.User
// @Router /users [get]
func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Println("Get Users request received")
	users, err := s.Queries.ListUsers(ctx)
	if err != nil {
		fmt.Println("error on retrieving users list: ", err)
		http.Error(w, "Error on Returing All Users", http.StatusInternalServerError)
		return
	}

	fmt.Println("users list retrieved", users)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		fmt.Println("error on retrieving users list: ", err)
		http.Error(w, "Error on Returing All Users", http.StatusInternalServerError)
		return
	}
}

// @Summary Create a New User
// @Description  Create a New User
// @Accept json
// @Produce json
// @Param UserInput body dto.User true "User Details for Creation"
// @Success 200 {object} dto.User
// @Router /users [post]
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user dto.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, userError, httpstatus := services.CreateUser(ctx, user, s.Queries)
	if httpstatus != http.StatusCreated {
		fmt.Println(userError)
		http.Error(w, userError, httpstatus)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatus)
	fmt.Println("User Creating with name: " + user.Firstname)
	err = json.NewEncoder(w).Encode(dbUser)
	if err != nil {
		fmt.Println("error on  Creating User: ", err)
		http.Error(w, "Error on Creating User", http.StatusInternalServerError)
		return
	}
}

// @Summary Get single user
// @Description Retrieve a user
// @Produce json
// @Success 200 {object} dto.User
// @Router /users/id [get]
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")

	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Println("error on Returing User: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := s.Queries.GetUser(ctx, int32(id))
	if err != nil {
		fmt.Println("error on Returing User: ", err)
		http.Error(w, "Error on Returing User with id: ", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		fmt.Println("error on Returing User: ", err)
		http.Error(w, "Error on Returing User", http.StatusInternalServerError)
		return
	}
}

// @Summary Update existing User
// @Description Update existing User
// @Accept json
// @Produce json
// @Param UserInput body dto.User true "User Details for Update"
// @Success 200 {object} dto.User
// @Router /users/id [patch]
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	fmt.Println("Updating User with id: " + userId)

	id, strErr := strconv.Atoi(userId)
	if strErr != nil {
		fmt.Println("error on Updating User: ", strErr)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var user dto.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateErr, httpstatus := services.UpdateUser(ctx, id, user, s.Queries)

	if httpstatus != http.StatusOK {
		fmt.Println("Error on updating user: " + updateErr)
		http.Error(w, updateErr, httpstatus)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("User Updated with id: " + userId)
	err = json.NewEncoder(w).Encode("User Updated with id: " + userId)
	if err != nil {
		fmt.Println("error on Updating User: ", err)
		http.Error(w, "Error on Updating User", http.StatusInternalServerError)
		return
	}
}

// @Summary Delete existing User
// @Description Delete existing User
// @Accept json
// @Produce json
// @Success 200
// @Router /users/id [delete]
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")

	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id, strErr := strconv.Atoi(userId)
	if strErr != nil {
		fmt.Println("error on Updating User: ", strErr)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := s.Queries.DeleteUser(ctx, int32(id))
	if err != nil {
		fmt.Println("error on Deleting User: ", strErr)
		http.Error(w, "Error on Deleting User with id: "+userId, http.StatusNotFound)
		return
	}

	fmt.Println("Deleting User with id: " + userId)
	err = json.NewEncoder(w).Encode("Deleting User with id: " + userId)
	if err != nil {
		fmt.Println("error on Deleting User: ", err)
		http.Error(w, "Error on Deleting User", http.StatusInternalServerError)
		return
	}
}
