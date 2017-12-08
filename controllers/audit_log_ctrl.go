package controllers

// import (
// 	"cattle-prism/models"
// 	"encoding/json"
// 	"github.com/astaxie/beego/httplib"
// 	"log"
// )

type BsAuditLogController struct {
	AppController
}

// func (this *BsAuditLogController) Prepare() {
// 	this.GetUserInfo()
// }

// func (this *BsAuditLogController) Get() {

// 	//验证身份
// 	if this.UserInfo.CompanyIdNum == 0 {
// 		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
// 	}

// 	//声明局部变量
// 	var resource models.CattleResource

// 	stackRequest := httplib.Get(`http://` + RancherEndpointHost + `/v2-beta/projects/1a5/auditlogs?companyId=` + this.UserInfo.CompanyId)
// 	for headerName, _ := range this.Ctx.Request.Header {
// 		if headerName != "User-Agent" {
// 			stackRequest.Header(headerName, this.Ctx.Input.Header(headerName))
// 		}
// 	}

// 	log.Println("header:", this.Ctx.Request.Header)

// 	body, err := stackRequest.Bytes()
// 	if err != nil {
// 		this.ServeError(500, err, "Internal Server Error")
// 	}
// 	err = json.Unmarshal(body, &resource)
// 	if err != nil {
// 		this.ServeError(500, err, "Internal Server Error")
// 	}

// 	this.Data["json"] = resource
// 	this.ServeJSON()
// }

// func (this *BsAuditLogController) Put() {

// }

// func (this *BsAuditLogController) Post() {

// }

// func (this *BsAuditLogController) Delete() {

// }
