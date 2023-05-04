package handler

import (
	"fmt"
	"log"
	"medauth/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	mod "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

type authHandler struct {
	app core.App
}

func AuthHandler(app core.App, e *echo.Echo) {
	api := &authHandler{app: app}

	e.POST("/token", api.token(), apis.LoadCollectionContext(app, mod.CollectionTypeAuth), apis.ActivityLogger(app))
	e.POST("/auth", api.authorize(), apis.LoadCollectionContext(app, mod.CollectionTypeAuth), apis.ActivityLogger(app))

}

func (ah *authHandler) authorize() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input *models.Authorize

		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "cant bind input")
		}

		if strings.TrimSpace(input.ClientId) == "" || strings.TrimSpace(input.ResponseType) == "" {
			return c.JSON(http.StatusBadRequest, "please insert one field")
		}

		recordClient, err := ah.app.Dao().FindFirstRecordByData("clients", "client_id", input.ClientId)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		} else {
			//input.Code = "BRG"

			authCodeCollection, err := ah.app.Dao().FindFirstRecordByData("auth_codes", "client_id", recordClient.Get("id"))

			// authCodeCollection, err := app.Dao().FindCollectionByNameOrId("auth_codes")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

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

func (ah *authHandler) token() echo.HandlerFunc {
	return func(c echo.Context) error {

		userColl, errCol := ah.app.Dao().FindCollectionByNameOrId("users")
		if errCol != nil {
			log.Println("error di collection " + errCol.Error())
			return c.JSON(http.StatusInternalServerError, "error di collection "+errCol.Error())
		}
		request := models.NewClientLogin(ah.app, userColl)
		r := request.App.Dao()

		if readErr := c.Bind(request); readErr != nil {
			log.Println("error di Bind Data" + readErr.Error())
			return c.JSON(http.StatusBadRequest, "error di Bind Data "+readErr.Error())
		}

		if strings.TrimSpace(request.ClientId) == "" || strings.TrimSpace(request.ClientSecret) == "" || strings.TrimSpace(request.Username) == "" || strings.TrimSpace(request.Password) == "" {
			return c.JSON(http.StatusBadRequest, "invalid request")
		}

		clientColl, errClinet := r.FindFirstRecordByData("clients", "client_id", request.ClientId)
		if errClinet != nil {
			return c.JSON(http.StatusNotFound, errClinet.Error())
		}
		fmt.Println(clientColl)

		//TIDAK BISA DI GANTIIIIIII
		eventReq := new(core.RecordAuthWithPasswordEvent)
		eventReq.HttpContext = c
		eventReq.Collection = userColl
		eventReq.Password = request.Password
		eventReq.Identity = request.Username

		_, submitErr := request.Submit(func(next models.InterceptorNextFunc[*mod.Record]) models.InterceptorNextFunc[*mod.Record] {
			return func(record *mod.Record) error {
				eventReq.Record = record

				token, tokenErr := tokens.NewRecordAuthToken(ah.app, record)
				if tokenErr != nil {
					log.Println("error di Token " + tokenErr.Error())
					return c.JSON(http.StatusBadRequest, "error di token "+tokenErr.Error())
				}

				return ah.app.OnRecordBeforeAuthWithPasswordRequest().Trigger(eventReq, func(e *core.RecordAuthWithPasswordEvent) error {
					if errNext := next(e.Record); errNext != nil {
						log.Println("error di Next " + errNext.Error())
						return c.JSON(http.StatusBadRequest, "error di Next "+errNext.Error())
					} else {

						result := map[string]any{
							"acess_token": token,
							"token_type":  "Bearrer",
							"expires_in":  3000,
						}

						return c.JSON(http.StatusOK, result)
						//return apis.RecordAuthResponse(app, e.HttpContext, e.Record, result)
					}
				})
			}
		})
		if submitErr == nil {
			if errSubmit := ah.app.OnRecordAfterAuthWithPasswordRequest().Trigger(eventReq); errSubmit != nil && ah.app.IsDebug() {
				log.Println("error di Next " + errSubmit.Error())
				return errSubmit
			}
		}

		return nil

	}
}
