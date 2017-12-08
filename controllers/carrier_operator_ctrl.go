package controllers

import (
	"cattle-prism/models"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
)

type CarrierOperatorController struct {
	AppController
}

func (this *CarrierOperatorController) Prepare() {

	this.GetUserInfo()
}

func (this *CarrierOperatorController) Get() {

}

func (this *CarrierOperatorController) Put() {

}

func (this *CarrierOperatorController) Post() {

	//验证身份
	if this.UserInfo.CompanyIdNum != 1 {
		this.ServeErrorWithDetail(403, nil, "Forbidden", "Invalid Identity")
	}

	//声明局部变量
	orm.Debug = true
	o := orm.NewOrm()
	var carrierOperator models.BsCarrierOperator

	//验证requestbody合法
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &carrierOperator); err != nil {
		this.ServeError(404, err, "Not Found")
	}

	//验证carrier_operator名称合法
	cnt, err := o.QueryTable("bs_carrier_operator").Filter("carrier_operator_name", carrierOperator.CarrierOperatorName).Count()
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}
	if carrierOperator.CarrierOperatorName == "" || cnt != 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid CarrierOperator Name")
	}
	re := regexp.MustCompile(`^[0-9a-zA-Z\x{4e00}-\x{9fa5}]{2,20}$`)
	if matched := re.MatchString(carrierOperator.CarrierOperatorName); !matched {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid CarrierOperator Name")
	}

	//写入数据库
	if _, err = o.Insert(&carrierOperator); err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//返回值构造
	var cattleResourceData []interface{} = make([]interface{}, 1)
	cattleResourceData[0] = carrierOperator

	cattleResource := models.CattleResource{
		Id:           "",
		Type:         "resource",
		ResourceType: "BsCarrierOperator",
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

func (this *CarrierOperatorController) Delete() {

}
