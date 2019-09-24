package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
	"github.com/davidrenne/heredoc"
)

func main() {
	outputNewRoles := true

	//MoreAccountTabs
	//
	//if session_functions.CheckRoleAccess(context(), constants.FEATURE_ROLE_VIEW) {
	//	c = viewModel.SETTINGS_CONST_ROLE
	//	vm.ButtonBar.Config.VisibleTabs[c] = c
	//	vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
	//	vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ROLELIST)
	//	vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
	//	vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, true)
	//	vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, "")
	//}
	//c = viewModel.SETTINGS_CONST_ROLE_ADD
	//vm.ButtonBar.Config.VisibleTabs[c] = c
	//vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
	//vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ROLELIST)
	//vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
	//vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, false)
	//vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, viewModel.SETTINGS_CONST_ROLE)
	//
	//c = viewModel.SETTINGS_CONST_ROLE_MODIFY
	//vm.ButtonBar.Config.VisibleTabs[c] = c
	//vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
	//vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ROLELIST)
	//vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
	//vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, false)
	//vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, viewModel.SETTINGS_CONST_ROLE)

	//Manually set this if you want your table to include foreign keys as "fields" for import and widget lists
	app.Initialize("src/github.com/DanielRenne/goCoreAppTemplate", "webConfig.json")
	settings.Initialize()
	dbServices.Initialize()
	allowIdFieldsToBeShown := false
	capCamel := os.Args[1]
	lowerCamel := os.Args[2]
	capCamelPlural := os.Args[3]
	var err error
	db := model.ResolveEntity(capCamel)
	reflectedFields := db.Reflect()
	hasView := false
	var viewFields []model.Field
	var inputFields []model.Field
	for i := 0; i < len(reflectedFields); i++ {
		if reflectedFields[i].Name == "Id" {
			continue
		}
		if reflectedFields[i].IsView {
			hasView = true
			viewFields = append(viewFields, reflectedFields[i])
		} else {
			inputFields = append(inputFields, reflectedFields[i])
		}
	}

	if !hasView {
		for i := 0; i < len(reflectedFields); i++ {
			viewFields = append(viewFields, reflectedFields[i])
		}
	}

	t, err := model.Transactions.New(constants.APP_CONSTANTS_CRONJOB_ID)

	template := capCamel + " Related"
	var featureGroups []model.FeatureGroup
	count, _ := model.FeatureGroups.Query().Filter(model.Q(model.FIELD_FEATUREGROUP_NAME, template)).Count(&featureGroups)
	if count < 1 && outputNewRoles {
		var fg model.FeatureGroup
		fg = model.FeatureGroup{
			Name: capCamel + " Related",
		}
		fg.SaveWithTran(t)

		f1 := model.Feature{
			FeatureGroupId: fg.Id.Hex(),
			Key:            strings.ToUpper(lowerCamel) + "_VIEW",
			Name:           "View " + capCamelPlural,
			Description:    "Ability to view " + capCamelPlural,
		}
		f1.SaveWithTran(t)

		f2 := model.Feature{
			FeatureGroupId: fg.Id.Hex(),
			Key:            strings.ToUpper(lowerCamel) + "_ADD",
			Name:           "Add " + capCamelPlural,
			Description:    "Ability to add " + capCamelPlural,
		}
		f2.SaveWithTran(t)

		f3 := model.Feature{
			FeatureGroupId: fg.Id.Hex(),
			Key:            strings.ToUpper(lowerCamel) + "_MODIFY",
			Name:           "Modify " + capCamelPlural,
			Description:    "Ability to modify " + capCamelPlural,
		}
		f3.SaveWithTran(t)

		f4 := model.Feature{
			FeatureGroupId: fg.Id.Hex(),
			Key:            strings.ToUpper(lowerCamel) + "_DELETE",
			Name:           "Delete " + capCamelPlural,
			Description:    "Ability to delete " + capCamelPlural,
		}
		f4.SaveWithTran(t)

		f5 := model.Feature{
			FeatureGroupId: fg.Id.Hex(),
			Key:            strings.ToUpper(lowerCamel) + "_EXPORT",
			Name:           "Export " + capCamelPlural,
			Description:    "Ability to export " + capCamelPlural,
		}
		f5.SaveWithTran(t)

		f6 := model.Feature{
			FeatureGroupId: fg.Id.Hex(),
			Key:            strings.ToUpper(lowerCamel) + "_COPY",
			Name:           "Copy " + lowerCamel,
			Description:    "Ability to copy a " + lowerCamel,
		}
		f6.SaveWithTran(t)
		t.Commit()

		var all []model.FeatureGroup
		model.FeatureGroups.Query().All(&all)

		for i, _ := range all {
			all[i].BootstrapMeta = &model.BootstrapMeta{
				AlwaysUpdate: true,
			}
			all[i].LastUpdateId = constants.APP_CONSTANTS_USERS_ANONYMOUS_ID
		}

		strjson, _ := json.MarshalIndent(all, "", "\t")
		err = ioutil.WriteFile(serverSettings.APP_LOCATION+"/db/bootstrap/featureGroups/featureGroups.json", []byte(strjson), 0644)
		if err != nil {
			return
		}

		var row model.Feature
		systemRoleFile := serverSettings.APP_LOCATION + "/constants/systemRoles.go"
		model.Features.Query().ById(f1.Id, &row)
		utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		model.Features.Query().ById(f2.Id, &row)
		utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		model.Features.Query().ById(f3.Id, &row)
		utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		model.Features.Query().ById(f4.Id, &row)
		utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		model.Features.Query().ById(f5.Id, &row)
		utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		model.Features.Query().ById(f6.Id, &row)
		utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")

		var allFeature []model.Feature
		model.Features.Query().All(&allFeature)

		for i, _ := range allFeature {
			allFeature[i].BootstrapMeta = &model.BootstrapMeta{
				AlwaysUpdate: true,
			}
			allFeature[i].LastUpdateId = constants.APP_CONSTANTS_USERS_ANONYMOUS_ID
		}
		strjson2, _ := json.MarshalIndent(allFeature, "", "\t")
		err = ioutil.WriteFile(serverSettings.APP_LOCATION+"/db/bootstrap/features/features.json", []byte(strjson2), 0644)

		var features []model.Feature
		_ = model.Features.Query().In(model.Q(model.FIELD_FEATURE_ID, utils.Array(f1.Id.Hex(), f2.Id.Hex(), f3.Id.Hex(), f4.Id.Hex(), f5.Id.Hex(), f6.Id.Hex()))).All(&features)
		var roles []model.Role
		_ = model.Roles.Query().All(&roles)

		t, _ := model.Transactions.New(constants.APP_CONSTANTS_CRONJOB_ID)
		for _, feature := range features {
			for _, role := range roles {
				rf := model.RoleFeature{
					FeatureId: feature.Id.Hex(),
					RoleId:    role.Id.Hex(),
				}
				rf.SaveWithTran(t)
			}
		}
		t.Commit()

		var allRf []model.RoleFeature
		model.RoleFeatures.Query().All(&allRf)
		for i, _ := range allRf {
			allRf[i].BootstrapMeta = &model.BootstrapMeta{
				AlwaysUpdate: true,
			}
			allRf[i].LastUpdateId = constants.APP_CONSTANTS_USERS_ANONYMOUS_ID
		}
		strjson2222, _ := json.MarshalIndent(allRf, "", "\t")
		ioutil.WriteFile(serverSettings.APP_LOCATION+"/db/bootstrap/roleFeatures/roleFeatures.json", []byte(strjson2222), 0644)

		importPopup := serverSettings.APP_LOCATION + "/web/app/javascript/components/addOrImportComponent.js"
		utils.ReplaceTokenInFile(importPopup, "//AdditionalPages", heredoc.Docf(`
      if (Action == "CONTROLLER_%sADD" && !this.globs.HasRole("%s_ADD")) {
        return <div><Component/></div>;
      }
      //AdditionalPages`, strings.ToUpper(lowerCamel), strings.ToUpper(lowerCamel)))
	}

	serverSettings.Initialize(os.Getenv("GOPATH")+"/src/github.com/DanielRenne/goCoreAppTemplate", "webConfig.json")
	post := serverSettings.APP_LOCATION + "/controllers/" + capCamelPlural + "PostController.go"
	vm := serverSettings.APP_LOCATION + "/viewModel/" + capCamelPlural + "ImportViewModel.go"
	list := serverSettings.APP_LOCATION + "/web/app/javascript/pages/" + lowerCamel + "List/" + lowerCamel + "ListComponents.js"
	modify := serverSettings.APP_LOCATION + "/web/app/javascript/pages/" + lowerCamel + "Modify/" + lowerCamel + "ModifyComponents.js"
	br := serverSettings.APP_LOCATION + "/controllers/" + lowerCamel + ".go"

	editFields := ""
	editEvents := ""
	hasNameField := false
	for i := 0; i < len(inputFields); i++ {
		if inputFields[i].Name == "Name" {
			hasNameField = true
		}
		if inputFields[i].Name != "LastUpdateId" && inputFields[i].Name != "CreateDate" && inputFields[i].Name != "UpdateDate" && (allowIdFieldsToBeShown || (!allowIdFieldsToBeShown && strings.Index(inputFields[i].Name, "Id") == -1)) && inputFields[i].Name != "Slug" && inputFields[i].Name != "RWMutex" {
			required := ""
			if inputFields[i].Validation != nil && inputFields[i].Validation.Required {
				required = "\"* \" + "
			}
			editFields += heredoc.Docf(`

                <TextField
                  floatingLabelText={%swindow.pageContent.%sModify%s || window.pageContent.%sAdd%s}
                  hintText={%swindow.pageContent.%sModify%s || window.pageContent.%sAdd%s}
                  fullWidth={true}
                  onChange={this.handle%sChange}
                  errorText={this.globs.translate(this.state.%s.Errors.%s)}
                  value={this.state.%s.%s}
                />
                <br />`, required, capCamel, inputFields[i].Name, capCamel, inputFields[i].Name, required, capCamel, inputFields[i].Name, capCamel, inputFields[i].Name, inputFields[i].Name, capCamel, inputFields[i].Name, capCamel, inputFields[i].Name)
			editEvents += heredoc.Docf(`

    this.handle%sChange = (event) => {
      this.setComponentState({%s: {
        %s: event.target.value,
        Errors: {%s: ""}
      }});
    };

`, inputFields[i].Name, capCamel, inputFields[i].Name, inputFields[i].Name)

		}
		os.Setenv("TMP_CONTROLLER", lowerCamel+"Modify")
		os.Setenv("TMP_TRANSLATIONKEY", capCamel+"Modify"+inputFields[i].Name)
		os.Setenv("TMP_TRANSLATION", inputFields[i].Label)
		_, err = exec.Command("bash", "-c", "add_env_translation").Output()
		if err != nil {
			log.Fatal(err)
		}
		os.Setenv("TMP_CONTROLLER", lowerCamel+"Add")
		os.Setenv("TMP_TRANSLATIONKEY", capCamel+"Add"+inputFields[i].Name)
		os.Setenv("TMP_TRANSLATION", inputFields[i].Label)
		_, err = exec.Command("bash", "-c", "add_env_translation").Output()
		if err != nil {
			log.Fatal(err)
		}
	}

	utils.ReplaceTokenInFile(modify, "-TEXT_FIELDS-", editFields)
	utils.ReplaceTokenInFile(modify, "-HANDLE_METHODS-", editEvents)
	if hasNameField {
		utils.ReplaceTokenInFile(br, "-COPY_NAME-", `
		// Copy Name
		replacements := queries.TagReplacements{
			Tag1: queries.Q("old_name_of_row", copyRowVm.`+capCamel+`.Name),
		}
		copyRowVm.`+capCamel+`.Name = queries.AppContent.GetTranslationWithReplacements(context, "CopyRow", &replacements)
		`)
		utils.ReplaceTokenInFile(br, "-QUERIES-", `"github.com/DanielRenne/goCoreAppTemplate/queries"`)
	} else {
		utils.ReplaceTokenInFile(br, "-COPY_NAME-", "")
		utils.ReplaceTokenInFile(br, "-QUERIES-", "")
	}

	listFields := ""
	fieldsCSV := ""
	importFields := ""
	importFieldSetters := ""
	var fieldsInlineCSV []string
	outputField := true
	for i := 0; i < len(inputFields); i++ {
		outputField = true
		if inputFields[i].Name == "CreateDate" {
			outputField = false

		} else if inputFields[i].Name == "UpdateDate" {
			listFields += `
                {
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdated,
                  stateKey: "Views.UpdateFromNow",
                  sortOn: "UpdateDate",
                  tooltipKey: "Views.UpdateDate"
                },
`
		} else if inputFields[i].Name == "LastUpdateId" {
			listFields += heredoc.Docf(`
                {
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdatedBy,
                  sortable: false,
                  stateKey: "Joins.LastUpdateUser.Views.FullName"
                },
`)
		} else if inputFields[i].Name != "LastUpdateId" && (allowIdFieldsToBeShown || (!allowIdFieldsToBeShown && strings.Index(inputFields[i].Name, "Id") == -1)) && inputFields[i].Name != "Slug" && inputFields[i].Name != "BootstrapMeta" && inputFields[i].Name != "RWMutex" {
			fieldsCSV += "vm." + capCamel + "." + inputFields[i].Name + " = row[i." + inputFields[i].Name + ".Idx]\n"
			importFields += "\t" + inputFields[i].Name + "       ImportField\n"
			csvFieldTranslationKey := "CSV" + capCamel + "Field" + inputFields[i].Name
			requiredStr := ""
			if inputFields[i].Validation != nil && inputFields[i].Validation.Required {
				requiredStr = "\n\nthis." + inputFields[i].Name + ".Required = IMPORT_REQUIRED\n"
			}

			os.Setenv("TMP_CONTROLLER", "app")
			os.Setenv("TMP_TRANSLATIONKEY", csvFieldTranslationKey)
			os.Setenv("TMP_TRANSLATION", inputFields[i].Name)
			_, err = exec.Command("bash", "-c", "add_env_translation").Output()
			if err != nil {
				log.Fatal(err)
			}

			// stub out the help version.  i have no clue what someone wants for help but i will comment it out
			os.Setenv("TMP_CONTROLLER", "app")
			os.Setenv("TMP_TRANSLATIONKEY", csvFieldTranslationKey+"Help")
			os.Setenv("TMP_TRANSLATION", inputFields[i].Name+" Add Your Help Here")
			_, err = exec.Command("bash", "-c", "add_env_translation").Output()
			if err != nil {
				log.Fatal(err)
			}

			importFieldSetters += heredoc.Docf(`

	i += 1
	this.%s.CsvHeader = queries.AppContent.GetTranslation(context, "%s")
	this.%s.Idx = i%s
	//this.%s.CsvHelp = queries.AppContent.GetTranslation(context, "%s")
	returnFields = append(returnFields, this.%s)

			`, inputFields[i].Name, csvFieldTranslationKey, inputFields[i].Name, requiredStr, inputFields[i].Name, csvFieldTranslationKey+"Help", inputFields[i].Name)
			fieldsInlineCSV = append(fieldsInlineCSV, "row."+inputFields[i].Name)
			listFields += heredoc.Docf(`
                {
                  tooltip: window.pageContent.%sListToolTip%s,
                  headerDisplay: window.pageContent.%sListHeader%s,
                  sortable: true,
                  stateKey: "%s"
                },
`, capCamel, inputFields[i].Name, capCamel, inputFields[i].Name, inputFields[i].Name)
		} else {
			outputField = false
		}
		if outputField {
			os.Setenv("TMP_CONTROLLER", lowerCamel+"List")
			os.Setenv("TMP_TRANSLATIONKEY", capCamel+"ListToolTip"+inputFields[i].Name)
			os.Setenv("TMP_TRANSLATION", inputFields[i].Label)
			_, err = exec.Command("bash", "-c", "add_env_translation").Output()
			if err != nil {
				log.Fatal(err)
			}

			os.Setenv("TMP_CONTROLLER", lowerCamel+"List")
			os.Setenv("TMP_TRANSLATIONKEY", capCamel+"ListHeader"+inputFields[i].Name)
			os.Setenv("TMP_TRANSLATION", inputFields[i].Label)
			_, err = exec.Command("bash", "-c", "add_env_translation").Output()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	utils.ReplaceTokenInFile(list, "-LIST_FIELDS-", listFields)
	utils.ReplaceTokenInFile(post, "-CSVFIELDSGOLANGSETTERS-", fieldsCSV)
	utils.ReplaceTokenInFile(post, "-CSVFIELDSALL-", "record := []string{row.Id.Hex(), "+strings.Join(fieldsInlineCSV, ", ")+"}")
	utils.ReplaceTokenInFile(vm, "-IMPORT_FIELDS-", importFields)
	utils.ReplaceTokenInFile(vm, "-IMPORT_FIELD_SETTERS-", importFieldSetters)
	println("All Done with Field Injections and Translation Creation!")
}
