package handler

import (
	"net/http"

	"github.com/SaikatDeb12/storeX/internal/database/dbhelper"
	"github.com/SaikatDeb12/storeX/internal/utils"
	"github.com/go-chi/chi/v5"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	role := query.Get("role")
	employment := query.Get("employment")
	assetStatus := query.Get("status")

	userDetails, err := dbhelper.FetchUsers(name, role, employment, assetStatus)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"users": userDetails,
	})
}

func GetUserInfoByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userDetails, err := dbhelper.FetchUserByID(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch user details")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"user": userDetails,
	})
}

// func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
// 	userID := chi.URLParam(r, "id")
// }
