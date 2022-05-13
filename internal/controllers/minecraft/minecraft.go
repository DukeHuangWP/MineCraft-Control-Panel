package minecraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"minecraft-control-panel/internal/common"
	"minecraft-control-panel/internal/global"
	"os/exec"
	"strings"
	"sync"
)

func init() { //app啟動後會自動將下載清單複製到網頁/download目錄作暫存

	errList := CreatFileTransList()
	if len(errList) > 0 {
		log.Panicf("Minecraft掃描伺服器資料異常: %v", errList)
	}
	err := removeZIPFiles()
	if err != nil {
		log.Panicf("Minecraft清除ZIP暫存檔案異常: %v", err)
	}
	return
}

const (
	ZIPFolderName = "ZIPs"
)

//對照whitelist.json結構
type Whitelist struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

var TransDownload *fileList

//var TransDownloadZIP *fileList
var TransUpload *fileList
var TransDelete *fileList

type fileList struct {
	//typeMap           map[string]map[string]bool   //[Tag][RelativePath]isFolder
	allTagList        []string
	fileFullPathMap   map[string]map[string]string //[Tag][RelativePath]FileFullPath
	folderFullPathMap map[string]map[string]string //[Tag][RelativePath]FolderFullPath

	allFileList   []fileListInfo
	allFolderList []fileListInfo
}

type fileListInfo struct {
	Tag          string //設定檔標記名稱
	RelativePath string //檔案或資料夾相對路徑
	FullPath     string //檔案或資料夾絕對路徑
}

func (thisList *fileList) IsExsit(tag, relativePath string) bool {
	_, isExsit := thisList.fileFullPathMap[tag][relativePath]
	if isExsit == true {
		return true
	}
	_, isExsit = thisList.folderFullPathMap[tag][relativePath]
	return isExsit
}

func (thisList *fileList) IsFileExsit(tag, relativePath string) bool {
	_, isExsit := thisList.fileFullPathMap[tag][relativePath]
	return isExsit
}

func (thisList *fileList) IsFolderExsit(tag, relativePath string) bool {
	_, isExsit := thisList.folderFullPathMap[tag][relativePath]
	return isExsit
}

func (thisList *fileList) GetFileList() (count int, tags, relativePaths, fullPaths []string) {
	count = len(thisList.allFileList)
	for index := 0; index < count; index++ {
		tags = append(tags, thisList.allFileList[index].Tag)
		relativePaths = append(relativePaths, thisList.allFileList[index].RelativePath)
		fullPaths = append(fullPaths, thisList.allFileList[index].FullPath)
	}
	return
}

func (thisList *fileList) GetFilePath(tag, relativePaths string) (fullPath string, isExsit bool) {
	fullPath, isExsit = thisList.fileFullPathMap[tag][relativePaths]
	return
}

func (thisList *fileList) GetFolderList() (count int, tags, relativePaths, fullPaths []string) {
	count = len(thisList.allFolderList)
	for index := 0; index < count; index++ {
		tags = append(tags, thisList.allFolderList[index].Tag)
		relativePaths = append(relativePaths, thisList.allFolderList[index].RelativePath)
		fullPaths = append(fullPaths, thisList.allFolderList[index].FullPath)
	}
	return
}

func (thisList *fileList) GetFolderPath(tag, relativePaths string) (fullPath string, isExsit bool) {
	fullPath, isExsit = thisList.folderFullPathMap[tag][relativePaths]
	return
}

