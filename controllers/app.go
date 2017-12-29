package controllers

import (
	"cattle-prism/models"

	"cattle-prism/utils"
	"cattle-prism/utils/wsutil"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"
	// "net/http"
	"strings"

	"cattle-prism/dao"
	"encoding/base64"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	// "strconv"
)

// Operations about Users
type AppController struct {
	beego.Controller
	Token    string
	UserInfo models.TokenDataItemUserIdentity
}

var RancherEndpointHost string
var UserInfoCache, _ = cache.NewCache("memory", `{"interval":60}`)
var RancherAdminApiKey string
var RancherAdminSecretKey string

func init() {
	if RancherEndpointHost = beego.AppConfig.String("RancherEndpointHost"); RancherEndpointHost == "" {
		RancherEndpointHost = "127.0.0.1:8080"
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)

	dbConfig := dao.InitConfig()
	dbConfigString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.DbName)
	orm.RegisterDataBase("default", "mysql", dbConfigString)

	orm.RegisterModel(new(models.BsGroup))
	orm.RegisterModel(new(models.BsIdc))
	orm.RegisterModel(new(models.BsCarrierOperator))
	orm.RegisterModel(new(models.BsArea))
	orm.RegisterModel(new(models.BsUserGroupIdcMap))
	orm.RegisterModel(new(models.BsUserResourceInstanceMap))
	// orm.RegisterModel(new(models.BsIdcHostMap))
	orm.RegisterModel(new(models.Company))
	orm.RegisterModel(new(models.BsUserResourceTotal))
	orm.RegisterModel(new(models.BsContainerType))

	orm.RegisterModel(new(models.Environment))
	orm.RegisterModel(new(models.Instance))
	orm.RegisterModel(new(models.Service))

	go CreateGlobleSocket()
}

func (this *AppController) ServeErrorWithDetail(status int, err error, message string, detail string) {

	if err == nil || detail != "" {

		cattleErr := &models.CattleError{
			Type:     "error",
			BaseType: "error",
			Status:   status,
			Code:     message,
			Message:  message,
			Detail:   detail,
		}
		logs.Error(cattleErr)
		this.Data["json"] = cattleErr

	} else {
		logs.Error(err)

		cattleErr := &models.CattleError{
			Type:     "error",
			BaseType: "error",
			Status:   status,
			Code:     message,
			Message:  message,
			Detail:   err.Error(),
		}

		this.Data["json"] = cattleErr
	}

	this.Ctx.Output.SetStatus(status)
	this.ServeJSON()
	this.StopRun()
}

func (this *AppController) ServeError(status int, err error, message string) {

	if err == nil {

		cattleErr := &models.CattleError{
			Type:     "error",
			BaseType: "error",
			Status:   status,
			Code:     message,
			Message:  message,
			Detail:   "",
		}
		logs.Error(cattleErr)
		this.Data["json"] = cattleErr

		this.Ctx.Output.SetStatus(status)
		this.ServeJSON()
		this.StopRun()

	} else {
		this.ServeErrorWithDetail(status, err, message, "")
	}
}

func (this *AppController) GetUserInfo() {
	this.Token = this.Ctx.Input.Cookie("token")
	if this.Token != "" {
		var tokenData models.TokenResource
		cacheKey := fmt.Sprintf("userinfo_%s", this.Token)
		cacheData := UserInfoCache.Get(cacheKey)
		if cacheData == nil {
			tokenRequest := httplib.Get(`http://` + RancherEndpointHost + `/v2-beta/token`)
			for headerName, _ := range this.Ctx.Request.Header {
				if headerName != "User-Agent" {
					tokenRequest.Header(headerName, this.Ctx.Input.Header(headerName))
				}
			}
			body, err := tokenRequest.Bytes()
			if err != nil {
				this.ServeError(500, err, "Internal Server Error")
			}
			err = json.Unmarshal(body, &tokenData)
			if err != nil {
				this.ServeError(500, err, "Internal Server Error")
			}
			if len(tokenData.Data) > 0 {
				this.UserInfo = tokenData.Data[0].UserIdentity
				this.UserInfo.Cached = false
				if this.UserInfo.Id != "" {
					if err := UserInfoCache.Put(cacheKey, body, 60*time.Second); err != nil {
						this.ServeError(500, err, "Internal Server Error")
					}
				}
				/*
					this.Data["json"] = this.UserInfo
					this.ServeJSON()
					this.StopRun()
				*/
			}
		} else {
			if cacheBytes, ok := cacheData.([]byte); ok {
				err := json.Unmarshal(cacheBytes, &tokenData)
				if err != nil {
					UserInfoCache.Delete(cacheKey)
					this.ServeError(500, err, "Internal Server Error")
				}
				if len(tokenData.Data) > 0 {
					this.UserInfo = tokenData.Data[0].UserIdentity
					this.UserInfo.Cached = true
					/*
						this.Data["json"] = this.UserInfo
						this.ServeJSON()
						this.StopRun()
					*/
				}
			} else {
				UserInfoCache.Delete(cacheKey)
				this.ServeError(500, errors.New("Decoding Cache Data Error"), "Internal Server Error")
			}
		}
	}
}

