package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Tokens struct to describe tokens object.
type Tokens struct {
	Access  string
	Refresh string
}

// GenerateNewTokens func for generate a new Access & Refresh tokens.
func GenerateNewTokens(id, email string) (*Tokens, error) {
	// Generate JWT Access token.
	accessToken, err := generateNewAccessToken(id, email)
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	// Generate JWT Refresh token.
	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func generateNewAccessToken(id, email string) (string, error) {

	envFile, _ := godotenv.Read(".env")
	// if err != nil {
	// 	return nil, fmt.Errorf("can't read env file")
	// }

	// Set secret key from .env file.
	secret := envFile["JWT_SECRET_KEY"]

	// Set expires minutes count for secret key from .env file.
	minutesCount, _ := strconv.Atoi(envFile["JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"])

	// Create a new claims.
	claims := jwt.MapClaims{}

	// Set public claims:
	claims["id"] = id
	claims["email"] = email
	claims["expires"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()
	claims["post:create"] = true
	// claims["post:create"] = false
	// claims["post:update"] = false
	// claims["book:delete"] = false

	// Set private token credentials:
	// for _, credential := range credentials {
	// 	claims[credential] = true
	// }

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}

func generateNewRefreshToken() (string, error) {
	// Create a new SHA256 hash.
	hash := sha256.New()

	envFile, err := godotenv.Read(".env")
	if err != nil {
		return "", fmt.Errorf("can't read env file")
	}

	// Create a new now date and time string with salt.
	refresh := envFile["JWT_REFRESH_KEY"] + time.Now().String()

	// See: https://pkg.go.dev/io#Writer.Write
	hashed, err := hash.Write([]byte(refresh))
	if err != nil {
		fmt.Println(hashed)
		// Return error, it refresh token generation failed.
		return "", err
	}

	//fmt.Println("hhhh", hashed)
	// Set expires hours count for refresh key from .env file.
	hoursCount, _ := strconv.Atoi(envFile["JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"])

	// Set expiration time.
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())

	// Create a new refresh token (sha256 string with salt + expire time).
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime

	return t, nil
}

// ParseRefreshToken func for parse second argument from refresh token.
func ParseRefreshToken(refreshToken string) (int64, error) {
	if refreshToken == "" {
		return 0, fmt.Errorf("index out of range [1] with length 1")
	}
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}
