package main

import (
	"ARTeamViewService/conf"
	_ "ARTeamViewService/docs"
	"ARTeamViewService/global"
	"ARTeamViewService/models"
	"ARTeamViewService/tasks"
	"ARTeamViewService/utils"
	"ARTeamViewService/web/controllers"
	"ARTeamViewService/web/middleware"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	_ "github.com/mattn/go-sqlite3"
	json "github.com/wulinlw/jsonAmbiguity"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// @title anyRTC arteamview server api server
// @version 1.0.0
// @description This is a arteamview server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 192.168.1.140:12681
// @BasePath /arapi/v1

/**
 * 当前版本
 */
func version() (ver string) {
	ver = "2.1.0"
	fmt.Printf("Welcome to use anyRTC arteamview server\r\n"+
		"Current version is %s, thanks you very much!!!\r\n",
		ver)
	return
}

/**
 * 获取当前路径
 */
func curDir() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	fmt.Println("OS is:", runtime.GOOS)
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return "", errors.New("can not find / or \\")
	}
	fmt.Println("The path is:", path)
	fmt.Println("The dir is:", string(path[0:i+1]))
	return string(path[0 : i+1]), nil
}

/**
 * 获取配置信息
 */
func getConfigToUse(conf *conf.AllConfig) {
	fmt.Println("Config to use is:", conf.UseConf)
	switch conf.UseConf {
	case "d":
		global.GConfig = conf.DConfig
		break
	case "r":
		global.GConfig = conf.RConfig
		break
	}
}

/**
 * 打开数据库
 */
func openDb() (err error) {
	var dbPath = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?collation=utf8mb4_general_ci",
		global.GConfig.Mysql.User, global.GConfig.Mysql.Password, global.GConfig.Mysql.Host, global.GConfig.Mysql.Port, global.GConfig.Mysql.Database)
	driverName := "mysql"

	global.GDb, err = sql.Open(driverName, dbPath)
	if err != nil {
		global.GLogger.Error("open mysql failed", dbPath, err)
		return err
	}

	if err = global.GDb.Ping(); err != nil {
		global.GLogger.Error("ping mysql failed", dbPath, err)
		return err
	}
	global.GLogger.Info("open mysql success")

	//连接最长存活期，超过这个时间连接将不再被复用
	global.GDb.SetConnMaxLifetime(3600 * time.Second)
	//最大空闲连接数
	global.GDb.SetMaxIdleConns(10)
	//最大连接数
	global.GDb.SetMaxOpenConns(50)

	return nil
}

/**
 * 关闭数据库
 */
func closeDb() {
	if global.GDb != nil {
		_ = global.GDb.Close()
	}
}

/**
 * 初始化context
 */
func initContext() {
	global.Context = context.Background()
}

/**
 * 打开redis
 */
func openRedis() (err error) {
	addr := fmt.Sprintf("%s:%d", global.GConfig.Redis.Host, global.GConfig.Redis.Port)

	global.GRedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: global.GConfig.Redis.Password,
		DB:       global.GConfig.Redis.Db,
	})
	pong, err := global.GRedisClient.Ping(global.Context).Result()
	if err != nil {
		global.GLogger.Error("Open redis failed, err:", err)
		return err
	}
	global.GLogger.Info("Redis pong:", pong)
	return nil
}

/**
 * 关闭redis
 */
func closeRedis() {
	if global.GRedisClient != nil {
		_ = global.GRedisClient.Close()
	}
}


