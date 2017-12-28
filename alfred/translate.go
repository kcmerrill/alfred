package alfred

import (
	"bytes"
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
	te := template.Must(template.New("template." + context.TaskName).Funcs(fmap).Parse(raw))
	var b bytes.Buffer
	context.Lock.Lock()
	err := te.Execute(&b, context)
	context.Lock.Unlock()
	if err != nil {
		context.Ok = false
		outFail("template", "Invalid Argument(s)", context)
		// TODO: chew on this some, should we fail? If not, how do we handle this gracefully?
		os.Exit(42)
	}
	return b.String()
}
