package helper

import (
	"medauth/domain"

	"github.com/pocketbase/pocketbase/tools/hook"
)

type Meda interface {
	OnRecordBeforeAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*domain.RecordAuthWithPasswordEvent]

	OnRecordAfterAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*domain.RecordAuthWithPasswordEvent]
}
