package controllers

import (
	"cattle-prism/models"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type ContainerTypeController struct {
	AppController
}

func (this *ContainerTypeController) Prepare() {

	this.GetUserInfo()
}

func (this *ContainerTypeController) Get() {

	//验证身份
	if this.UserInfo.CompanyIdNum == 0 {
		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
	}

	orm.Debug = true
	o := orm.NewOrm()

	//验证query合法
	querys := this.Ctx.Input.Params()
	query, ok := querys[":service_id"]

	if !ok || query == "" {

		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Service Id")

	} else {

		//某service所属的container_type

		//读数据库
		var containerTypes []models.BsContainerType

		_, err := o.Raw("SELECT * FROM `bs_container_type` WHERE `id` IN (SELECT `container_type_id` FROM `bs_user_resource` WHERE `id` IN (SELECT `user_resource_id` FROM `bs_user_resource_instance_map` WHERE `company_id` = ? AND `service_id` = ?))", this.UserInfo.CompanyIdNum, query).QueryRows(&containerTypes)

		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}

		//返回值构造
		var cattleResourceData []interface{} = make([]interface{}, len(containerTypes))

		for i, ct := range containerTypes {
			cattleResourceData[i] = ct
		}

		cattleResource := models.CattleResource{
			Id:           "",
			Type:         "resource",
			ResourceType: "BsContainerType",
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
}
