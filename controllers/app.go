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
	sqlQueryString = fmt.Sprintf("UPDATE `bs_user_resource_total` SET `occupied` = `occupied` + %d WHERE `company_id` = %d AND `container_type_id` = %s AND `idc_id` IN (", serviceRequestBody.Scale, this.UserInfo.CompanyIdNum, containerTypeId)

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
			var subscribe models.SubscribeResource
			if err := json.Unmarshal(msg, &subscribe); err == nil {

				if subscribe.Data.Resource.CompanyId != "" && subscribe.Data.Resource.Type == "container" && subscribe.Data.Resource.State == "starting" {
					//使用规格 创建 容器时，相应的数据表内容写入及更新

					instanceIdNum := utils.IdStringToIdNumber(subscribe.Data.Resource.Id)
					serviceIdNum := utils.IdStringToIdNumber(subscribe.Data.Resource.ServiceIds[0])
					companyIdNum := utils.IdStringToIdNumber(subscribe.Data.Resource.CompanyId)

					var stackIdNum int64
					var userResourceIdNum int64
					var containerTypeIdNum int64

					orm.Debug = true
					o := orm.NewOrm()

					//获取containerTypeId
					sqlQueryString := fmt.Sprintf("SELECT `data` FROM `service` WHERE `id` = %d", serviceIdNum)
					var serviceDatas []string
					_, err := o.Raw(sqlQueryString).QueryRows(&serviceDatas)
					if err != nil {
						log.Println("[globel socket db err]: \n", err)
						continue
					} else {
						var serviceDataStruct models.ServiceData
						if err := json.Unmarshal([]byte(serviceDatas[0]), &serviceDataStruct); err == nil {
							containerTypeIdNum = utils.IdStringToIdNumber(serviceDataStruct.Fields.Metadata["containerTypeId"])
						} else {
							log.Println("[globel socket db err]: \n", err)
							continue
						}
					}

					//查看bs_user_resource_instance_map里面是否已经写入过数据，说明资源分配后已经被记录

					sqlQueryString = fmt.Sprintf("SELECT `id` FROM `bs_user_resource_instance_map` WHERE `company_id` = %d AND `instance_id` = %d", companyIdNum, instanceIdNum)
					var ids []int64

					num, err := o.Raw(sqlQueryString).QueryRows(&ids)
					if err != nil {
						log.Println("[globel socket db err]: \n", err)
						continue
					}
					if num == 0 {
						//没写入过，则收集所需数据，写入bs_user_resource_instance_map
						sqlQueryString := fmt.Sprintf("SELECT `environment_id` FROM `service` WHERE `id` = %d", serviceIdNum)
						var stackIds []int64

						num, err := o.Raw(sqlQueryString).QueryRows(&stackIds)

						if err != nil || num == 0 {
							log.Println("[globel socket db err]: \n", err)
							continue
						} else {
							stackIdNum = stackIds[0]
						}

						//查出可用资源的Id
						var userResources []models.BsUserResource
						sqlQueryString = fmt.Sprintf("SELECT * FROM `bs_user_resource` WHERE `company_id` = %d AND `container_type_id` = %d AND `can_use` = 1 AND `idc_id` IN (SELECT `idc_id` FROM `bs_idc_host_map` WHERE `host_id` IN (SELECT `host_id` FROM `instance_host_map` WHERE `instance_id` = %d))", companyIdNum, containerTypeIdNum, instanceIdNum)

						if _, err = o.Raw(sqlQueryString).QueryRows(&userResources); err != nil {
							log.Println("[globel socket db err]: \n", err)
							continue
						} else {
							userResourceIdNum = userResources[0].Id
						}

						//写入bs_user_resource_instance_map，更新bs_user_resource和bs_user_resource_total
						o1 := orm.NewOrm()
						sqlQueryString = fmt.Sprintf("INSERT INTO `bs_user_resource_instance_map` (`company_id`,`user_resource_id`,`instance_id`,`service_id`,`stack_id`) VALUES(%d,%d,%d,%d,%d)", companyIdNum, userResourceIdNum, instanceIdNum, serviceIdNum, stackIdNum)
						if _, err = o1.Raw(sqlQueryString).Exec(); err != nil {
							log.Println("[globel socket db err]: \n", err)
							continue
							//TODO: 插入retry队列
						}

						o2 := orm.NewOrm()
						sqlQueryString = fmt.Sprintf("UPDATE `bs_user_resource` SET `can_use` = 0 WHERE `id`= %d", userResourceIdNum)
						if _, err = o2.Raw(sqlQueryString).Exec(); err != nil {
							o1.Rollback()
							log.Println("[globel socket db err]: \n", err)
							continue
							//TODO: 插入retry队列
						}

						o3 := orm.NewOrm()
						sqlQueryString = fmt.Sprintf("UPDATE `bs_user_resource_total` SET `used` = `used` + 1, `free` = `free` - 1 WHERE (`company_id`= %d AND `container_type_id` = %d AND `idc_id` = %d)", companyIdNum, containerTypeIdNum, userResources[0].IdcId)
						if _, err = o3.Raw(sqlQueryString).Exec(); err != nil {
							o2.Rollback()
							o1.Rollback()
							log.Println("[globel socket db err]: \n", err)
							continue
							//TODO: 插入retry队列
						}

						//更新group里stack_count和instance_count信息

						var groupStackCount int = 0
						var groupInstanceCount int = 0

						var groupIds []int64
						sqlQueryString = fmt.Sprintf("SELECT `group_id` FROM `bs_user_group_idc_map` WHERE `idc_id` = %d AND `company_id` = %d", userResources[0].IdcId, companyIdNum)
						if _, err = o.Raw(sqlQueryString).QueryRows(&groupIds); err != nil {
							//this.ServeError(500, err, "Internal Server Error")
							//TODO: error log
						}

						if len(groupIds) != 0 {

							//分别处理instance对应的idc所属于的每个分组（这个产品逻辑很bug）
							for _, groupId := range groupIds {

								//查值
								var userResourceIds []int64
								sqlQueryString = fmt.Sprintf("SELECT `id` FROM `bs_user_resource` WHERE `company_id` = %d AND `can_use` = 0 AND `idc_id` IN (SELECT `idc_id` FROM `bs_user_group_idc_map` WHERE `company_id` = %d AND `group_id` = %d)", companyIdNum, companyIdNum, groupId)
								if _, err = o.Raw(sqlQueryString).QueryRows(&userResourceIds); err != nil {
									log.Println("[globel socket db err]: \n", err)
									continue
								}

								groupInstanceCount = len(userResourceIds)

								if groupInstanceCount != 0 {

									var stackIds []int64
									sqlQueryString = fmt.Sprintf("SELECT `stack_id` FROM `bs_user_resource_instance_map` WHERE `company_id` = %d AND `user_resource_id` IN (SELECT `id` FROM `bs_user_resource` WHERE `company_id` = %d AND `can_use` = 0 AND `idc_id` IN (SELECT `idc_id` FROM `bs_user_group_idc_map` WHERE `company_id` = %d AND `group_id` = %d))", companyIdNum, companyIdNum, companyIdNum, groupId)
									if _, err = o.Raw(sqlQueryString).QueryRows(&stackIds); err != nil {
										log.Println("[globel socket db err]: \n", err)
										continue
									}

									sidMap := make(map[int64]int)
									for _, sid := range stackIds {
										sidMap[sid] = 1
									}

									groupStackCount = len(sidMap)
								}

								//写库
								sqlQueryString = fmt.Sprintf("UPDATE `bs_group` SET `stack_count` = %d, `instance_count` = %d WHERE `id`= %d", groupStackCount, groupInstanceCount, groupId)
								if _, err = o.Raw(sqlQueryString).Exec(); err != nil {
									log.Println("[globel socket db err]: \n", err)
									continue
								}
							}
						}

					} else {
						//写入过则什么都不做
					}

				} else if subscribe.Data.Resource.CompanyId != "" && subscribe.Data.Resource.Type == "container" && subscribe.Data.Resource.State == "removed" {
					//删除 容器时，相应的数据表内容写入及更新
					instanceIdNum := utils.IdStringToIdNumber(subscribe.Data.Resource.Id)
					companyIdNum := utils.IdStringToIdNumber(subscribe.Data.Resource.CompanyId)

					orm.Debug = true
					o := orm.NewOrm()

					//更新bs_user_resource中的can_use
					var userResourceIds []int64
					sqlQueryString := fmt.Sprintf("SELECT `user_resource_id` FROM `bs_user_resource_instance_map` WHERE `instance_id` = %d", instanceIdNum)
					if _, err := o.Raw(sqlQueryString).QueryRows(&userResourceIds); err != nil {
						log.Println("[globel socket db err]: \n", err)
						continue
					}

					//没删过才删
					if len(userResourceIds) != 0 {

						o1 := orm.NewOrm()
						sqlQueryString = fmt.Sprintf("UPDATE `bs_user_resource` SET `can_use` = 1 WHERE `id` = %d", userResourceIds[0])
						if _, err = o1.Raw(sqlQueryString).Exec(); err != nil {
							log.Println("[globel socket db err]: \n", err)
							continue
						}

						//更新bs_user_resource_total中的used和free
						var userResourceTotalIds []int64
						sqlQueryString = fmt.Sprintf("SELECT `urt`.`id` FROM `bs_user_resource_total` urt, `bs_user_resource` ur WHERE `urt`.`company_id` = `ur`.`company_id` AND `urt`.`container_type_id` = `ur`.`container_type_id` AND `urt`.`idc_id` = `ur`.`idc_id` AND `ur`.`id` = %d", userResourceIds[0])
						if _, err := o.Raw(sqlQueryString).QueryRows(&userResourceTotalIds); err != nil {
							log.Println("[globel socket db err]: \n", err)
							continue
						}

						o2 := orm.NewOrm()
						sqlQueryString = fmt.Sprintf("UPDATE `bs_user_resource_total` SET `used`=`used`-1, `free`=`free`+1 WHERE `id` = %d", userResourceTotalIds[0])
						if _, err = o2.Raw(sqlQueryString).Exec(); err != nil {
							o1.Rollback()
							log.Println("[globel socket db err]: \n", err)
							continue
						}

						//删除bs_user_resource_instance_map中的记录
						o3 := orm.NewOrm()
						sqlQueryString = fmt.Sprintf("DELETE FROM `bs_user_resource_instance_map` WHERE `instance_id` = %d", instanceIdNum)
						if _, err = o3.Raw(sqlQueryString).Exec(); err != nil {
							o2.Rollback()
							o1.Rollback()
							log.Println("[globel socket db err]: \n", err)
							continue
						}

						//更新group的stack_count和instance_count
						var userResources []models.BsUserResource
						sqlQueryString = fmt.Sprintf("SELECT * FROM `bs_user_resource` WHERE `id` = %d", userResourceIds[0])

						if _, err = o.Raw(sqlQueryString).QueryRows(&userResources); err != nil {
							log.Println("[globel socket db err]: \n", err)
							continue
						}

						var groupStackCount int = 0
						var groupInstanceCount int = 0

						var groupIds []int64
						sqlQueryString = fmt.Sprintf("SELECT `group_id` FROM `bs_user_group_idc_map` WHERE `idc_id` = %d AND `company_id` = %d", userResources[0].IdcId, companyIdNum)
						if _, err = o.Raw(sqlQueryString).QueryRows(&groupIds); err != nil {
							log.Println("[globel socket db err]: \n", err)
							continue
						}

						if len(groupIds) != 0 {

							//分别处理instance对应的idc所属于的每个分组（这个产品逻辑很bug）
							for _, groupId := range groupIds {

								//查值
								var userResourceIds []int64
								sqlQueryString = fmt.Sprintf("SELECT `id` FROM `bs_user_resource` WHERE `company_id` = %d AND `can_use` = 0 AND `idc_id` IN (SELECT `idc_id` FROM `bs_user_group_idc_map` WHERE `company_id` = %d AND `group_id` = %d)", companyIdNum, companyIdNum, groupId)
								if _, err = o.Raw(sqlQueryString).QueryRows(&userResourceIds); err != nil {
									log.Println("[globel socket db err]: \n", err)
									continue
								}

								groupInstanceCount = len(userResourceIds)

								if groupInstanceCount != 0 {

									var stackIds []int64
									sqlQueryString = fmt.Sprintf("SELECT `stack_id` FROM `bs_user_resource_instance_map` WHERE `company_id` = %d AND `user_resource_id` IN (SELECT `id` FROM `bs_user_resource` WHERE `company_id` = %d AND `can_use` = 0 AND `idc_id` IN (SELECT `idc_id` FROM `bs_user_group_idc_map` WHERE `company_id` = %d AND `group_id` = %d))", companyIdNum, companyIdNum, companyIdNum, groupId)
									if _, err = o.Raw(sqlQueryString).QueryRows(&stackIds); err != nil {
										log.Println("[globel socket db err]: \n", err)
										continue
									}

									sidMap := make(map[int64]int)
									for _, sid := range stackIds {
										sidMap[sid] = 1
									}

									groupStackCount = len(sidMap)
								}

								//写库
								sqlQueryString = fmt.Sprintf("UPDATE `bs_group` SET `stack_count` = %d, `instance_count` = %d WHERE `id`= %d", groupStackCount, groupInstanceCount, groupId)
								if _, err = o.Raw(sqlQueryString).Exec(); err != nil {
									log.Println("[globel socket db err]: \n", err)
									continue
								}
							}
						}
					}
				}
			}
		}
	}
}

func SocketConnect(url string, header http.Header, subscribeMessage chan []byte) {

	connectionMessage := make(chan bool)

	wsClient, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Println("[globel socket connect err]: ", err)
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
					log.Println("[globel socket connect err]: ", err)
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

				log.Println("[globel socket connect err]: ", err)

				connectionMessageChan <- true

				return

			} else {
				log.Println("[globel socket err]: ", err)
				return
			}
		}
		messageChan <- message
	}
}
