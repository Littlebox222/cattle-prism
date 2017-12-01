package controllers

import (
	"cattle-prism/models"
	// "encoding/json"
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
	_, err := o.Raw("SELECT `rst`.`id`,`rst`.`company_id`,`rst`.`idc_id`,`rst`.`total`,`rst`.`used`,`rst`.`free`,`rst`.`container_type_id`,`ct`.`name`,`ct`.`cpu`,`ct`.`memory`,`ct`.`storage` FROM `bs_user_resource_total` rst LEFT JOIN `bs_container_type` ct ON `rst`.`container_type_id`=`ct`.`id` WHERE `rst`.`company_id` = ?", this.UserInfo.CompanyIdNum).QueryRows(&resourceTotal)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	var cattleResourceData []interface{} = make([]interface{}, len(resourceTotal))

	for i, resource := range resourceTotal {

		cattleResourceData[i] = models.BsUserResourceTotalResponseData{
			Id:        resource.Id,
			CompanyId: resource.CompanyId,
			ContainerType: models.BsContainerType{
				Id:      resource.ContainerTypeId,
				Name:    resource.Name,
				Cpu:     resource.Cpu,
				Memory:  resource.Memory,
				Storage: resource.Storage,
			},
			IdcId: resource.IdcId,
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

}

func (this *BsUserResourceTotalController) Delete() {

}