//addMap["tag"]=["fullPath"],hidenText要隱藏的路徑前墜文字
func (thisList *fileList) addToList(addMap map[string]string, hidenText string) (errList []error) {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock() //上鎖並免thisList存取發生崩潰
	if len(thisList.allFileList) <= 0 {
		thisList.fileFullPathMap = make(map[string]map[string]string)
		thisList.allFileList = []fileListInfo{}
	}

	if len(thisList.allFolderList) <= 0 {
		thisList.folderFullPathMap = make(map[string]map[string]string)
		thisList.allFolderList = []fileListInfo{}
	}

	for tag, fullPath := range addMap {
		if _, isExsit := thisList.fileFullPathMap[tag]; isExsit == false {
			thisList.fileFullPathMap[tag] = make(map[string]string)
		}

		if common.IsFileExist(fullPath) == true {
			relativePath := strings.Replace(fullPath, hidenText, "", 1)
			thisList.fileFullPathMap[tag][relativePath] = fullPath
			thisList.allFileList = append(thisList.allFileList, fileListInfo{Tag: tag, RelativePath: relativePath, FullPath: fullPath})
			continue
		}

		if _, isExsit := thisList.folderFullPathMap[tag]; isExsit == false {
			thisList.folderFullPathMap[tag] = make(map[string]string)
		}

		if common.IsFolderExist(fullPath) == true {
			subFileNameList, err := common.GetFileList(fullPath)
			if err != nil {
				errList = append(errList, err)
				continue
			}

			if len(subFileNameList) > 0 {

				relativePath := strings.Replace(fullPath, hidenText, "", 1)
				thisList.folderFullPathMap[tag][relativePath] = fullPath
				thisList.allFolderList = append(thisList.allFolderList, fileListInfo{Tag: tag, RelativePath: relativePath, FullPath: fullPath})

				for index := 0; index < len(subFileNameList); index++ {
					if _, isExsit := global.Minecraft_TransmissionList_Excludes[subFileNameList[index]]; isExsit {
						continue
					} //忽略檔名清單

					subFullPath := fullPath + "/" + subFileNameList[index]
					relativePath := strings.Replace(subFullPath, hidenText, "", 1)
					thisList.fileFullPathMap[tag][relativePath] = subFullPath
					thisList.allFileList = append(thisList.allFileList, fileListInfo{Tag: tag, RelativePath: relativePath, FullPath: subFullPath})
				}
			}
		}
	}
	return

}

//創建檔案傳輸清單,會檢查檔案是否存在
func CreatFileTransList() (errList []error) {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock() //上鎖並免檔案清單存取發生崩潰
	transFileMap := make(map[string]string)
	for tag, value := range global.Minecraft_TransmissionList_Paths {
		transFileMap[tag] = fmt.Sprint(value)
	}
	TransDownload = &fileList{}
	errList = append(errList, TransDownload.addToList(transFileMap, global.Minecraft_TransmissionList_HideText)...)
	TransUpload = TransDownload
	TransDelete = TransDownload

	// fmt.Println("--------------")
	// fmt.Println(TransDownload.Map)
	// fmt.Println(TransDownload.FileList)
	// fmt.Println(TransDownload.FolderList)
	// fmt.Println("--------------")
	// fmt.Println(TransUpload)
	// fmt.Println("--------------")
	// fmt.Println(transZipMap)
	return
}

func removeZIPFiles() (err error) {

	webFolderPath := global.WebPath_Public + "/" + global.WebPath_Downloads + "/" + ZIPFolderName
	if common.IsFolderExist(webFolderPath) == false {
		if err := common.CreatAFolder(webFolderPath); err != nil {
			return err
		}
	} else {
		if err := common.RemoveContents(webFolderPath); err != nil {
			return err
		}
	}
	return nil
}

func ReadLogFile() (fileText string, err error) {

	bytes, err := ioutil.ReadFile(global.MCSVPath_LastLog)
	return string(bytes), err
}

func RunScript(scriptPath string) (output string, err error) {
	if strings.Contains(scriptPath, ";") || strings.Contains(scriptPath, "\n") || strings.Contains(scriptPath, "\r") {
		return "", fmt.Errorf("不安全的腳本路徑")
	}
	outBytes, err := exec.Command("/bin/sh", scriptPath).Output()
	return string(outBytes), err
}

func GetPlayerAllowList() (allowMap map[string]string, err error) {

	bytes, err := ioutil.ReadFile(global.MCSVPath_PlayerAllowList)
	if err != nil {
		return nil, err
	}

	var cacheList []*Whitelist
	json.Unmarshal(bytes, &cacheList)
	allowMap = make(map[string]string, len(cacheList))
	for _, player := range cacheList {
		allowMap[player.Uuid] = player.Name
	}

	return allowMap, nil
}

func AddPlayerAllowList(playerName, playerUUID string) (err error) {

	allowMap, err := GetPlayerAllowList()
	if err != nil {
		return err
	}
	allowMap[playerUUID] = playerName

	allowSlice := []*Whitelist{}

	for uuid, name := range allowMap {
		allowSlice = append(allowSlice, &Whitelist{Uuid: uuid, Name: name})
	}

	// for _, value := range allowSlice {
	// 	fmt.Printf("'%v'(%v)\n", value.Name,len(value.Name))
	// 	fmt.Printf("'%v'(%v)\n", value.Uuid,len(value.Uuid))
	// }

	bytes, _ := json.MarshalIndent(allowSlice, "", "    ")
	if err != nil {
		return err
	}

	err = common.TruncateFile(global.MCSVPath_PlayerAllowList, bytes)
	if err != nil {
		return err
	}

	return nil
}
