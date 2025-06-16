package helpers

import (
	"net/http"

	"github.com/Reewd/WASAproject/service/api/reqcontext"
)

func HandleInternalServerError(ctx reqcontext.RequestContext, w http.ResponseWriter, err error, message string) {
	ctx.Logger.WithError(err).Error(message)
	http.Error(w, "An unexpected error occurred. Please try again later.", http.StatusInternalServerError)
}
