package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SaikatDeb12/storeX/internal/database"
	"github.com/SaikatDeb12/storeX/internal/database/dbhelper"
	"github.com/SaikatDeb12/storeX/internal/middleware"
	"github.com/SaikatDeb12/storeX/internal/models"
	"github.com/SaikatDeb12/storeX/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAssetRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "payload validation failed")
		return
	}

	err := database.Tx(func(tx *sqlx.Tx) error {
		assetID, err := dbhelper.CreateAsset(tx, req)
		if err != nil {
			return err
		}

		assetType := req.Type
		switch assetType {
		case "laptop":
			err = dbhelper.InsertLaptopDetails(tx, assetID, *req.Laptop)
		case "keyboard":
			err = dbhelper.InsertKeyboardDetails(tx, assetID, *req.Keyboard)
		case "mouse":
			err = dbhelper.InsertMouseDetails(tx, assetID, *req.Mouse)
		case "mobile":
			err = dbhelper.InsertMobileDetails(tx, assetID, *req.Mobile)
		default:
			err = errors.New("invalid asset type")
		}

		return err
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create asset")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "asset added",
	})
}

func FetchAssets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	brand := query.Get("brand")
	model := query.Get("model")
	assetType := query.Get("type")
	status := query.Get("status")
	owner := query.Get("owner")
	serialNumber := query.Get("serial_number")

	limit := 10
	page := 1
	var err error

	if limitInput := query.Get("limit"); limitInput != "" {
		limit, err = strconv.Atoi(limitInput)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err, "invalid limit")
			return
		}
	}

	if pageInput := query.Get("page"); pageInput != "" {
		page, err = strconv.Atoi(pageInput)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err, "invalid limit")
			return
		}
	}

	page = max(page, 1)
	offset := (page - 1) * limit

	allAssets, err := dbhelper.FetchAssets(brand, model, assetType, serialNumber, status, owner, limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch assets")
		return
	}

	assetsCount, err := dbhelper.GettingAssetsCount()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get asset count")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"dashboard": assetsCount,
		"assets":    allAssets,
	})
}

func AssignAssets(w http.ResponseWriter, r *http.Request) {
	var req models.AssetAssignRequest

	if parseErr := utils.ParseBody(r.Body, &req); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed parsing body")
		return
	}
	userCtx, _ := middleware.UserContext(r)
	currectUserID := userCtx.UserID

	err := dbhelper.AssignAssets(req.AssetID, currectUserID, req.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to assign assets")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "successfully assigned",
	})
}

func UpdateAsset(w http.ResponseWriter, r *http.Request) {
	assetId := chi.URLParam(r, "id")
	if assetId == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid asset id")
		return
	}
	var req models.UpdateAssetRequest
	err := utils.ParseBody(r.Body, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid body")
		return
	}
	validateErr := utils.ValidateStruct(&req)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}
	warrantyStart, err := time.Parse(time.RFC3339, req.WarrantyStart)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid warrantyStart")
		return
	}

	warrantyEnd, err := time.Parse(time.RFC3339, req.WarrantyEnd)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid warrantyEnd")
		return
	}

	if warrantyEnd.Before(warrantyStart) {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid warranty range")
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelper.UpdateAsset(tx, assetId, req.Brand, req.Model, req.SerialNumber, req.Type, req.Owner, warrantyStart, warrantyEnd)
		if err != nil {
			return err
		}
		switch req.Type {

		case "laptop":
			if req.Laptop == nil {
				return fmt.Errorf("laptop details required")
			}
			return dbhelper.UpdateLaptop(tx, assetId, req.Laptop)

		case "mouse":
			if req.Mouse == nil {
				return fmt.Errorf("mouse details required")
			}
			return dbhelper.UpdateMouse(tx, assetId, req.Mouse)
		case "keyboard":
			if req.Keyboard == nil {
				return fmt.Errorf("keyboard details required")
			}
			return dbhelper.UpdateKeyboard(tx, assetId, req.Keyboard)
		case "mobile":
			if req.Mobile == nil {
				return fmt.Errorf("mobile details required")
			}
			return dbhelper.UpdateMobile(tx, assetId, req.Mobile)

		default:
			return fmt.Errorf("unsupported asset type")
		}
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusBadRequest, txErr, "fail to update asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "asset updated",
	})
}

func SentToService(w http.ResponseWriter, r *http.Request) {
	assetId := chi.URLParam(r, "id")
	if assetId == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid id")
		return
	}

	var req models.SentServiceRequest
	err := utils.ParseBody(r.Body, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid body")
		return
	}
	err = utils.ValidateStruct(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "validation failed")
		return
	}

	serviceStart, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid service date")
		return
	}
	serviceEnd, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid service end date")
		return
	}
	if serviceEnd.Before(serviceStart) {
		utils.RespondError(w, http.StatusBadRequest, nil, "end date must be after start date")
		return
	}
	err = dbhelper.SentToService(assetId, serviceStart, serviceEnd)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "fail to sent for service")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "asset sent for service successfully",
	})
}