func (this *AppController) Prepare() {
	this.GetUserInfo()

	//websocket过滤
	if wsutil.IsWebSocketRequest(this.Ctx.Request) && this.Ctx.Input.IsGet() {

		re := regexp.MustCompile(`^/v2-beta/projects/[a-z0-9]+/subscribe$`)
		if matched := re.MatchString(this.Ctx.Input.URL()); matched {
			// fmt.Println(this.Ctx.Input.URL())
			this.Subscribe()
			return
		}
	}

	//put、delete方法加身份验证
	if this.Ctx.Input.IsPut() || this.Ctx.Input.IsDelete() {
		re := regexp.MustCompile(`^/v2-beta/projects/[a-z0-9]+/(stacks|services|instances)/[a-z0-9]+$`)
		if matched := re.MatchString(this.Ctx.Input.URL()); matched {
			if ok := this.IdendityCheck(); !ok {
				return
			}
		}
	}

	//资源上限判断
	if this.Ctx.Input.IsPost() {
		re := regexp.MustCompile(`^/v2-beta/projects/[a-z0-9]+/service`)
		if matched := re.MatchString(this.Ctx.Input.URL()); matched {
			if ok := this.ResourceLimitCheck(); !ok {
				return
			}
		}
	}

	this.Proxy()
}

func (this *AppController) Proxy() {
	if this.UserInfo.CompanyId != "" {
		this.Ctx.Request.Header.Set("X-Company-Id", this.UserInfo.CompanyId)
		this.Ctx.ResponseWriter.Header().Set("X-Company-Id", this.UserInfo.CompanyId)
	}

	/*
		if this.UserInfo.Cached {
			this.Ctx.ResponseWriter.Header().Set("X-Userinfo-Cached", "True")
		} else {
			this.Ctx.ResponseWriter.Header().Set("X-Userinfo-Cached", "False")
		}
	*/

	// this.Ctx.Request.SetBasicAuth("E37CE8E5038A794B25FC", "wQCfoBTrT8PBFUobT4oQJRHEDaBkuv3wwZi4fVCb")
	if wsutil.IsWebSocketRequest(this.Ctx.Request) {
		remoteWs := &url.URL{
			Scheme: "ws://",
			Host:   RancherEndpointHost,
		}
		proxyWs := wsutil.NewSingleHostReverseProxy(remoteWs)
		proxyWs.ServeHTTP(this.Ctx.ResponseWriter, this.Ctx.Request)
	} else {
		remoteHttp, err := url.Parse(`http://` + RancherEndpointHost)
		if err != nil {
			// panic(err)
			logs.Error(err)
		}
		proxyHttp := httputil.NewSingleHostReverseProxy(remoteHttp)
		proxyHttp.ServeHTTP(this.Ctx.ResponseWriter, this.Ctx.Request)
	}
}

