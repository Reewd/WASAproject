package api

import (
	"net/http"
	"strconv"

	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) idVerifierMiddleware(next httpRouterHandler) httpRouterHandler {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
		// Check if the request has a valid bearer token
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		id, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Check if the user ID exists in the database
		exists, err := rt.db.UserExistsById(id)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to check user ID existence")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx.UserID = id
		ctx.Username, err = rt.db.GetUsername(id)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to get username")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Here you would typically validate the token and extract user information
		// For simplicity, we assume the token is valid and proceed

		next(w, r, ps, ctx)
	}
}
