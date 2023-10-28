package http_router

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type MiddleWare func(w http.ResponseWriter, r *http.Request, next func(), abort func(code int))

type InvalidPathParameterErr struct {
	Message string
}

func (r *InvalidPathParameterErr) Error() string {
	return r.Message
}

func HandleError(err error, writer http.ResponseWriter, statusCode int) (proceed bool) {
	if err == nil {
		return true
	}
	http.Error(writer, err.Error(), statusCode)
	return false
}

func Vars(request *http.Request) url.Values {
	query := request.URL.Query()
	return query
}

func Params(request *http.Request) map[string]any {
	params := request.Context().Value("params").(map[string]any)
	return params
}

type Router interface {
	Get(relativePath string, handler http.HandlerFunc)
	Post(relativePath string, handler http.HandlerFunc)
	Put(relativePath string, handler http.HandlerFunc)
	Delete(relativePath string, handler http.HandlerFunc)
	Patch(relativePath string, handler http.HandlerFunc)
	Use(handler MiddleWare)
	http.Handler
}

type GinRouter struct {
	Engine *gin.Engine
}

func New() Router {
	router := GinRouter{Engine: gin.Default()}
	return router
}

func overwriteContext(request *http.Request, params gin.Params) *http.Request {
	mapParams := make(map[string]any)

	for _, param := range params {
		mapParams[param.Key] = param.Value
	}

	ctx := context.WithValue(request.Context(), "params", mapParams)
	return request.WithContext(ctx)
}

func (r GinRouter) Get(relativePath string, handler http.HandlerFunc) {
	r.Engine.GET(relativePath, func(c *gin.Context) {
		handler(c.Writer, overwriteContext(c.Request, c.Params))
	})
}

func (r GinRouter) Post(relativePath string, handler http.HandlerFunc) {
	r.Engine.POST(relativePath, func(c *gin.Context) {
		handler(c.Writer, overwriteContext(c.Request, c.Params))
	})
}

func (r GinRouter) Put(relativePath string, handler http.HandlerFunc) {
	r.Engine.PUT(relativePath, func(c *gin.Context) {
		handler(c.Writer, overwriteContext(c.Request, c.Params))
	})
}
func (r GinRouter) Patch(relativePath string, handler http.HandlerFunc) {
	r.Engine.PATCH(relativePath, func(c *gin.Context) {
		handler(c.Writer, overwriteContext(c.Request, c.Params))
	})
}

func (r GinRouter) Delete(relativePath string, handler http.HandlerFunc) {
	r.Engine.DELETE(relativePath, func(c *gin.Context) {
		handler(c.Writer, overwriteContext(c.Request, c.Params))
	})
}

func (r GinRouter) Use(handler MiddleWare) {
	r.Engine.Use(func(c *gin.Context) {
		handler(c.Writer, c.Request, c.Next, c.AbortWithStatus)
	})
}

func (r GinRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.Engine.ServeHTTP(writer, request)
}
