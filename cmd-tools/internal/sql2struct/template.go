package sql2struct

import (
	"fmt"
	"html/template"
	"os"

	"github.com/summerKK/go-code-snippet-library/cmd-tools/internal/word"
)

const strcutTpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{.Name | ToCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model *{{.TableName | ToCamelCase}}) TableName() string {
	return "{{.TableName}}"
}
`

type StructTemplate struct {
	structTpl string
}

type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type StructTempateDb struct {
	TableName string
	Columns   []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{structTpl: strcutTpl}
}

func (t *StructTemplate) AssemblyColumns(tableColumns []*TableColumn) []*StructColumn {
	structColumns := make([]*StructColumn, 0, len(tableColumns))
	for _, column := range tableColumns {
		structColumns = append(structColumns, &StructColumn{
			Name:    column.ColumnName,
			Type:    DbTypeToStructType[column.DataType],
			Tag:     fmt.Sprintf("`json:"+"%s"+"`", column.ColumnName),
			Comment: column.ColumnComment,
		})
	}

	return structColumns
}

func (t *StructTemplate) Generate(tableName string, structColumns []*StructColumn) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase": word.UnderscoreToUpperCamelCase,
	}).Parse(t.structTpl))

	structTempateDb := StructTempateDb{
		TableName: tableName,
		Columns:   structColumns,
	}

	err := tpl.Execute(os.Stdout, structTempateDb)

	return err
}
