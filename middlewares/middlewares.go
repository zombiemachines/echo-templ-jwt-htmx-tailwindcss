package middlewares

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/zombiemachines/echo-templ-htmx-tailwindcss/controllers"
	"github.com/zombiemachines/echo-templ-htmx-tailwindcss/models"
)

func AuthCookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		e := echo.New().Logger
		cookie, err := c.Cookie(controllers.AccessTokenCookieName)
		if err != nil {
			e.Infof("\n----Invalid Cookie----")
			c.Response().Header().Set("Location", "/login")
			return c.Redirect(301, "/login")
		}
		token := cookie.Value
		if token == "" {
			e.Infof("\n----Empty token----")
			c.Response().Header().Set("Location", "/login")
			return c.Redirect(308, "/login")
		}
		if isValid := controllers.ValidateToken(token); !isValid {
			e.Infof("\n----Invalid token----")
			c.Response().Header().Set("Location", "/login")
			return c.Redirect(301, "/login")
		}
		if token == "" {
			e.Infof("\n----Redirect To Login Page----")
			c.Response().Header().Set("Location", "/login")
			return c.Redirect(301, "/login")
		}

		return next(c)
	}
}

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := echo.New().Logger

		cookie, err := c.Cookie(controllers.AccessTokenCookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				logger.Infof("\n----ERROR NO COOKIE >> LOGIN----")
				c.Response().Header().Set("Location", "/login")
				return c.Redirect(301, "/login")
			}
			logger.Errorf("Error retrieving cookie: %v", err)
			return c.String(http.StatusBadRequest, "INVALID COOKIE >> LOGIN")
		}

		tokenStr := cookie.Value
		claims := &models.CustomClaims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return controllers.JwtSecretKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.String(http.StatusUnauthorized, "SIGNATURE INVALID >> LOGIN")
			}
			logger.Errorf("Error parsing token: %v", err)
			return c.String(http.StatusBadRequest, "INTERNAL ERROR ^^^")
		}

		if !tkn.Valid {
			return c.String(http.StatusUnauthorized, "INVALID TOKEN")
		}

		// Check time remaining for token to expire
		if time.Until(time.Unix(claims.ExpiresAt, 0)) < 15*time.Minute {
			// Make new token
			expirationTime := time.Now().Add(time.Minute * 15)
			claims.ExpiresAt = expirationTime.Unix()

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(controllers.JwtSecretKey)

			if err != nil {
				logger.Errorf("Error creating new token: %v", err)
				return c.String(http.StatusBadRequest, "FAILED CONVERTING TOKEN TO STRING")
			}

			// todo : call the SetTokenCookie() instead
			c.SetCookie(
				&http.Cookie{
					Name:    controllers.AccessTokenCookieName,
					Value:   tokenString,
					Expires: expirationTime,

					HttpOnly: true,
					Secure:   true,
					Path:     "/",
				})
		}

		return next(c)
	}
}
