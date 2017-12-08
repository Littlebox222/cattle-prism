package controllers

// import (
// 	"cattle-prism/models"
// 	"encoding/json"
// 	"github.com/astaxie/beego/httplib"
// )

type BsServiceController struct {
	AppController
}

// func (this *BsServiceController) Prepare() {
// 	this.GetUserInfo()
// }

// func (this *BsServiceController) Get() {

// 	//验证身份
// 	if this.UserInfo.CompanyIdNum == 0 {
// 		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
// 	}

// 	stackRequest := httplib.Get(`http://` + CattlePrismHttpHost + `/v2-beta/projects/1a5/services`)
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

// func (this *BsServiceController) Put() {

// }

// func (this *BsServiceController) Post() {

// }

// func (this *BsServiceController) Delete() {

// }
