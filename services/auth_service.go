package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"aura-erp/backend/config"
	"aura-erp/backend/models"
)

var jwtSecret = []byte(os.Getenv("JWT_Secret"))

type TokenPayload struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func init() {
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("your-secret-key-change-in-production")
	}
}

// Hash password, sha256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// Generate token
func GenerateToken(user models.User) string {
	payload := TokenPayload{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		Exp:   time.Now().Add(time.Hour*24).Unix() * 1000, //24hrs
	}

	payloadBytes, _ := json.Marshal(payload)
	payloadB64 := base64.StdEncoding.EncodeToString(payloadBytes)

	// Create HMAC signature
	h := hmac.New(sha256.New, jwtSecret)
	h.Write(payloadBytes)
	signature := fmt.Sprintf("%x", h.Sum(nil))

	return payloadB64 + "." + signature
}

func VerifyToken(token string) (*TokenPayload, error) {
	parts := strings.Split(token, ".")

	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	payloadB64 := parts[0]
	signature := parts[1]

	//Decode payload
	payloadBytes, err := base64.StdEncoding.DecodeString(payloadB64)

	if err != nil {
		return nil, errors.New("invalid token encoding")
	}

	//Verify Signature
	h := hmac.New(sha256.New, jwtSecret)
	h.Write(payloadBytes)
	expectedSignature := fmt.Sprintf("%x", h.Sum(nil))

	if signature != expectedSignature {
		return nil, errors.New("invalid token signature")
	}

	//parse payload
	var payload TokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, errors.New("invalid token payload")
	}

	//check expiration
	if payload.Exp < time.Now().Unix()*1000 {
		return nil, errors.New("token expired")
	}

	return &payload, nil
}

func Login(email, password string) (*LoginResponse, error) {
	hashedPassword := HashPassword(password)

	var user models.User
	query := `SELECT id, name, email, role FROM users WHERE email = $1 AND password_hash = $2`

	err := config.DB.QueryRow(query, email, hashedPassword).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role,
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token := GenerateToken(user)

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func VerifyTokenAndGetUser(token string) (*models.User, error) {
	payload, err := VerifyToken(token)

	if err != nil {
		return nil, err
	}

	var user models.User
	query := `SELECT id, name, email, role, created_at FROM users WHERE id = $1`

	err = config.DB.QueryRow(query, payload.ID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
