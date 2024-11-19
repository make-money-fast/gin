package gin

import (
	"net/http"
	"strings"
)

type muxHandler func(path string, h HandlersChain)

func (e *Engine) getMuxHandlers() map[string]muxHandler {
	return map[string]muxHandler{
		http.MethodGet:     e.muxGet,
		http.MethodPost:    e.muxPost,
		http.MethodPut:     e.muxPut,
		http.MethodDelete:  e.muxDelete,
		http.MethodPatch:   e.muxPatch,
		http.MethodOptions: e.muxOptions,
		http.MethodHead:    e.muxHead,
		http.MethodConnect: e.muxConnect,
		http.MethodTrace:   e.muxTrace,
	}
}

func (e *Engine) addMuxRouter(method string, path string, h HandlersChain) {
	if strings.HasSuffix(path, "/*") {
		path = strings.TrimSuffix(path, "*")
		e.muxPrefix(method, path, h)
		return
	}

	mh, ok := e.getMuxHandlers()[method]
	if ok {
		mh(path, h)
		return
	}

	panic("unsupport method: " + method + ", path: " + path)
}

func (e *Engine) muxPrefix(method string, path string, h HandlersChain) {
	e.muxRouter.Methods(method).PathPrefix(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxGet(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodGet).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxPost(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodPost).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxPut(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodPut).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxDelete(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodDelete).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxPatch(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodPatch).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxOptions(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodOptions).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxHead(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodHead).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxConnect(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodConnect).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) muxTrace(path string, h HandlersChain) {
	e.muxRouter.Methods(http.MethodTrace).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.handlerRequest(w, r, h)
	})
}

func (e *Engine) handlerRequest(w http.ResponseWriter, r *http.Request, h HandlersChain) {
	c := e.pool.Get().(*Context)
	c.reset()

	c.writermem.reset(w)
	c.Request = r
	c.handlers = h
	c.Next()
	e.pool.Put(c)
}
