package gin

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/golang-module/carbon/v2"
	"reflect"
	"time"
)

var (
	globalFunctions = make(map[string]jet.Func)
	contextBuilders = make([]FunctionsBuilder, 0)
)

func init() {
	fn := new(TemplateFunc)
	carbon.SetDefault(carbon.Default{
		Layout:   carbon.DateTimeFormat,
		Timezone: carbon.Local,
		Locale:   "zh-CN",
	})
	RegisterContextFunctionBuilders(
		fn.Route(),
		fn.CurrentURI(),
		fn.AbsUrl(),
		fn.HasQuery(),
		fn.Carbon(),
		fn.Date(),
		fn.DateTime(),
	)
}

func RegisterGlobalFunctions(name string, p jet.Func) {
	globalFunctions[name] = p
}

func RegisterContextFunctionBuilders(b ...FunctionsBuilder) {
	contextBuilders = append(contextBuilders, b...)
}

type FunctionsBuilder func(c *Context) (string, jet.Func)

func of(a any) reflect.Value {
	return reflect.ValueOf(a)
}

type TemplateFunc struct {
}

// === url 辅助函数

// Route 用法: {{ route "/index" }} , {{ route "/index" "a=1&b=2" }}
func (c *TemplateFunc) Route() FunctionsBuilder {
	return func(c *Context) (string, jet.Func) {
		name := "route"
		return name, func(a jet.Arguments) reflect.Value {
			a.RequireNumOfArguments(name, 1, 2)
			path := a.Get(0).String()

			var query string
			if a.IsSet(1) {
				query = a.Get(1).String()
			}
			if len(query) == 0 {
				return of(fmt.Sprintf("%s", path))
			}
			return of(fmt.Sprintf("%s?%s", path, query))
		}
	}
}

// CurrentURI 获取当前路劲
func (c *TemplateFunc) CurrentURI() FunctionsBuilder {
	name := "currentURI"
	return func(c *Context) (string, jet.Func) {
		return name, func(a jet.Arguments) reflect.Value {
			return reflect.ValueOf(c.Request.RequestURI)
		}
	}
}

// AbsUrl 获取相对url
func (c *TemplateFunc) AbsUrl() FunctionsBuilder {
	name := "absUrl"
	return func(c *Context) (string, jet.Func) {
		return name, func(a jet.Arguments) reflect.Value {
			a.RequireNumOfArguments(name, 1, 1)
			newUrl, _ := c.Request.URL.Parse(a.Get(0).String())
			if newUrl != nil {
				return of(newUrl.String())
			}
			return of("")
		}
	}
}

// HasQuery 获取相对url
func (c *TemplateFunc) HasQuery() FunctionsBuilder {
	name := "hasQuery"
	return func(c *Context) (string, jet.Func) {
		return name, func(a jet.Arguments) reflect.Value {
			a.RequireNumOfArguments(name, 1, 1)
			key := a.Get(0).String()
			val := c.Query(key)
			if val == "" {
				return of(false)
			}
			return of(true)
		}
	}
}

// Query 获取 query 参数.
func (c *TemplateFunc) Query() FunctionsBuilder {
	name := "query"
	return func(c *Context) (string, jet.Func) {
		return name, func(a jet.Arguments) reflect.Value {
			a.RequireNumOfArguments(name, 1, 1)
			key := a.Get(0).String()
			val := c.Query(key)
			return of(val)
		}
	}
}

// == 时间函数
func (c *TemplateFunc) carbonBase(a jet.Arguments) time.Time {
	var (
		ti time.Time
	)
	if a.IsSet(0) {
		arg := a.Get(0)
		numVal, ok := numberValue(arg)
		if ok {
			ti = time.Unix(numVal, 0)
		} else {
			if arg.Kind() == reflect.Struct { // time.time
				ti = arg.Interface().(time.Time)
			}
			if arg.Kind() == reflect.String { // string .
				var layout string
				if a.IsSet(1) {
					// layout
					layout = a.Get(1).String()
				} else {
					layout = time.RFC3339
				}
				t, err := time.Parse(layout, arg.Interface().(string))
				if err != nil {
					panic("failed to parse time: " + err.Error())
				}
				ti = t
			}
		}
	}
	return ti
}

func (f *TemplateFunc) Carbon() FunctionsBuilder {
	name := "carbon"
	return func(c *Context) (string, jet.Func) {
		return name, func(a jet.Arguments) reflect.Value {
			return of(carbon.CreateFromStdTime(f.carbonBase(a)).DiffForHumans(carbon.Now()))
		}
	}
}

func (f *TemplateFunc) Date() FunctionsBuilder {
	name := "date"
	return func(c *Context) (string, jet.Func) {
		return name, func(a jet.Arguments) reflect.Value {
			return of(carbon.CreateFromStdTime(f.carbonBase(a)).ToDateString())
		}
	}
}

func (f *TemplateFunc) DateTime() FunctionsBuilder {
	name := "datetime"
	return func(c *Context) (string, jet.Func) {
		return name, func(a jet.Arguments) reflect.Value {
			return of(carbon.CreateFromStdTime(f.carbonBase(a)).ToDateTimeString())
		}
	}
}

func numberValue(value reflect.Value) (int64, bool) {
	switch value.Kind() {
	case reflect.Int:
		v := value.Interface().(int)
		return int64(v), true
	case reflect.Int8:
		v := value.Interface().(int8)
		return int64(v), true
	case reflect.Int16:
		v := value.Interface().(int16)
		return int64(v), true
	case reflect.Int32:
		v := value.Interface().(int32)
		return int64(v), true
	case reflect.Int64:
		v := value.Interface().(int64)
		return int64(v), true
	case reflect.Uint:
		v := value.Interface().(uint)
		return int64(v), true
	case reflect.Uint8:
		v := value.Interface().(uint8)
		return int64(v), true
	case reflect.Uint16:
		v := value.Interface().(uint16)
		return int64(v), true
	case reflect.Uint32:
		v := value.Interface().(uint32)
		return int64(v), true
	case reflect.Uint64:
		v := value.Interface().(uint64)
		return int64(v), true
	}
	return 0, false
}
