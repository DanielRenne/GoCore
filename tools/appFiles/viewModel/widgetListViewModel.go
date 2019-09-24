package viewModel

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"gopkg.in/mgo.v2/bson"
)

type WidgetListUserControlsViewModel struct {
	PerPage         int      `json:"PerPage"`
	Page            int      `json:"Page"`
	SearchFields    []string `json:"SearchFields"`
	SortBy          string   `json:"SortBy"`
	SortDirection   string   `json:"SortDirection"`
	Criteria        string   `json:"Criteria"`
	DataKey         string   `json:"DataKey"`
	CustomCriteria  string   `json:"CustomCriteria"`
	ListTitle       string   `json:"ListTitle"`
	IsDefaultFilter bool     `json:"IsDefaultFilter"`
}

func (this *WidgetListUserControlsViewModel) LoadDefaultState() {
	this.PerPage = 10
	this.Page = 1
	this.SearchFields = utils.Array()
	this.SortBy = "CreateDate"
	this.SortDirection = "-"
	this.Criteria = ""
	this.ListTitle = ""
	this.IsDefaultFilter = true
}

func (self *WidgetListUserControlsViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}

func InitWidgetList() WidgetListUserControlsViewModel {
	widgetList := WidgetListUserControlsViewModel{}
	widgetList.LoadDefaultState()
	return widgetList
}

func InitWidgetListWithParams(uriParams map[string]string) WidgetListUserControlsViewModel {
	widgetList := WidgetListUserControlsViewModel{}
	widgetList.LoadDefaultState()
	var page string
	var perPage string
	perPage, _ = uriParams["PerPage"]
	perPageSplit := strings.Split(perPage, ".")
	if len(perPageSplit) > 0 {
		widgetList.PerPage = extensions.StringToInt(perPageSplit[0])
	} else {
		widgetList.PerPage = extensions.StringToInt(perPage)
	}
	page, _ = uriParams["Page"]
	pSplit := strings.Split(page, ".")
	if len(pSplit) > 0 {
		widgetList.Page = extensions.StringToInt(pSplit[0])
	} else {
		widgetList.Page = extensions.StringToInt(page)
	}

	criteria, ok := uriParams["CustomCriteria"]
	if ok && len(criteria) > 0 {
		widgetList.CustomCriteria = uriParams["CustomCriteria"]
	}
	widgetList.ListTitle, _ = uriParams["ListTitle"]
	widgetList.SortBy, _ = uriParams["SortBy"]
	widgetList.SortDirection, _ = uriParams["SortDirection"]
	widgetList.Criteria, _ = uriParams["Criteria"]

	return widgetList
}

func FilterWidgetList(vm WidgetListUserControlsViewModel, q *model.Query) {
	FilterWidgetListLogic(vm, q, true)

}

func FilterWidgetListLogic(vm WidgetListUserControlsViewModel, q *model.Query, withLimit bool) {
	//.Filter()
	if withLimit {
		q.Limit(vm.PerPage)
	}
	if vm.Page > 1 {
		q.Skip((vm.Page - 1) * vm.PerPage)
	}
	if vm.SortBy != "" {
		q.Sort(vm.SortDirection + vm.SortBy)
	}
	andOrFilters := q.GetAndOr()
	if len(andOrFilters) > 0 {
		if vm.Criteria != "" && len(vm.SearchFields) > 0 {
			if vm.Criteria != "" && len(vm.SearchFields) > 0 {
				for _, field := range vm.SearchFields {
					if field == "Id" && len(vm.Criteria) == 24 {
						field = "_id"
						q.OrFilter(0, model.Q(field, bson.ObjectIdHex(vm.Criteria)))
					} else {
						q.OrFilter(0, model.Q(field, bson.M{"$regex": bson.RegEx{`.*` + regexp.QuoteMeta(vm.Criteria) + `.*`, "i"}}))
					}
				}
			}
		} else {
			q.OrFilter(0, model.Q("_id", bson.M{"$exists": true}))
		}
	}
}

func FilterWidgetListNoLimit(vm WidgetListUserControlsViewModel, q *model.Query) {
	FilterWidgetListLogic(vm, q, false)
}
