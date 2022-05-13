package web

import (
	"bytes"
	"fmt"
	"log"
	"minecraft-control-panel/internal/common"
	"minecraft-control-panel/internal/controllers/compress"
	"minecraft-control-panel/internal/controllers/encrypt"
	"minecraft-control-panel/internal/controllers/googleOauth2"
	"minecraft-control-panel/internal/controllers/minecraft"
	"minecraft-control-panel/internal/global"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//Design Pattern : Factory
// type Handlers interface {
// 	Action(account *googleOauth2.GoogleAcc, w http.ResponseWriter, r *http.Request)
// }

const (
	mc_DownloadFile = "mc-downland-file"
	mc_ZIPFolders   = "mc-zip-folders"
	mc_DeleteFile   = "mc-delete-file"
	mc_RefreshFiles = "mc-refresh-files"
	mc_PathChar     = "?" //callback?state= 分隔符號
)

type Handlers interface {
	Action(ginCTX *gin.Context, account *googleOauth2.GoogleAcc)
}

//callback通過帳號認證後,所執行func
func NewHandler(state string) (Handlers, error) {
	var stateType string
	var statePlus string
	postFrom := make(map[string]string) //post資料處理
	if strings.Count(state, "~") == 1 {
		decode, err := encrypt.Decode(state[strings.Index(state, "~")+1:])
		if err != nil {
			log.Printf("Post解碼錯誤 : %v (%v)", decode, err)
		} else {
			for _, lineStr := range strings.Split(decode, "\n") {
				if strings.Count(lineStr, "\t") != 1 {
					log.Printf("Post解碼異常 : %v (%v)", decode, lineStr)
					break
				}
				postFrom[lineStr[:strings.Index(lineStr, "\t")]] = lineStr[strings.Index(lineStr, "\t")+1:]
			}
		}

	}

	if strings.Contains(state, ":") == false {
		stateType = state
	} else {

		stateType = state[:strings.Index(state, ":")]
		statePlus = state[strings.Index(state, ":")+1:]
	}

	switch stateType {

	case "test": //測試callback認證功能
		return &HandlerTestAcc{}, nil
	case global.APIIndexMenuValue: //轉入首頁選單
		return &HandlerIndexMenu{}, nil

	case "mc-showLogs": //Mincraft伺服器監控 - 查看log
		return &HandlerMinecraftShowLogs{}, nil
	case "mc-showControl": //Mincraft伺服器操作 - 伺服器操作
		return &HandlerMinecraftShowControl{StateData: statePlus}, nil

	case "mc-showDownland": //Mincraft伺服器檔案 - 顯示下載檔案頁面
		return &HandlerMinecraftDownloadList{}, nil
	case "mc-showUpload": //Mincraft伺服器檔案 - 顯示下載檔案頁面
		return &HandlerMinecraftUploadMenu{}, nil
	case "mc-showDelete": //Mincraft伺服器檔案 - 顯示下載檔案頁面
		return &HandlerMinecraftDeleteList{}, nil

	case mc_RefreshFiles: //Mincraft伺服器檔案 - 刷新伺服器檔案目錄
		return &HandlerMinecraftRefreshFiles{}, nil
	case mc_ZIPFolders: //Mincraft伺服器檔案 - 認證後並讓客戶端下載檔案
		return &HandlerMinecraftZIPDownloadFolder{}, nil
	case mc_DownloadFile: //Mincraft伺服器檔案 - 認證後並讓客戶端下載檔案
		return &HandlerMinecraftDownloadFile{StateData: statePlus}, nil
	case mc_DeleteFile: //Mincraft伺服器檔案 - 認證後並讓客戶端刪除檔案
		return &HandlerMinecraftDeleteFile{StateData: statePlus}, nil

	case "mc-showAllowList": //Mincraft伺服器管理 - 顯示白名單
		return &HandlerMinecraftAllowList{}, nil
	case "mc-addAllowList": //Mincraft伺服器管理 - 添加白名單
		return &HandlerMinecraftAllowList{FromData: postFrom}, nil

	default:
		return nil, fmt.Errorf("get url query state 設置錯誤！")
	}

}

//測試
type HandlerTestAcc struct{}

func (_ *HandlerTestAcc) Action(ginCTX *gin.Context, account *googleOauth2.GoogleAcc) {
	//fmt.Printf("Client was login : %v\n", account.Email)
	ginCTX.Writer.WriteString("Is working. \n")
	//ginCTX.Writer.WriteString(fmt.Sprintf("Your Email: %v\n", account.Email))
	ginCTX.Abort()
}

type HandlerIndexMenu struct{}

func (_ *HandlerIndexMenu) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {
	ginCTX.Redirect(http.StatusFound, "./"+global.WebPath_Pages+"/"+"index-Menu.html")
	ginCTX.Abort()
}

type HandlerMinecraftShowLogs struct{}

func (_ *HandlerMinecraftShowLogs) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {
	logText, err := minecraft.ReadLogFile()
	if err != nil {
		logText = err.Error()
	}

	ginCTX.HTML(http.StatusOK, "minecraft-ShowLogs.html", gin.H{
		global.H5Value_LastUpdate: global.LastOnLineTime,
		global.H5Value_Title:      "Minecraft - Log監控",
		global.H5Page_Menu:        "./" + global.WebPath_Pages + "/" + "index-Menu.html",
		"HTMLValue_OutputLogText": logText,
	})

	ginCTX.Abort()
}

type HandlerMinecraftShowControl struct {
	StateData string
}

const (
	scriptRestart            = "restart-minecraft"
	scriptStop               = "stop-minecraft"
	scriptDeleteMods         = "delete-mods"
	scriptDeleteMapOverworld = "delete-map-overworld"
	scriptDeleteMapNether    = "delete-map-nether"
	scriptDeleteMapEnd       = "delete-map-end"
	scriptCustom1            = "script-custom1"
	scriptCustom2            = "script-custom2"
	scriptCustom3            = "script-custom3"
	scriptCustom4            = "script-custom4"
	scriptCustom5            = "script-custom5"
)

func (feild *HandlerMinecraftShowControl) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {

	var outputMessage string
	var outputReturn string

	if feild.StateData != "" {
		var err error
		var scriptPath string
		switch feild.StateData {
		case scriptRestart:
			scriptPath = global.MCSVPath_RestartScript
		case scriptStop:
			scriptPath = global.MCSVPath_StopScript
		case scriptDeleteMods:
			scriptPath = global.MCSVPath_DeleteModsScript
		case scriptDeleteMapOverworld:
			scriptPath = global.MCSVPath_DeleteMapOverWorldScript
		case scriptDeleteMapNether:
			scriptPath = global.MCSVPath_DeleteMapNetherScript
		case scriptDeleteMapEnd:
			scriptPath = global.MCSVPath_DeleteMapEndScript
		case scriptCustom1:
			scriptPath = global.MCSVPath_CustomScript1
		case scriptCustom2:
			scriptPath = global.MCSVPath_CustomScript2
		case scriptCustom3:
			scriptPath = global.MCSVPath_CustomScript3
		case scriptCustom4:
			scriptPath = global.MCSVPath_CustomScript4
		case scriptCustom5:
			scriptPath = global.MCSVPath_CustomScript5
		default:
			println()
			fmt.Println(feild.StateData)
			ginCTX.Abort() //跑到這邊肯定客戶端故意為之
			return
		}

		outputMessage, err = minecraft.RunScript(scriptPath)
		if err != nil {
			outputMessage = "Script path : " + global.MCSVPath_RestartScript + "\r\n" + err.Error()
		}
		outputReturn = "執行腳本時間:" + time.Now().Format("2006-01-02_15:04:05")
	}

	ginCTX.HTML(http.StatusOK, "minecraft-ShowControl.html", gin.H{
		global.H5Value_LastUpdate:         global.LastOnLineTime,
		global.H5Value_Title:              "Minecraft - 伺服器操作",
		global.H5Page_Menu:                "./" + global.WebPath_Pages + "/" + "index-Menu.html",
		"HTMLValue_RunRestartScript":      "." + global.APIName_Login + "?state=mc-showControl:" + scriptRestart,
		"HTMLValue_RunStopScript":         "." + global.APIName_Login + "?state=mc-showControl:" + scriptStop,
		"HTMLValue_RunDeleteModsScript":   "." + global.APIName_Login + "?state=mc-showControl:" + scriptDeleteMods,
		"HTMLValue_RunDeleteOverWorldMap": "." + global.APIName_Login + "?state=mc-showControl:" + scriptDeleteMapOverworld,
		"HTMLValue_RunDeleteNetherMap":    "." + global.APIName_Login + "?state=mc-showControl:" + scriptDeleteMapNether,
		"HTMLValue_RunDeleteEndMap":       "." + global.APIName_Login + "?state=mc-showControl:" + scriptDeleteMapEnd,
		"HTMLValue_RunCustomScript1":      "." + global.APIName_Login + "?state=mc-showControl:" + scriptCustom1,
		"HTMLValue_RunCustomScript2":      "." + global.APIName_Login + "?state=mc-showControl:" + scriptCustom2,
		"HTMLValue_RunCustomScript3":      "." + global.APIName_Login + "?state=mc-showControl:" + scriptCustom3,
		"HTMLValue_RunCustomScript4":      "." + global.APIName_Login + "?state=mc-showControl:" + scriptCustom4,
		"HTMLValue_RunCustomScript5":      "." + global.APIName_Login + "?state=mc-showControl:" + scriptCustom5,

		"HTMLValue_OutputMessage": outputMessage,
		"HTMLValue_OutputReturn":  outputReturn,
	})
	ginCTX.Abort()
}

type HandlerMinecraftAllowList struct {
	FromData map[string]string
}

func (field *HandlerMinecraftAllowList) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {

	var outputReturn string
	if len(field.FromData) > 0 {
		minecraft.AddPlayerAllowList(field.FromData["HTMLInput_AllowUUID"], field.FromData["HTMLInput_AllowName"])
		outputReturn = "白名單寫入完成:" + time.Now().Format("2006-01-02_15:04:05")
	}

	allowMap, err := minecraft.GetPlayerAllowList()
	var outputMessage string
	if err != nil {
		outputMessage = err.Error()
	} else {
		allowSlice := []string{}
		for uuid, name := range allowMap {
			allowSlice = append(allowSlice, name+"	:	"+uuid)
		}
		sort.Strings(allowSlice)
		for _, value := range allowSlice {
			outputMessage = outputMessage + value + "\n"
		}
	}

	ginCTX.HTML(http.StatusOK, "minecraft-ShowAllowList.html", gin.H{
		global.H5Value_LastUpdate: global.LastOnLineTime,
		global.H5Value_Title:      "Minecraft - 添加白名單",
		global.H5Page_Menu:        "./" + global.WebPath_Pages + "/" + "index-Menu.html",

		"HTMLValue_AddAllowList":     "." + global.APIName_Login + "?state=mc-addAllowList",
		"HTMLValue_OutputPlayerText": outputMessage,
		"HTMLValue_OutputReturn":     outputReturn,
	})

	ginCTX.Abort()
}

type HandlerMinecraftDownloadList struct{}

func (_ *HandlerMinecraftDownloadList) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {

	var downloadLinks string
	pathBase := global.Oauth2CallbackName + "?state=" + mc_DownloadFile
	buffer := bytes.NewBuffer([]byte{})

	count, tags, relativePaths, fullPaths := minecraft.TransDownload.GetFolderList()
	if len(tags) != count || len(relativePaths) != count || len(fullPaths) != count {
		log.Println("獲取打包下載清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.String(http.StatusBadGateway, "獲取打包下載清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.Abort()
		return
	}

	// if count > 0 {
	// 	buffer.WriteString(fmt.Sprintf(`<a href="%v" >%v</a><br>`, pathBase+":"+encode, "[@打包下載] "+tags[index]+" ▶ "+relativePaths[index]))
	// }

	for index := 0; index < count; index++ {
		encode, err := encrypt.Encode(tags[index] + mc_PathChar + relativePaths[index])
		if err != nil {
			log.Println("製作下載清單加密內容發生錯誤 :", err)
		} else {
			buffer.WriteString(fmt.Sprintf(`<a href="%v" >%v</a><br>`, pathBase+":"+encode, "[@打包下載] "+tags[index]+" ▶ "+relativePaths[index]))
		}
	}

	buffer.WriteString("\n<br>\n")

	count, tags, relativePaths, fullPaths = minecraft.TransDownload.GetFileList()
	if len(tags) != count || len(relativePaths) != count || len(fullPaths) != count {
		log.Println("獲取下載清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.String(http.StatusBadGateway, "獲取下載清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.Abort()
		return
	}

	for index := 0; index < count; index++ {
		encode, err := encrypt.Encode(tags[index] + mc_PathChar + relativePaths[index])
		if err != nil {
			log.Println("製作下載清單加密內容發生錯誤 :", err)
		} else {
			buffer.WriteString(fmt.Sprintf(`<a href="%v" >%v</a><br>`, pathBase+":"+encode, tags[index]+" ▶ "+relativePaths[index]))
		}
	}

	downloadLinks = buffer.String()
	ginCTX.HTML(http.StatusOK, "minecraft-ShowDownland.html", gin.H{
		global.H5Value_LastUpdate:       global.LastOnLineTime,
		global.H5Value_Title:            "Minecraft - 伺服器操作",
		global.H5Page_Menu:              "./" + global.WebPath_Pages + "/" + "index-Menu.html",
		"HTMLValue_StartRefreshFiles":   "." + global.APIName_Login + "?state=" + mc_RefreshFiles,
		"HTMLValue_StartZIPFolders":     "." + global.APIName_Login + "?state=" + mc_ZIPFolders,
		"HTMLValue_OutputDownloadLinks": downloadLinks,
		"HTMLValue_OutputReturn":        "檔案掃描時間:" + time.Now().Format("2006-01-02_15:04:05"),
	})

	ginCTX.Abort()
}

type HandlerMinecraftZIPDownloadFolder struct{}

func (feild *HandlerMinecraftZIPDownloadFolder) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {

	count, tags, relativePaths, fullPaths := minecraft.TransDownload.GetFolderList()
	if len(tags) != count || len(relativePaths) != count || len(fullPaths) != count {
		log.Println("更新打包下載清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.String(http.StatusBadGateway, "更新打包下載上傳清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.Abort()
		return
	}

	for index := 0; index < count; index++ {
		ZIPFullPath := global.WebPath_Public + "/" + global.WebPath_Downloads + "/" + minecraft.ZIPFolderName + "/" + tags[index] + ".zip"
		if err := compress.RecursiveZip(fullPaths[index], ZIPFullPath); err != nil {
			errMessage := fmt.Sprintf("壓縮打包下載檔案出現錯誤(%v)(%v) : %v\n", tags[index], fullPaths[index], err)
			log.Println(errMessage)
			ginCTX.String(http.StatusBadGateway, errMessage)
			ginCTX.Abort()
			return
		}
	}

	var buffer bytes.Buffer
	for index := 0; index < count; index++ {
		buffer.WriteString(tags[index] + "\n")
	}
	messageZIPs := buffer.String()
	ginCTX.String(http.StatusOK, fmt.Sprintf("壓縮打包完成! 請回下載頁面進行下載 : \n%v", messageZIPs))
	ginCTX.Abort()
	encrypt.NewAesKey()
	return

}

type HandlerMinecraftDownloadFile struct {
	StateData string
}

func (feild *HandlerMinecraftDownloadFile) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {

	if feild.StateData == "" {
		ginCTX.Abort() //來到這裡應該是惡意使用者,還破解了googleOauth2
		return
	}

	requestPath, err := encrypt.Decode(feild.StateData)
	if err != nil {
		ginCTX.Writer.WriteString(err.Error())
		log.Println(err.Error())
		ginCTX.Abort()
		return
	}

	downloadPath := strings.Split(requestPath, mc_PathChar)
	if len(downloadPath) != 2 {
		ginCTX.Writer.WriteString("解密後字串失效! 疑似有管理員上傳/刪除新的檔案,請刷新頁面獲取新的Token")
		return
	}

	var fullPath string
	fullFilePath, isFileExsit := minecraft.TransDownload.GetFilePath(downloadPath[0], downloadPath[1])
	isFolderExsit := minecraft.TransDownload.IsFolderExsit(downloadPath[0], downloadPath[1])

	if isFileExsit == false && isFolderExsit == false {
		ginCTX.Writer.WriteString("該檔案" + requestPath + " ,不存在於檔案傳輸清單當中!")
		return
	} else {

		if isFolderExsit == true && isFileExsit == false {
			ZIPFullPath := global.WebPath_Public + "/" + global.WebPath_Downloads + "/" + minecraft.ZIPFolderName + "/" + downloadPath[0] + ".zip"
			if common.IsFileExist(ZIPFullPath) == false {
				ginCTX.Writer.WriteString("該檔案" + requestPath + " 為資料夾形式, 請重新至下載頁面進行壓縮打包後再進行下載!")
				return
			} else {
				fullPath = ZIPFullPath
			}
		} else {
			fullPath = fullFilePath
		}

		fileName, _, _, _ := common.FullPathToFileInfo(fullPath)
		ginCTX.Header("Content-Type", "application/octet-stream")              //強制瀏覽器下載
		ginCTX.Header("Content-Disposition", "attachment; filename="+fileName) //瀏覽器下載或預覽
		ginCTX.Header("Content-Disposition", "inline;filename="+fileName)
		ginCTX.Header("Content-Transfer-Encoding", "binary")
		ginCTX.Header("Cache-Control", "no-cache")
		ginCTX.File(fullPath)
		return
	}
}

type HandlerMinecraftDeleteList struct{}

func (_ *HandlerMinecraftDeleteList) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {

	var deleteLinks string
	pathBase := global.Oauth2CallbackName + "?state=" + mc_DeleteFile
	buffer := bytes.NewBuffer([]byte{})

	count, tags, relativePaths, fullPaths := minecraft.TransDownload.GetFileList()
	if len(tags) != count || len(relativePaths) != count || len(fullPaths) != count {
		log.Println("獲取刪除清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.String(http.StatusBadGateway, "獲取刪除清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.Abort()
		return
	}

	deleteFormat := `<input type="button" onclick="let input = prompt('請輸入「刪除」後將會開始執行刪除 %v'); if(input=='刪除'){ location.href='%v' };" value="刪除 %v" /><br>`

	for index := 0; index < count; index++ {
		encode, err := encrypt.Encode(tags[index] + mc_PathChar + relativePaths[index])
		if err != nil {
			log.Println("製作刪除清單加密內容發生錯誤 :", err)
		} else {
			buffer.WriteString(fmt.Sprintf(deleteFormat, tags[index]+" ▶ "+relativePaths[index], pathBase+":"+encode, tags[index]+" ▶ "+relativePaths[index]))
		}
	}

	deleteLinks = buffer.String()
	ginCTX.HTML(http.StatusOK, "minecraft-ShowDelete.html", gin.H{
		global.H5Value_LastUpdate:     global.LastOnLineTime,
		global.H5Value_Title:          "Minecraft - 伺服器操作",
		global.H5Page_Menu:            "./" + global.WebPath_Pages + "/" + "index-Menu.html",
		"HTMLValue_OutputDeleteLinks": deleteLinks,
		"HTMLValue_OutputReturn":      "檔案掃描時間:" + time.Now().Format("2006-01-02_15:04:05"),
	})

	ginCTX.Abort()
}

type HandlerMinecraftDeleteFile struct {
	StateData string
}

func (feild *HandlerMinecraftDeleteFile) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {

	if feild.StateData == "" {
		ginCTX.Abort() //來到這裡應該是惡意使用者,還破解了googleOauth2
		return
	}

	requestPath, err := encrypt.Decode(feild.StateData)
	if err != nil {
		ginCTX.Writer.WriteString(err.Error())
		log.Println(err.Error())
		ginCTX.Abort()
		return
	}

	deletePath := strings.Split(requestPath, mc_PathChar)
	if len(deletePath) != 2 {
		ginCTX.Writer.WriteString("解密後字串失效! 疑似有管理員上傳/刪除新的檔案,請刷新頁面獲取新的Token")
		return
	}

	if fullPath, isExsit := minecraft.TransDelete.GetFilePath(deletePath[0], deletePath[1]); isExsit != true {
		ginCTX.Writer.WriteString("該檔案" + requestPath + " ,不存在於檔案傳輸清單當中!")
		return
	} else {

		if err := common.RemoveAFile(fullPath); err != nil {
			ginCTX.Writer.WriteString(fmt.Sprintf("刪除檔案失敗 %v(%v) : %v", deletePath[0], deletePath[1], err.Error()))
			log.Println(err.Error())
			ginCTX.Abort()
			return
		} else {
			ginCTX.String(http.StatusOK, fmt.Sprintf("檔案刪除完成! : %v(%v)", deletePath[0], deletePath[1]))
			minecraft.CreatFileTransList()
			encrypt.NewAesKey()
			ginCTX.Abort()
		}

		return
	}
}

type HandlerMinecraftRefreshFiles struct{}

func (_ *HandlerMinecraftRefreshFiles) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {
	errList := minecraft.CreatFileTransList()
	if len(errList) <= 0 {
		ginCTX.String(http.StatusOK, fmt.Sprintf("Minecraft檔案清單已刷新!"))
	} else {
		ginCTX.String(http.StatusOK, fmt.Sprintf("Minecraft檔案清單已刷新失敗 : %v", errList))
	}
	encrypt.NewAesKey()
	ginCTX.Abort()
	return
}

type HandlerMinecraftUploadMenu struct{}

func (_ *HandlerMinecraftUploadMenu) Action(ginCTX *gin.Context, account *googleOauth2.GoogleAcc) {

	token, err := encrypt.Encode(fmt.Sprint(time.Now().Unix()) + ":" + account.Email)
	if err != nil {
		ginCTX.HTML(http.StatusOK, "minecraft-ShowUpload.html", gin.H{
			global.H5Value_LastUpdate: global.LastOnLineTime,
			global.H5Value_Title:      "Minecraft - 伺服器操作",
			global.H5Page_Menu:        "./" + global.WebPath_Pages + "/" + "index-Menu.html",
			"HTMLValue_UploadMenu":    err.Error(),
			"HTMLValue_OutputReturn":  "處理時間:" + time.Now().Format("2006-01-02_15:04:05"),
		})
		return
	}

	uploadMenuH5_Top := `
	<form enctype="multipart/form-data" action="/minecraft/upload" method="POST">
	<input type="hidden" name="%v" value="%v">
	<input type="hidden" name="%v" value="%v">` + "\n"
	uploadMenuH5_Top = fmt.Sprintf(uploadMenuH5_Top, global.H5Output_Token, token, global.H5Output_Upload_Type, "allowList")

	uploadLineFile := `%v<input type="file" name="%v" />` + "<div></div>\n"
	uploadLineFolder := `%v<input type="file" name="%v" multiple />` + "<div></div>\n"
	var buffer bytes.Buffer

	count, tags, relativePaths, fullPaths := minecraft.TransUpload.GetFolderList()
	if len(tags) != count || len(relativePaths) != count || len(fullPaths) != count {
		log.Println("獲取上傳清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.String(http.StatusBadGateway, "獲取上傳清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.Abort()
		return
	}

	for index := 0; index < count; index++ {
		buffer.WriteString(fmt.Sprintf(uploadLineFolder, tags[index]+" ◀ "+relativePaths[index], tags[index]+mc_PathChar+relativePaths[index]))
	}
	buffer.WriteString("<br>")

	count, tags, relativePaths, fullPaths = minecraft.TransUpload.GetFileList()
	if len(tags) != count || len(relativePaths) != count || len(fullPaths) != count {
		log.Println("獲取上傳清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.String(http.StatusBadGateway, "獲取上傳清單異常數量", count, len(tags), len(relativePaths), len(fullPaths))
		ginCTX.Abort()
		return
	}
	for index := 0; index < count; index++ {
		buffer.WriteString(fmt.Sprintf(uploadLineFile, tags[index]+" ◀ "+relativePaths[index], tags[index]+mc_PathChar+relativePaths[index]))
	}

	uploadMenuH5_End := `<br>
	<input type="submit" value="開始上傳！！！！" />
	</form>` + "\n"

	uploadMenuAll := uploadMenuH5_Top + buffer.String() + uploadMenuH5_End

	ginCTX.HTML(http.StatusOK, "minecraft-ShowUpload.html", gin.H{
		global.H5Value_LastUpdate: global.LastOnLineTime,
		global.H5Value_Title:      "Minecraft - 伺服器操作",
		global.H5Page_Menu:        "./" + global.WebPath_Pages + "/" + "index-Menu.html",
		"HTMLValue_UploadMenu":    uploadMenuAll,
		"HTMLValue_OutputReturn":  "檔案掃描時間:" + time.Now().Format("2006-01-02_15:04:05"),
	})

	ginCTX.Abort()
}

type FileUpdater interface {
	Action(ginCTX *gin.Context)
}

func NewFileUpdater(uploadType, email string) (FileUpdater, error) {

	switch uploadType {

	case "allowList":
		return &FileUpdaterAllowList{Email: email}, nil
	default:
		return nil, fmt.Errorf("get url query uploadType 設置錯誤！")
	}

}

type FileUpdaterAllowList struct {
	Email string
}

func (thisField *FileUpdaterAllowList) Action(ginCTX *gin.Context) {

	muiltFilesForm, _ := ginCTX.MultipartForm()
	var getFileList []string
	for fromName, muiltFilesHeaders := range muiltFilesForm.File {
		for _, headerCache := range muiltFilesHeaders {

			fromNameSpit := strings.Split(fromName, mc_PathChar)
			if len(fromNameSpit) != 2 {
				log.Printf("解析 %v 的上傳檔案失敗 : %v(%v)(%v)", thisField.Email, headerCache.Filename, fromName, headerCache.Size)
				break
			} else {
				log.Printf("正在解析 %v 的上傳檔案... : %v(%v)(%v)", thisField.Email, headerCache.Filename, fromName, headerCache.Size)
			}

			if fileFullPath, isExsit := minecraft.TransUpload.GetFilePath(fromNameSpit[0], fromNameSpit[1]); isExsit == true {
				getFileList = append(getFileList, headerCache.Filename)
				ginCTX.SaveUploadedFile(headerCache, fileFullPath) //上傳並儲存到目錄
				log.Printf("成功接收 %v 上傳檔案%v", thisField.Email, fileFullPath)
				continue
			}

			if fileFolderPath, isExsit := minecraft.TransUpload.GetFolderPath(fromNameSpit[0], fromNameSpit[1]); isExsit == true {

				getFileList = append(getFileList, headerCache.Filename)
				savePath := fileFolderPath + "/" + headerCache.Filename
				ginCTX.SaveUploadedFile(headerCache, savePath) //上傳並儲存到目錄
				log.Printf("成功接收上傳檔案%v", savePath)
				continue
			}

		}
	}

	var outputReturn bytes.Buffer
	for _, fileName := range getFileList {
		outputReturn.WriteString(fileName + "\n")
	}
	outputReturn.WriteString(fmt.Sprintf("%d files uploaded!", len(getFileList)))
	ginCTX.String(http.StatusOK, outputReturn.String())
	//minecraft.CopyToDownload()
	//minecraft.ZipTheMods()
	minecraft.CreatFileTransList()
	encrypt.NewAesKey()
	return
}
