package main

import (
	"log"
	"medauth/handler"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	mod "github.com/pocketbase/pocketbase/models"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// add new "GET /hello" route to the app router (echo)
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/hello",
			Handler: func(c echo.Context) error {
				return c.String(200, "Hello world!")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		e.Router.POST("/register", handler.Register(app), apis.LoadCollectionContext(app, mod.CollectionTypeAuth), apis.ActivityLogger(app))
		e.Router.POST("/login", handler.Login(app), apis.LoadCollectionContext(app, mod.CollectionTypeAuth), apis.ActivityLogger(app))

		e.Router.GET("/authorize", handler.Authorize(app), apis.LoadCollectionContext(app, mod.CollectionTypeAuth), apis.ActivityLogger(app))

		e.Router.POST("/token", handler.Token(app), apis.LoadCollectionContext(app, mod.CollectionTypeAuth), apis.ActivityLogger(app))
		//e.Router.POST("/token")
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
