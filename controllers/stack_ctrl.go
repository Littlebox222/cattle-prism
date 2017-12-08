package controllers

// import (
// 	"cattle-prism/models"
// 	"encoding/json"
// 	"github.com/astaxie/beego/httplib"
// 	"log"
// 	"regexp"
// )

type BsStackController struct {
	AppController
}

// func (this *BsStackController) Prepare() {
// 	this.GetUserInfo()
// }

// func (this *BsStackController) Get() {

// 	//验证身份
// 	if this.UserInfo.CompanyIdNum == 0 {
// 		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
// 	}

// 	//请求数据
// 	stackRequest := httplib.Get(`http://` + CattlePrismHttpHost + `/v2-beta/projects/1a5/stacks`)
// 	for headerName, _ := range this.Ctx.Request.Header {
// 		if headerName != "User-Agent" {
// 			stackRequest.Header(headerName, this.Ctx.Input.Header(headerName))
// 		}
// 	}

// 	body, err := stackRequest.Bytes()
// 	if err != nil {
// 		this.ServeError(500, err, "Internal Server Error")
// 	}

// 	//返回值构造
// 	var resource models.CattleResource

// 	err = json.Unmarshal(body, &resource)
// 	if err != nil {
// 		this.ServeError(500, err, "Internal Server Error")
// 	}

// 	this.Data["json"] = resource
// 	this.ServeJSON()
// }

// func (this *BsStackController) Put() {

// 	//验证身份
// 	if this.UserInfo.CompanyIdNum == 0 {
// 		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
// 	}

// 	//验证query合法
// 	querys := this.Ctx.Input.Params()
// 	query, ok := querys[":id"]

// 	if !ok || query == "" {
// 		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Stack Id")
// 	}

// 	//验证requestbody合法
// 	var stackRequestData models.BsStackRequestData
// 	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &stackRequestData); err != nil {
// 		this.ServeError(404, err, "Not Found")
// 	}

// 	log.Printf("---------------------- stackRequestData = %s", stackRequestData)

// 	//验证stackName名称合法
// 	re := regexp.MustCompile(`^[0-9a-zA-Z\x{4e00}-\x{9fa5}]{2,20}$`)
// 	if matched := re.MatchString(stackRequestData.StackName); !matched {
// 		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Stack Name")
// 	}

// 	//请求数据
// 	stackRequest := httplib.Put(`http://` + CattlePrismHttpHost + `/v2-beta/projects/1a5/stacks/` + query)
// 	for headerName, _ := range this.Ctx.Request.Header {
// 		if headerName != "User-Agent" {
// 			stackRequest.Header(headerName, this.Ctx.Input.Header(headerName))
// 		}
// 	}

// 	aaa, _ := json.Marshal(stackRequestData)

// 	ss := make(map[string]string, 1)
// 	ss["name"] = "qwe123"

// 	stackRequest.GetRequest().Body = aaa

// 	stackRequest.Body("{\"name\":\"sssssss\"}")

// 	// stackRequest.Param("name", "wer1wer1")

// 	log.Printf("---------------------- aaa = %s, stackRequest = %s", aaa, stackRequest)

// 	body, err := stackRequest.Bytes()
// 	if err != nil {
// 		this.ServeError(500, err, "Internal Server Error")
// 	}

// 	//返回值构造
// 	var resource models.CattleResource

// 	err = json.Unmarshal(body, &resource)
// 	if err != nil {
// 		this.ServeError(500, err, "Internal Server Error")
// 	}

// 	this.Data["json"] = resource
// 	this.ServeJSON()

// }

// func (this *BsStackController) Post() {

// }

// func (this *BsStackController) Delete() {

// }
