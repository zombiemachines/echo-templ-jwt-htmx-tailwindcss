package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	//
	"github.com/zombiemachines/echo-templ-htmx-tailwindcss/controllers"
	"github.com/zombiemachines/echo-templ-htmx-tailwindcss/middlewares"

	//

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
)

func main() {
	port := flag.String("port", "4000", "HTTP network port for http://localhost:")
	flag.Parse()

	e := echo.New()
	e.Logger.SetLevel(echoLog.INFO)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Static("/static"))
	e.Static("/static", "static")

	e.Use(middleware.Secure())
	e.Any("/login", controllers.LoginHandler)
	// e.POST("/logout", controllers.LogoutHandler)
	e.GET("/", controllers.GoHomeHandler)

	guardedRoutes := e.Group("/v1")
	guardedRoutes.Use(middlewares.AuthCookieMiddleware)
	guardedRoutes.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:   controllers.JwtSecretKey,
		TokenLookup:  "cookie:auth", // "<source>:<name>"
		ErrorHandler: controllers.JWTErrorChecker,
	}))
	guardedRoutes.Use(middlewares.TokenRefresherMiddleware)

	guardedRoutes.GET("", controllers.HomeHandler).Name = "HOME"

	{
		go func() {

			if err := e.StartTLS(":"+*port, "tls/cert.pem", "tls/key.pem"); err != nil && err != http.ErrServerClosed {
				e.Logger.Fatal("•• Shutting down server ◤•_•_•〓〓")
			}
		}()

		e.Logger.Infof("〓〓•_•_•◥ ••• Server Listening on https://localhost:" + *port)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		e.Logger.Printf("〓〓•_•_•◥ •• Shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
		e.Logger.Infof("• Server gracefully stopped ◤•_•_•〓〓")
	}

}
