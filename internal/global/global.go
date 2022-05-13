package global

import (
	"fmt"
	"minecraft-control-panel/internal/common/customVar"
	"minecraft-control-panel/internal/global/config"
	"strings"
)

const (
	WebURLRoot        = "minecraft" //網頁根目錄 http://localhost:8080/xxx
	WebPath_Public    = "./website" //網頁公用檔案目錄
	WebPath_Images    = "images"    //網頁公用目錄名稱 - 圖片
	WebPath_Pages     = "pages"     //網頁公用目錄名稱 - 純html頁面
	WebPath_Templs    = "templates" //網頁公用目錄名稱 - 版模
	WebPath_Downloads = "downloads" //網頁公用目錄名稱 - 下載暫存
	WebPath_Others    = "others"    //網頁公用目錄名稱 - 其他檔案分類

	//ConfigPath     = "./configs"
	ConfigFilePath = "./configs/webstie.conf" //環境變數設定檔

	APIName_Login  = "/login"  //登陸api名稱
	APIName_Upload = "/upload" //上傳檔案api名稱

	APIPost_UploadType = "upload_type" //上傳檔案Post 標籤名稱
	APIIndexMenuValue  = "indexMenu"   //首頁api1標籤ㄋ

	H5Page_Menu          = "HTMLPage_Menu"
	H5Value_Title        = "HTMLValue_Title"
	H5Value_LastUpdate   = "HTMLValue_LastUpdate"
	H5Value_Token        = "HTMLValue_Token"
	H5Value_UploadType   = "HTMLValue_Upload_Type"
	H5Output_Token       = "HTMLOutput_Token"
	H5Output_Upload_Type = "HTMLOutput_Upload_Type"
)

var (
	LastOnLineTime string //伺服器啟動時間

	WebURLPage = WebURLRoot + "/" + WebPath_Pages //直接顯示頁面

	Oauth2ClientID     string   //= "32054680383-r6adfs24i84edvqtfr3eesskums6l2n1.apps.googleusercontent.com"
	Oauth2ClientSecret string   //= "_cha_wE8IT-pMGeGUDrbcn9j"
	Oauth2CallbackName string   //= "callback"
	Oauth2EmailList    []string //准許的google email 清單

	HostScheme     string //= "http://"
	HostDomainName string //= "localhost"
	HostPort       string //= "8080"

	APIProxyHeaderClientIP string //獲取代理服務器之ClientIP之header名稱
	APICallbackServiceName string //API回調服務名稱
	APICallbackMaxNum      int    // API傳送回調最高上限次數

	APILimiterMaxCount          uint64   //短時間連線次數上限
	APILimiterRateSecs          int64    //定義短時間秒數
	APILimiterClientIPUnbanSecs int64    //遭禁ClientIP自動解封時間
	APILimiterTimerUnbanSecs    int      //自動解封時間執行周期
	APILimiterAllowIPList       []string //限制器白名單IP

	MCSVPath_LastLog                  string //= "/Users/Duke/+WorkDir/Minecraft Server/logs/latest.log"
	MCSVPath_PlayerAllowList          string //= "/Users/Duke/+WorkDir/Minecraft Server/whitelist.json"
	MCSVPath_ServerProperties         string //= "/Users/Duke/+WorkDir/Minecraft Server/server.properties"
	MCSVPath_RestartScript            string //重啟伺服器腳本 //= "./configs/restart_script.sh"
	MCSVPath_StopScript               string //暫停伺服器腳本 //= "./configs/stop_script.sh"
	MCSVPath_DeleteModsScript         string //刪除模組檔案腳本 //= "./configs/delete_mods_script.sh"
	MCSVPath_DeleteMapOverWorldScript string //刪除模組主世界地圖腳本 //= "./configs/delete_map_overworld_script.sh"
	MCSVPath_DeleteMapNetherScript    string //刪除模組地獄地圖腳本 //= "./configs/delete_map_nether_script.sh"
	MCSVPath_DeleteMapEndScript       string //刪除模組終界地圖腳本 //= "./configs/delete_map_end_script.sh"
	MCSVPath_CustomScript1            string //刪除模組檔案腳本 //= "./configs/custom_script1.sh"
	MCSVPath_CustomScript2            string //刪除模組檔案腳本 //= "./configs/custom_script2.sh"
	MCSVPath_CustomScript3            string //刪除模組檔案腳本 //= "./configs/custom_script3.sh"
	MCSVPath_CustomScript4            string //刪除模組檔案腳本 //= "./configs/custom_script4.sh"
	MCSVPath_CustomScript5            string //刪除模組檔案腳本 //= "./configs/custom_script5.sh"

	Minecraft_TransmissionList_HideText     string                 //前端網頁隱藏目錄文字
	Minecraft_TransmissionList_Excludes     map[string]struct{}    //下載/上傳/刪除清單中要忽略的檔案名稱
	Minecraft_TransmissionList_Paths     map[string]interface{} //下載/上傳/刪除清單設定
	//Minecraft_TransmissionList_ZIP_Paths map[string]interface{} //下載打包成ZIP清單設定

	SettingMap = config.SettingMap //參照package config

)

