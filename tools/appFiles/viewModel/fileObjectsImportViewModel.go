package viewModel

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"strings"
)

type FileObjectImport struct {
	Id           ImportField
	Name         ImportField
	Content      ImportField
	Size         ImportField
	Type         ImportField
	ModifiedUnix ImportField
	Modified     ImportField
	MD5          ImportField
}

func (this *FileObjectImport) ValidateRows(context session_functions.RequestContext, rows [][]string) ([]string, [][]string, [][]string) {
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

func (this *FileObjectImport) ValidateCustom(context session_functions.RequestContext, row []string) (returnErrors []string, newRow []string) {
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

func (this *FileObjectImport) LoadSchema(context session_functions.RequestContext) (returnFields []ImportField) {
	this.Id = importFieldId(context)
	returnFields = append(returnFields, this.Id)

	var i int

	i += 1
	this.Name.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldName")
	this.Name.Idx = i
	//this.Name.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldNameHelp")
	returnFields = append(returnFields, this.Name)

	i += 1
	this.Content.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldContent")
	this.Content.Idx = i
	//this.Content.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldContentHelp")
	returnFields = append(returnFields, this.Content)

	i += 1
	this.Size.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldSize")
	this.Size.Idx = i
	//this.Size.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldSizeHelp")
	returnFields = append(returnFields, this.Size)

	i += 1
	this.Type.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldType")
	this.Type.Idx = i
	//this.Type.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldTypeHelp")
	returnFields = append(returnFields, this.Type)

	i += 1
	this.ModifiedUnix.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldModifiedUnix")
	this.ModifiedUnix.Idx = i
	//this.ModifiedUnix.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldModifiedUnixHelp")
	returnFields = append(returnFields, this.ModifiedUnix)

	i += 1
	this.Modified.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldModified")
	this.Modified.Idx = i
	//this.Modified.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldModifiedHelp")
	returnFields = append(returnFields, this.Modified)

	i += 1
	this.MD5.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldMD5")
	this.MD5.Idx = i
	//this.MD5.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFileObjectFieldMD5Help")
	returnFields = append(returnFields, this.MD5)

	return returnFields
}

func (this *FileObjectImport) LoadSchemaAndParseFile(context session_functions.RequestContext, fileContent string) ([][]string, error) {
	this.LoadSchema(context)
	contents, err := parseFile(context, this, fileContent)
	return contents, err
}
