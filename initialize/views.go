package initialize

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/LINQQ1212/common/global"
	"io"
	"path"
	"reflect"
)

func NewViews() {
	global.View = jet.NewSet(
		jet.NewOSFileSystemLoader(path.Join(global.CONFIG.System.MainDir, "views")),
		jet.InDevelopmentMode(), // remove in production
		jet.WithSafeWriter(func(w io.Writer, b []byte) {
			w.Write(b)
		}),
	)
	global.View.AddGlobalFunc("ListDesc", func(arg jet.Arguments) reflect.Value {
		desc := ""
		num := 100
		if arg.IsSet(0) {
			desc = arg.Get(0).String()
		}
		if arg.IsSet(1) {
			num = int(arg.Get(1).Float())
		}
		if len(desc) > num {
			return reflect.ValueOf(desc[0:num] + "...")
		}
		return reflect.ValueOf(desc)
	})
}
