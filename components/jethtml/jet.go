package jethtml

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/make-money-fast/gin"
	"github.com/make-money-fast/gin/render"
	"net/http"
)

type JetInstantRender struct {
	views *jet.Set
	funcs map[string]jet.Func
}

type JetInstantRenderOption struct {
	extensions []string
	leftDelim  string
	rightDelim string
	functions  []jet.Func
	debug      bool
}

type Options func(*JetInstantRenderOption)

func Extensions(s ...string) Options {
	return func(r *JetInstantRenderOption) {
		r.extensions = s
	}
}

func Debug() Options {
	return func(option *JetInstantRenderOption) {
		option.debug = true
	}
}

func Functions(functions ...jet.Func) Options {
	return func(o *JetInstantRenderOption) {
		o.functions = functions
	}
}

func Delims(left string, right string) Options {
	return func(option *JetInstantRenderOption) {
		option.leftDelim = left
		option.rightDelim = right
	}
}

func NewJetRender(directory string, options ...Options) *JetInstantRender {
	var jetOption JetInstantRenderOption
	for _, opt := range options {
		opt(&jetOption)
	}

	var opts []jet.Option
	if jetOption.debug {
		opts = append(opts, jet.InDevelopmentMode())
	}
	if jetOption.leftDelim != "" || jetOption.rightDelim != "" {
		opts = append(opts, jet.WithDelims(jetOption.leftDelim, jetOption.rightDelim))
	}

	if len(jetOption.extensions) != 0 {
		opts = append(opts, jet.WithTemplateNameExtensions(jetOption.extensions))
	}

	views := jet.NewSet(
		jet.NewOSFileSystemLoader(directory),
		opts...,
	)

	return &JetInstantRender{
		views: views,
	}
}

var (
	contextKey = struct{}{}
)

func NewContext(name string, ctx *gin.Context) *render.Context {
	c := &render.Context{
		Name: name,
	}
	c.Set(contextKey, ctx)
	return c
}

func (j *JetInstantRender) Instance(ctx *render.Context) render.Render {
	return &JetHtmlRender{
		views: j.views,
		name:  ctx.Name,
		ctx:   ctx.ContextValue(contextKey).(*gin.Context),
	}
}

type JetHtmlRender struct {
	views *jet.Set
	name  string
	ctx   *gin.Context
}

func (j *JetHtmlRender) Render(writer http.ResponseWriter) error {
	t, err := j.views.GetTemplate(j.name)
	if err != nil {
		return err
	}
	data := j.ctx.Assigned()
	variables := make(jet.VarMap)
	for key, item := range data {
		variables.Set(key, item)
	}

	if err := t.Execute(writer, variables, nil); err != nil {
		return err
	}
	return nil
}

func (j *JetHtmlRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}
