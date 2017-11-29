package controllers

import (
	"cattle-prism/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type GroupController struct {
	AppController
}

func (this *GroupController) Prepare() {

	this.GetUserInfo()
}

func (this *GroupController) Get() {

	orm.Debug = true
	o := orm.NewOrm()

	var groups []models.BsGroup

	num, err := o.Raw("SELECT * FROM bs_group WHERE company_id = ?", this.UserInfo.CompanyIdNum).QueryRows(&groups)

	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	fmt.Println(num)
	fmt.Println(groups)

	this.Data["json"] = groups
	this.ServeJSON()

}

func (this *GroupController) Put() {

	var group models.BsGroup
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &group)

	//检验参数有无及合法性
	if err != nil {
		this.ServeError(404, err, "Not Found")
	}

	if group.Name == "" {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Name")
	}

	if group.Id == 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Id")
	}

	orm.Debug = true
	o := orm.NewOrm()

	var groups []models.BsGroup
	num, err := o.Raw("SELECT * FROM bs_group WHERE company_id = ? AND name = ?", this.UserInfo.CompanyIdNum, group.Name).QueryRows(&groups)
	if num > 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Name")
	} else if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//检验权限
	var groupQuery models.BsGroup
	err = o.QueryTable("bs_group").Filter("id", group.Id).One(&groupQuery)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	if groupQuery.CompanyId != this.UserInfo.CompanyIdNum {
		this.ServeError(401, err, "Unauthorized")
	}

	//更新表
	group.Updated = time.Now()
	_, err = o.QueryTable("bs_group").Filter("id", group.Id).Update(orm.Params{
		"name":    group.Name,
		"updated": group.Updated,
	})

	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	groupQuery.Name = group.Name
	groupQuery.Updated = group.Updated

	this.Data["json"] = groupQuery
	this.ServeJSON()

}

func (this *GroupController) Post() {

	var group models.BsGroup
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &group)

	if err != nil {
		this.ServeError(404, err, "Not Found")
	}

	if group.Name == "" {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Name")
	}

	group.Created = time.Now() //.Format("2006-01-02 15:04:05")
	group.CompanyId = this.UserInfo.CompanyIdNum

	orm.Debug = true
	o := orm.NewOrm()

	var groups []models.BsGroup
	num, err := o.Raw("SELECT * FROM bs_group WHERE company_id = ? AND name = ?", this.UserInfo.CompanyIdNum, group.Name).QueryRows(&groups)
	if num > 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Name")
	} else if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	_, err = o.Insert(&group)

	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	this.Data["json"] = group
	this.ServeJSON()
}

func (this *GroupController) Delete() {

	var group models.BsGroup
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &group)

	if err != nil {
		this.ServeError(404, err, "Not Found")
	}

	if group.Id == 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Id")
	}

	orm.Debug = true
	o := orm.NewOrm()

	var groupQuery models.BsGroup
	err = o.QueryTable("bs_group").Filter("Id", group.Id).One(&groupQuery)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	if groupQuery.CompanyId != this.UserInfo.CompanyIdNum {
		this.ServeErrorWithDetail(403, nil, "Forbidden", "Invalid Identity")
	}

	_, err = o.QueryTable("bs_group").Filter("id", group.Id).Delete()
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	this.Data["json"] = groupQuery
	this.ServeJSON()

}
