package controllers

import (
	"cattle-prism/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"time"
)

type IdcController struct {
	AppController
}

func (this *IdcController) Prepare() {

	this.GetUserInfo()
}

func (this *IdcController) Get() {

}

func (this *IdcController) Put() {

}

func (this *IdcController) Post() {

	//验证身份
	if this.UserInfo.CompanyIdNum != 1 {
		this.ServeErrorWithDetail(403, nil, "Forbidden", "Invalid Identity")
	}

	fmt.Printf("验证身份  成功\n")

	//声明局部变量
	orm.Debug = true
	o := orm.NewOrm()
	var idc models.BsIdc

	//验证requestbody合法
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &idc); err != nil {
		this.ServeError(404, err, "Not Found")
	}

	fmt.Printf("验证requestbody合法  成功\n")

	//验证idc名称合法
	cnt, err := o.QueryTable("bs_idc").Filter("idc_name", idc.IdcName).Count()
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}
	if idc.IdcName == "" || cnt != 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Name")
	}
	re := regexp.MustCompile(`^[0-9a-zA-Z\x{4e00}-\x{9fa5}]{2,20}$`)
	if matched := re.MatchString(idc.IdcName); !matched {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Name")
	}

	fmt.Printf("验证idc名称合法  成功\n")

	//验证idc运营商合法
	if idc.CarrierOperatorId <= 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc CarrierOperator")
	} else {
		cnt, err := o.QueryTable("bs_carrier_operator").Filter("id", idc.CarrierOperatorId).Count()
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
		if cnt == 0 {
			this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc CarrierOperator")
		}
	}

	fmt.Printf("验证idc运营商合法  成功\n")

	//验证idc区域合法
	if idc.AreaId <= 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Area")
	} else {
		cnt, err := o.QueryTable("bs_area").Filter("id", idc.AreaId).Count()
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
		if cnt == 0 {
			this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Area")
		}
	}

	fmt.Printf("验证idc区域合法  成功\n")

	//验证idc状态描述合法
	switch idc.State {
	case "":
		idc.State = "normal"
		break
	case "normal":
		break
	case "warning":
		break
	case "error":
		break
	default:
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc State Description")
	}

	fmt.Printf("验证idc状态描述合法  成功\n")

	//写入数据库
	idc.Created = time.Now()

	if _, err = o.Insert(&idc); err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//返回值构造
	var cattleResourceData []interface{} = make([]interface{}, 1)
	cattleResourceData[0] = idc

	cattleResource := models.CattleResource{
		Id:           "",
		Type:         "resource",
		ResourceType: "BsIdc",
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

func (this *IdcController) Delete() {

}
