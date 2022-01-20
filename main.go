package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jblim0125/redredis-expire/common"
	"github.com/jblim0125/redredis-expire/common/appdata"
	"github.com/jblim0125/redredis-expire/infrastructures/datastore"
	"github.com/jblim0125/redredis-expire/infrastructures/router"
	"github.com/jblim0125/redredis-expire/injectors"
	"github.com/jblim0125/redredis-expire/internal/redis"
	"github.com/jblim0125/redredis-expire/models"
)

// Context context of main
type Context struct {
	Env       *appdata.Environment
	Conf      *appdata.Configuration
	Log       *common.Logger
	CM        *common.ConfigManager
	Datastore *datastore.DataStore
	Router    *router.Router
	Redis     *redis.Manager
}

// InitLog Initialize logger
func (c *Context) InitLog() {
	log := common.Logger{}.GetInstance()
	log.SetLogLevel(logrus.DebugLevel)
	c.Log = log
	c.Log.Start()
}

// ReadEnv Read value of the environment
func (c *Context) ReadEnv() error {
	c.Env = new(appdata.Environment)
	// Get Home
	homePath := os.Getenv("APP_ROOT")
	if len(homePath) > 0 {
		c.Log.Errorf("APP_ROOT: %s", homePath)
		c.Env.Home = homePath
	} else {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}
		c.Env.Home = dir
	}
	// Get Profile
	profile := os.Getenv("PROFILE")
	if len(profile) > 0 {
		c.Log.Errorf("ENV PROFILE : %s", profile)
		c.Env.Profile = profile
	}
	// Get Log Level
	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) > 0 {
		c.Log.Errorf("ENV LOG_LEVEL : %s", logLevel)
		_, err := appdata.CheckLogLevel(logLevel)
		if err != nil {
			return err
		}
		c.Env.LogLevel = logLevel
	}
	c.Log.Errorf("[ Env ] Read ...................................................................... [ OK ]")
	return nil
}

// ReadConfig Read Configuration File By Viper
func (c *Context) ReadConfig() error {
	c.CM = common.ConfigManager{}.New(c.Log.Logger)
	// Write Config File Info
	configPath := c.Env.Home + "/configs"
	var configName string
	if len(c.Env.Profile) <= 0 {
		configName = "prod"
	} else {
		configName = c.Env.Profile
	}
	configType := "yaml"
	// Config file struct
	conf := new(appdata.Configuration)
	// Read
	if err := c.CM.ReadConfig(configPath, configName, configType, conf); err != nil {
		return err
	}
	// Save
	c.Conf = conf

	// Set Watcher
	c.CM.SetOnChanged(configPath, configName, configType,
		func(conf interface{}) {
			newConf := conf.(*appdata.Configuration)
			c.Log.Infof("%+v\n", newConf)
		}, conf)
	c.Log.Errorf("[ Configuration ] Read ............................................................ [ OK ]")
	return nil
}

// SetLogger set log level, log output. and etc
func (c *Context) SetLogger() error {
	if len(c.Env.LogLevel) > 0 {
		c.Conf.Log.Level = c.Env.LogLevel
	}
	if len(c.Conf.Log.Level) <= 0 {
		c.Conf.Log.Level = "debug"
	}
	if err := c.Log.Setting(&c.Conf.Log); err != nil {
		return err
	}
	return nil
}

// InitDatastore Initialize datastore
func (c *Context) InitDatastore() error {
	// Create datastore
	ds, err := datastore.DataStore{}.New(c.Log.Logger)
	if err != nil {
		return err
	}
	// Connect
	if err := ds.Connect(&c.Conf.Datastore); err != nil {
		return err
	}

	// Migrate
	if err := ds.Migrate(); err != nil {
		return err
	}

	c.Datastore = ds
	c.Log.Errorf("[ DataStore ] Initialize .......................................................... [ OK ]")
	return nil
}

// InitRouter Initialize router
func (c *Context) InitRouter() error {
	// init echo framework
	r, err := router.Init(c.Log.Logger, c.Conf.Server.Debug)
	if err != nil {
		return err
	}
	c.Router = r
	c.Log.Errorf("[ Router ] Initialize ............................................................. [ OK ]")
	return nil
}

