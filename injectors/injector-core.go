package injectors

import (
	"github.com/jblim0125/redredis-expire/common"
	"github.com/jblim0125/redredis-expire/common/appdata"
	"github.com/jblim0125/redredis-expire/infrastructures/datastore"
	"github.com/jblim0125/redredis-expire/infrastructures/router"
)

// Injector web-server layer initializer : Dependency Injection )
type Injector struct {
	Router    *router.Router
	Datastore *datastore.DataStore
	Log       *common.Logger
	Conf      *appdata.Configuration
}

// New create Injector
func (Injector) New(r *router.Router, d *datastore.DataStore,
	l *common.Logger, c *appdata.Configuration) *Injector {
	return &Injector{
		Router:    r,
		Datastore: d,
		Log:       l,
		Conf:      c,
	}
}

// Init web-server layer interconnection create (web server layer)
func (injector *Injector) Init() error {
	// For Version
	ver := Version{}.Init(injector)
	injector.Router.GET("/version", ver.GetVersion)

	//// 캐시 서버가 처리하지 않는 요청을 Angora로 전달하기 위함.
	//angoraURL := fmt.Sprintf("http://%s:%d",
	//	injector.Conf.DestInfo.Angora.IP,
	//	injector.Conf.DestInfo.Angora.Port)
	//injector.RegisterProxy(angoraURL)
	//
	//// For Query(DSL, SID)
	//qry := QueryInjector{}.New(injector.Log)
	//injector.Router.GET("/angora/v2/query/jobs/:sid", qry.GetEventData)
	//injector.Router.POST("/angora/v2/query/jobs", qry.GetEventSid)
	//
	//// For Setting(Admin)
	//settingGrp := injector.Router.Group("/cache-server/api/v1/setting")
	//setting := Setting{}.Init(injector)
	//settingGrp.GET("/loglevel", setting.GetLogLevel)
	//settingGrp.POST("/loglevel", setting.UpdateLogLevel)
	//settingGrp.GET("/cache", setting.GetCacheSetting)
	//settingGrp.POST("/cache", setting.UpdateCacheSetting)
	//settingGrp.GET("/jobqueue", setting.GetJobQueueSetting)
	//settingGrp.POST("/jobqueue", setting.UpdateJobQueueSetting)
	//settingGrp.GET("/stat", setting.GetStatSetting)
	//settingGrp.POST("/stat", setting.UpdateStatSetting)

	// Sample
	//injector.Log.Errorf("[ PATH ] /api/v1/sample ........................................................... [ OK ]")
	//apiv1.GET("/sample/:id", sample.GetByID)
	//apiv1.POST("/sample", sample.Create)
	//apiv1.POST("/sample/update", sample.Update)
	//apiv1.DELETE("/sample/:id", sample.Delete)
	return nil
}

// RegisterProxy 캐시 서버는 Angora의 특정 API만을 처리 하므로 ReverseProxy로 Angora를 등록한다.
//func (injector *Injector) RegisterProxy(angoraURL string) error {
//	url1, err := url.Parse(angoraURL)
//	if err != nil {
//		injector.Log.Error(err)
//		return err
//	}
//	targets := []*middleware.ProxyTarget{
//		{
//			URL: url1,
//		},
//	}
//	// Proxy
//	proxyConf := middleware.DefaultProxyConfig
//	// 특정 API를 Cache서버가 처리하기 위해 Proxy처리 하지 않을 Method - Path를 처리
//	proxyConf.Skipper = injector.Skipper
//	proxyConf.Balancer = middleware.NewRoundRobinBalancer(targets)
//	proxy := middleware.ProxyWithConfig(proxyConf)
//
//	// proxy 등록
//	angora := injector.Router.Group("/angora")
//	angora.Use(proxy)
//	return nil
//}
//
//// Skipper 용
//const (
//	angoraPath   = "/angora/query/jobs"
//	angoraV2Path = "/angora/v2/query/jobs"
//)
//
//// Skipper 미들웨어로 처리하지 않을 Path와 Method 선별한다.
//func (injector *Injector) Skipper(echo echo.Context) bool {
//	method := echo.Request().Method
//	path := echo.Request().URL.Path
//	injector.Log.Debugf("Method[ %s ] Path[ %s ]", method, path)
//
//	if strings.Index(path, angoraPath) == 0 {
//		// DSL Query
//		if path == angoraPath && method == http.MethodPost {
//			return true
//		}
//		// SID
//		if strings.Count(path, "/") == 4 && method == http.MethodGet {
//			return true
//		}
//	}
//
//	if strings.Index(path, angoraV2Path) == 0 {
//		// DSL Query
//		if path == angoraV2Path && method == http.MethodPost {
//			return true
//		}
//		// SID
//		if strings.Count(path, "/") == 5 && method == http.MethodGet {
//			return true
//		}
//	}
//	return false
//}
