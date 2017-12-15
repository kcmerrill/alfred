package alfred

import (
	event "github.com/kcmerrill/hook"
)

func init() {
	event.Register("output", output)
	event.Register("setup", setup)
	event.Register("summary", summary)
	event.Register("watch", watch)
	event.Register("serve", serve)
	event.Register("result", result)
	event.Register("ok", ok)
	event.Register("fail", fail)
	event.Register("wait", wait)
}
