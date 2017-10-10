package viewModel

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"strings"
)

type FeatureImport struct {
	Id             ImportField
	Key            ImportField
	Name           ImportField
	Description    ImportField
	FeatureGroupId ImportField
}

func (this *FeatureImport) ValidateRows(context session_functions.RequestContext, rows [][]string) ([]string, [][]string, [][]string) {
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

func (this *FeatureImport) ValidateCustom(context session_functions.RequestContext, row []string) (returnErrors []string, newRow []string) {
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

func (this *FeatureImport) LoadSchema(context session_functions.RequestContext) (returnFields []ImportField) {
	this.Id = importFieldId(context)
	returnFields = append(returnFields, this.Id)

	var i int

	i += 1
	this.Key.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFeatureFieldKey")
	this.Key.Idx = i

	this.Key.Required = IMPORT_REQUIRED

	//this.Key.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFeatureFieldKeyHelp")
	returnFields = append(returnFields, this.Key)

	i += 1
	this.Name.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFeatureFieldName")
	this.Name.Idx = i

	this.Name.Required = IMPORT_REQUIRED

	//this.Name.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFeatureFieldNameHelp")
	returnFields = append(returnFields, this.Name)

	i += 1
	this.Description.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFeatureFieldDescription")
	this.Description.Idx = i

	this.Description.Required = IMPORT_REQUIRED

	//this.Description.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFeatureFieldDescriptionHelp")
	returnFields = append(returnFields, this.Description)

	i += 1
	this.FeatureGroupId.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFeatureFieldFeatureGroupId")
	this.FeatureGroupId.Idx = i

	this.FeatureGroupId.Required = IMPORT_REQUIRED

	//this.FeatureGroupId.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFeatureFieldFeatureGroupIdHelp")
	returnFields = append(returnFields, this.FeatureGroupId)

	return returnFields
}

func (this *FeatureImport) LoadSchemaAndParseFile(context session_functions.RequestContext, fileContent string) ([][]string, error) {
	this.LoadSchema(context)
	contents, err := parseFile(context, this, fileContent)
	return contents, err
}
