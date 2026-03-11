package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/SaikatDeb12/storeX/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// package-level variables are initialized before the execution of main
var (
	// SecretKey = GetEnvVariables("SECRET_KEY")
	SecretKey string
	validate  = validator.New()
)

func GetEnvVariables(key string) string {
	return os.Getenv(key)
}

func ParseBody(body io.Reader, out interface{}) error {
	return json.NewDecoder(body).Decode(out)
}

func EncodeBody(w http.ResponseWriter, out interface{}) error {
	return json.NewEncoder(w).Encode(out)
}

func ValidateStruct(payload interface{}) error {
	return validate.Struct(payload)
}

func RespondError(w http.ResponseWriter, statusCode int, err error, message string) {
	w.WriteHeader(statusCode)
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	NewError := models.ErrorModel{
		Message:    message,
		Error:      errStr,
		StatusCode: statusCode,
	}

	if err := EncodeBody(w, NewError); err != nil {
		fmt.Printf("error: %+v", err)
	}
}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		err := EncodeBody(w, body)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}

func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(userID, sessionID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"session_id": sessionID,
		"role":       role,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}
