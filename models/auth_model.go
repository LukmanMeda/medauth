package models

import (
	"database/sql"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	mod "github.com/pocketbase/pocketbase/models"
)

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

	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
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
		validation.Field(&form.Username, validation.Required, validation.Length(1, 255)),
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

	isEmail := is.EmailFormat.Validate(form.Username) == nil

	if isEmail {
		if authOptions.AllowEmailAuth {
			authRecord, fetchErr = form.Dao.FindAuthRecordByEmail(form.Collection.Id, form.Username)
		}
	} else if authOptions.AllowUsernameAuth {
		authRecord, fetchErr = form.Dao.FindAuthRecordByUsername(form.Collection.Id, form.Username)
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

type RecordClientLogin struct {
	App        core.App
	Dao        *daos.Dao
	Collection *mod.Collection

	Username     string `form:"username" json:"username"`
	Password     string `form:"password" json:"password"`
	ClientId     string `form:"client_id" json:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
}

func NewRecordClientLogin(app core.App, collection *mod.Collection) *RecordClientLogin {
	return &RecordClientLogin{
		App:        app,
		Dao:        app.Dao(),
		Collection: collection,
	}
}

func (form *RecordClientLogin) SetDao(dao *daos.Dao) {
	form.Dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordClientLogin) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Username, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the form.
// On success returns the authorized record model.
//
// You can optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
func (form *RecordClientLogin) Submit(interceptors ...InterceptorFunc[*mod.Record]) (*mod.Record, error) {
	//var form *RecordPasswordLogin

	if err := form.Validate(); err != nil {
		return nil, err
	}

	authOptions := form.Collection.AuthOptions()

	var authRecord *mod.Record
	var fetchErr error

	isEmail := is.EmailFormat.Validate(form.Username) == nil

	if isEmail {
		if authOptions.AllowEmailAuth {
			authRecord, fetchErr = form.Dao.FindAuthRecordByEmail(form.Collection.Id, form.Username)
		}
	} else if authOptions.AllowUsernameAuth {
		authRecord, fetchErr = form.Dao.FindAuthRecordByUsername(form.Collection.Id, form.Username)
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