func (this *AppController) Subscribe() {
	// fmt.Println("Subscribe")
	// remoteWs := &url.URL{
	// 	Scheme: "ws://",
	// 	Host: RancherEndpointHost,
	// }
	// proxyWs := wsutil.NewSingleHostReverseProxy(remoteWs)
	// proxyWs.ServeHTTP(this.Ctx.ResponseWriter, this.Ctx.Request)

	c, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 4096, 4096)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	subscribeMessage := make(chan []byte)

	filterUrl := `ws://` + RancherEndpointHost + this.Ctx.Input.URI()
	wsHeader := this.Ctx.Request.Header
	wsHeader.Del("Sec-Websocket-Version")
	wsHeader.Del("Connection")
	wsHeader.Del("Sec-Websocket-Key")
	wsHeader.Del("Sec-Websocket-Extensions")
	wsHeader.Del("Upgrade")
	wsClient, _, err := websocket.DefaultDialer.Dial(filterUrl, wsHeader)
	if err != nil {
		log.Println("read:", err)
		return
	}

	defer c.Close()
	defer wsClient.Close()

	go func() {
		for {
			_, message, err := wsClient.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			// 过滤
			if true {
				var subscribe models.SubscribeResource
				if err := json.Unmarshal(message, &subscribe); err == nil {
					// if subscribe.Name == "ping" || subscribe.ResourceType == "containerEvent" || subscribe.ResourceType == "host" {
					if subscribe.Name == "ping" {
						subscribeMessage <- message
						log.Printf("recv ping --------:\n %s", message)
					} else if subscribe.Data.Resource.CompanyId == this.UserInfo.CompanyId {

						// 过滤后塞ws消息
						stringMsg := string(message)
						stringMsg = strings.Replace(stringMsg, RancherEndpointHost, fmt.Sprintf("%s:%d", this.Ctx.Input.Host(), this.Ctx.Input.Port()), -1)
						subscribeMessage <- []byte(stringMsg)
						log.Printf("recv change --------:\n %s", message)

					} else {
						log.Printf("interception change --------:\n %s", message)
					}
				} else {
					log.Println("read json:", err)
				}
			} else {
				subscribeMessage <- message
				log.Printf("recv --------:\n %s", message)
			}
		}

	}()

	for {
		select {
		case msg := <-subscribeMessage:
			err = c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
		// mt, message, err := c.ReadMessage()
		// if err != nil {
		// 	log.Println("read:", err)
		// 	break
		// }
		// log.Printf("recv: %s", message)
		// time.Sleep(time.Second)
	}
}

func (this *AppController) IdendityCheck() bool {
	//取id和type
	params := this.Ctx.Input.Params()
	itemId := params["3"]
	itemIdNum := utils.IdStringToIdNumber(itemId)

	orm.Debug = true
	o := orm.NewOrm()

	var itemType string
	switch params["2"] {
	case "stacks":
		itemType = "environment"
		break
	case "services":
		itemType = "service"
		break
	case "instances":
		itemType = "instance"
		break
	default:
		itemType = ""
	}

	//对应type的数据库里取company_id

	sqlQueryString := fmt.Sprintf("SELECT `company_id` FROM %s WHERE id = %d", itemType, itemIdNum)
	var companyIds []int64
	if _, err := o.Raw(sqlQueryString).QueryRows(&companyIds); err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	if companyIds[0] != this.UserInfo.CompanyIdNum {
		this.ServeErrorWithDetail(403, nil, "Forbidden", "Invalid Identity")
		return false
	} else {
		return true
	}
}

