package controllers

import (
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

func FileObjectPostCommitHook(actionPerformed string, context session_functions.RequestContext, id string) bool {
	if utils.InArray(actionPerformed, utils.Array("CreateFileObject", "CopyFileObject", "UpdateFileObject")) && id != "" {
		var row model.FileObject
		err := model.FileObjects.Query().ById(id, &row)
		if err == nil {
			t, err := session_functions.StartTransaction(context())
			if err == nil {
				// Add custom business rules to clean up data based on relationships or changes here

				//if row.Joins.Buildings.Count > 1 {
				//	row.IsCampus = true
				//} else {
				//	row.IsCampus = false
				//}

				err = row.SaveWithTran(t)
				if err == nil {
					err = t.Commit()
					if err == nil {
						return true
					}
				}
			} else {
				session_functions.Dump("Desc->Error in LocationEntityPostSaveHook", err)
			}
		} else {
			session_functions.Dump("Desc->Error in LocationEntityPostSaveHook 29", err)
		}
		return false
	} else {
		return true
	}
}
