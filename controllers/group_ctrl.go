package controllers

import (
	"cattle-prism/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	// "log"
	"regexp"
	"strconv"
	"time"
)

type GroupController struct {
	AppController
}

func (this *GroupController) Prepare() {

	this.GetUserInfo()
}

func (this *GroupController) Get() {

	//验证身份
	if this.UserInfo.CompanyIdNum == 0 {
		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
	}

	orm.Debug = true
	o := orm.NewOrm()

	//验证query合法
	querys := this.Ctx.Input.Params()
	query, ok := querys[":group_id"]

	if !ok || query == "" {

		query, ok := querys[":service_id"]

		if !ok || query == "" {

			//group列表

			//读数据库
			var groups []models.BsGroup
			if _, err := o.Raw("SELECT * FROM bs_group WHERE company_id = ?", this.UserInfo.CompanyIdNum).QueryRows(&groups); err != nil {
				this.ServeError(500, err, "Internal Server Error")
			}

			//返回值构造
			var cattleResourceData []interface{} = make([]interface{}, len(groups))

			for i, group := range groups {
				cattleResourceData[i] = group
			}

			cattleResource := models.CattleResource{
				Id:           "",
				Type:         "resource",
				ResourceType: "BsGroup",
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

		} else {

			//某service所属的group

			//读数据库
			var groups []models.BsGroup

			_, err := o.Raw("SELECT * FROM `bs_group` WHERE `id` IN (SELECT `group_id` FROM `bs_user_group_idc_map` WHERE `company_id` = ? AND `idc_id` = (SELECT `idc_id` FROM `bs_user_resource` WHERE `id` = (SELECT `user_resource_id` FROM `bs_user_resource_instance_map` WHERE `service_id` = ?)))", this.UserInfo.CompanyIdNum, query).QueryRows(&groups)

			if err != nil {
				this.ServeError(500, err, "Internal Server Error")
			}

			//返回值构造
			var cattleResourceData []interface{} = make([]interface{}, len(groups))

			for i, group := range groups {
				cattleResourceData[i] = group
			}

			cattleResource := models.CattleResource{
				Id:           "",
				Type:         "resource",
				ResourceType: "BsGroup",
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

	} else {

		//某个group的详情

		//验证资源权限
		var groups []models.BsGroup
		num, err := o.Raw("SELECT * FROM bs_group WHERE id = ? AND company_id = ?", query, this.UserInfo.CompanyIdNum).QueryRows(&groups)
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
		if num == 0 {
			this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Id")
		}

		//读数据库
		var idcs []models.BsIdcResponseData
		_, err = o.Raw("SELECT `m`.`idc_id`,`idc`.`idc_name`,`idc`.`state`,`idc`.`created`,`idc`.`updated`,`idc`.`carrier_operator_id`,`co`.`carrier_operator_name`,`idc`.`area_id`,`area`.`area_name` FROM `bs_user_group_idc_map` m LEFT JOIN `bs_idc` idc ON `m`.`idc_id`=`idc`.`id` LEFT JOIN `bs_carrier_operator` co ON `idc`.`carrier_operator_id`=`co`.`id` LEFT JOIN `bs_area` area ON `idc`.`area_id`=`area`.`id` WHERE (group_id = ? AND company_id = ?)", query, this.UserInfo.CompanyIdNum).QueryRows(&idcs)
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}

		var groupUserResource []models.BsGroupUserResource

		_, err = o.Raw("SELECT `rst`.`idc_id`,`ct`.`id`,`ct`.`container_type_name`,`ct`.`cpu`,`ct`.`memory`,`ct`.`storage`,`rst`.`total`,`rst`.`used`,`rst`.`free` FROM `bs_user_resource_total` rst LEFT JOIN `bs_container_type` ct ON `rst`.`container_type_id`=`ct`.`id` WHERE `company_id` = ? AND `idc_id` IN (SELECT `m`.`idc_id` FROM `bs_user_group_idc_map` m WHERE (group_id = ? AND company_id = ?))", this.UserInfo.CompanyIdNum, query, this.UserInfo.CompanyIdNum).QueryRows(&groupUserResource)
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}

		//返回值构造
		var cattleResourceData []interface{} = make([]interface{}, 1)

		bsGroupResponseDataDetail := models.BsGroupResponseDataDetail{
			Group:             groups[0],
			Idcs:              idcs,
			GroupUserResource: groupUserResource,
		}

		cattleResourceData[0] = bsGroupResponseDataDetail

		cattleResource := models.CattleResource{
			Id:           "",
			Type:         "resource",
			ResourceType: "BsGroup",
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

func (this *GroupController) Put() {

	//验证身份
	if this.UserInfo.CompanyIdNum == 0 {
		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
	}

	//声明局部变量
	orm.Debug = true
	o := orm.NewOrm()
	var groupRequestData models.BsGroupRequestData

	//验证requestbody合法
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &groupRequestData)
	if err != nil {
		this.ServeError(404, err, "Not Found")
	}

	groupRequestData.CompanyId = this.UserInfo.CompanyIdNum

	//验证GroupId合法
	var groups []models.BsGroup
	num, err := o.Raw("SELECT * FROM bs_group WHERE company_id = ? AND id = ?", this.UserInfo.CompanyIdNum, groupRequestData.GroupId).QueryRows(&groups)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}
	if num == 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Id")
	}

	//验证GroupName合法
	group := groups[0]
	if group.GroupName != groupRequestData.GroupName {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Name")
	}

	//验证Idcs合法
	if len(groupRequestData.IdcIds) == 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Ids")
	}

	//验证每个idc都有用户购买过的资源
	var idcIdsFromDB []int64
	_, err = o.Raw("SELECT idc_id FROM bs_user_resource_total WHERE company_id = ?", this.UserInfo.CompanyIdNum).QueryRows(&idcIdsFromDB)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	for _, idcId := range groupRequestData.IdcIds {
		for j, id := range idcIdsFromDB {
			if idcId == id {
				break
			}

			if j == len(idcIdsFromDB)-1 {
				//该资源未授权
				this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Idc Ids")
			}
		}
	}

	//区分哪些需要新插入user_group_idc_map，哪些需要从该表中移除，哪些不处理
	var idcIdsFromMap []int64
	num, err = o.Raw("SELECT idc_id FROM bs_user_group_idc_map WHERE company_id = ? AND group_id = ?", this.UserInfo.CompanyIdNum, groupRequestData.GroupId).QueryRows(&idcIdsFromMap)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	toBeInsertIds := []int64{0}
	toBeDeleteIds := []int64{0}

	if num == 0 {
		toBeInsertIds = groupRequestData.IdcIds
	} else {
		for _, inputId := range groupRequestData.IdcIds {
			for j, existId := range idcIdsFromMap {
				if inputId == existId {
					break
				}

				if j == len(idcIdsFromMap)-1 {
					//新增的inputId
					toBeInsertIds = append(toBeInsertIds, inputId)
				}
			}
		}

		for _, existId := range idcIdsFromMap {
			for j, inputId := range groupRequestData.IdcIds {
				if existId == inputId {
					break
				}

				if j == len(groupRequestData.IdcIds)-1 {
					//要删除的existId
					toBeDeleteIds = append(toBeDeleteIds, existId)
				}
			}
		}

		toBeInsertIds = toBeInsertIds[1:]
		toBeDeleteIds = toBeDeleteIds[1:]
	}

	//写入表
	if len(toBeInsertIds) > 0 {

		var userGroupIdcMaps []models.BsUserGroupIdcMap = make([]models.BsUserGroupIdcMap, len(toBeInsertIds))

		for i, id := range toBeInsertIds {
			userGroupIdcMaps[i] = models.BsUserGroupIdcMap{
				CompanyId: groupRequestData.CompanyId,
				GroupId:   groupRequestData.GroupId,
				IdcId:     id,
			}
		}

		_, err = o.InsertMulti(len(userGroupIdcMaps), userGroupIdcMaps)
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
	}

	//删除行
	if len(toBeDeleteIds) > 0 {

		qs := fmt.Sprintf("DELETE FROM bs_user_group_idc_map WHERE company_id = ? AND group_id = ? AND (")

		for i, id := range toBeDeleteIds {
			qs = fmt.Sprintf("%s idc_id = %d ", qs, id)

			if i != len(toBeDeleteIds)-1 {
				qs = fmt.Sprintf("%s OR", qs)
			} else {
				qs = fmt.Sprintf("%s)", qs)
			}
		}

		var userGroupIdcMapDelete []models.BsUserGroupIdcMap
		_, err = o.Raw(qs, this.UserInfo.CompanyIdNum, groupRequestData.GroupId).QueryRows(&userGroupIdcMapDelete)
		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
	}

	//查user_resource_id -> user_resource_instance_map -> instance and stack -> update group's instance_count and stack_count

	//查user_resource_id
	var userResourceIds []int64

	qs := fmt.Sprintf("SELECT id FROM bs_user_resource WHERE company_id = ? AND can_use = 0 AND (")

	for i, id := range groupRequestData.IdcIds {
		qs = fmt.Sprintf("%s idc_id = %d ", qs, id)

		if i != len(groupRequestData.IdcIds)-1 {
			qs = fmt.Sprintf("%s OR", qs)
		} else {
			qs = fmt.Sprintf("%s)", qs)
		}
	}

	_, err = o.Raw(qs, this.UserInfo.CompanyIdNum).QueryRows(&userResourceIds)

	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//查user_resource_instance_map
	var userResourceInstanceMaps []models.BsUserResourceInstanceMap

	qs = fmt.Sprintf("SELECT * FROM bs_user_resource_instance_map WHERE (")

	for i, id := range userResourceIds {
		qs = fmt.Sprintf("%s user_resource_id = %d ", qs, id)

		if i != len(userResourceIds)-1 {
			qs = fmt.Sprintf("%s OR", qs)
		} else {
			qs = fmt.Sprintf("%s)", qs)
		}
	}

	num, err = o.Raw(qs).QueryRows(&userResourceInstanceMaps)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	if num == 0 {
		//直接改group表

		group.Updated = time.Now()
		group.InstanceCount = 0
		group.StackCount = 0
		group.GroupName = groupRequestData.GroupName
		_, err = o.Update(&group)

		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}

	} else {
		//统计容器和应用个数
		ins := make(map[int64]int64)
		sta := make(map[int64]int64)

		for _, urim := range userResourceInstanceMaps {
			ins[urim.InstanceId] = urim.Id
			sta[urim.StackId] = urim.Id
		}

		instanceCount := len(ins)
		stackCount := len(sta)

		group.Updated = time.Now()
		group.InstanceCount = instanceCount
		group.StackCount = stackCount
		group.GroupName = groupRequestData.GroupName
		_, err = o.Update(&group)

		if err != nil {
			this.ServeError(500, err, "Internal Server Error")
		}
	}

	//返回值构造
	var cattleResourceData []interface{} = make([]interface{}, 1)
	cattleResourceData[0] = models.BsGroupResponseData{
		Group:  group,
		IdcIds: groupRequestData.IdcIds,
	}

	cattleResource := models.CattleResource{
		Id:           "",
		Type:         "resource",
		ResourceType: "BsGroup",
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

func (this *GroupController) Post() {

	//验证身份
	if this.UserInfo.CompanyIdNum == 0 {
		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
	}

	//声明局部变量
	orm.Debug = true
	o := orm.NewOrm()
	var group models.BsGroup
	var groups []models.BsGroup

	//验证requestbody合法
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &group)
	if err != nil {
		this.ServeError(404, err, "Not Found")
	}

	//验证group.GroupName名称合法
	re := regexp.MustCompile(`^[0-9a-zA-Z\x{4e00}-\x{9fa5}]{2,20}$`)
	if matched := re.MatchString(group.GroupName); !matched {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Name")
	}

	num, err := o.Raw("SELECT * FROM bs_group WHERE company_id = ? AND group_name = ?", this.UserInfo.CompanyIdNum, group.GroupName).QueryRows(&groups)
	if num > 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Name")
	} else if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//相应参数构造
	group.Created = time.Now() //.Format("2006-01-02 15:04:05")
	group.CompanyId = this.UserInfo.CompanyIdNum

	//写入表
	_, err = o.Insert(&group)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//返回值构造
	var cattleResourceData []interface{} = make([]interface{}, 1)
	cattleResourceData[0] = group

	cattleResource := models.CattleResource{
		Id:           "",
		Type:         "resource",
		ResourceType: "BsGroup",
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

func (this *GroupController) Delete() {

	//验证身份
	if this.UserInfo.CompanyIdNum == 0 {
		this.ServeErrorWithDetail(401, nil, "Unauthorized", "Token Expired")
	}

	//验证query合法
	querys := this.Ctx.Input.Params()
	query, ok := querys[":group_id"]

	if !ok || query == "" {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Id")
	}

	re := regexp.MustCompile(`^[0-9]*$`)
	if matched := re.MatchString(query); !matched {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Id")
	}

	//验证资源权限
	orm.Debug = true
	o := orm.NewOrm()

	var groups []models.BsGroup
	num, err := o.Raw("SELECT * FROM bs_group WHERE id = ? AND company_id = ?", query, this.UserInfo.CompanyIdNum).QueryRows(&groups)

	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}
	if num == 0 {
		this.ServeErrorWithDetail(404, nil, "Not Found", "Invalid Group Id")
	}

	//删记录
	var userGroupIdcMaps []models.BsUserGroupIdcMap
	num, err = o.Raw("DELETE FROM bs_user_group_idc_map WHERE company_id = ? AND group_id = ?", this.UserInfo.CompanyIdNum, query).QueryRows(&userGroupIdcMaps)
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	_, err = o.Delete(&groups[0])
	if err != nil {
		this.ServeError(500, err, "Internal Server Error")
	}

	//返回值构造
	groups[0].Id, _ = strconv.ParseInt(query, 10, 64)

	var cattleResourceData []interface{} = make([]interface{}, 1)
	cattleResourceData[0] = groups[0]

	cattleResource := models.CattleResource{
		Id:           "",
		Type:         "resource",
		ResourceType: "BsGroup",
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
