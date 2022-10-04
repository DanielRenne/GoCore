// package gitWebHooks - deprecated use git actions instead
package gitWebHooks

import (
	"sync"
)

const (
	PROJECT_COLUMN = "ProjectColumn"
	PROJECT        = "Project"
	PROJECT_CARD   = "ProjectCard"
	ISSUES         = "Issues"
	PUSH_TYPE      = "Push"
	ISSUE_COMMENT  = "IssueComment"
)

var registry sync.Map

type webHooks interface {
	RunEvent(x interface{})
}

// RegisterHook will register a new store to the store registry.
func RegisterHook(typeHook string, x interface{}) {
	registry.Store(typeHook, x)
}

func RunEvent(key string, x interface{}) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	hook, ok := getRegistry(key)
	if !ok {
		return
	}

	hook.RunEvent(x)
	return
}

func getRegistry(key string) (x webHooks, ok bool) {

	obj, ok := registry.Load(key)
	if ok {
		x = obj.(webHooks)
	}
	return
}
