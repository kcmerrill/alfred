package alfred

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
)

func translateArgs(args []string, context *Context) []string {
	translated := make([]string, 0)
	for _, arg := range args {
		translated = append(translated, translate(arg, context))
	}
	return translated
}

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
		fmt.Println(translate("{{ .Text.Failure }}{{ .Text.FailureIcon }} Missing arguments{{ .Text.Reset }}", context))
		// TODO: chew on this some, should we fail? If not, how do we handle this gracefully?
		os.Exit(42)
	}
	return b.String()
}
