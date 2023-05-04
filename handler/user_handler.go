package handler

import (
	"log"
	"medauth/models"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	mod "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

type userHandler struct {
	app core.App
}

func UserHandler(app core.App, e *echo.Echo) {
	api := &userHandler{app: app}

	e.POST("/register", api.register(), apis.LoadCollectionContext(app, mod.CollectionTypeAuth), apis.ActivityLogger(app))
	e.POST("/login", api.login(), apis.LoadCollectionContext(app, mod.CollectionTypeAuth), apis.ActivityLogger(app))

}

func (uh *userHandler) register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input *models.User
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "cant bind input")
		}

		collection, err := uh.app.Dao().FindCollectionByNameOrId("users")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		record := mod.NewRecord(collection)

		record.Set("username", input.Username)
		record.Set("name", input.Name)
		record.Set("email", input.Email)
		record.Set("phone", input.Phone)
		record.Set("active", true)

		record.RefreshTokenKey()
		record.SetPassword(input.Password)
		record.ValidatePassword(input.PasswordConfirm)
		record.PasswordHash()

		if err := uh.app.Dao().SaveRecord(record); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, record)
	}
}

func (uh *userHandler) login() echo.HandlerFunc {
	return func(c echo.Context) error {

		collection, err := uh.app.Dao().FindCollectionByNameOrId("users")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		form := models.NewUserLogin(uh.app, collection)

		if readErr := c.Bind(form); readErr != nil {
			return c.JSON(http.StatusBadRequest, readErr)
		}

		event := new(core.RecordAuthWithPasswordEvent)
		event.HttpContext = c
		event.Collection = collection
		event.Password = form.Password
		event.Identity = form.Username

		_, submitErr := form.Submit(func(next models.InterceptorNextFunc[*mod.Record]) models.InterceptorNextFunc[*mod.Record] {
			return func(record *mod.Record) error {
				event.Record = record

				token, tokenErr := tokens.NewRecordAuthToken(uh.app, record)
				if tokenErr != nil {
					return c.JSON(http.StatusBadRequest, tokenErr)
				}

				return uh.app.OnRecordBeforeAuthWithPasswordRequest().Trigger(event, func(e *core.RecordAuthWithPasswordEvent) error {
					if err := next(e.Record); err != nil {
						return c.JSON(http.StatusBadRequest, err)
					} else {

						result := map[string]any{
							"acess_token": token,
							"token_type":  "bearer",
							"expires_in":  3000,
						}

						return c.JSON(http.StatusOK, result)
					}
				})
			}
		})

		if submitErr == nil {
			if err := uh.app.OnRecordAfterAuthWithPasswordRequest().Trigger(event); err != nil && uh.app.IsDebug() {
				log.Println(err)
			}
		}

		return submitErr
		//return c.JSON(http.StatusOK, submitErr)

	}
}
