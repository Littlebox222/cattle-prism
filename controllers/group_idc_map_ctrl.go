package controllers

import (
// "cattle-prism/models"
// "encoding/json"
// "fmt"
// "github.com/astaxie/beego/orm"
// _ "github.com/go-sql-driver/mysql"
// "time"
)

type GroupIdcMapController struct {
	AppController
}

func (this *GroupIdcMapController) Prepare() {

	this.GetUserInfo()
}

func (this *GroupIdcMapController) Get() {

}

func (this *GroupIdcMapController) Put() {

}

func (this *GroupIdcMapController) Post() {

	// //验证身份
	// if this.UserInfo.CompanyIdNum != 1 {
	// 	this.ServeErrorWithDetail(403, nil, "Forbidden", "Invalid Identity")
	// }

	// fmt.Printf("验证身份  成功\n")

	// //声明局部变量
	// orm.Debug = true
	// o := orm.NewOrm()
	// var groupIdcMapReqData models.BsGroupIdcMapReqData

	// //验证requestbody合法
	// if err := json.Unmarshal(this.Ctx.Input.RequestBody, &groupIdcMapReqData); err != nil {
	// 	this.ServeError(404, err, "Not Found")
	// }

	// fmt.Printf("验证requestbody合法  成功\n")

	// //验证companyId合法
	// cnt, err := o.QueryTable("company").Filter("id", groupIdcMapReqData.CompanyId).Count()
	// if err != nil {
	// 	this.ServeError(500, err, "Internal Server Error")
	// }
	// if cnt == 0 {
	// 	this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Company Id")
	// }

	// fmt.Printf("验证companyId合法  成功\n")

	// //验证idcs合法

	// //排重
	// for i, id := range groupIdcMapReqData.Idcs {
	// 	if i > 0 && groupIdcMapReqData.Idcs[i-1] == id {
	// 		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Id")
	// 	}
	// }

	// for _, id := range groupIdcMapReqData.Idcs {
	// 	cnt, err := o.QueryTable("bs_idc").Filter("id", id).Count()
	// 	if err != nil {
	// 		this.ServeError(500, err, "Not Found")
	// 	}
	// 	if cnt == 0 {
	// 		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Id")
	// 	}
	// }

	// fmt.Printf("验证idcs合法  成功\n")

	// //写入数据库
	// for _, id := range groupIdcMapReqData.Idcs {

	// 	var groupIdcMap models.BsGroupIdcMap
	// 	groupIdcMap.GroupId = 1
	// 	groupIdcMap.CompanyId = groupIdcMapReqData.CompanyId
	// 	groupIdcMap.IdcId = id
	// 	groupIdcMap.Created = time.Now()

	// 	if _, err = o.Insert(&groupIdcMap); err != nil {
	// 		this.ServeError(500, err, "Internal Server Error")
	// 	}
	// }

	// //返回值
	// groupIdcMapReqData.GroupId = 1
	// this.Data["json"] = groupIdcMapReqData
	// this.ServeJSON()
}

func (this *GroupIdcMapController) Delete() {

}
