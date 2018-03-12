package viewModel

type ViewModel interface {
	LoadDefaultState()
}

func GetViewModel(key string) ViewModel {

	switch key {
	case VIEWMODEL_LOGIN:
		var vm LoginViewModel
		return &vm

	case VIEWMODEL_APP:
		var vm AppViewModel
		return &vm

	case VIEWMODEL_SETTINGS:
		var vm SettingsViewModel
		return &vm

	case VIEWMODEL_USERS:
		var vm UsersViewModel
		return &vm

	case VIEWMODEL_ACCOUNTS:
		var vm AccountsViewModel
		return &vm

	case VIEWMODEL_ACCOUNTLIST:
		var vm AccountListViewModel
		return &vm

	case VIEWMODEL_ACCOUNTMODIFY:
		var vm AccountModifyViewModel
		return &vm

	case VIEWMODEL_USERMODIFY:
		var vm UserModifyViewModel
		return &vm

	case VIEWMODEL_SERVERSETTINGSMODIFY:
		var vm ServerSettingsModifyViewModel
		return &vm

	case VIEWMODEL_USERLIST:
		var vm UserListViewModel
		return &vm

	case VIEWMODEL_ACCOUNTADD:
		var vm AccountModifyViewModel
		return &vm

	case VIEWMODEL_APPERRORMODIFY:
		var vm AppErrorModifyViewModel
		return &vm

	case VIEWMODEL_APPERRORLIST:
		var vm AppErrorListViewModel
		return &vm

	case VIEWMODEL_FILEOBJECT:
		var vm FileObjectViewModel
		return &vm

	case VIEWMODEL_FEATUREMODIFY:
		var vm FeatureModifyViewModel
		return &vm

	case VIEWMODEL_FEATURELIST:
		var vm FeatureListViewModel
		return &vm

	case VIEWMODEL_ROLEFEATUREMODIFY:
		var vm RoleFeatureModifyViewModel
		return &vm

	case VIEWMODEL_ROLEFEATURELIST:
		var vm RoleFeatureListViewModel
		return &vm

	case VIEWMODEL_FEATUREGROUPMODIFY:
		var vm FeatureGroupModifyViewModel
		return &vm

	case VIEWMODEL_FEATUREGROUPLIST:
		var vm FeatureGroupListViewModel
		return &vm

	case VIEWMODEL_ROLEMODIFY:
		var vm RoleModifyViewModel
		return &vm

	case VIEWMODEL_ROLELIST:
		var vm RoleListViewModel
		return &vm

	case VIEWMODEL_FILEOBJECTMODIFY:
		var vm FileObjectModifyViewModel
		return &vm

	case VIEWMODEL_FILEOBJECTLIST:
		var vm FileObjectListViewModel
		return &vm

		//-DONT-REMOVE-NEW-CASE
	}

	var vm AppViewModel
	return &vm
}
