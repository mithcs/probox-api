package users

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Invalid User ID.",
		})

		return
	}

	err = compareUserIDWithToken(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusUnauthorized,
			Title:   "Unauthorized Action.",
			Details: err.Error(),
		})

		return
	}

	err = deleteUserByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not delete user.",
		})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	identifier := chi.URLParam(r, "identifier")
	if id, err := strconv.ParseInt(identifier, 10, 64); err == nil {
		GetUserByID(w, r, id)
	} else {
		GetUserByUsername(w, r, identifier)
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	user, err := getUserByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not get user.",
		})

		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list user.",
		})

		return
	}

	w.Write(response)
}

func GetUserByUsername(w http.ResponseWriter, r *http.Request, username string) {
	user, err := getUserByUsername(r.Context(), username)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not get user.",
		})

		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list user.",
		})

		return
	}

	w.Write(response)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := getAllUsers(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get all users.",
		})

		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list all users.",
		})

		return
	}

	w.Write(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := globals.ParseBody[CreateUserRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not parse body.",
		})

		return
	}

	err = user.Validate(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	hashedPass, err := hashPass(user.Password)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not hash your password.",
		})

		return
	}

	userID, err := storeUser(r.Context(), user.Username, hashedPass, user.FullName)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not store credentials.",
		})

		return
	}

	tokens, err := generateTokens(userID)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not generate tokens for you.",
		})

		return
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not give you tokens.",
		})

		return
	}

	w.Write(response)
}
