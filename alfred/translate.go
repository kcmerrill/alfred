package alfred

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
)

func translate(raw string, context *Context) string {
	if raw == "" {
		// Nothing to translate, move along
		return raw
	}
	fmap := sprig.TxtFuncMap()
	te := template.Must(template.New("template").Funcs(fmap).Parse(raw))
	var b bytes.Buffer
	err := te.Execute(&b, context)
	if err != nil {
		context.Ok = false
		fmt.Print(translate("{{ .Text.Failure }}Missing arguments{{ .Text.Reset }}", context))
	}
	return b.String()
}
