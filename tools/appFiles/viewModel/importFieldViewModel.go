package viewModel

import (
	"encoding/csv"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/davidrenne/reflections"
	"github.com/go-errors/errors"
	"io"
	"strings"
)

type ImportField struct {
	CsvHeader string
	CsvHelp   string
	Idx       int
	Required  bool
}

func importFieldId(context session_functions.RequestContext) ImportField {
	var ret ImportField
	ret.CsvHeader = queries.AppContent.GetTranslation(context, "CSVFieldId")
	ret.Idx = 0
	ret.Required = IMPORT_NOT_REQUIRED
	ret.CsvHelp = queries.AppContent.GetTranslation(context, "CSVFieldIdHelp")

	return ret
}

func checkForMissingFields(context session_functions.RequestContext, obj interface{}, rows [][]string) (errors []string, invalidRows [][]string) {

	var err error
	var fields []string
	fields, err = reflections.FieldsDeep(obj)

	for lineNumber, row := range rows {
		if len(row) > 0 {
			var fieldPointer interface{}
			var valueIndex int64
			var valueRequired bool
			var valueName string
			var valueHelp string
			realLine := lineNumber + 2
			var allFieldErrors []string
			if err == nil {
				for _, field := range fields {
					fieldPointer, _ = reflections.GetField(obj, field)
					valueRequired, _ = reflections.GetFieldAsBool(fieldPointer, "Required")
					valueIndex, _ = reflections.GetFieldAsInt(fieldPointer, "Idx")
					valueName, _ = reflections.GetFieldAsString(fieldPointer, "CsvHeader")
					valueHelp, _ = reflections.GetFieldAsString(fieldPointer, "CsvHelp")
					if valueRequired && row[valueIndex] == "" {
						var help string
						var required string
						if valueHelp != "" {
							help = " (" + valueHelp + ")"
						}
						if valueRequired {
							required = queries.AppContent.GetTranslation(context, "CSVRequired") + " "
						}
						allFieldErrors = append(allFieldErrors, "\""+required+valueName+help+"\"")
					}

				}

				if len(allFieldErrors) > 0 {
					var requiredStr string
					if len(allFieldErrors) == 1 {
						requiredStr = "CSVLineNumberRequiredSingular"
					} else {
						requiredStr = "CSVLineNumberRequiredPlural"
					}
					replacements := queries.TagReplacements{
						Tag1: queries.Q("line_number", extensions.IntToString(realLine)),
						Tag2: queries.Q("csv_fields_errored", strings.Join(allFieldErrors, ", ")),
					}
					errors = append(errors, queries.AppContent.GetTranslationWithReplacements(context, requiredStr, &replacements))
					invalidRows = append(invalidRows, row)
				}
			} else {
				replacements := queries.TagReplacements{
					Tag1: queries.Q("line_number", extensions.IntToString(realLine)),
				}
				errors = append(errors, queries.AppContent.GetTranslationWithReplacements(context, "CSVLineNumberRrBackendFailure", &replacements))
				invalidRows = append(invalidRows, row)
			}

			if len(allFieldErrors) == 0 {
				errors = append(errors, "")
			}
		}
	}
	return errors, invalidRows
}

func GetCSVHeaderArray(context session_functions.RequestContext, schema []ImportField) (headers []string) {
	for _, field := range schema {
		var header string
		if field.Required {
			header += queries.AppContent.GetTranslation(context, "CSVFieldRequired") + " "
		}
		header += field.CsvHeader
		if field.CsvHelp != "" {
			header += " (" + field.CsvHelp + ")"
		}
		headers = append(headers, header)
	}
	return headers
}

func GetCSVTemplate(context session_functions.RequestContext, schema []ImportField) (headers string) {
	for _, field := range schema {
		headers += "\""
		if field.Required {
			headers += queries.AppContent.GetTranslation(context, "CSVFieldRequired") + " "
		}
		headers += field.CsvHeader
		if field.CsvHelp != "" {
			headers += " (" + field.CsvHelp + ")"
		}
		headers += "\","
	}
	return headers[:len(headers)-1]
}

func parseFile(context session_functions.RequestContext, obj interface{}, fileContents string) ([][]string, error) {
	fileContents = strings.Replace(fileContents, "\r", "", -1)
	lines := strings.Split(fileContents, "\n")
	var csvRows [][]string

	r := csv.NewReader(strings.NewReader(string(lines[0])))
	var header []string
	var err error
	header, err = r.Read()
	if err != nil {
		return csvRows, err
	}

	var fields []string
	fields, err = reflections.FieldsDeep(obj)

	if len(fields) != len(header) {
		return csvRows, errors.New(queries.AppContent.GetTranslation(context, "CSVOutdated"))
	}

	r = csv.NewReader(strings.NewReader(strings.Join(lines[1:], "\n")))
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = len(fields)

	for {
		record, err := r.Read()
		if len(record) > 0 {
			if len(header) > len(record) {
				for i := len(record) + 1; i <= r.FieldsPerRecord; i++ {
					record = append(record, "")
				}
			}
			csvRows = append(csvRows, record)
		}
		if err == io.EOF {
			break
		}
		if err != nil && strings.Index(err.Error(), "wrong number of fields") == -1 {
			return csvRows, err
		}
	}

	return csvRows, nil
}
