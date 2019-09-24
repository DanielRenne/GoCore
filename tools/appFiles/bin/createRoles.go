package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
)

func main() {
	log.SetOutput(ioutil.Discard)
	os.Chdir(os.Getenv("GOPATH"))
	outputNewRoles := true
	app.Initialize("src/github.com/DanielRenne/goCoreAppTemplate", "webConfig.json")
	settings.Initialize()
	dbServices.Initialize()
	log.SetOutput(ioutil.Discard)

	var capCamel string
	var lowerCamel string
	var capCamelPlural string

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Object Name Singular [i.e. 'Manufacturer']: ")
		capCamel, _ = reader.ReadString('\n')
		capCamel = strings.Trim(strings.Trim(capCamel, "\n"), "\r")
		break
	}

	lowerCamel = strings.ToLower(capCamel)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Plural [i.e. 'manufacturers']: ")
		capCamelPlural, _ = reader.ReadString('\n')
		capCamelPlural = strings.Trim(strings.ToLower(strings.Trim(capCamelPlural, "\n")), "\r")
		break
	}

	fieldMapping := make(map[string]string, 0)
	fieldMapping["view"] = "y"
	fieldMapping["add"] = "y"
	fieldMapping["modify"] = "y"
	fieldMapping["delete"] = "y"
	fieldMapping["export"] = "y"
	fieldMapping["copy"] = "y"
	for k, _ := range fieldMapping {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Allow role [" + k + "] defaults to \"y\": ")
			tmpType, _ := reader.ReadString('\n')
			tmpType = strings.Trim(strings.Trim(tmpType, "\n"), "\r")
			ok := false
			if tmpType == "" {
				tmpType = "y"
			}
			ok = true
			fieldMapping[k] = tmpType
			if ok {
				break
			}
		}
	}
	logger.Message(fieldMapping["copy"], logger.GREEN)
	logger.Message(fieldMapping["add"], logger.GREEN)
	logger.Message(fieldMapping["modify"], logger.GREEN)
	logger.Message(fieldMapping["view"], logger.GREEN)
	logger.Message(fieldMapping["delete"], logger.GREEN)
	logger.Message(fieldMapping["export"], logger.GREEN)
	type Custom struct {
		add         string
		action      string
		key         string
		name        string
		description string
	}
	fieldMappingCustom := make(map[string]*Custom, 0)
	fieldMappingCustom["custom1"] = &Custom{
		add: "n",
	}
	fieldMappingCustom["custom2"] = &Custom{
		add: "n",
	}
	fieldMappingCustom["custom3"] = &Custom{
		add: "n",
	}
	fieldMappingCustom["custom4"] = &Custom{
		add: "n",
	}
	fieldMappingCustom["custom5"] = &Custom{
		add: "n",
	}
	for k, _ := range fieldMappingCustom {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("[" + k + "] defaults to \"n\": ")
			tmpType, _ := reader.ReadString('\n')
			tmpType = strings.Trim(tmpType, "\n")
			if tmpType == "" {
				tmpType = "n"
			}
			ok := false
			if tmpType == "y" || tmpType == "n" {
				ok = true
				fieldMappingCustom[k].add = tmpType
				if fieldMappingCustom[k].add == "y" {
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Type the action or feature: ")
					tmpType, _ := reader.ReadString('\n')
					tmpType = strings.Trim(tmpType, "\n")
					fieldMappingCustom[k].action = strings.ToLower(tmpType)
					fieldMappingCustom[k].key = strings.ToUpper(strings.Replace(lowerCamel, " ", "", -1)) + "_" + strings.ToUpper(fieldMappingCustom[k].action)
					fieldMappingCustom[k].name = strings.Title(fieldMappingCustom[k].action) + " " + capCamelPlural
					fieldMappingCustom[k].description = "Ability to " + fieldMappingCustom[k].action + " " + capCamelPlural
				}
			} else {
				fmt.Println("Invalid type 'mongo' or 'bolt'")
			}
			if ok {
				break
			}
		}
	}
	t, err := model.Transactions.New(constants.APP_CONSTANTS_CRONJOB_ID)

	template := capCamel + " Related"
	var x []model.FeatureGroup
	count, _ := model.FeatureGroups.Query().Filter(model.Q(model.FIELD_FEATUREGROUP_NAME, template)).Count(&x)

	var f1 model.Feature
	var f2 model.Feature
	var f3 model.Feature
	var f4 model.Feature
	var f5 model.Feature
	var f6 model.Feature
	var f11 model.Feature
	var f12 model.Feature
	var f13 model.Feature
	var f14 model.Feature
	var f15 model.Feature
	var fg model.FeatureGroup
	upperSingular := strings.ToUpper(strings.Replace(lowerCamel, " ", "", -1))
	if outputNewRoles {
		if count == 0 {
			fg = model.FeatureGroup{
				Name: capCamel + " Related",
			}
			fg.SaveWithTran(t)
			t.Commit()
		} else {
			model.FeatureGroups.Query().Filter(model.Q("Name", capCamel+" Related")).One(&fg)
		}
		if fieldMapping["view"] == "y" {
			f1 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            upperSingular + "_VIEW",
				Name:           "View " + capCamelPlural,
				Description:    "Ability to view " + capCamelPlural,
			}
			f1.Save()
		}

		if fieldMapping["add"] == "y" {
			f2 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            upperSingular + "_ADD",
				Name:           "Add " + capCamelPlural,
				Description:    "Ability to add " + capCamelPlural,
			}
			err = f2.Save()
		}

		if fieldMapping["modify"] == "y" {
			f3 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            upperSingular + "_MODIFY",
				Name:           "Modify " + capCamelPlural,
				Description:    "Ability to modify " + capCamelPlural,
			}
			f3.Save()
		}

		if fieldMapping["delete"] == "y" {
			f4 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            upperSingular + "_DELETE",
				Name:           "Delete " + capCamelPlural,
				Description:    "Ability to delete " + capCamelPlural,
			}
			f4.Save()
		}

		if fieldMapping["export"] == "y" {
			f5 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            upperSingular + "_EXPORT",
				Name:           "Export " + capCamelPlural,
				Description:    "Ability to export " + capCamelPlural,
			}
			f5.Save()
		}

		if fieldMapping["copy"] == "y" {
			f6 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            upperSingular + "_COPY",
				Name:           "Copy " + lowerCamel,
				Description:    "Ability to copy a " + lowerCamel,
			}
			f6.Save()
		}
		var customkey string
		customkey = "custom1"
		if fieldMappingCustom[customkey].add == "y" {
			f11 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            fieldMappingCustom[customkey].key,
				Name:           fieldMappingCustom[customkey].name,
				Description:    fieldMappingCustom[customkey].description,
			}
			f11.Save()
		}
		customkey = "custom2"
		if fieldMappingCustom[customkey].add == "y" {
			f12 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            fieldMappingCustom[customkey].key,
				Name:           fieldMappingCustom[customkey].name,
				Description:    fieldMappingCustom[customkey].description,
			}
			f12.Save()
		}
		customkey = "custom3"
		if fieldMappingCustom[customkey].add == "y" {
			f13 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            fieldMappingCustom[customkey].key,
				Name:           fieldMappingCustom[customkey].name,
				Description:    fieldMappingCustom[customkey].description,
			}
			f13.Save()
		}
		customkey = "custom4"
		if fieldMappingCustom[customkey].add == "y" {
			f14 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            fieldMappingCustom[customkey].key,
				Name:           fieldMappingCustom[customkey].name,
				Description:    fieldMappingCustom[customkey].description,
			}
			f14.Save()
		}
		customkey = "custom5"
		if fieldMappingCustom[customkey].add == "y" {
			f15 = model.Feature{
				FeatureGroupId: fg.Id.Hex(),
				Key:            fieldMappingCustom[customkey].key,
				Name:           fieldMappingCustom[customkey].name,
				Description:    fieldMappingCustom[customkey].description,
			}
			f15.Save()
		}

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
		allIds := utils.Array()
		if fieldMapping["view"] == "y" {
			model.Features.Query().Where("Key", upperSingular+"_VIEW").One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		if fieldMapping["add"] == "y" {
			model.Features.Query().Where("Key", upperSingular+"_ADD").One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		if fieldMapping["modify"] == "y" {
			model.Features.Query().Where("Key", upperSingular+"_MODIFY").One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		if fieldMapping["delete"] == "y" {
			model.Features.Query().Where("Key", upperSingular+"_DELETE").One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		if fieldMapping["export"] == "y" {
			model.Features.Query().Where("Key", upperSingular+"_EXPORT").One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		if fieldMapping["copy"] == "y" {
			model.Features.Query().Where("Key", upperSingular+"_COPY").One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		customkey = "custom1"
		if fieldMappingCustom[customkey].add == "y" {
			model.Features.Query().Where("Key", fieldMappingCustom[customkey].key).One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		customkey = "custom2"
		if fieldMappingCustom[customkey].add == "y" {
			model.Features.Query().Where("Key", fieldMappingCustom[customkey].key).One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		customkey = "custom3"
		if fieldMappingCustom[customkey].add == "y" {
			model.Features.Query().Where("Key", fieldMappingCustom[customkey].key).One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		customkey = "custom4"
		if fieldMappingCustom[customkey].add == "y" {
			model.Features.Query().Where("Key", fieldMappingCustom[customkey].key).One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}
		customkey = "custom5"
		if fieldMappingCustom[customkey].add == "y" {
			model.Features.Query().Where("Key", fieldMappingCustom[customkey].key).One(&row)
			allIds = append(allIds, row.Id.Hex())
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
		}

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
		_ = model.Features.Query().In(model.Q(model.FIELD_FEATURE_ID, allIds)).All(&features)
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

	}
	log.SetOutput(os.Stdout)
	log.Println("Done with role creation")
}
