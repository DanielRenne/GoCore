package queries

import (
	"encoding/json"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
	"github.com/davidrenne/reflections"
)

type AppContentJson struct {
	InInvalidIP                          string `json:"InInvalidIP"`
	StartIPHigher                        string `json:"StartIPHigher"`
	EndIPHigher                          string `json:"EndIPHigher"`
	IPAddressInRangeOf                   string `json:"IPAddressInRangeOf"`
	IPAddressInRange                     string `json:"IPAddressInRange"`
	IPAddressCantStart                   string `json:"IPAddressCantStart"`
	SkippedXCSVRows                      string `json:"SkippedXCSVRows"`
	CopyRow                              string `json:"CopyRow"`
	RestartDone                          string `json:"RestartDone"`
	CSVFeatureGroupFieldName             string `json:"CSVFeatureGroupFieldName"`
	CSVFeatureFieldKey                   string `json:"CSVFeatureFieldKey"`
	CSVFeatureFieldName                  string `json:"CSVFeatureFieldName"`
	CSVFeatureFieldDescription           string `json:"CSVFeatureFieldDescription"`
	CSVFeatureFieldFeatureGroupId        string `json:"CSVFeatureFieldFeatureGroupId"`
	CSVFileObjectFieldName               string `json:"CSVFileObjectFieldName"`
	CSVFileObjectFieldContent            string `json:"CSVFileObjectFieldContent"`
	CSVFileObjectFieldSize               string `json:"CSVFileObjectFieldSize"`
	CSVFileObjectFieldType               string `json:"CSVFileObjectFieldType"`
	CSVFileObjectFieldModifiedUnix       string `json:"CSVFileObjectFieldModifiedUnix"`
	CSVFileObjectFieldModified           string `json:"CSVFileObjectFieldModified"`
	CSVFileObjectFieldMD5                string `json:"CSVFileObjectFieldMD5"`
	CSVFieldId                           string `json:"CSVFieldId"`
	CSVFieldIdHelp                       string `json:"CSVFieldIdHelp"`
	CSVRequired                          string `json:"CSVRequired"`
	CSVLineNumberRrBackendFailure        string `json:"CSVLineNumberRrBackendFailure"`
	CSVOutdated                          string `json:"CSVOutdated"`
	CSVRoleFeatureFieldRoleId            string `json:"CSVRoleFeatureFieldRoleId"`
	CSVRoleFeatureFieldFeatureId         string `json:"CSVRoleFeatureFieldFeatureId"`
	CSVRoleFieldName                     string `json:"CSVRoleFieldName"`
	CSVRoleFieldShortName                string `json:"CSVRoleFieldShortName"`
	SideBarAccessAccountName             string `json:"SideBarAccessAccountName"`
	ValidationFieldSpecificEmailRequired string `json:"ValidationFieldSpecificEmailRequired"`
	ErrorAccountSave                     string `json:"ErrorAccountSave"`
	ErrorAccountDelete                   string `json:"ErrorAccountDelete"`
	AccountSaveSuccess                   string `json:"AccountSaveSuccess"`
	AccountDeleteSuccessful              string `json:"AccountDeleteSuccessful"`
	ErrorAccountNotFound                 string `json:"ErrorAccountNotFound"`
	ErrorAccountNoId                     string `json:"ErrorAccountNoId"`
	ErrorAccountRelated                  string `json:"ErrorAccountRelated"`
	ErrorAccountRoleSave                 string `json:"ErrorAccountRoleSave"`
	ErrorAccountRoleDelete               string `json:"ErrorAccountRoleDelete"`
	AccountRoleSaveSuccess               string `json:"AccountRoleSaveSuccess"`
	AccountRoleDeleteSuccessful          string `json:"AccountRoleDeleteSuccessful"`
	ErrorAccountRoleNotFound             string `json:"ErrorAccountRoleNotFound"`
	ErrorAccountRoleNoId                 string `json:"ErrorAccountRoleNoId"`
	ErrorAccountUserNotSelected          string `json:"ErrorAccountUserNotSelected"`
	ImportSuccessful                     string `json:"ImportSuccessful"`
	ErrorUnauthorizedAccountRoleChange   string `json:"ErrorUnauthorizedAccountRoleChange"`
	ErrorAppErrorSave                    string `json:"ErrorAppErrorSave"`
	ErrorAppErrorDelete                  string `json:"ErrorAppErrorDelete"`
	AppErrorSaveSuccess                  string `json:"AppErrorSaveSuccess"`
	AppErrorDeleteSuccessful             string `json:"AppErrorDeleteSuccessful"`
	ErrorAppErrorNotFound                string `json:"ErrorAppErrorNotFound"`
	ErrorAppErrorNoId                    string `json:"ErrorAppErrorNoId"`
	ErrorAppErrorCopy                    string `json:"ErrorAppErrorCopy"`
	AppErrorDeleteMany                   string `json:"AppErrorDeleteMany"`
	AppErrorCopySuccessful               string `json:"AppErrorCopySuccessful"`
	InvalidAuth                          string `json:"InvalidAuth"`
	LoginPageLockLimitReached            string `json:"LoginPageLockLimitReached"`
	Authorized                           string `json:"Authorized"`
	LoggedOut                            string `json:"LoggedOut"`
	DialogMessage                        string `json:"DialogMessage"`
	DialogOpen                           string `json:"DialogOpen"`
	DialogTitle                          string `json:"DialogTitle"`
	DialogSuccess                        string `json:"DialogSuccess"`
	DialogWarning                        string `json:"DialogWarning"`
	DialogError                          string `json:"DialogError"`
	ErrorUserRelatedAccounts             string `json:"ErrorUserRelatedAccounts"`
	GenericDatabaseErrorOccurred         string `json:"GenericDatabaseErrorOccurred"`
	SuccessUserSaveDefaultAccount        string `json:"SuccessUserSaveDefaultAccount"`
	ErrorAccountCompanyFailed            string `json:"ErrorAccountCompanyFailed"`
	GlobalNoCountriesFound               string `json:"GlobalNoCountriesFound"`
	ErrorTranFailedRetrieve              string `json:"ErrorTranFailedRetrieve"`
	ErrorTranCommitFailed                string `json:"ErrorTranCommitFailed"`
	PostHookFailure                      string `json:"PostHookFailure"`
	ErrorTranRollbackFailed              string `json:"ErrorTranRollbackFailed"`
	ErrorRolesFailedRetrieve             string `json:"ErrorRolesFailedRetrieve"`
	ErrorAccountRoleFailedRetrieve       string `json:"ErrorAccountRoleFailedRetrieve"`
	ErrorAccountFailedRetrieve           string `json:"ErrorAccountFailedRetrieve"`
	ErrorCompanyFailedRetrieve           string `json:"ErrorCompanyFailedRetrieve"`
	ErrorsDeletingCurrentAccount         string `json:"ErrorsDeletingCurrentAccount"`
	ErrorUnauthorizedDelete              string `json:"ErrorUnauthorizedDelete"`
	ErrorDeleteAccountRole               string `json:"ErrorDeleteAccountRole"`
	ErrorAccountInvitation               string `json:"ErrorAccountInvitation"`
	ErrorPasswordReset                   string `json:"ErrorPasswordReset"`
	ErrorPasswordResetEmail              string `json:"ErrorPasswordResetEmail"`
	ErrorRequiredField                   string `json:"ErrorRequiredField"`
	SaveChanges                          string `json:"SaveChanges"`
	PasswordResetSuccess                 string `json:"PasswordResetSuccess"`
	ErrorFeatureSave                     string `json:"ErrorFeatureSave"`
	ErrorFeatureDelete                   string `json:"ErrorFeatureDelete"`
	FeatureSaveSuccess                   string `json:"FeatureSaveSuccess"`
	FeatureDeleteSuccessful              string `json:"FeatureDeleteSuccessful"`
	ErrorFeatureNotFound                 string `json:"ErrorFeatureNotFound"`
	ErrorFeatureNoId                     string `json:"ErrorFeatureNoId"`
	FeatureDeleteMany                    string `json:"FeatureDeleteMany"`
	ErrorFeatureCopy                     string `json:"ErrorFeatureCopy"`
	FeatureCopySuccessful                string `json:"FeatureCopySuccessful"`
	ErrorFeatureGroupSave                string `json:"ErrorFeatureGroupSave"`
	ErrorFeatureGroupDelete              string `json:"ErrorFeatureGroupDelete"`
	FeatureGroupSaveSuccess              string `json:"FeatureGroupSaveSuccess"`
	FeatureGroupDeleteSuccessful         string `json:"FeatureGroupDeleteSuccessful"`
	ErrorFeatureGroupNotFound            string `json:"ErrorFeatureGroupNotFound"`
	ErrorFeatureGroupNoId                string `json:"ErrorFeatureGroupNoId"`
	FeatureGroupDeleteMany               string `json:"FeatureGroupDeleteMany"`
	ErrorFeatureGroupCopy                string `json:"ErrorFeatureGroupCopy"`
	FeatureGroupCopySuccessful           string `json:"FeatureGroupCopySuccessful"`
	ErrorFileObjectSave                  string `json:"ErrorFileObjectSave"`
	ErrorFileObjectDelete                string `json:"ErrorFileObjectDelete"`
	FileObjectSaveSuccess                string `json:"FileObjectSaveSuccess"`
	FileObjectDeleteSuccessful           string `json:"FileObjectDeleteSuccessful"`
	ErrorFileObjectNotFound              string `json:"ErrorFileObjectNotFound"`
	ErrorFileObjectNoId                  string `json:"ErrorFileObjectNoId"`
	FileObjectDeleteMany                 string `json:"FileObjectDeleteMany"`
	ErrorFileObjectCopy                  string `json:"ErrorFileObjectCopy"`
	FileObjectCopySuccessful             string `json:"FileObjectCopySuccessful"`
	ErrorPasswordSave                    string `json:"ErrorPasswordSave"`
	ErrorPasswordDelete                  string `json:"ErrorPasswordDelete"`
	PasswordSaveSuccess                  string `json:"PasswordSaveSuccess"`
	PasswordDeleteSuccessful             string `json:"PasswordDeleteSuccessful"`
	ErrorPasswordNotFound                string `json:"ErrorPasswordNotFound"`
	ErrorPasswordNoId                    string `json:"ErrorPasswordNoId"`
	ErrorRoleSave                        string `json:"ErrorRoleSave"`
	ErrorRoleDelete                      string `json:"ErrorRoleDelete"`
	RoleSaveSuccess                      string `json:"RoleSaveSuccess"`
	RoleDeleteSuccessful                 string `json:"RoleDeleteSuccessful"`
	ErrorRoleNotFound                    string `json:"ErrorRoleNotFound"`
	ErrorRoleNoId                        string `json:"ErrorRoleNoId"`
	RoleDeleteMany                       string `json:"RoleDeleteMany"`
	ErrorRoleCopy                        string `json:"ErrorRoleCopy"`
	RoleCopySuccessful                   string `json:"RoleCopySuccessful"`
	ErrorRoleFeatureSave                 string `json:"ErrorRoleFeatureSave"`
	ErrorRoleFeatureDelete               string `json:"ErrorRoleFeatureDelete"`
	RoleFeatureSaveSuccess               string `json:"RoleFeatureSaveSuccess"`
	RoleFeatureDeleteSuccessful          string `json:"RoleFeatureDeleteSuccessful"`
	ErrorRoleFeatureNotFound             string `json:"ErrorRoleFeatureNotFound"`
	ErrorRoleFeatureNoId                 string `json:"ErrorRoleFeatureNoId"`
	RoleFeatureDeleteMany                string `json:"RoleFeatureDeleteMany"`
	ErrorRoleFeatureCopy                 string `json:"ErrorRoleFeatureCopy"`
	RoleFeatureCopySuccessful            string `json:"RoleFeatureCopySuccessful"`
	ServerSettingSaveSuccess             string `json:"ServerSettingSaveSuccess"`
	ServerSettingSaveFail                string `json:"ServerSettingSaveFail"`
	ErrorSiteSave                        string `json:"ErrorSiteSave"`
	ErrorSiteDelete                      string `json:"ErrorSiteDelete"`
	SiteSaveSuccess                      string `json:"SiteSaveSuccess"`
	SiteDeleteSuccessful                 string `json:"SiteDeleteSuccessful"`
	ErrorSiteNotFound                    string `json:"ErrorSiteNotFound"`
	ErrorSiteNoId                        string `json:"ErrorSiteNoId"`
	ErrorTransactionSave                 string `json:"ErrorTransactionSave"`
	ErrorTransactionDelete               string `json:"ErrorTransactionDelete"`
	TransactionSaveSuccess               string `json:"TransactionSaveSuccess"`
	TransactionDeleteSuccessful          string `json:"TransactionDeleteSuccessful"`
	ErrorTransactionNotFound             string `json:"ErrorTransactionNotFound"`
	ErrorTransactionNoId                 string `json:"ErrorTransactionNoId"`
	ErrorUserSave                        string `json:"ErrorUserSave"`
	ErrorUserDelete                      string `json:"ErrorUserDelete"`
	UserSaveSuccess                      string `json:"UserSaveSuccess"`
	UserDeleteSuccessful                 string `json:"UserDeleteSuccessful"`
	ErrorUserNotFound                    string `json:"ErrorUserNotFound"`
	ErrorUserNoId                        string `json:"ErrorUserNoId"`
	UserAddEditEmail                     string `json:"UserAddEditEmail"`
	UserAddEditFirstName                 string `json:"UserAddEditFirstName"`
	UserAddEditLastName                  string `json:"UserAddEditLastName"`
	UserAddEditPhone                     string `json:"UserAddEditPhone"`
	UserAddEditExtension                 string `json:"UserAddEditExtension"`
	UserAddEditMobile                    string `json:"UserAddEditMobile"`
	UserAddEditJobTitle                  string `json:"UserAddEditJobTitle"`
	UserAddEditOfficeName                string `json:"UserAddEditOfficeName"`
	UserAddEditDepartment                string `json:"UserAddEditDepartment"`
	UserAddEditUserBio                   string `json:"UserAddEditUserBio"`
	UserAddEditSkypeId                   string `json:"UserAddEditSkypeId"`
	UserAddEditUserProfileIcon           string `json:"UserAddEditUserProfileIcon"`
	UserAddEditSetDefaultAccount         string `json:"UserAddEditSetDefaultAccount"`
	ErrorEmailExists                     string `json:"ErrorEmailExists"`
	ErrorEnterValidEmail                 string `json:"ErrorEnterValidEmail"`
	ErrorEnterFirst                      string `json:"ErrorEnterFirst"`
	ErrorEnterLast                       string `json:"ErrorEnterLast"`
	ErrorEnterPassAndConfirm             string `json:"ErrorEnterPassAndConfirm"`
	ErrorPasswordsDoNotMatch             string `json:"ErrorPasswordsDoNotMatch"`
	ErrorPasswordNotValid                string `json:"ErrorPasswordNotValid"`
	UserAddSuccess                       string `json:"UserAddSuccess"`
	ErrorUserListFailedParse             string `json:"ErrorUserListFailedParse"`
	ErrorUserListFailedRevoke            string `json:"ErrorUserListFailedRevoke"`
	UserListRevokedUser                  string `json:"UserListRevokedUser"`

	//AdditionalConstructs
}

