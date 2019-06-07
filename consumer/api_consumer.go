package ogiconsumer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	golenv "github.com/abhishekkr/gol/golenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/pseidemann/finish"

	logger "github.com/OpenChaos/ogi/logger"
	ogitransformer "github.com/OpenChaos/ogi/transformer"
)

var (
	ListenAt          = golenv.OverrideIfEnv("CONSUMER_API_LISTENAT", ":8080")
	ApisServed        = golenv.OverrideIfEnv("CONSUMER_API_PATHS_CSV", "^/$")
	BasicAuthEnabled  = golenv.OverrideIfEnv("CONSUMER_API_BASICAUTH_ENABLED", "false")
	BasicAuthUsername = golenv.OverrideIfEnv("CONSUMER_API_BASICAUTH_USERNAME", "")
	BasicAuthPassword = golenv.OverrideIfEnv("CONSUMER_API_BASICAUTH_PASSWORD", "")
)

func NewHttpApiConsumer() Consumer {
	var apiServer APIServer
	return &apiServer
}

type Handler func(*Context)

type Route struct {
	Pattern *regexp.Regexp
	Handler Handler
}

type APIServer struct {
	Routes       []Route
	DefaultRoute Handler
}

type APIRequest struct {
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (apiServer *APIServer) Consume() {
	logger.Infof("listening at: %s\n", ListenAt)

	api := NewAPIServer()

	http.Handle("^/metrics$", promhttp.Handler())

	svr := &http.Server{Addr: ListenAt, Handler: api}
	fin := finish.New()
	fin.Add(svr)
	go func() {
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	fin.Wait()
}

type Context struct {
	http.ResponseWriter
	*http.Request
	Params []string
}

func (c *Context) Json(code int, body string) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.WriteHeader(code)

	io.WriteString(c.ResponseWriter, body)
}

func NewAPIServer() *APIServer {
	api := &APIServer{
		DefaultRoute: func(ctx *Context) {
			ctx.Json(http.StatusNotFound, `{"error": "wrong api", "success": false}`)
		},
		Routes: []Route{},
	}

	api.Handle("^/ping$",
		func(ctx *Context) {
			ctx.Json(200, `{"status": "pong", "success": true}`)
		},
	)

	for _, apiPath := range strings.Split(ApisServed, ",") {
		api.Handle(strings.Trim(apiPath, " "), apiHandler)
	}

	return api
}

func (ctx *Context) IsAllowed() bool {
	username, password, ok := ctx.Request.BasicAuth()
	if authEnabled, err := strconv.ParseBool(BasicAuthEnabled); err == nil && authEnabled == false {
		return true
	}
	if ok == true && BasicAuthUsername == username && BasicAuthPassword == password {
		return true
	}
	return false
}

func apiHandler(ctx *Context) {
	if ctx.IsAllowed() == false {
		ctx.Json(http.StatusForbidden, `{"error": "auth-failure", "success": false}`)
		return
	}

	var apiRequest APIRequest
	requestBytes, err := apiRequest.Marshal(ctx.Request)
	if err != nil {
		logger.Errorf("json consume error: %s", ctx.Request.URL.Path)
		ctx.Json(http.StatusOK, `{"status": "consume-error", "success": false}`)
		return
	}
	ogitransformer.Transform(requestBytes)
	logger.Infof("consumed: %s", ctx.Request.URL.RawPath)
	ctx.Json(http.StatusOK, `{"status": "consumed", "success": true}`)
}

func (a *APIServer) Handle(pattern string, handler Handler) {
	re := regexp.MustCompile(pattern)
	route := Route{Pattern: re, Handler: handler}

	a.Routes = append(a.Routes, route)
}

func (a *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{Request: r, ResponseWriter: w}

	for _, route := range a.Routes {
		if matches := route.Pattern.FindStringSubmatch(ctx.URL.Path); len(matches) > 0 {
			if len(matches) > 1 {
				ctx.Params = matches[1:]
			}

			route.Handler(ctx)
			return
		}
	}

	a.DefaultRoute(ctx)
}

func (apiRequest *APIRequest) Marshal(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	if err != nil {
		return []byte{}, err
	}

	apiRequest.Path = req.URL.Path
	apiRequest.Method = req.Method
	apiRequest.Headers = map[string]string{}
	apiRequest.Body = string(body)

	for header, val := range req.Header {
		apiRequest.Headers[header] = strings.Join(val, " ")
	}
	return json.Marshal(apiRequest)
}

func (apiRequest *APIRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &apiRequest)
}