func main() {

	//// get version and name
	_ = version()

	//// get current path
	curDir, err := curDir()
	if err != nil {
		fmt.Println("get cur dir failed", err)
	}

	//// load config
	confPath := curDir + "conf/config.toml"
	fmt.Println("--->", os.Args, len(os.Args))
	if len(os.Args) >= 2 {
		confPath = os.Args[1]
	}
	fmt.Println("the conf path is:", confPath)
	config := conf.NewAllConfig()
	//conf.LoadJsonConfig(confPath)
	config.LoadTomlConfig(confPath)
	getConfigToUse(config)

	// init log
	global.GLogger = utils.GetLogger(global.GConfig.Logger)

	byteConfig, configErr := json.Marshal(global.GConfig)
	if nil != configErr {
		global.GLogger.Error("the config to configErr:", configErr)
	}

	global.GLogger.Info("the conf path is:", confPath)
	global.GLogger.Info("the config is: ", string(byteConfig))

	title := global.GConfig.ProcessName

	global.GLogger.Info("the process name is:", title)
	//// write pid
	pid := os.Getpid()
	pidDir := global.GConfig.PidDir
	if pidDir == "" {
		pidDir = "."
	}
	pidPath := pidDir + "/" + title + ".pid"
	err = ioutil.WriteFile(pidPath, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		global.GLogger.Fatal("Write pid file failed,err:", err)
	}

	//// open db
	err = openDb()
	if err != nil {
		global.GLogger.Fatal(err)
	}

	////初始化context
	initContext()

	//// init redis client
	err = openRedis()
	if err != nil {
		global.GLogger.Fatal(err)
	}

	//// init iris
	app := iris.New()

	//app.RegisterView(iris.HTML("./web/views", ".html"))

	//// init docs
	//// TODO:: when build release, the bin will smaller if comment this line
	//app.Get("/swagger/*arapi", swagger.WrapHandler(swaggerFiles.Handler))

	swaggerConfig := &swagger.Config{
		URL: "http://192.168.1.140:12681/swagger/doc.json", //The url pointing to API definition
	}
	// use swagger middleware to
	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(swaggerConfig, swaggerFiles.Handler))

	//// handle error
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		_, _ = ctx.JSON(models.ApiJsonResp(404, "404 not found", nil))
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		_, _ = ctx.JSON(models.ApiJsonResp(500, "Oups something went wrong, try again", nil))
	})

	iris.RegisterOnInterrupt(func() {
		global.GLogger.Fatalln("Error: RegisterOnInterrupt!!!")
		closeDb()
		closeRedis()
	})

	//// handle cors
	irisCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	//// config router, platform
	apiPrefix := "/arapi/v1"
	v1 := app.Party(apiPrefix, irisCors).AllowMethods(iris.MethodOptions)
	{

		v1.PartyFunc("/teamview/", func(svr router.Party) {
			//用户登录
			svr.Post("signIn", controllers.GetInsApiUser().SignIn)

			//消息回调
			svr.Post("teamViewRtcNotify", controllers.GetInsApiUser().TeamViewRtcNotify)
			//录像回调
			svr.Post("teamViewVodNotify", controllers.GetInsApiUser().TeamViewVodNotify)
		})

		v1.PartyFunc("/teamview/", func(svr router.Party) {
			svr.Use(middleware.JwtHandler().Serve, middleware.JwtAuthToken)

			//创建房间
			svr.Post("insertRoom", controllers.GetInsApiUser().InsertRoom)
			//获取用户信息
			svr.Post("getUserInfo", controllers.GetInsApiUser().GetUserInfo)
			//加入房间
			svr.Post("joinRoom", controllers.GetInsApiUser().JoinRoom)
			//用户离开房间
			svr.Post("leaveRoom", controllers.GetInsApiUser().LeaveRoom)
			//获取房间列表
			svr.Post("getRoomList", controllers.GetInsApiUser().GetRoomList)
			//获取专家列表
			svr.Post("getSpecialist", controllers.GetInsApiUser().GetSpecialist)
			//记录用户在线心跳包信息
			svr.Post("insertUserOnlineInfo", controllers.GetInsApiUser().InsertUserOnlineInfo)
		})
	}

	//// run loop
	go tasks.ServeCheckRunLoop()

	//// start & listen server
	addr := fmt.Sprintf("%s:%d", global.GConfig.Host, global.GConfig.Port)
	err = app.Run(

		iris.Addr(addr),
		iris.WithOptimizations,
		//设置此选项为了body的多次消费,可以二次获取body
		iris.WithoutBodyConsumptionOnUnmarshal,
	)
	if err != nil {
		global.GLogger.Error(err)
		panic(err)
	}
}
