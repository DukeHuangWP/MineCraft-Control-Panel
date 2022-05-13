package config

import (
	"log"
	"minecraft-control-panel/internal/common/customVar"
)

func init() {
	CreatSettingMap(false)
	SettingMap["cpu_core"] = &AddSetings{DefaultValue: 0, OutputPointer: &Cfgs.CpuCore, Custom: &customVar.Uint8Type{}}
	//SettingMap["XXX"] = &AddSetings{初始值, 輸出變數指標, 自訂設定值類型}
	return
}

//初始化環境變數設定值
func InitLoad(configPath string) {
	Cfgs.ReadFile(configPath, true)
	return
}

//重新讀取環境變數設定值
func Reload() {
	CreatSettingMap(false)
	Cfgs.ReadFile("", false)
	log.Printf("重新讀取環境設定值完畢 (%v)", Cfgs.ConfigPath)
	return
}
