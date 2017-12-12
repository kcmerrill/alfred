package alfred

import (
	event "github.com/kcmerrill/hook"
)

func init() {
	event.Register("speak", speak)
	event.Register("setup", setup)
	event.Register("summary", summary)
	event.Register("serve", serve)
	event.Register("result", result)
	event.Register("ok", ok)
	event.Register("fail", fail)
	event.Register("wait", wait)
}
