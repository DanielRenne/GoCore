package queries

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/pkg/errors"
)

const (
	ACTIVITY_ACCOUNT_ENTITY = "Account"
	ACTIVITY_ACTION_ACCESS  = "Access"
	ACTIVITY_ACTION_VIEW    = "View"
)

type queryActivityLogs struct{}

func (self queryActivityLogs) queryAccountActivityContext(context session_functions.RequestContext) (q *model.Query, err error) {
	user, err := getUser(context)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	filter := make(map[string]interface{}, 3)
	filter[model.FIELD_ACTIVITYLOG_USERID] = user.Id.Hex()
	filter[model.FIELD_ACTIVITYLOG_ENTITY] = ACTIVITY_ACCOUNT_ENTITY
	filter[model.FIELD_ACTIVITYLOG_ACTION] = ACTIVITY_ACTION_ACCESS

	q = model.ActivityLogs.Query().Sort("UpdateDate").Filter(filter)
	return
}

func (self queryActivityLogs) AccountActivity(context session_functions.RequestContext, limit int) (logs []model.ActivityLog, err error) {
	q, err := self.queryAccountActivityContext(context)
	err = q.Limit(limit).All(&logs)
	return
}