func (this *AppController) ResourceLimitCheck() bool {

	var serviceRequestBody models.ServiceRequestBody

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &serviceRequestBody); err != nil {
		this.ServeError(404, err, "Not Found")
		return false
	}

	containerTypeId := serviceRequestBody.Metadata["containerTypeId"]
	if containerTypeId == "" {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid ContainerTypeId")
		return false
	}

	idcIds := serviceRequestBody.IdcIds
	if len(idcIds) == 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid IdcIds")
		return false
	}

	//判断逻辑：scale < (total - occupied) 则ok
	orm.Debug = true
	o := orm.NewOrm()

	sqlQueryString := fmt.Sprintf("SELECT `total`-`occupied` FROM `bs_user_resource_total` WHERE `company_id` = %d AND `container_type_id` = %s AND `idc_id` IN (", this.UserInfo.CompanyIdNum, containerTypeId)

	for i, id := range idcIds {
		sqlQueryString = fmt.Sprintf("%s%d", sqlQueryString, id)

		if i != len(idcIds)-1 {
			sqlQueryString = fmt.Sprintf("%s,", sqlQueryString)
		} else {
			sqlQueryString = fmt.Sprintf("%s)", sqlQueryString)
		}
	}

	var availables []int

	if _, err := o.Raw(sqlQueryString).QueryRows(&availables); err != nil {
		this.ServeError(500, err, "Internal Server Error")
		return false
	}

	if len(availables) == 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Available User Resource Not Found")
	}

	for _, a := range availables {
		if a == 0 || a < serviceRequestBody.Scale {
			this.ServeErrorWithDetail(404, nil, "Not Found", "Available User Resource Not Found")
			return false
		}
	}

	//给occupid赋值占位
	sqlQueryString = fmt.Sprintf("UPDATE `bs_user_resource_total` SET `occupied` = `occupied` + %d WHERE (`company_id`= %d AND `container_type_id` = %s AND `idc_id` IN (", serviceRequestBody.Scale, this.UserInfo.CompanyIdNum, containerTypeId)

	for i, id := range idcIds {
		sqlQueryString = fmt.Sprintf("%s%d", sqlQueryString, id)

		if i != len(idcIds)-1 {
			sqlQueryString = fmt.Sprintf("%s,", sqlQueryString)
		} else {
			sqlQueryString = fmt.Sprintf("%s)", sqlQueryString)
		}
	}

	if _, err := o.Raw(sqlQueryString).Exec(); err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	return true
}

func (this *AppController) Finish() {

}

func CreateGlobleSocket() {

	RancherAdminApiKey = beego.AppConfig.String("RancherAdminApiKey")
	RancherAdminSecretKey = beego.AppConfig.String("RancherAdminSecretKey")
	url := `ws://` + RancherEndpointHost + `/v2-beta/projects/1a5/subscribe?eventNames=resource.change&limit=-1&sockId=1`
	header := http.Header{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(RancherAdminApiKey+":"+RancherAdminSecretKey))}}

	subscribeMessage := make(chan []byte)

	go SocketConnect(url, header, subscribeMessage)

	for {
		select {
		case msg := <-subscribeMessage:
			// log.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~ \n %s", msg)

			var subscribe models.SubscribeResource
			if err := json.Unmarshal(msg, &subscribe); err == nil {

				if subscribe.Data.Resource.CompanyId != "" && subscribe.Data.Resource.Type == "container" && subscribe.Data.Resource.State == "starting" {
					//使用规格 创建 容器时，相应的数据表内容写入及更新

				} else if subscribe.Data.Resource.CompanyId != "" && subscribe.Data.Resource.Type == "container" && subscribe.Data.Resource.State == "removed" {
					//删除 容器时，相应的数据表内容写入及更新

				}

			}
		}
	}
}

func SocketConnect(url string, header http.Header, subscribeMessage chan []byte) {

	connectionMessage := make(chan bool)

	wsClient, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Println("globel socket connect err: ", err)
	}

	defer wsClient.Close()

	go ReceiveSocketMessage(subscribeMessage, wsClient, connectionMessage)

	//监听socket连接断的消息
	for {
		select {
		case shouldConnect := <-connectionMessage:

			if shouldConnect {

				wsClient, _, err := websocket.DefaultDialer.Dial(url, header)
				if err != nil {
					log.Println("globel socket connect err: ", err)
					continue
				}

				defer wsClient.Close()

				go ReceiveSocketMessage(subscribeMessage, wsClient, connectionMessage)
			}
		}
	}
}

func ReceiveSocketMessage(messageChan chan []byte, wsClient *websocket.Conn, connectionMessageChan chan bool) {

	for {
		_, message, err := wsClient.ReadMessage()

		if err != nil {
			//错误码详见：https://github.com/gorilla/websocket/blob/23059f29570f0e13fca80ef6aea0f04c11daaa4d/conn.go
			if websocket.IsUnexpectedCloseError(err, 1000, 1001, 1002, 1003, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1015) ||
				websocket.IsCloseError(err, 1000, 1001, 1002, 1003, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1015) {

				log.Printf("error: %v", err)

				connectionMessageChan <- true

				return

			} else {
				log.Println("read:", err)
				return
			}
		}
		messageChan <- message
	}
}
