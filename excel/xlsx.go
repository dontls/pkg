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

type excel struct {
	filename string
	file     *xlsx.File
	sheet    *xlsx.Sheet
}

// 自定义title接口
type ITitles interface {
	SheetTitles() []string
}

func Sheet(sheet string, v any) *excel {
	r := &excel{}
	r.filename = fmt.Sprintf("%s_%s.xlsx", sheet, time.Now().Format("2006-01-02 150405"))
	r.file = xlsx.NewFile()
	r.sheet, _ = r.file.AddSheet(sheet)
	titleRow := r.sheet.AddRow()
	value, titles := scanTitles(v)
	for _, v := range titles {
		titleRow.AddCell().Value = v
	}
	r.setData(value)
	return r
}

func structTitles(typOf reflect.Type) []string {
	var titles []string
	for i := 0; i < typOf.NumField(); i++ {
		f := typOf.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			titles = append(titles, structTitles(f.Type)...)
		} else {
			titles = append(titles, f.Name)
		}
	}
	return titles
}

func scanTitles(v any) (reflect.Value, []string) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	vtype := value.Type()
	if value.Kind() == reflect.Slice {
		vtype = value.Type().Elem()
	}
	if t, ok := reflect.New(vtype).Interface().(ITitles); ok {
		return value, t.SheetTitles()
	}
	return value, structTitles(vtype)
}

func (o *excel) setData(value reflect.Value) *excel {
	switch value.Kind() {
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			o.addRow(o.sheet.AddRow(), value.Index(i), -1)
		}
	default:
		o.addRow(o.sheet.AddRow(), value, -1)
	}
	return o
}

func (o *excel) Write(c *gin.Context) error {
	var buffer bytes.Buffer
	o.file.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, o.filename))
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeContent(c.Writer, c.Request, o.filename, time.Now(), content)
	return nil
}

func (o *excel) WriteFile(dir string) error {
	return o.file.Save(filepath.Join(dir, o.filename))
}
