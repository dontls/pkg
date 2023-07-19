package excel

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

type excel struct {
	filename string
	file     *xlsx.File
	sheet    *xlsx.Sheet
}

func Sheet(v interface{}, sheet string) *excel {
	r := &excel{}
	typOf := reflect.TypeOf(v)
	if typOf.Kind() == reflect.Ptr {
		typOf = typOf.Elem()
	}
	r.filename = fmt.Sprintf("%s_%s.xlsx", sheet, time.Now().Format("2006-01-02 150405"))
	r.file = xlsx.NewFile()
	r.sheet, _ = r.file.AddSheet(sheet)
	titleRow := r.sheet.AddRow()
	titles := r.scanTitles(typOf)
	for _, v := range titles {
		titleRow.AddCell().Value = v
	}
	return r
}

func (o *excel) scanTitles(typOf reflect.Type) []string {
	var titles []string
	for i := 0; i < typOf.NumField(); i++ {
		f := typOf.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			titles = append(titles, o.scanTitles(f.Type)...)
		} else {
			titles = append(titles, f.Name)
		}
	}
	return titles
}

// Writes a struct to row r. Accepts a pointer to struct type 'e',
// and the number of columns to write, `cols`. If 'cols' is < 0,
// the entire struct will be written if possible. Returns -1 if the 'e'
// doesn't point to a struct, otherwise the number of columns written
func (o *excel) addRow(r *xlsx.Row, v reflect.Value, cols int) int {
	if cols == 0 {
		return cols
	}
	n := v.NumField() // number of fields in struct
	if cols < n && cols > 0 {
		n = cols
	}

	var k int
	for i := 0; i < n; i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.String, reflect.Int, reflect.Int8,
			reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float64, reflect.Float32:
			cell := r.AddCell()
			cell.SetValue(f.Interface())
		case reflect.Bool:
			cell := r.AddCell()
			cell.SetBool(f.Bool())
		default:
			if v.Type().Field(i).Anonymous {
				o.addRow(r, f, -1)
			} else {
				cell := r.AddCell()
				cell.SetValue(f.Interface())
			}
		}
	}

	return k
}

func (o *excel) Write(v interface{}) *excel {
	valOf := reflect.ValueOf(v)
	if valOf.Kind() == reflect.Ptr {
		valOf = valOf.Elem()
	}
	switch valOf.Kind() {
	case reflect.Slice:
		for i := 0; i < valOf.Len(); i++ {
			o.addRow(o.sheet.AddRow(), valOf.Index(i), -1)
		}
	default:
		o.addRow(o.sheet.AddRow(), valOf, -1)
	}
	return o
}

func (o *excel) To(c *gin.Context) {
	var buffer bytes.Buffer
	o.file.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, o.filename))
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeContent(c.Writer, c.Request, o.filename, time.Now(), content)
}

func (o *excel) ToFile(dir string) {
	o.file.Save(filepath.Join(dir, o.filename))
}
