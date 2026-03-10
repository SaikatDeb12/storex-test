package handler

import (
	"net/http"

	"github.com/SaikatDeb12/storeX/internal/database"
	"github.com/SaikatDeb12/storeX/internal/database/dbhelper"
	"github.com/SaikatDeb12/storeX/internal/middleware"
	"github.com/SaikatDeb12/storeX/internal/models"
	"github.com/SaikatDeb12/storeX/internal/utils"
	"github.com/jmoiron/sqlx"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "payload validation error")
		return
	}

	isEmailExists, err := dbhelper.CheckUserExistsByEmail(req.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error while checking email exists or not")
		return
	}

	if isEmailExists {
		utils.RespondError(w, http.StatusUnauthorized, nil, "user already exists")
		return
	}

	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "password hashing failed")
		return
	}

	var token string
	err = database.Tx(func(tx *sqlx.Tx) error {
		userID, err := dbhelper.CreateUser(req.Name, req.Email, req.PhoneNumber, req.Role, req.Employment, hashedPassword)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
			return err
		}

		sessionID, err := dbhelper.CreateSession(userID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to create session")
			return err
		}

		role := req.Role
		token, err = utils.GenerateJWT(userID, sessionID, role)

		return err
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error while generating token")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "user register successfully",
		"token":   token,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	user, err := dbhelper.GetUserAuthByEmail(req.Email)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid credentials")
		return
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid credentials")
		return
	}

	userID := user.ID
	sessionID, err := dbhelper.CreateSession(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create session")
		return
	}

	role := user.Role
	token, err := utils.GenerateJWT(userID, sessionID, role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error while generating token")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "login successfull",
		"token":   token,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userContext, ok := middleware.UserContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	sessionID := userContext.SessionID
	if err := dbhelper.ValidateUserSession(sessionID); err != nil {
		utils.RespondError(w, http.StatusForbidden, err, "no active session found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
