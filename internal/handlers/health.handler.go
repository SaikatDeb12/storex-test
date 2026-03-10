package handler

import (
	"net/http"

	"github.com/SaikatDeb12/storeX/internal/utils"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"status": "server is running",
	})
}
