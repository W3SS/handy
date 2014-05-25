package lib

import "github.com/go4r/handy"

type r_defer chan func()

var (
	_ = handy.Server.Context.SetProvider(
		"defer",
		func(c *handy.Context) func() interface{} {

			var rdefer = make(r_defer)

			go func() {

				var defereds []func()

				for deferedCall := range rdefer {
					defereds = append(defereds, deferedCall)
				}

				for _, v := range defereds {
					v()
				}

			}()

			c.CleanupFunc(func() {
				close(rdefer)
			})

			return func() interface{} {
				return rdefer
			}
		})
)

func Defer(r interface{}, call func()) {
	handy.CContext(r).Get("defer").(r_defer)<-call
}
