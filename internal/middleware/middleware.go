package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/SaikatDeb12/storeX/internal/models"
	"github.com/SaikatDeb12/storeX/internal/utils"
	"github.com/golang-jwt/jwt"
)

type ContextKeys struct{}

var RequestContextKey = ContextKeys{}

func UserContext(r *http.Request) (models.RequestContext, bool) {
	user, ok := r.Context().Value(RequestContextKey).(models.RequestContext)
	return user, ok
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusNotFound, nil, "missing authorization header")
			return
		}

		const brearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, brearerPrefix) {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, brearerPrefix)
		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(utils.SecretKey), nil
		})

		if !token.Valid || parseErr != nil {
			utils.RespondError(w, http.StatusUnauthorized, parseErr, "invalid token")
			return
		}

		claimValues, err := token.Claims.(jwt.MapClaims)
		if !err {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token claims")
			return
		}

		userID, err := claimValues["user_id"].(string)
		if !err {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid user id")
			return
		}

		sessionID, err := claimValues["session_id"].(string)
		if !err {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid session id")
			return
		}

		role := claimValues["role"].(string)
		requestContext := models.RequestContext{
			UserID:    userID,
			SessionID: sessionID,
			Role:      role,
		}

		ctx := context.WithValue(r.Context(), RequestContextKey, requestContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func CheckUserRole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		userCtx, _ := UserContext(r)
// 		userRole := userCtx.Role
//
// 		if userRole == "admin" || userRole == "asset_manager" {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 		utils.RespondError(w, http.StatusUnauthorized, nil, "not authorized")
// 	})
// }

func CheckUserRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx, _ := UserContext(r)
		role := userCtx.Role

		if role == "admin" || role == "asset_manager" {
			next.ServeHTTP(w, r)
			return
		}
		utils.RespondError(w, http.StatusUnauthorized, nil, "access denied")
	})
}
