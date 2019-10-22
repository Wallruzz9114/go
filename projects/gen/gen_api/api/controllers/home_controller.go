package controllers

import (
	"net/http"

	responses "github.com/Wallruzz9114/gen_api/api/responses"
)

// Home ...
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
