package viewModel

import (
	"encoding/json"
	"time"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type AppViewModel struct {
	Routes                             RoutesViewModel           `json:"routes" js:"routes"`
	DisplayVersion                     string                    `json:"displayVersion" js:"displayVersion"`
	ProductName                        string                    `json:"productName" js:"productName"`
	LoggedIn                           bool                      `json:"loggedIn" js:"loggedIn"`
	Version                            int                       `json:"version" js:"version"`
	UserName                           string                    `json:"username" js:"username"`
	CopyrightYear                      int                       `json:"CopyrightYear" js:"CopyrightYear"`
	HTTPPort                           int                       `json:"HTTPPort" js:"HTTPPort"`
	AccountId                          string                    `json:"AccountId" js:"AccountId"`
	AccountName                        string                    `json:"AccountName" js:"AccountName"`
	AccountUsername                    string                    `json:"AccountUsername" js:"AccountUsername"`
	AccountRoleId                      string                    `json:"AccountRoleId" js:"AccountRoleId"`
	IsSystemAccount                    bool                      `json:"IsSystemAccount" js:"IsSystemAccount"`
	UserInitials                       string                    `json:"UserInitials" js:"UserInitials"`
	UserFirst                          string                    `json:"UserFirst" js:"UserFirst"`
	UserLast                           string                    `json:"UserLast" js:"UserLast"`
	UserEmail                          string                    `json:"UserEmail" js:"UserEmail"`
	UserPrimaryAccount                 string                    `json:"UserPrimaryAccount"`
	UserLanguage                       string                    `json:"UserLanguage" js:"UserLanguage"`
	UserEnforcePasswordChange          bool                      `json:"UserEnforcePasswordChange" js:"UserEnforcePasswordChange"`
	UserPreferences                    []model.UsersPreference   `json:"UserPreferences" js:"UserPreferences"`
	HasRole                            map[string]bool           `json:"HasRole"`
	UserId                             string                    `json:"UserId" js:"UserId"`
	ConfirmPassword                    string                    `json:"ConfirmPassword"`
	ConfirmPasswordErrors              string                    `json:"ConfirmPasswordErrors"`
	AccountTypeShort                   string                    `json:"AccountTypeShort" js:"AccountTypeShort"`
	DialogOpen                         bool                      `json:"DialogOpen" js:"DialogOpen"`
	DialogTranslationTitle             string                    `json:"DialogTitle" js:"DialogTitle"`
	DialogMessage                      string                    `json:"DialogMessage" js:"DialogMessage"`
	DialogOpen2                        bool                      `json:"DialogOpen2" js:"DialogOpen2"`
	DialogTranslationTitle2            string                    `json:"DialogTitle2" js:"DialogTitle2"`
	ShowDialogSubmitBug2               bool                      `json:"ShowDialogSubmitBug2" js:"ShowDialogSubmitBug2"`
	DialogMessage2                     string                    `json:"DialogMessage2" js:"DialogMessage2"`
	DialogGenericOpen                  bool                      `json:"DialogGenericOpen" js:"DialogGenericOpen"`
	DialogGenericTranslationTitle      string                    `json:"DialogGenericTitle" js:"DialogGenericTitle"`
	DialogGenericMessage               string                    `json:"DialogGenericMessage" js:"DialogGenericMessage"`
	SnackBarUndoTransactionId          string                    `json:"SnackBarUndoTransactionId" js:"SnackBarUndoTransactionId"`
	SnackBarAutoHideDuration           int                       `json:"SnackbarAutoHideDuration" js:"SnackbarAutoHideDuration"`
	SnackBarMessage                    string                    `json:"SnackbarMessage" js:"SnackbarMessage"`
	SnackBarOpen                       bool                      `json:"SnackbarOpen" js:"SnackbarOpen"`
	PopupErrorSubmit                   bool                      `json:"PopupErrorSubmit" js:"PopupErrorSubmit"`
	SnackBarType                       string                    `json:"SnackbarType" js:"SnackbarType"`
	Banner                             Banner                    `json:"Banner"`
	SideBarMenu                        SideBarViewModel          `json:"SideBarMenu"`
	DeveloperMode                      bool                      `json:"DeveloperMode"`
	DeveloperLogTheseObjects           []string                  `json:"DeveloperLogTheseObjects"`
	DeveloperSuppressTheseObjects      []string                  `json:"DeveloperSuppressTheseObjects"`
	DeveloperSuppressThesePages        []string                  `json:"DeveloperSuppressThesePages"`
	DeveloperLogStateChangePerformance bool                      `json:"DeveloperLogStateChangePerformance"`
	DeveloperLogState                  bool                      `json:"DeveloperLogState"`
	DeveloperLogReact                  bool                      `json:"DeveloperLogReact"`
}

func (self *AppViewModel) LoadDefaultState() {
	self.Routes = RoutesViewModel{}
	self.Routes.LoadDefaultState()
	self.DisplayVersion = "0.0.1"
	self.Version = 1
	self.CopyrightYear = time.Now().Year()
	self.SnackBarAutoHideDuration = 10000
	self.DialogOpen = false
	self.DialogTranslationTitle = ""
	self.DialogMessage = ""
	self.Banner = Banner{}
	self.Banner.LoadDefaultState()
}

func (self *AppViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
