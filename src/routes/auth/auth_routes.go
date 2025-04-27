package auth

import (
	"backend/data"
	"backend/errors"
	"backend/models"
	"backend/utils"
	"backend/validators"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// ValidateRefreshToken verifies whether a refresh token is valid
//
// meta:operation POST /api/auth/validate-refresh
func ValidateRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.InvalidRequestMethod, http.StatusMethodNotAllowed)
		return
	}

	var reqBody models.RefreshBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		http.Error(w, errors.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	_, err = utils.ValidateToken(reqBody.RefreshToken)

	if err != nil {
		http.Error(w, errors.InvalidToken, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := map[string]string{
		"refresh_token": reqBody.RefreshToken,
	}

	err = json.NewEncoder(w).Encode(res)
}

// RefreshAuthToken regenerates a new JWT from the refresh token
//
// meta:operation POST /api/auth/refresh-token
func RefreshAuthToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.InvalidRequestMethod, http.StatusMethodNotAllowed)
		return
	}

	var reqBody models.RefreshBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		http.Error(w, errors.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	token, err := utils.ValidateToken(reqBody.RefreshToken)

	if token == nil || err != nil {
		http.Error(w, errors.InvalidToken, http.StatusUnauthorized)
		return
	}

	user, err := data.GetUserByEmail(token.Email)

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	if user == nil {
		http.Error(w, errors.InvalidToken, http.StatusUnauthorized)
		return
	}

	newToken, err := utils.CreateToken(user.Email, utils.RefreshJWTDuration)

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	res := map[string]string{
		"token": newToken,
	}

	err = json.NewEncoder(w).Encode(res)
}

// CreateUser creates a new user object and returns a JWT
//
// meta:operation POST /api/auth/users
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.InvalidRequestMethod, http.StatusMethodNotAllowed)
		return
	}

	var user models.NewUser
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, errors.InvalidRequestBody, http.StatusBadRequest)
		return
	}

	if err = validators.ValidateEmail(user.Email); err != nil {
		http.Error(w, errors.InvalidRequestBody, http.StatusBadRequest)
		log.Printf("Error: %v", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	user.Password = string(hashedPassword)

	authToken, err := utils.CreateToken(user.Email, utils.AuthJWTDuration)

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	refreshToken, err := utils.CreateToken(user.Email, utils.RefreshJWTDuration)

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	err = data.AddUser(&user)

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	res := map[string]string{
		"token":         authToken,
		"refresh_token": refreshToken,
	}

	err = json.NewEncoder(w).Encode(res)
}

// VerifyUser validates a user's credentials or JWT
//
// meta:operation POST /api/auth/verify
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.InvalidRequestMethod, http.StatusMethodNotAllowed)
		return
	}

	token, _ := utils.ExtractTokenFromHeader(r)
	var useCredentialAuth = token == ""
	var reqUserCredentials models.UserCredentials

	if useCredentialAuth {
		var credentials models.UserCredentials
		err := json.NewDecoder(r.Body).Decode(&credentials)

		if err != nil {
			http.Error(w, errors.InvalidRequestBody, http.StatusBadRequest)
			return
		}

		reqUserCredentials = credentials
	} else {
		claims, err := utils.ValidateToken(token)

		if err != nil {
			http.Error(w, errors.InvalidToken, http.StatusUnauthorized)
			return
		}

		reqUserCredentials.Email = claims.Email
	}

	dbUser, err := data.GetUserByEmail(reqUserCredentials.Email)

	if dbUser == nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	if useCredentialAuth {
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(reqUserCredentials.Password))

		if err != nil {
			http.Error(w, errors.InvalidCredentials, http.StatusUnauthorized)
			return
		}
	}

	userJSON, err := json.Marshal(struct {
		Id            string  `json:"id"`
		Email         string  `json:"email"`
		DisplayName   string  `json:"display_name"`
		ProfileImgx64 *string `json:"profile_img_x64"`
	}{
		Id:            dbUser.Id,
		Email:         dbUser.Email,
		DisplayName:   dbUser.DisplayName,
		ProfileImgx64: dbUser.ProfileImgx64,
	})

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	authToken, err := utils.CreateToken(dbUser.Email, utils.AuthJWTDuration)

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	refreshToken, err := utils.CreateToken(dbUser.Email, utils.RefreshJWTDuration)

	if err != nil {
		http.Error(w, errors.InternalServerError, http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := map[string]string{
		"token":         authToken,
		"refresh_token": refreshToken,
		"user":          string(userJSON),
	}

	_ = json.NewEncoder(w).Encode(res)
}
