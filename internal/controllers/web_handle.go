package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"minecraft-control-panel/internal/controllers/encrypt"
	"minecraft-control-panel/internal/controllers/googleOauth2"
	"minecraft-control-panel/internal/global"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

//認證 使用者是否在gooogle email 名單內
func IsAllowEmail(email string) bool {
	for _, allowEmail := range global.Oauth2EmailList {
		if allowEmail == email {
			return true
		}
	}
	return false
}

// func HandleDebug(ginCTX *gin.Context) {

// 	log.Printf("Connected with > %v ,return the index Page\n", ginCTX.Request.RemoteAddr)
// 	ginCTX.Writer.WriteString("<a href=\"http://localhost:8080/minecraft/login?state=test\">http://localhost:8080/minecraft/login?state=test</a><br>")
// 	ginCTX.Writer.WriteString("<a href=\"http://localhost:8080/minecraft/login?state=loginTest\">http://localhost:8080/minecraft/login?state=loginTest</a><br>")
// 	ginCTX.Writer.WriteString("<a href=\"http://localhost:8080/minecraft/login?state=upload\">http://localhost:8080/minecraft/login?state=upload</a><br>")

// 	ginCTX.Writer.WriteString("<br>")
// 	ginCTX.Writer.WriteString("<a href=\"http://localhost:8080/minecraft/login?state=mc-showLogs\">http://localhost:8080/minecraft/login?state=mc-showLogs</a><br>")

// 	return

// }

//轉向首頁選單
func HandleIndexMenu(ginCTX *gin.Context) {

	ginCTX.Redirect(http.StatusTemporaryRedirect, global.APIName_Login[1:]+"?state="+global.APIIndexMenuValue)
	//http://localhost:8080/minecraft/login?state=indexMenu
	ginCTX.Abort() //首頁

}

//API /login 接口
func HandleLogin(ginCTX *gin.Context) {

	ginCTX.Request.ParseForm()
	//fmt.Println(ginCTX.Request.PostForm)
	buffer := bytes.NewBuffer([]byte{})
	for key, value := range ginCTX.Request.PostForm {
		buffer.WriteString(key + "\t" + value[len(value)-1] + "\n")
	} //post from data

	var state string
	if buffer.Len() == 0 {
		buffer = nil //CG優化
		state = ginCTX.Query("state")
	} else {
		cacheStr := buffer.String()
		encode, err := encrypt.Encode(cacheStr[:len(cacheStr)-1])
		if err != nil {
			log.Printf("Post編碼錯誤 : %v (%v)", ginCTX.Request.PostForm, err)
		}
		state = ginCTX.Query("state") + "~" + encode
	}

	//state為googleOauth2准許自定義的query的自訂值
	googleOauthURL := googleOauth2.OauthConfig.AuthCodeURL(state) //URI Query : state
	ginCTX.Redirect(http.StatusTemporaryRedirect, googleOauthURL)
	return

}

//API /callback 接口
func HandleCallback(ginCTX *gin.Context) {
	//http://localhost:8080/start?state=test&code=4%2F0AY0e-g42CDKUTrW7IG_0k4nI6tCxjILz776Z9zizcAcU4x0BTf1RfmbZWfTrbV0D3_UILQ&scope=email+openid+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&authuser=0&prompt=none
	token, err := googleOauth2.OauthConfig.Exchange(oauth2.NoContext, ginCTX.Query("code"))
	if err != nil {
		//fmt.Fprintf(w, "code exchange failed: %s", err.Error())
		state := ginCTX.Query("state")
		//http.Redirect(w, r, APIName_Login+"state="+state, http.StatusTemporaryRedirect) //force login again
		ginCTX.Redirect(http.StatusTemporaryRedirect, "login?state="+state)
		return
	} //將client接收到的callback資訊,回傳給google伺服器作認證,若認證失敗則要求client重新登陸

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		ginCTX.Writer.WriteString(fmt.Sprintf("Failed getting user info: %s", err.Error()))
		return
	} //認證失敗,發生原因情境狀況不明,待確認

	defer response.Body.Close()

	var account googleOauth2.GoogleAcc
	json.NewDecoder(response.Body).Decode(&account) //save the google account to struct
	if err != nil {
		ginCTX.Writer.WriteString(fmt.Sprintf("Failed reading response body: %s", err.Error()))
		return
	} //googleOauth2回傳格式為json

	if !IsAllowEmail(account.Email) {
		ginCTX.Writer.WriteString("Account fail.")
		log.Printf("有個傢伙亂登入伺服器注意一下: %v", account.Email)
		return
	}

	state := ginCTX.Query("state")
	handler, err := NewHandler(state)
	if err != nil {
		ginCTX.Writer.WriteString(fmt.Sprintf("Somthing error : %v)", err))
		return
	} //依據query內的state參數設置,轉發給api_model作實際的業務處理
	handler.Action(ginCTX, &account)
	return
}

//API /upload 接口
func HandleUploadFile(ginCTX *gin.Context) {

	token := ginCTX.PostForm(global.H5Output_Token)
	var errList []error
	decode, err := encrypt.Decode(token)
	if err != nil {
		errList = append(errList, err)
	}

	cache := strings.Split(decode, ":")
	if len(cache) == 2 {
		tokenTime, _ := strconv.ParseInt(cache[0], 10, 64)
		if (time.Now().Unix() - tokenTime) > 300 { //token超過5分鐘
			errList = append(errList, fmt.Errorf("token超過規定時間"))
		} else if !IsAllowEmail(cache[1]) {
			errList = append(errList, fmt.Errorf("Gmail帳號不准許使用！"))
		}
	}

	if len(errList) > 0 {
		log.Println("注意一下token解析異常,留意可能的惡意使用者！")
		ginCTX.Writer.WriteString(fmt.Sprintf("token error : %v)", errList))
		return
	}

	uploadType := ginCTX.PostForm(global.H5Output_Upload_Type)
	fileUpdater, err := NewFileUpdater(uploadType, cache[1])
	if err != nil {
		ginCTX.Writer.WriteString(fmt.Sprintf("Somthing error : %v)", err))
		return
	}

	fileUpdater.Action(ginCTX)
	return

}
