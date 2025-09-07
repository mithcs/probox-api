package tokens

import (
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

func CreateTokens(w http.ResponseWriter, r *http.Request) {
	creds, err := globals.ParseBody[CreateTokensRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	userID, err := creds.Verify(r.Context())
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

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	userID, err := retrieveUserID(r.Context())
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
