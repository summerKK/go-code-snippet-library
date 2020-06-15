package controllers

import (
	"net/http"

	"github.com/summerKK/go-code-snippet-library/champion-go/api/responses"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome")
}
