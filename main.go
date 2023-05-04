package main

import (
	"log"
	"medauth/handler"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		//call API method in Handler

		handler.UserHandler(app, e.Router)
		handler.AuthHandler(app, e.Router)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
