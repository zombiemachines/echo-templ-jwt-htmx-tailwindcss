package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/zombiemachines/echo-templ-htmx-tailwindcss/models"
	"github.com/zombiemachines/echo-templ-htmx-tailwindcss/views"
)

func HomeHandler(c echo.Context) error {

	context := c.Request().Context()
	writer := c.Response().Writer

	// get the username from the context
	//CAUTION: c.Set("user") is already used by JWT to save the token
	//so we did set c.Set("username") instead of c.Set("user") that JWT uses
	user := c.Get("username").(string)

	return views.IndexHome(user).Render(context, writer)
}

func LoginHandler(c echo.Context) error {

	context := c.Request().Context()
	writer := c.Response().Writer

	m := c.Request().Method
	switch m {
	case "GET":
		user, ok := c.Get("username").(string)
		if !ok {
			return views.IndexLogin(user).Render(context, writer)
		}
		return views.IndexLogin("").Render(context, writer)
	case "POST":
		user := models.User{}
		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if user.Username != "admin" || user.Password != "1234" {
			return echo.ErrUnauthorized
		}
		token, err := CreateToken(user.Username, []string{"admin"}, c)
		if err != nil {
			return err
		}

		expTime := time.Now().Add(time.Minute * 30) // * 24* 10
		SetTokenCookie(AccessTokenCookieName, token, expTime, c)

		c.Response().Header().Set("Location", "/v1")
		return c.Redirect(301, c.Echo().Reverse("HOME"))
	}

	// return c.NoContent(400)
	// code 400  http.StatusBadRequest
	return echo.ErrBadRequest
}

func GoHomeHandler(c echo.Context) error {

	c.Response().Header().Set("Location", "/v1")
	return c.Redirect(301, c.Echo().Reverse("HOME"))
}

func LogoutHandler(c echo.Context) error {
	// Clear or expire the user cookie
	rc := &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		Path:     "/"}

	c.SetCookie(rc)
	// If the token is expired, remove the user from the context
	c.Set("username", nil)
	c.Set("user", nil)
	fmt.Println("+++++++++++++++++LOGOUT+++++++++++++++++++")

	//go login
	c.Response().Header().Set("Location", "/login")   //
	return c.Redirect(302, c.Echo().Reverse("LOGIN")) //

	// return views.IndexLogin("").Render(c.Request().Context(), c.Response().Writer)
}

// func Restricted(c echo.Context) error {
// 	user := c.Get("username").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	name := claims["name"].(string)
// 	roles := claims["roles"].([]interface{})
// 	return c.String(http.StatusOK,
// 		fmt.Sprintf("Welcome %s, you have roles: %v\n", name, roles))
// }

// func SkipperFn(skipURLs []string) func(echo.Context) bool {
// 	return func(context echo.Context) bool {
// 		for _, url := range skipURLs {
// 			if url == context.Request().URL.String() {
// 				return true
// 			}
// 		}
// 		return false
// 	}
// }

func HelloPostHandler(c echo.Context) error {
	context := c.Request().Context()
	writer := c.Response().Writer
	var v = models.Visitor{LikesDueling: true}
	v.Name = c.FormValue("name")
	if v.Name != "" {
		return views.Card(v.Name).Render(context, writer)
	}

	return htmx.NewResponse().
		Retarget("#toTarget").Reswap(htmx.SwapInnerHTML.Transition(false)). //disable transition api
		RenderTempl(context, writer, views.WftMsg("Need a name!!"))
	// return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid name")
}

func FormHandler(c echo.Context) error {
	context := c.Request().Context()
	writer := c.Response().Writer
	return views.LoginForm().Render(context, writer)
}
