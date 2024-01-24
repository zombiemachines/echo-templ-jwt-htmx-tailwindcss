package controllers

import (
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
	return views.IndexHome("Lyoko").Render(context, writer)
}

func LoginHandler(c echo.Context) error {

	context := c.Request().Context()
	writer := c.Response().Writer

	m := c.Request().Method
	switch m {
	case "GET":
		return views.IndexLogin("Lyoko").Render(context, writer)
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

// func LogoutHandler(c echo.Context) error {
// 	context := c.Request().Context()
// 	writer := c.Response().Writer
// 	return views.LoginPage().Render(context, writer)
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
