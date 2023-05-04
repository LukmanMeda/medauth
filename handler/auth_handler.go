package handler

import (
	"fmt"
	"log"
	"medauth/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
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

func Token(app core.App) echo.HandlerFunc {
	return func(c echo.Context) error {

		userColl, errCol := app.Dao().FindCollectionByNameOrId("users")
		if errCol != nil {
			log.Println("error di collection " + errCol.Error())
			return c.JSON(http.StatusInternalServerError, "error di collection "+errCol.Error())
		}
		form := models.NewRecordClientLogin(app, userColl)

		if readErr := c.Bind(form); readErr != nil {
			log.Println("error di Bind Data" + readErr.Error())
			return c.JSON(http.StatusBadRequest, "error di Bind Data "+readErr.Error())
		}
		clientColl, errClinet := app.Dao().FindFirstRecordByData("clients", "client_id", form.ClientId)
		if errClinet != nil {
			return c.JSON(http.StatusNotFound, errClinet.Error())
		}
		fmt.Println(clientColl)

		//TIDAK BISA DI GANTIIIIIII
		event := new(core.RecordAuthWithPasswordEvent)
		event.HttpContext = c
		event.Collection = userColl
		event.Password = form.Password
		event.Identity = form.Username

		_, submitErr := form.Submit(func(next models.InterceptorNextFunc[*mod.Record]) models.InterceptorNextFunc[*mod.Record] {
			return func(record *mod.Record) error {
				event.Record = record

				token, tokenErr := tokens.NewRecordAuthToken(app, record)
				if tokenErr != nil {
					log.Println("error di Token " + tokenErr.Error())
					return c.JSON(http.StatusBadRequest, "error di token "+tokenErr.Error())
				}

				return app.OnRecordBeforeAuthWithPasswordRequest().Trigger(event, func(e *core.RecordAuthWithPasswordEvent) error {
					if errNext := next(e.Record); errNext != nil {
						log.Println("error di Next " + errNext.Error())
						return c.JSON(http.StatusBadRequest, "error di Next "+errNext.Error())
					} else {

						result := map[string]any{
							"acess_token": token,
							"token_type":  form.ClientId,
							"expires_in":  3000,
						}

						return c.JSON(http.StatusOK, result)
						//return apis.RecordAuthResponse(app, e.HttpContext, e.Record, result)
					}
				})
			}
		})
		if submitErr == nil {
			if errSubmit := app.OnRecordAfterAuthWithPasswordRequest().Trigger(event); errSubmit != nil && app.IsDebug() {
				log.Println("error di Next " + errSubmit.Error())
				return errSubmit
			}
		}

		return nil

	}
}
