package handler

import (
	"log"
	"medauth/models"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	mod "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

const ContextCollectionKey string = "collection"

func Register(app core.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input *models.User
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "cant bind input")
		}

		collection, err := app.Dao().FindCollectionByNameOrId("users")
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

		if err := app.Dao().SaveRecord(record); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, record)
	}
}

// func Login(app core.App) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		collection, err := app.Dao().FindCollectionByNameOrId("users")
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, err.Error())
// 		}

// 		form := forms.NewRecordPasswordLogin(app, collection)
// 		if readErr := c.Bind(form); readErr != nil {
// 			return c.JSON(http.StatusBadRequest, readErr)
// 		}

// 		event := new(core.RecordAuthWithPasswordEvent)
// 		event.HttpContext = c
// 		event.Collection = collection
// 		event.Password = form.Password
// 		event.Identity = form.Identity

// 		_, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*mod.Record]) forms.InterceptorNextFunc[*mod.Record] {
// 			return func(record *mod.Record) error {
// 				event.Record = record

// 				token, tokenErr := tokens.NewRecordAuthToken(app, record)
// 				if tokenErr != nil {
// 					return c.JSON(http.StatusBadRequest, tokenErr)
// 				}

// 				return app.OnRecordBeforeAuthWithPasswordRequest().Trigger(event, func(e *core.RecordAuthWithPasswordEvent) error {
// 					if err := next(e.Record); err != nil {
// 						return c.JSON(http.StatusBadRequest, err)
// 					}

// 					result := map[string]any{
// 						"acess_token": token,
// 						"token_type":  "bearer",
// 						"expires_in":  3000,
// 					}

// 					return c.JSON(http.StatusOK, result)
// 					//return apis.RecordAuthResponse(app, e.HttpContext, e.Record, nil)
// 				})
// 			}
// 		})

// 		if submitErr == nil {
// 			if err := app.OnRecordAfterAuthWithPasswordRequest().Trigger(event); err != nil && app.IsDebug() {
// 				log.Println(err)
// 			}
// 		}

// 		return submitErr
// 		//return c.JSON(http.StatusOK, submitErr)

// 	}
// }

func Login(app core.App) echo.HandlerFunc {
	return func(c echo.Context) error {

		collection, err := app.Dao().FindCollectionByNameOrId("users")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		form := models.NewRecordPasswordLogin(app, collection)

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

				token, tokenErr := tokens.NewRecordAuthToken(app, record)
				if tokenErr != nil {
					return c.JSON(http.StatusBadRequest, tokenErr)
				}

				return app.OnRecordBeforeAuthWithPasswordRequest().Trigger(event, func(e *core.RecordAuthWithPasswordEvent) error {
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
			if err := app.OnRecordAfterAuthWithPasswordRequest().Trigger(event); err != nil && app.IsDebug() {
				log.Println(err)
			}
		}

		return submitErr
		//return c.JSON(http.StatusOK, submitErr)

	}
}
