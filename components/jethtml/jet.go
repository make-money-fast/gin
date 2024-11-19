package jethtml

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/clearcodecn/gin/render"
	"net/http"
)

type JetInstantRender struct {
	views *jet.Set
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

func (j *JetInstantRender) Instance(s string, a any) render.Render {
	return &JetHtmlRender{
		views: j.views,
		name:  s,
		data:  a,
	}
}

type JetHtmlRender struct {
	views *jet.Set
	name  string
	data  any
}

func (j *JetHtmlRender) Render(writer http.ResponseWriter) error {
	t, err := j.views.GetTemplate(j.name)
	if err != nil {
		return err
	}
	if err := t.Execute(writer, nil, j.data); err != nil {
		return err
	}
	return nil
}

func (j *JetHtmlRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}
