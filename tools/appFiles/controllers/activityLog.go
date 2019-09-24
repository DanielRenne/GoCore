package controllers

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/pkg/errors"
)

var ActivityLog ActivityLogController

type ActivityLogController struct {
}

func (self ActivityLogController) UpsertActivityByContext(context session_functions.RequestContext, entity string, entityId string, action string, value string) (err error) {

	user, err := queries.Common.GetUser(context)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	account, err := session_functions.GetSessionAccount(context())

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	err = self.UpsertActivity(t, account.Id.Hex(), user.Id.Hex(), entity, entityId, action, value)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	err = t.Commit()

	return
}

func (self ActivityLogController) UpsertActivity(t *model.Transaction, acctId string, userId string, entity string, entityId string, action string, value string) (err error) {

	filter := make(map[string]interface{}, 6)
	filter[model.FIELD_ACTIVITYLOG_ACCOUNTID] = acctId
	filter[model.FIELD_ACTIVITYLOG_USERID] = userId
	filter[model.FIELD_ACTIVITYLOG_ENTITY] = entity
	filter[model.FIELD_ACTIVITYLOG_ENTITYID] = entityId
	filter[model.FIELD_ACTIVITYLOG_ACTION] = action
	filter[model.FIELD_ACTIVITYLOG_VALUE] = value

	var log model.ActivityLog
	err = model.ActivityLogs.Query().Filter(filter).GetOrCreate(&log, t)

	log.SaveWithTran(t)

	return
}
