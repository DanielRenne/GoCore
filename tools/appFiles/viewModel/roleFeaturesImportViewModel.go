package viewModel

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"strings"
)

type RoleFeatureImport struct {
	Id        ImportField
	RoleId    ImportField
	FeatureId ImportField
}

func (this *RoleFeatureImport) ValidateRows(context session_functions.RequestContext, rows [][]string) ([]string, [][]string, [][]string) {
	var validRows [][]string
	errors, invalidRows := checkForMissingFields(context, this, rows)

	for lineNumber, row := range rows {
		if len(row) > 0 {
			realLine := lineNumber + 2
			errorMessages, newRow := this.ValidateCustom(context, row)
			rows[lineNumber] = newRow
			if len(errorMessages) > 0 {
				var lineDesc string
				errorMessage := lineDesc + strings.Join(errorMessages, ". ")
				if errors[lineNumber] == "" {
					invalidRows = append(invalidRows, row)
					lineDesc = "Line (" + extensions.IntToString(realLine) + "): "
					errors = append(errors, lineDesc+errorMessage)
				} else {
					errors[lineNumber] = errors[lineNumber] + ".  " + errorMessage
				}
			} else if errors[lineNumber] == "" {
				validRows = append(validRows, row)
			}
		}
	}

	var finalErrors []string
	for _, row := range errors {
		if row != "" {
			finalErrors = append(finalErrors, row)
		}
	}

	return finalErrors, invalidRows, validRows
}

func (this *RoleFeatureImport) ValidateCustom(context session_functions.RequestContext, row []string) (returnErrors []string, newRow []string) {
	newRow = row
	//if newRow[this.SubmitInvitation.Idx] == "1" && newRow[this.RoleType.Idx] == "" {
	//	replacements := queries.TagReplacements{
	//		Tag1: queries.Q("csv_header", this.RoleType.CsvHeader),
	//		Tag2: queries.Q("csv_help", this.RoleType.CsvHelp),
	//	}
	//	returnErrors = append(returnErrors, queries.AppContent.GetTranslationWithReplacements(context, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", &replacements))
	//}
	return returnErrors, newRow
}

func (this *RoleFeatureImport) LoadSchema(context session_functions.RequestContext) (returnFields []ImportField) {
	this.Id = importFieldId(context)
	returnFields = append(returnFields, this.Id)

	var i int

	i += 1
	this.RoleId.CsvHeader = queries.AppContent.GetTranslation(context, "CSVRoleFeatureFieldRoleId")
	this.RoleId.Idx = i

	this.RoleId.Required = IMPORT_REQUIRED

	//this.RoleId.CsvHelp = queries.AppContent.GetTranslation(context, "CSVRoleFeatureFieldRoleIdHelp")
	returnFields = append(returnFields, this.RoleId)

	i += 1
	this.FeatureId.CsvHeader = queries.AppContent.GetTranslation(context, "CSVRoleFeatureFieldFeatureId")
	this.FeatureId.Idx = i

	this.FeatureId.Required = IMPORT_REQUIRED

	//this.FeatureId.CsvHelp = queries.AppContent.GetTranslation(context, "CSVRoleFeatureFieldFeatureIdHelp")
	returnFields = append(returnFields, this.FeatureId)

	return returnFields
}

func (this *RoleFeatureImport) LoadSchemaAndParseFile(context session_functions.RequestContext, fileContent string) ([][]string, error) {
	this.LoadSchema(context)
	contents, err := parseFile(context, this, fileContent)
	return contents, err
}
