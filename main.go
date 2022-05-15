package main

import (
	"context"
	"fmt"
	"log"
	web "minecraft-control-panel/internal/controllers"
	"minecraft-control-panel/internal/controllers/middleware/antiflood"
	"minecraft-control-panel/internal/global"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	var err error

	//gin.DisableConsoleColor()

	ginRouter := gin.Default()
	global.LastOnLineTime = fmt.Sprint(time.Now().Format("2006-01-02_15:04:05")) //紀錄Web Server 啟動時間

	//--------初始化anti SYN-FOOD模組------------------{
	antiflood.StartTimer(&antiflood.Cfgs{
		MaxCount:          global.APILimiterMaxCount,
		RateSecs:          global.APILimiterRateSecs,
		ClientIPUnbanSecs: global.APILimiterClientIPUnbanSecs,
		TimerUnbanSecs:    global.APILimiterTimerUnbanSecs,
		ProxyHeaderKey:    global.APIProxyHeaderClientIP,
	}) //啟動反SYN-Flood計時器

	if antiflood.SetAllowIPRule(global.APILimiterAllowIPList) != nil {
		log.Fatal("IP白名設定失敗 > ", err.Error())
	} else if len(global.APILimiterAllowIPList) > 0 && (len(antiflood.ClientIPAllowList)+len(antiflood.ClientIPRuleList)) > 0 {
		log.Printf("%v : 載入開放無限制使用API > 准許IP清單:%v 准許IP規則:%v\n", antiflood.Service_Name, antiflood.ClientIPAllowList, antiflood.ClientIPRuleList)
	}
	//------------END--------------------------------}

	//ginRouter.Use中每個路徑都會先加載func
	ginRouter.Use(func(ginCTX *gin.Context) {
		antiflood.TriggerHandler(ginCTX) //反SYN-Flood
	})

	ginRouter.LoadHTMLGlob(global.WebPath_Public + "/" + global.WebPath_Templs + "/*.html")         //.LoadHTMLGlob()僅會在板模後傳給客戶端，該html檔案並不會被公開
	ginRouter.StaticFS(global.WebURLPage, http.Dir(global.WebPath_Public+"/"+global.WebPath_Pages)) //公開靜態頁面

	ginGroup1 := ginRouter.Group("/" + global.WebURLRoot)
	{
		//StaticFS()可指定目錄內哪些檔案可以被公開
		//ginRouter.Static("BoswerPath", "./Public") //Static()會將目錄內檔案公開出去，較不安全
		//ginRouter.StaticFS("BoswerPath", http.Dir("./Public")) //StaticFS()可指定目錄內哪些檔案可以被公開
		//ginRouter.StaticFile(global.WebPath_Images, "./Public/gabo.png") //只能公開一個檔案
		ginGroup1.StaticFS(global.WebPath_Images, http.Dir(global.WebPath_Public+"/"+global.WebPath_Images))
		ginGroup1.StaticFS(global.WebPath_Others, http.Dir(global.WebPath_Public+"/"+global.WebPath_Others))
		ginGroup1.StaticFS(global.WebPath_Downloads, http.Dir(global.WebPath_Public+"/"+global.WebPath_Downloads))

		//ginGroup1.GET("/debug", api.HandleDebug)
		ginGroup1.GET("/", web.HandleIndexMenu)                      //根目錄
		ginGroup1.GET(global.APIName_Login, web.HandleLogin)         //GET 登陸api接口,包含所有服務需通過接口
		ginGroup1.POST(global.APIName_Login, web.HandleLogin)        //POST 登陸api接口,包含所有服務接需通過接口
		ginGroup1.GET(global.Oauth2CallbackName, web.HandleCallback) //GET 驗證客戶端獲得的GoogleOauth2認證碼,並藉由Query:'state' 作為Param後開始後續使用服務
		ginGroup1.POST(global.APIName_Upload, web.HandleUploadFile)  //POST 上傳檔案接口,上傳過程包含Token驗證
	}

	//--------初始化API Router(gin)模組-----------------{
	gin.SetMode("release") //debug mode會有較多的性能資訊
	apiServ := &http.Server{
		Addr:           ":" + global.HostPort, //host port
		Handler:        ginRouter,             //ginRouter
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20, //計算原理不明,限制Request接收長度
	}

	go func() {
		if err := apiServ.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP端口監聽失敗 > ", err.Error())
		}
	}()
	//------------END---------------------------------}

	//========================往下為關機後觸發func========================================{
	quitChan := make(chan os.Signal, 10) //Notify：系統訊號轉將發至channel
	signal.Notify(quitChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quitChan
	//永久等待收到系統發出訊號類型,若符合則往下執行
	close(quitChan)

	timeoutSec := 4
	timeoutCTX, timeoutClose := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	//透過context.WithTimeout產生一個新的子context，它的特性是有生命週期，只要超過x秒就會自動發出Done()的訊息
	if err := apiServ.Shutdown(timeoutCTX); err != nil {
		log.Fatal("API關閉服務過程中發生錯誤:", err)
	}
	//========================往下為timeout後觸發func======================================={
	go func() {
		<-timeoutCTX.Done() //timeout觸發
		timeoutClose()
		log.Printf("超時%v秒, API服務強制關閉！\n", timeoutSec)
		os.Exit(3)
	}()
	log.Printf("\n正在關閉API服務... \n")
	//====================================================================================}

	//========================往下可帶入timeoutCTX執行自定義需要觸發關機程序的func============================={
	//func(timeoutCTX)
	//====================================================================================}

	log.Printf("API服務正常已關閉！\n")
	return
}
