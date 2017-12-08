package controllers

import (
	"cattle-prism/models"
	"encoding/json"
	// "fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	// "time"
)

type BsUserResourceTotalController struct {
	AppController
}

func (this *BsUserResourceTotalController) Prepare() {
	this.GetUserInfo()
}

func (this *BsUserResourceTotalController) Get() {

	//验证身份
	if this.UserInfo.CompanyIdNum == 0 {
		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
	}

	//声明局部变量
	orm.Debug = true
	o := orm.NewOrm()
	var resourceTotal []models.BsUserResourceTotalQueryData

	//查表
	_, err := o.Raw("SELECT `rst`.`id`,`rst`.`company_id`,`rst`.`idc_id`,`rst`.`total`,`rst`.`used`,`rst`.`free`,`rst`.`container_type_id`,`ct`.`container_type_name`,`ct`.`cpu`,`ct`.`memory`,`ct`.`storage`,`idc`.`idc_name`,`idc`.`state`,`idc`.`created`,`idc`.`updated`,`idc`.`carrier_operator_id`,`co`.`carrier_operator_name`,`idc`.`area_id`,`area`.`area_name` FROM `bs_user_resource_total` rst LEFT JOIN `bs_container_type` ct ON `rst`.`container_type_id`=`ct`.`id` LEFT JOIN `bs_idc` idc ON `rst`.`idc_id`=`idc`.`id` LEFT JOIN `bs_carrier_operator` co ON `idc`.`carrier_operator_id`=`co`.`id` LEFT JOIN `bs_area` area ON `idc`.`area_id`=`area`.`id` WHERE `rst`.`company_id` = ?", this.UserInfo.CompanyIdNum).QueryRows(&resourceTotal)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//返回值构造
	var cattleResourceData []interface{} = make([]interface{}, len(resourceTotal))

	for i, resource := range resourceTotal {

		cattleResourceData[i] = models.BsUserResourceTotalResponseData{
			Id:        resource.Id,
			CompanyId: resource.CompanyId,
			ContainerType: models.BsContainerType{
				Id:                resource.ContainerTypeId,
				ContainerTypeName: resource.ContainerTypeName,
				Cpu:               resource.Cpu,
				Memory:            resource.Memory,
				Storage:           resource.Storage,
			},
			Idc: models.BsIdcResponseData{
				IdcId:   resource.IdcId,
				IdcName: resource.IdcName,
				CarrierOperator: models.BsCarrierOperatorResponseData{
					CarrierOperatorId:   resource.CarrierOperatorId,
					CarrierOperatorName: resource.CarrierOperatorName,
				},
				Area: models.BsAreaResponseData{
					AreaId:   resource.AreaId,
					AreaName: resource.AreaName,
				},
				State:   resource.State,
				Created: resource.Created,
				Updated: resource.Updated,
			},
			Total: resource.Total,
			Used:  resource.Used,
			Free:  resource.Free,
		}
	}

	cattleResource := models.CattleResource{
		Id:           "",
		Type:         "resource",
		ResourceType: "BsUserResourceTotal",
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

func (this *BsUserResourceTotalController) Put() {

}

func (this *BsUserResourceTotalController) Post() {

	//验证身份
	if this.UserInfo.CompanyIdNum <= 0 {
		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
	}
	if this.UserInfo.CompanyIdNum != 1 {
		this.ServeErrorWithDetail(403, nil, "Forbidden", "Invalid Identity")
	}

	//声明局部变量
	orm.Debug = true
	o := orm.NewOrm()
	var userResourceTotal models.BsUserResourceTotal

	//验证requestbody合法
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &userResourceTotal); err != nil {
		this.ServeError(404, err, "Not Found")
	}

	//验证CompanyId合法性
	if userResourceTotal.CompanyId <= 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Company Id")
	} else {
		cnt, err := o.QueryTable("company").Filter("id", userResourceTotal.CompanyId).Count()
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
		if cnt == 0 {
			this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Company Id")
		}
	}

	//验证ContainerTypeId合法性
	if userResourceTotal.ContainerTypeId <= 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Container Type Id")
	} else {
		cnt, err := o.QueryTable("bs_container_type").Filter("id", userResourceTotal.ContainerTypeId).Count()
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
		if cnt == 0 {
			this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Container Type Id")
		}
	}

	//验证IdcId合法性
	if userResourceTotal.IdcId <= 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Id")
	} else {
		cnt, err := o.QueryTable("bs_idc").Filter("id", userResourceTotal.IdcId).Count()
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
		if cnt == 0 {
			this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Id")
		}
	}

	//验证total合法性
	if userResourceTotal.Total <= 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Total")
	} else {
		//TODO: 验证资源总量是否超过可分配资源总数
	}

	//写入数据库
	userResourceTotal.Used = 0
	userResourceTotal.Free = userResourceTotal.Total

	if _, err := o.Insert(&userResourceTotal); err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//返回值构造

	var containerType models.BsContainerType
	err := o.QueryTable("bs_container_type").Filter("id", userResourceTotal.ContainerTypeId).One(&containerType)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	var idc models.BsIdc
	err = o.QueryTable("bs_idc").Filter("id", userResourceTotal.IdcId).One(&idc)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	var carrierOperator models.BsCarrierOperator
	err = o.QueryTable("bs_carrier_operator").Filter("id", idc.CarrierOperatorId).One(&carrierOperator)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	var area models.BsArea
	err = o.QueryTable("bs_area").Filter("id", idc.AreaId).One(&area)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	var idcResponseData models.BsIdcResponseData
	idcResponseData = models.BsIdcResponseData{
		IdcId:   idc.Id,
		IdcName: idc.IdcName,
		CarrierOperator: models.BsCarrierOperatorResponseData{
			CarrierOperatorId:   carrierOperator.Id,
			CarrierOperatorName: carrierOperator.CarrierOperatorName,
		},
		Area: models.BsAreaResponseData{
			AreaId:   area.Id,
			AreaName: area.AreaName,
		},
		State:   idc.State,
		Created: idc.Created,
		Updated: idc.Updated,
	}

	var cattleResourceData []interface{} = make([]interface{}, 1)

	cattleResourceData[0] = models.BsUserResourceTotalResponseData{
		Id:            userResourceTotal.Id,
		CompanyId:     userResourceTotal.CompanyId,
		ContainerType: containerType,
		Idc:           idcResponseData,
		Total:         userResourceTotal.Total,
		Used:          userResourceTotal.Used,
		Free:          userResourceTotal.Free,
	}

	cattleResource := models.CattleResource{
		Id:           "",
		Type:         "resource",
		ResourceType: "BsUserResourceTotal",
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

func (this *BsUserResourceTotalController) Delete() {

}
