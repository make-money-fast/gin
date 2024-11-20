// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"golang.org/x/net/context"
	"html/template"
	"net/http"
)

// Delims represents a set of Left and Right delimiters for HTML template rendering.
type Delims struct {
	// Left delimiter, defaults to {{.
	Left string
	// Right delimiter, defaults to }}.
	Right string
}

type Context struct {
	RequestContext context.Context
	Name           string
	DataMap        map[string]any
}

func (c *Context) ContextValue(key interface{}) any {
	return c.RequestContext.Value(key)
}

func (c *Context) Set(key interface{}, val interface{}) {
	if c.RequestContext == nil {
		c.RequestContext = context.Background()
	}
	c.RequestContext = context.WithValue(c.RequestContext, key, val)
}

// HTMLRender interface is to be implemented by HTMLProduction and HTMLDebug.
type HTMLRender interface {
	// Instance returns an HTML instance.
	Instance(*Context) Render
}

// HTMLProduction contains template reference and its delims.
type HTMLProduction struct {
	Template *template.Template
	Delims   Delims
}

// HTMLDebug contains template delims and pattern and function with file list.
type HTMLDebug struct {
	Files   []string
	Glob    string
	Delims  Delims
	FuncMap template.FuncMap
}

// HTML contains template reference and its name with given interface object.
type HTML struct {
	Template *template.Template
	Name     string
	Data     any
}

var htmlContentType = []string{"text/html; charset=utf-8"}

// Instance (HTMLProduction) returns an HTML instance which it realizes Render interface.
func (r HTMLProduction) Instance(ctx *Context) Render {
	return HTML{
		Template: r.Template,
		Name:     ctx.Name,
		Data:     ctx.DataMap,
	}
}

// Instance (HTMLDebug) returns an HTML instance which it realizes Render interface.
func (r HTMLDebug) Instance(ctx *Context) Render {
	return HTML{
		Template: r.loadTemplate(),
		Name:     ctx.Name,
		Data:     ctx.DataMap,
	}
}
func (r HTMLDebug) loadTemplate() *template.Template {
	if r.FuncMap == nil {
		r.FuncMap = template.FuncMap{}
	}
	if len(r.Files) > 0 {
		return template.Must(template.New("").Delims(r.Delims.Left, r.Delims.Right).Funcs(r.FuncMap).ParseFiles(r.Files...))
	}
	if r.Glob != "" {
		return template.Must(template.New("").Delims(r.Delims.Left, r.Delims.Right).Funcs(r.FuncMap).ParseGlob(r.Glob))
	}
	panic("the HTML debug render was created without files or glob pattern")
}

// Render (HTML) executes template and writes its result with custom ContentType for response.
func (r HTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

// WriteContentType (HTML) writes HTML ContentType.
func (r HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
