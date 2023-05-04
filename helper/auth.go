package helper

import (
	"medauth/models"

	"github.com/pocketbase/pocketbase/tools/hook"
)

type Meda interface {
	OnRecordBeforeAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*models.RecordAuthWithPasswordEvent]

	OnRecordAfterAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*models.RecordAuthWithPasswordEvent]
}