type queryAppContent struct{}

type TagReplacements struct {
	Tag1  map[string]string
	Tag2  map[string]string
	Tag3  map[string]string
	Tag4  map[string]string
	Tag5  map[string]string
	Tag6  map[string]string
	Tag7  map[string]string
	Tag8  map[string]string
	Tag9  map[string]string
	Tag10 map[string]string
}

func Q(k string, v string) map[string]string {
	var ret map[string]string
	ret = make(map[string]string)
	ret[k] = v
	return ret
}

func (self *queryAppContent) GetTranslationWithReplacementsFromUser(user model.User, key string, tags *TagReplacements) (originalReplacement string) {
	originalReplacement = self.GetTranslationFromUser(user, key)
	if originalReplacement != "" && tags != nil {
		for k, v := range tags.Tag1 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag2 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag3 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag4 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag5 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag6 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag7 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag8 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag9 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag10 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
	}
	return originalReplacement
}

func (self *queryAppContent) GetTranslationWithReplacements(context session_functions.RequestContext, key string, tags *TagReplacements) (originalReplacement string) {
	user, err := getUser(context)
	if err != nil {
		return
	}
	return self.GetTranslationWithReplacementsFromUser(user, key, tags)
}

func (self *queryAppContent) GetTranslationFromUser(user model.User, key string) (translatedText string) {
	contentData, err := extensions.ReadFile(settings.WebRoot + "/globalization/translations/app/" + user.Language + "/" + user.Language + ".json")
	if err != nil {
		contentData, err = extensions.ReadFile(settings.WebRoot + "/globalization/translations/app/en/US.json")
	}
	var appContent AppContentJson
	err = json.Unmarshal(contentData, &appContent)
	if err != nil {
		return ""
	}
	translatedText, _ = reflections.GetFieldAsString(appContent, key)
	return translatedText
}

func (self *queryAppContent) GetTranslation(context session_functions.RequestContext, key string) (translatedText string) {

	user, err := getUser(context)
	if err != nil { //Just get the English Language for no context
		var usr model.User
		usr.Language = "en"
		translatedText = self.GetTranslationFromUser(usr, key)
		return
	}
	return self.GetTranslationFromUser(user, key)
}
