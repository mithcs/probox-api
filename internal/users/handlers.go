package users

import (
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = compareUserIDWithToken(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = deleteUserByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, err)
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
		globals.WriteErrorResponse(w, err)
		return
	}

	response, err := globals.EncodeJson(user)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}

func GetUserByUsername(w http.ResponseWriter, r *http.Request, username string) {
	user, err := getUserByUsername(r.Context(), username)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	response, err := globals.EncodeJson(user)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := getAllUsers(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	response, err := globals.EncodeJson(users)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := globals.ParseBody[CreateUserRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = user.Validate(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	hashedPass, err := hashPass(user.Password)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	userID, err := storeUser(r.Context(), user.Username, hashedPass, user.FullName)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	tokens, err := generateTokens(userID)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	response, err := globals.EncodeJson(tokens)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}
