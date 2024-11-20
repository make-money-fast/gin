package binding

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

type varBinding struct {
}

func (v varBinding) Name() string {
	return "var"
}

func (v varBinding) Bind(req *http.Request, obj any) error {
	vars := mux.Vars(req)
	varForms := toFormValues(vars)
	if err := mapForm(obj, varForms); err != nil {
		return err
	}
	return validate(obj)
}

type varFormBinding struct {
}

func (varFormBinding) Name() string {
	return "varForm"
}

func (varFormBinding) Bind(req *http.Request, obj any) error {
	vars := mux.Vars(req)
	varForms := toFormValues(vars)

	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}
	values := mergeUrlValues(req.Form, varForms)
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

type varQueryBinding struct{}

func (varQueryBinding) Name() string {
	return "varQuery"
}

func (varQueryBinding) Bind(req *http.Request, obj any) error {
	vars := mux.Vars(req)
	varForms := toFormValues(vars)
	values := mergeUrlValues(varForms, req.URL.Query())
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

func toFormValues(vars map[string]string) url.Values {
	var val = url.Values{}
	for k, v := range vars {
		val.Set(k, v)
	}
	return val
}

func mergeUrlValues(values ...url.Values) url.Values {
	var base = make(url.Values)
	for _, val := range values {
		for k := range val {
			base.Set(k, val.Get(k))
		}
	}
	return base
}
