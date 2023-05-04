package domain

import (
	"github.com/labstack/echo/v5"
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

type RecordAuthWithPasswordEvent struct {
	BaseCollectionEvent

	HttpContext echo.Context
	Record      *mod.Record
	Username    string
	Password    string
}
