package handler

import (
	"log"
	"medauth/helper"
	"medauth/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	mod "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

func Authorize(app core.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input *models.Authorize

		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "cant bind input")
		}

		if strings.TrimSpace(input.ClientId) == "" || strings.TrimSpace(input.ResponseType) == "" {
			return c.JSON(http.StatusBadRequest, "please insert one field")
		}

		recordClient, err := app.Dao().FindFirstRecordByData("clients", "client_id", input.ClientId)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		} else {
			//input.Code = "BRG"

			authCodeCollection, err := app.Dao().FindFirstRecordByData("auth_codes", "client_id", recordClient.Get("id"))

			// authCodeCollection, err := app.Dao().FindCollectionByNameOrId("auth_codes")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

			// authCodeRecord := mod.NewRecord(authCodeCollection)
			// authCodeRecord.Set("code", input.Code)
			// authCodeRecord.Set("redirect_uri", recordClient.Get("redirect_uri"))
			// authCodeRecord.Set("client_id", recordClient.Get("client_id"))

			// if err := app.Dao().SaveRecord(authCodeRecord); err != nil {
			// 	return c.JSON(http.StatusInternalServerError, err.Error())
			// }

			result := map[string]any{
				"redirect_uri":  recordClient.Get("redirect_uri"),
				"response_type": input.ResponseType,
				"code":          authCodeCollection.Get("code"),
			}
			return c.JSON(http.StatusOK, result)
		}

		//return nil

	}
}

// func Token(app core.App) echo.HandlerFunc {
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

func Token(app core.App) echo.HandlerFunc {
	return func(c echo.Context) error {

		var meda helper.Meda

		collection, errCol := app.Dao().FindCollectionByNameOrId("users")
		if errCol != nil {
			log.Println("error di collection " + errCol.Error())
			return c.JSON(http.StatusInternalServerError, "error di collection "+errCol.Error())
		}

		form := models.NewRecordPasswordLogin(app, collection)

		if readErr := c.Bind(form); readErr != nil {
			log.Println("error di Bind Data" + readErr.Error())
			return c.JSON(http.StatusBadRequest, "error di Bind Data "+readErr.Error())
		}

		recordClient, err := app.Dao().FindFirstRecordByData("clients", "client_id", form.ClientId)
		if err != nil {
			log.Println("error di recordClient " + err.Error())
			return c.JSON(http.StatusNotFound, "error di recordClinet"+err.Error())
		}

		event := new(models.RecordAuthWithPasswordEvent)
		event.HttpContext = c
		event.Collection = collection
		event.Password = form.Password
		event.Identity = form.Identity
		event.ClientId = recordClient.Get("client_id").(string)
		event.ClientSecret = recordClient.Get("client_secret").(string)

		_, submitErr := form.Submit(func(next models.InterceptorNextFunc[*mod.Record]) models.InterceptorNextFunc[*mod.Record] {
			return func(record *mod.Record) error {
				event.Record = record

				token, tokenErr := tokens.NewRecordAuthToken(app, record)
				if tokenErr != nil {
					log.Println("error di Token " + tokenErr.Error())
					return c.JSON(http.StatusBadRequest, "error di token "+tokenErr.Error())
				}

				return meda.OnRecordBeforeAuthWithPasswordRequest().Trigger(event, func(e *models.RecordAuthWithPasswordEvent) error {
					if errNext := next(e.Record); errNext != nil {
						log.Println("error di Next " + errNext.Error())
						return c.JSON(http.StatusBadRequest, "error di Next "+errNext.Error())
					} else {

						result := map[string]any{
							"acess_token": token,
							"massage":     "Selamat datang" + " " + event.Identity + " " + "Anda berhasil login",
						}

						return apis.RecordAuthResponse(app, e.HttpContext, e.Record, result)
					}

					//return c.JSON(http.StatusOK, result)
				})
			}
		})
		if submitErr == nil {
			if errSubmit := meda.OnRecordAfterAuthWithPasswordRequest().Trigger(event); errSubmit != nil && app.IsDebug() {
				log.Println("error di Next " + errSubmit.Error())
				return errSubmit
			}
		}

		return nil

	}
}
