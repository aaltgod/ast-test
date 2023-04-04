package main

import (
	"fmt"
	"html/template"
	"strings"
)

type builder struct {
	strings.Builder
}

type builderHeader struct {
	PackageName string
	Imports     []string
	Types       []string
}

func newBuilder() *builder {
	return &builder{
		strings.Builder{},
	}
}

func (b *builder) Render() string {
	return b.String()
}
func (b *builder) RenderHeader(builderHeader builderHeader) error {
	importsTemplate := template.Must(template.New("header").Parse(`
package {{.PackageName}}_test

import (
{{range $_, $val := .Imports}}
	{{ $val }}	
{{end}}
)

type (
{{range $_, $val := .Types}}
	{{ $val }}	
{{end}}
)
`))

	if err := importsTemplate.Execute(b, builderHeader); err != nil {
		return err
	}

	return nil
}

func (b *builder) BuildTestCase() {

}

type frame struct {
	strings.Builder
}

func newFrame() *frame {
	return &frame{
		strings.Builder{},
	}
}

func (f *frame) BuildTest(funcName string) {
	f.buildTestHead(funcName)
	f.buildSuccessCase()
	f.buildTestFooter()
}

func (f *frame) View() string {
	return f.String()
}

func (f *frame) buildTestHead(funcName string) {
	f.WriteString(fmt.Sprintf(`
func Test_%s(t *testing.T) {
	t.Parallel()

`, funcName,
	))
}

func (f *frame) buildTestFooter() {
	f.WriteString("}")
}

func (f *frame) buildSuccessCase() {
	f.buildSuccessCaseHead()
	f.buildCaseFooter()
}

func (f *frame) buildSuccessCaseHead() {
	f.WriteString(`
	t.Run("success", func (t *testing.T) {
		t.Parallel()
		
		// arrange
		f := setUp(t)
		defer f.TearDown(t)

`)
}

func (f *frame) buildCaseFooter() {
	f.WriteString(`
	})`)
}