// Initialize env/config load and sub moduel init
func Initialize() (*Context, error) {
	c := new(Context)
	c.Conf = new(appdata.Configuration)

	// 환경 변수, 컨피그를 읽어 들이는 과정에서 로그 출력을 위해
	// 아주 기초적인 부분만을 초기화 한다.
	c.InitLog()

	// Env
	if err := c.ReadEnv(); err != nil {
		return nil, err
	}

	// Read Config
	if err := c.ReadConfig(); err != nil {
		return nil, err
	}

	// Datastore
	if err := c.InitDatastore(); err != nil {
		return nil, err
	}

	// Setting Log(from env and conf)
	if err := c.SetLogger(); err != nil {
		return nil, err
	}

	// Echo Framework Init
	if err := c.InitRouter(); err != nil {
		return nil, err
	}

	// TODO: Other Module Init
	redisManager := redis.Manager{}.CreateInstance(c.Log)
	if err := redisManager.ClientInit(c.Conf.RedisServer); err != nil {
		return nil, err
	}
	if err := c.Redis.SetNotifyEvent(true); err != nil {
		return nil, err
	}
	c.Redis = redisManager

	c.Log.Errorf("[ ALL ] Initialize ................................................................ [ OK ]")
	return c, nil
}

// InitDepencyInjection sub model Dependency injection and path regi to server
func (c *Context) InitDepencyInjection() error {
	injector := injectors.Injector{}.New(c.Router, c.Datastore, c.Log, c.Conf)
	injector.Init()
	return nil
}

// StartSubModules Start SubModule And Waiting Signal / Main Loop
func (c *Context) StartSubModules() {
	// Signal
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	c.Log.Errorf("[ Signal ] Listener Start ......................................................... [ OK ]")

	// Echo Framework
	echoServerErr := make(chan error)
	listenAddr := fmt.Sprintf("%s:%d", c.Conf.Server.Host, c.Conf.Server.Port)
	go func() {
		if err := c.Router.Run(listenAddr); err != nil {
			echoServerErr <- err
		}
	}()
	c.Log.Errorf("[ Router ] Listener Start ......................................................... [ OK ]")

	// TODO : Start Other Sub Modules
	go c.Redis.StartEventReceiver()

	err := c.Redis.WriteData(&models.SomeData{
		Key: "key_jbtest",
		SomeDataValue: models.SomeDataValue{
			Value1: "value1_jb",
			Value2: "value2_jb",
		},
	})
	if err != nil {
		c.Log.Errorf("Redis Insert Data Error[ %s ]", err.Error())
	}
	for {
		select {
		case err := <-echoServerErr:
			c.Log.Errorf("[ SERVER ] ERROR[ %s ]", err.Error())
			c.StopSubModules()
			return
		case sig := <-signalChannel:
			c.Log.Errorf("[ SIGNAL ] Receive [ %s ]", sig.String())
			c.StopSubModules()
			return
		case <-time.After(time.Second * 5):
			// 메인 Goroutine에서 주기적으로 무언가를 동작하게 하고 싶다면? 다음을 이용
			// 예 : 성능 시험과 같이 로그가 아닌 통계적인 부분만으로 상태를 체크해야 한다면?
			//c.Log.Errorf("I'm Alive...")

			//printConf, _ := yaml.Marshal(&c.Conf.Configuration)
			//c.Log.Infof("--- conf dump:\n%s\n\n", string(printConf))
		}
	}
}

// StopSubModules Stop Submodules
func (c *Context) StopSubModules() {
	if err := c.Datastore.Shutdown(); err != nil {
		c.Log.Errorf("[ DataStore ] Shutdown .......................................................... [ Fail ]")
		c.Log.Errorf("[ ERROR ] %s", err.Error())
	} else {
		c.Log.Errorf("[ DataStore ] Shutdown ............................................................ [ OK ]")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15*time.Second))
	defer cancel()
	c.Router.Shutdown(ctx)
	c.Log.Errorf("[ Router ] Shutdown ............................................................... [ OK ]")

	// TODO : 사용하는 서브 모듈(Goroutine)들이 안전하게 종료 될 수 있도록 종료 코드를 추가한다.
	c.Redis.Destroy()
	//localdata.JobManager{}.Destroy()
	//c.Log.Errorf("[ JobQueue ] Shutdown ............................................................. [ OK ]")
}

// @title Cache Server API
// @version 1.0.0
// @description This is a cache server.

// @contact.name API Support
// @contact.url http://mobigen.com
// @contact.email irisdev@mobigen.com

// @host localhost:8080
// @BashPath /cache-server/api/v1
func main() {
	// Initialize Sub module And Read Env, Config
	c, err := Initialize()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Initialization Interconnection of WebServer Layer
	// controller - application - domain - repository - infrastructures
	c.InitDepencyInjection()

	// Start sub module and main loop
	c.StartSubModules()

	// Bye Bye
	c.Log.Shutdown()
}
