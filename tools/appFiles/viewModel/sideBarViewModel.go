package viewModel

import (
	"errors"

	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/gin-gonic/gin"
)

func GetSideBarViewModel(c *gin.Context) SideBarViewModel {
	var sbv SideBarViewModel
	sbv.LoadDefaultState(c)
	return sbv
}

type SideBarViewModel struct {
	Items []SideBarMenuItem `json:"Items"`
}

func (self *SideBarViewModel) LoadDefaultState(c *gin.Context) {
	var dashboardItem SideBarMenuItem
	dashboardItem.Title = "DashboardTitle"
	dashboardItem.URL = "/#/home"
	dashboardItem.Icon = "Dashboard"

	var allUsersItem SideBarMenuItem
	allUsersItem.Title = "UsersTitle"
	allUsersItem.Icon = "Users"
	allUsersItem.URL = "/#/userList"

	var allAccountsItem SideBarMenuItem
	allAccountsItem.Title = "-99999"
	allAccountsItem.URL = "/#/accountList"

	var accountsItem SideBarMenuItem
	accountsItem.Title = "AccountsTitle"
	accountsItem.Icon = "Account"
	accountsItem.Items = append(accountsItem.Items, allAccountsItem)

	self.Items = append(self.Items, dashboardItem)

	if session_functions.CheckRoleAccess(c, constants.FEATURE_ACCOUNT_VIEW) {
		self.Items = append(self.Items, accountsItem)
	}
	if session_functions.CheckRoleAccess(c, constants.FEATURE_USER_VIEW) {
		self.Items = append(self.Items, allUsersItem)
	}
}

func (self *SideBarViewModel) GetIndex(c *gin.Context, page string) (index int) {
	if page == "ACCOUNT" {
		if session_functions.CheckRoleAccess(c, constants.FEATURE_ACCOUNT_VIEW) {
			index++
		}
	}
	if page == "USER" {
		if session_functions.CheckRoleAccess(c, constants.FEATURE_ACCOUNT_VIEW) {
			index++
		}
		if session_functions.CheckRoleAccess(c, constants.FEATURE_USER_VIEW) {
			index++
		}
	}
	return index
}

func (self *SideBarViewModel) RenderApp(c *gin.Context, account model.Account, user model.User) (err error) {
	if len(self.Items) == 0 {
		self.LoadDefaultState(c)
	}

	err = self.renderAccounts(c, account, user)
	return
}

func (self *SideBarViewModel) renderAccounts(c *gin.Context, account model.Account, user model.User) (err error) {

	if account.Id.Hex() == "" || user.Id.Hex() == "" {
		err = errors.New("Failed to Render Sidebar Menu.  User or Account not valid.")
		return
	}
	q, err := queries.AccountRoles.QueryByUser(user)

	if err != nil {
		return err
	}
	var accountRoles []model.AccountRole
	err = q.Join("Account").All(&accountRoles)

	if err != nil {
		return err
	}
	if len(self.Items[self.GetIndex(c, "ACCOUNT")].Items) == 1 {
		for _, acr := range accountRoles {
			if acr.AccountId == account.Id.Hex() {
				continue
			}
			var item SideBarMenuItem

			replacements := queries.TagReplacements{
				Tag1: queries.Q("account_name", acr.Joins.Account.AccountName),
			}
			item.Title = queries.AppContent.GetTranslationWithReplacements(session_functions.PassContext(c), "SideBarAccessAccountName", &replacements)
			item.URL = "javascript:window.globals.handleAccountSwitch('" + acr.Joins.Account.Id.Hex() + "')"
			item.RightIcon = "AccessAccount"
			if acr.AccountId == account.Id.Hex() {
				item.Selected = true
			}
			self.Items[self.GetIndex(c, "ACCOUNT")].Items = append(self.Items[self.GetIndex(c, "ACCOUNT")].Items, item)
		}
	}
	return
}

type SideBarMenuItem struct {
	Title         string            `json:"Title"`
	URL           string            `json:"URL"`
	Icon          string            `json:"Icon"`
	RightIcon     string            `json:"RightIcon"`
	RightIconLink string            `json:"RightIconLink"`
	Expanded      bool              `json:"Expanded"`
	Hidden        bool              `json:"Hidden"`
	Selected      bool              `json:"Selected"`
	Items         []SideBarMenuItem `json:"Items"`
}
