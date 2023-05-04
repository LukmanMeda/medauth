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
	Code         string `json:"code"`
}

type UserLogin struct {
	App        core.App
	Dao        *daos.Dao
	Collection *mod.Collection

	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func NewUserLogin(app core.App, collection *mod.Collection) *UserLogin {
	return &UserLogin{
		App:        app,
		Dao:        app.Dao(),
		Collection: collection,
	}
}

type ClientLogin struct {
	App        core.App
	Dao        *daos.Dao
	Collection *mod.Collection

	Username     string `form:"username" json:"username"`
	Password     string `form:"password" json:"password"`
	ClientId     string `form:"client_id" json:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
}

func NewClientLogin(app core.App, collection *mod.Collection) *ClientLogin {
	return &ClientLogin{
		App:        app,
		Dao:        app.Dao(),
		Collection: collection,
	}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (ul *UserLogin) Validate() error {
	return validation.ValidateStruct(ul,
		validation.Field(&ul.Username, validation.Required, validation.Length(1, 255)),
		validation.Field(&ul.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the form.
// On success returns the authorized record model.
//
// You can optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
func (ul *UserLogin) Submit(interceptors ...InterceptorFunc[*mod.Record]) (*mod.Record, error) {
	//var form *RecordPasswordLogin

	if err := ul.Validate(); err != nil {
		return nil, err
	}

	authOptions := ul.Collection.AuthOptions()

	var authRecord *mod.Record
	var fetchErr error

	isEmail := is.EmailFormat.Validate(ul.Username) == nil

	if isEmail {
		if authOptions.AllowEmailAuth {
			authRecord, fetchErr = ul.Dao.FindAuthRecordByEmail(ul.Collection.Id, ul.Username)
		}
	} else if authOptions.AllowUsernameAuth {
		authRecord, fetchErr = ul.Dao.FindAuthRecordByUsername(ul.Collection.Id, ul.Username)
	}

	// ignore not found errors to allow custom fetch implementations
	if fetchErr != nil && !errors.Is(fetchErr, sql.ErrNoRows) {
		return nil, fetchErr
	}

	interceptorsErr := runInterceptors(authRecord, func(m *mod.Record) error {
		authRecord = m

		if authRecord == nil || !authRecord.ValidatePassword(ul.Password) {
			return errors.New("error: invalid credentials")
		}

		return nil
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return authRecord, nil
}

// func (form *RecordClientLogin) SetDao(dao *daos.Dao) {
// 	form.Dao = dao
// }

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (cl *ClientLogin) Validate() error {
	return validation.ValidateStruct(cl,
		validation.Field(&cl.Username, validation.Required, validation.Length(1, 255)),
		validation.Field(&cl.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the form.
// On success returns the authorized record model.
//
// You can optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
func (cl *ClientLogin) Submit(interceptors ...InterceptorFunc[*mod.Record]) (*mod.Record, error) {
	//var form *RecordPasswordLogin

	if err := cl.Validate(); err != nil {
		return nil, err
	}

	authOptions := cl.Collection.AuthOptions()

	var authRecord *mod.Record
	var fetchErr error

	isEmail := is.EmailFormat.Validate(cl.Username) == nil

	if isEmail {
		if authOptions.AllowEmailAuth {
			authRecord, fetchErr = cl.Dao.FindAuthRecordByEmail(cl.Collection.Id, cl.Username)
		}
	} else if authOptions.AllowUsernameAuth {
		authRecord, fetchErr = cl.Dao.FindAuthRecordByUsername(cl.Collection.Id, cl.Username)
	}

	// ignore not found errors to allow custom fetch implementations
	if fetchErr != nil && !errors.Is(fetchErr, sql.ErrNoRows) {
		return nil, fetchErr
	}

	interceptorsErr := runInterceptors(authRecord, func(m *mod.Record) error {
		authRecord = m

		if authRecord == nil || !authRecord.ValidatePassword(cl.Password) {
			return errors.New("error: invalid credentials")
		}

		return nil
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return authRecord, nil
}