func init() {

	Minecraft_TransmissionList_Excludes = map[string]struct{}{
		".DS_Store": {},
		"Thumbs.db": {},
	}

	//環境變數設定值解析 SettingMap["XXX"] = &config.AddSetings{初始值, 輸出變數指標, 自訂設定值類型}
	SettingMap["host_scheme"] = &config.AddSetings{DefaultValue: "http://", OutputPointer: &HostScheme, Custom: &customVar.StringType{}}
	SettingMap["host_domain_name"] = &config.AddSetings{DefaultValue: "", OutputPointer: &HostDomainName, Custom: &customVar.StringType{}}
	SettingMap["host_port"] = &config.AddSetings{DefaultValue: "8080", OutputPointer: &HostPort, Custom: &customVar.Uint16Type{}}

	SettingMap["api_proxy_header_clientip"] = &config.AddSetings{DefaultValue: "Cf-Connecting-Ip", OutputPointer: &APIProxyHeaderClientIP, Custom: &customVar.StringType{}}
	SettingMap["api_limter_max_count"] = &config.AddSetings{DefaultValue: 3000, OutputPointer: &APILimiterMaxCount, Custom: &customVar.Uint64Type{}}
	SettingMap["api_limter_rate_secs"] = &config.AddSetings{DefaultValue: 10, OutputPointer: &APILimiterRateSecs, Custom: &customVar.SecondsInADay{}}
	SettingMap["api_limter_client_unban_sec"] = &config.AddSetings{DefaultValue: 600, OutputPointer: &APILimiterClientIPUnbanSecs, Custom: &customVar.SecondsInADay{}}
	SettingMap["api_limter_timer_unban_sec"] = &config.AddSetings{DefaultValue: 1800, OutputPointer: &APILimiterTimerUnbanSecs, Custom: &customVar.SecondsInADay{}}
	SettingMap["api_limter_allow_IPList_add"] = &config.AddSetings{DefaultValue: nil, OutputPointer: &APILimiterAllowIPList, Custom: &customVar.BaseSlice{}}

	SettingMap["google_oauth2_clientID"] = &config.AddSetings{DefaultValue: "", OutputPointer: &Oauth2ClientID, Custom: &customVar.StringType{}}
	SettingMap["google_oauth2_secret_code"] = &config.AddSetings{DefaultValue: "", OutputPointer: &Oauth2ClientSecret, Custom: &customVar.StringType{}}
	SettingMap["google_oauth2_callback_apiname"] = &config.AddSetings{DefaultValue: "callback", OutputPointer: &Oauth2CallbackName, Custom: &customVar.StringType{}}
	SettingMap["add_google_email"] = &config.AddSetings{DefaultValue: "callback", OutputPointer: &Oauth2EmailList, Custom: &customVar.BaseSlice{}}

	SettingMap["minecraft_path_last_log"] = &config.AddSetings{DefaultValue: "/Minecraft Server/logs/latest.log", OutputPointer: &MCSVPath_LastLog, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_whitelist"] = &config.AddSetings{DefaultValue: "/Minecraft Server/whitelist.json", OutputPointer: &MCSVPath_PlayerAllowList, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_server_properties"] = &config.AddSetings{DefaultValue: "/Minecraft Server/server.properties", OutputPointer: &MCSVPath_ServerProperties, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_restart_script"] = &config.AddSetings{DefaultValue: "./configs/restart_script.sh", OutputPointer: &MCSVPath_RestartScript, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_stop_script"] = &config.AddSetings{DefaultValue: "./configs/stop_script.sh", OutputPointer: &MCSVPath_StopScript, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_delete_mods_script"] = &config.AddSetings{DefaultValue: "./configs/delete_mods_script.sh", OutputPointer: &MCSVPath_DeleteModsScript, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_delete_map_overworld_script"] = &config.AddSetings{DefaultValue: "./configs/delete_map_overworld_script.sh", OutputPointer: &MCSVPath_DeleteMapOverWorldScript, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_delete_map_nether_script"] = &config.AddSetings{DefaultValue: "./configs/delete_map_nether_script.sh", OutputPointer: &MCSVPath_DeleteMapNetherScript, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_delete_map_end_script"] = &config.AddSetings{DefaultValue: "./configs/delete_map_end_script.sh", OutputPointer: &MCSVPath_DeleteMapEndScript, Custom: &customVar.StringType{}}

	SettingMap["minecraft_path_custom1_script"] = &config.AddSetings{DefaultValue: "./configs/custom1_script.sh", OutputPointer: &MCSVPath_CustomScript1, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_custom2_script"] = &config.AddSetings{DefaultValue: "./configs/custom2_script.sh", OutputPointer: &MCSVPath_CustomScript2, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_custom3_script"] = &config.AddSetings{DefaultValue: "./configs/custom3_script.sh", OutputPointer: &MCSVPath_CustomScript3, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_custom4_script"] = &config.AddSetings{DefaultValue: "./configs/custom4_script.sh", OutputPointer: &MCSVPath_CustomScript4, Custom: &customVar.StringType{}}
	SettingMap["minecraft_path_custom5_script"] = &config.AddSetings{DefaultValue: "./configs/custom5_script.sh", OutputPointer: &MCSVPath_CustomScript5, Custom: &customVar.StringType{}}

	SettingMap["minecraft_transmission_list_hidetext"] = &config.AddSetings{DefaultValue: "", OutputPointer: &Minecraft_TransmissionList_HideText, Custom: &customVar.StringType{}}
	SettingMap["add_minecraft_transmission_list"] = &config.AddSetings{DefaultValue: nil, OutputPointer: &Minecraft_TransmissionList_Paths, Custom: &filesMapType{}}
	//SettingMap["add_minecraft_transmission_zip_list"] = &config.AddSetings{DefaultValue: nil, OutputPointer: &Minecraft_TransmissionList_ZIP_Paths, Custom: &filesMapType{}}

	config.InitLoad(ConfigFilePath)
	// fmt.Printf("host_scheme=%v\n", HostScheme)
	// fmt.Printf("host_domain_name=%v\n", HostDomainName)
	// fmt.Printf("host_port=%v\n", HostPort)
	// fmt.Printf("google_oauth2_clientID=%v\n", Oauth2ClientID)
	// fmt.Printf("google_oauth2_secret_code=%v\n", Oauth2ClientSecret)
	// fmt.Printf("google_oauth2_callback_apiname=%v\n", Oauth2CallbackName)
	// fmt.Printf("add_google_email=%v\n", Oauth2EmailList)
	// fmt.Printf("minecraft_path_last_log=%v\n", MCSVPath_LastLog)
	// fmt.Printf("minecraft_path_whitelist=%v\n", MCSVPath_PlayerAllowList)
	// fmt.Printf("minecraft_path_server_properties=%v\n", MCSVPath_ServerProperties)
	// fmt.Printf("minecraft_path_restart_script=%v\n", MCSVPath_RestartScript)
	// fmt.Printf("add_minecraft_download_list=%v\n", Minecraft_DownloadList_Settings)

}

//自定義格式 key|value
type filesMapType struct{}

func (_ filesMapType) GetValue(inputValue string) (output interface{}, err error) {

	cacheSlice := strings.Split(inputValue, "|")
	if len(cacheSlice) < 1 {
		return nil, fmt.Errorf("'%v' > ' | ' 分隔號數量不可低於1", inputValue)
	} else if cacheSlice[0] == "" || cacheSlice[1] == "" {
		return nil, fmt.Errorf("'%v' > 部份數值不可為空", inputValue)
	}

	cacheMap := make(map[string]interface{})
	cacheMap[cacheSlice[0]] = cacheSlice[1]
	return cacheMap, nil
}
