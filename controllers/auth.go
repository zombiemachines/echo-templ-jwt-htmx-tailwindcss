package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/zombiemachines/echo-templ-htmx-tailwindcss/models"
)

var (
	AccessTokenCookieName = "auth"
	JwtSecretKey          = []byte("372de7cb-1f4e-4fde-bd77-7ae1f8f2f879")
	// RefreshTokenCookieName = "refresh-auth"
	// JwtRefreshSecretKey    = "c0bee718-d567-4bd3-8d4c-9f523a83ec4c"
)

func ValidateToken(tokenString string) (bool, string) {
	// Define the secret key used to sign the tokens
	secret := JwtSecretKey
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return false, ""
	}
	// Check if the token is valid
	claims, ok := token.Claims.(*models.CustomClaims)
	if ok && token.Valid {
		isExpired := claims.VerifyExpiresAt(time.Now().Unix(), true) //maybe should be isNotExpired
		if isExpired {                                               //if isNotExpired{....
			fmt.Println("Token has not expired.")
			fmt.Printf("\n-----------Token is valid for user: %s +++++++++++\n", claims.Username)
			//GET USERNAME FROM CLAIMS
			return true, claims.Username
		}
	}
	fmt.Println("-----------Invalid token//Token has expired.+++++++++++")

	return false, ""
}

func CreateToken(name string, roles []string, c echo.Context) (string, error) {
	// token := jwt.New(jwt.SigningMethodHS256)
	// claims := jwt.MapClaims{}
	// token := jwt.New(jwt.SigningMethodHS256)
	claims := &models.CustomClaims{}
	claims.Username = name
	claims.ExpiresAt = time.Now().Add(time.Minute * 30).Unix() //* 24 * 10
	// claims2 := &models.CustomClaims{
	// 	Username:       name,
	// 	StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24 * 10).Unix()},
	// }
	// claims := token.Claims.(jwt.MapClaims)
	// claims["username"] = name //i think this working and appearently case insensitive
	// claims["authorized"] = true
	// claims["roles"] = roles
	// claims["exp"] = time.Now().Add(time.Hour * 24 * 10).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //shit
	tokenString, err := token.SignedString(JwtSecretKey)

	return tokenString, err
}

// Here we are creating a new cookie, which will store the valid JWT token.
func SetTokenCookie(name string, token string, expiration time.Time, c echo.Context) {

	cookie := &http.Cookie{
		Name:     name,
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Expires:  expiration,
		Path:     "/"}

	c.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(c echo.Context, err error) error {
	// Redirects to the signIn form.

	c.Response().Header().Set("Location", "/v1")
	return c.Redirect(301, c.Echo().Reverse("HOME"))
}
