package models

import (
	"medauth/domain"
	"regexp"

	"github.com/pocketbase/pocketbase/tools/hook"
)

// base ID value regex pattern
var idRegex = regexp.MustCompile(`^[^\@\#\$\&\|\.\,\'\"\\\/\s]+$`)

// InterceptorNextFunc is a interceptor handler function.
// Usually used in combination with InterceptorFunc.
type InterceptorNextFunc[T any] func(t T) error

// InterceptorFunc defines a single interceptor function that
// will execute the provided next func handler.
type InterceptorFunc[T any] func(next InterceptorNextFunc[T]) InterceptorNextFunc[T]

// runInterceptors executes the provided list of interceptors.
func runInterceptors[T any](
	data T,
	next InterceptorNextFunc[T],
	interceptors ...InterceptorFunc[T],
) error {
	for i := len(interceptors) - 1; i >= 0; i-- {
		next = interceptors[i](next)
	}

	return next(data)
}

type BaseAuth struct {
	onRecordBeforeAuthWithPasswordRequest *hook.Hook[*domain.RecordAuthWithPasswordEvent]

	onRecordAfterAuthWithPasswordRequest *hook.Hook[*domain.RecordAuthWithPasswordEvent]
}

func (base *BaseAuth) OnRecordBeforeAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*domain.RecordAuthWithPasswordEvent] {
	return hook.NewTaggedHook(base.onRecordBeforeAuthWithPasswordRequest, tags...)
}

func (base *BaseAuth) OnRecordAfterAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*domain.RecordAuthWithPasswordEvent] {
	return hook.NewTaggedHook(base.onRecordAfterAuthWithPasswordRequest, tags...)
}
