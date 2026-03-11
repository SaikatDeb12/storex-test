package handler

import (
	"net/http"

	"github.com/SaikatDeb12/storeX/internal/database"
	"github.com/SaikatDeb12/storeX/internal/database/dbhelper"
	"github.com/SaikatDeb12/storeX/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
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

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid user id")
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelper.DeleteUser(tx, userID)
		if err != nil {
			return err
		}

		err = dbhelper.UnassignAssets(tx, userID)
		if err != nil {
			return err
		}

		err = dbhelper.DeleteUserSession(tx, userID)
		if err != nil {
			return err
		}

		return err
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to delete user")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"status": "user deleted successfully",
	})
}
