package models

import (
	"database/sql"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	mod "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/hook"
)

var (
	//	_ hook.Tagger = (*BaseModelEvent)(nil)
	_ hook.Tagger = (*BaseCollectionEvent)(nil)
)

type BaseCollectionEvent struct {
	Collection *mod.Collection
}

// Tags implements hook.Tagger
func (e *BaseCollectionEvent) Tags() []string {
	if e.Collection == nil {
		return nil
	}

	tags := make([]string, 0, 2)

	if e.Collection.Id != "" {
		tags = append(tags, e.Collection.Id)
	}

	if e.Collection.Name != "" {
		tags = append(tags, e.Collection.Name)
	}

	return tags
}

type Authorize struct {
	ClientId     string `db:"client_id" json:"client_id"`
	ClientSecret string `db:"client_secret" json:"client_secret"`
	ResponseType string `db:"response_type" json:"response_type"`
	Code         string ` json:"code"`
}

type PasswordLogin struct {
	ClientId     string `db:"client_id" json:"client_id"`
	ClientSecret string `db:"client_secret" json:"client_secret"`
}

type RecordPasswordLogin struct {
	App        core.App
	Dao        *daos.Dao
	Collection *mod.Collection

	Identity     string `form:"identity" json:"identity"`
	Password     string `form:"password" json:"password"`
	ClientId     string `form:"client_id" json:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
}

type RecordClientLogin struct {
	ClientId     string `form:"client_id" json:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
}

type RecordAuthWithPasswordEvent struct {
	BaseCollectionEvent

	HttpContext  echo.Context
	Record       *mod.Record
	Identity     string
	Password     string
	ClientId     string
	ClientSecret string
}

func NewRecordPasswordLogin(app core.App, collection *mod.Collection) *RecordPasswordLogin {
	return &RecordPasswordLogin{
		App:        app,
		Dao:        app.Dao(),
		Collection: collection,
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordPasswordLogin) SetDao(dao *daos.Dao) {
	form.Dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordPasswordLogin) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Identity, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the form.
// On success returns the authorized record model.
//
// You can optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
func (form *RecordPasswordLogin) Submit(interceptors ...InterceptorFunc[*mod.Record]) (*mod.Record, error) {
	//var form *RecordPasswordLogin

	if err := form.Validate(); err != nil {
		return nil, err
	}

	authOptions := form.Collection.AuthOptions()

	var authRecord *mod.Record
	var fetchErr error

	isEmail := is.EmailFormat.Validate(form.Identity) == nil

	if isEmail {
		if authOptions.AllowEmailAuth {
			authRecord, fetchErr = form.Dao.FindAuthRecordByEmail(form.Collection.Id, form.Identity)
		}
	} else if authOptions.AllowUsernameAuth {
		authRecord, fetchErr = form.Dao.FindAuthRecordByUsername(form.Collection.Id, form.Identity)
	}

	// ignore not found errors to allow custom fetch implementations
	if fetchErr != nil && !errors.Is(fetchErr, sql.ErrNoRows) {
		return nil, fetchErr
	}

	interceptorsErr := runInterceptors(authRecord, func(m *mod.Record) error {
		authRecord = m

		if authRecord == nil || !authRecord.ValidatePassword(form.Password) {
			return errors.New("error: invalid credentials")
		}

		return nil
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return authRecord, nil
}
