package gin

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"html/template"
	"io"
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
	functions  map[string]jet.Func
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

func Functions(functions map[string]jet.Func) Options {
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

	opts = append(opts, htmlEscaper)

	views := jet.NewSet(
		jet.NewOSFileSystemLoader(directory),
		opts...,
	)

	if len(jetOption.functions) != 0 {
		for name, fn := range jetOption.functions {
			views.AddGlobalFunc(name, fn)
		}
	}

	return &JetInstantRender{
		views: views,
	}
}

func htmlEscaper(w io.Writer, b []byte) {
	template.HTMLEscape(w, b)
}

func (j *JetInstantRender) Instance(ctx *RenderContext) Render {
	return &JetHtmlRender{
		views: j.views,
		name:  ctx.Name,
		ctx:   ctx.GinContext,
	}
}

type JetHtmlRender struct {
	views *jet.Set
	name  string
	ctx   *Context
}

func (j *JetHtmlRender) Render(writer http.ResponseWriter) error {
	t, err := j.views.GetTemplate(j.name)
	if err != nil {
		return fmt.Errorf("parse template failed: %+w", err)
	}
	data := j.ctx.Assigned()
	variables := make(jet.VarMap)
	for key, item := range data {
		variables.Set(key, item)
	}

	for _, functionBuilder := range contextBuilders {
		name, fn := functionBuilder(j.ctx)
		variables.Set(name, fn)
	}

	if err := t.Execute(writer, variables, nil); err != nil {
		return err
	}
	return nil
}

func (j *JetHtmlRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}
