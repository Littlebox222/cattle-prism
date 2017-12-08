package controllers

import (
	"cattle-prism/models"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
)

type AreaController struct {
	AppController
}

func (this *AreaController) Prepare() {

	this.GetUserInfo()
}

func (this *AreaController) Get() {

}

func (this *AreaController) Put() {

}

func (this *AreaController) Post() {

	//验证身份
	if this.UserInfo.CompanyIdNum != 1 {
		this.ServeErrorWithDetail(403, nil, "Forbidden", "Invalid Identity")
	}

	//声明局部变量
	orm.Debug = true
	o := orm.NewOrm()
	var area models.BsArea

	//验证requestbody合法
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &area); err != nil {
		this.ServeError(404, err, "Not Found")
	}

	//验证area名称合法
	cnt, err := o.QueryTable("bs_area").Filter("area_name", area.AreaName).Count()
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}
	if area.AreaName == "" || cnt != 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Area Name")
	}
	re := regexp.MustCompile(`^[0-9a-zA-Z\x{4e00}-\x{9fa5}]{2,20}$`)
	if matched := re.MatchString(area.AreaName); !matched {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Area Name")
	}

	//写入数据库
	if _, err = o.Insert(&area); err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//返回值构造
	var cattleResourceData []interface{} = make([]interface{}, 1)
	cattleResourceData[0] = area

	cattleResource := models.CattleResource{
		Id:           "",
		Type:         "resource",
		ResourceType: "BsArea",
		Links:        nil,
		SortLinks:    nil,
		Actions:      nil,
		CreateTypes:  nil,
		Data:         cattleResourceData,
		Pagination:   models.CattlePagination{},
		Sort:         models.CattleSort{},
		Filters:      nil,
	}

	this.Data["json"] = cattleResource

	this.ServeJSON()
}

func (this *AreaController) Delete() {

}
