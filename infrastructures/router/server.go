package router

import (
	stdContext "context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	// For Swagger
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Router echo.Echo
type Router struct {
	*echo.Echo
	Debug bool
	Stats *Stats
}

// Stats .. for server stats
type Stats struct {
	RequestCount map[string]uint64     `json:"requestCount"`
	Statuses     map[string]StatStatus `json:"statuses"`
	mutex        sync.RWMutex
}

// StatStatus ...
type StatStatus struct {
	Status map[string]uint64 `json:"status"`
}

// Init Echo Framework Initialize
func Init(log *logrus.Logger, debug bool) (r *Router, err error) {
	r = &Router{Echo: echo.New(), Debug: debug}

	// Recover returns a middleware which recovers from panics anywhere in the chain
	// and handles the control to the centralized HTTPErrorHandler.
	r.Use(middleware.Recover())

	// CORS default
	// 모든 원격지에서 오는 모든 메서드를 허용합니다.
	r.Use(middleware.CORS())

	// CORS restricted
	// `https://jblim.mobigen.com`과 `https://mobigen.com`로부터 오는 요청 중
	// GET, PUT, POST or DELETE 메서드를 허용합니다.
	//r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"https://jblim.mobigen.com", "https://mobigen.com"},
	//	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	//}))

	// ${id}: HeaderXRequestID
	// ${remote_ip} : RealIP
	// ${host} : Host
	// ${uri} : RequestURI
	// ${method} : Method
	// ${path} : Path
	// ${protocol} : Proto
	// ${referer} : req.Referer()
	// ${user_agent} : req.UserAgent()
	// ${status} : response status
	// ${error} : golang error string
	// ${latency} :
	// ${latency_human} :
	// ${bytes_out} : response size
	// ${header} : ..
	// ${query} : ..
	// ${form} : ..
	// ${cookie} : ..

	// Customize Log Format Sample
	logConfig := middleware.LoggerConfig{
		Skipper: r.LoggerSkipper,
		Format: "${time_custom} [DEBU] [echo-framework   :  - ] [ Router ] " +
			"${method} ${uri} ${status} Laency[ ${latency_human} ]\n",
		CustomTimeFormat: "2006-01-02 15:04:05.000",
		Output:           log.Out,
	}
	r.Use(middleware.LoggerWithConfig(logConfig))

	// Stats
	stats := new(Stats)
	stats.RequestCount = make(map[string]uint64)
	stats.Statuses = make(map[string]StatStatus)
	r.Use(stats.Process)
	r.GET("/stats", stats.StatsHandle)
	r.Stats = stats

	// Swager
	r.GET("/swagger/*", echoSwagger.WrapHandler)

	r.HideBanner = true
	r.HidePort = true

	return r, nil
}

// EnableDebug debug mode on
func (r *Router) EnableDebug() {
	r.Debug = true
}

// DisableDebug disable debug
func (r *Router) DisableDebug() {
	r.Debug = false
}

// LoggerSkipper .. logger skipper
func (r *Router) LoggerSkipper(e echo.Context) bool {
	if r.Debug {
		// 설정에 따라 skip 되지 못하도록 하거나,
		return false
	}
	// skip 되도록 한다.
	return true
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		// Get Request URI
		uri := c.Request().RequestURI
		// Get Response Status
		status := strconv.Itoa(c.Response().Status)

		// Request Stat
		if _, ok := s.RequestCount[uri]; ok {
			s.RequestCount[uri]++
		} else {
			s.RequestCount[uri] = 1
		}
		// Response Stat
		if val, ok := s.Statuses[uri]; ok {
			if _, ok := val.Status[status]; ok {
				s.Statuses[uri].Status[status]++
			} else {
				s.Statuses[uri].Status[status] = 1
			}
		} else {
			ss := make(map[string]uint64)
			ss[status] = 1
			s.Statuses[uri] = StatStatus{
				Status: ss,
			}
		}
		return nil
	}
}

// StatsHandle send stats
func (s *Stats) StatsHandle(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}

// Run echo framework
func (r *Router) Run(listenAddr string) error {
	if r == nil {
		return fmt.Errorf("ERROR. Router Not Initialize")
	}
	return r.Start(listenAddr)
}

// Stop echo framework
func (r *Router) Stop(ctx stdContext.Context) error {
	return r.Shutdown(ctx)
}
