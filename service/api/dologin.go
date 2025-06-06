package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Reewd/WASAproject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type LoginRequest struct {
	Name string `json:"name"`
}

// doLogin allows the user to login using their username in a bearer auth fashion
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the username
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id, err := rt.db.Login(req.Name)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	_, _ = w.Write([]byte("Logged in with ID: " + strconv.FormatInt(id, 10)))

}
