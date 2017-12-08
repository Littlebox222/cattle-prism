// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"cattle-prism/controllers"
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/context"
	// "net/http/httputil"
	// "net/url"
	// "github.com/yhat/wsutil"
)

// func init() {
// 	ns := beego.NewNamespace("/v1",
// 		beego.NSNamespace("/object",
// 			beego.NSInclude(
// 				&controllers.ObjectController{},
// 			),
// 		),
// 		beego.NSNamespace("/user",
// 			beego.NSInclude(
// 				&controllers.UserController{},
// 			),
// 		),
// 	)
// 	beego.AddNamespace(ns)
// }

func init() {

	beego.Router("/v2-beta/projects/:project_id/groups/?:group_id", &controllers.GroupController{})
	beego.Router("/v2-beta/projects/:project_id/idcs", &controllers.IdcController{})
	beego.Router("/v2-beta/projects/:project_id/areas", &controllers.AreaController{})
	beego.Router("/v2-beta/projects/:project_id/carrieroperators", &controllers.CarrierOperatorController{})
	beego.Router("/v2-beta/projects/:project_id/resource/total", &controllers.BsUserResourceTotalController{})

	beego.Router("/v2-beta/projects/:project_id/services/:service_id/groups", &controllers.GroupController{})
	beego.Router("/v2-beta/projects/:project_id/services/:service_id/containertypes", &controllers.ContainerTypeController{})

	// beego.Router("/bsgroupidcmaps", &controllers.GroupIdcMapController{})
	// beego.Router("/bs-v1/stacks/?:id", &controllers.BsStackController{})
	// beego.Router("/bs-v1/services", &controllers.BsServiceController{})
	// beego.Router("/bs-v1/instances", &controllers.BsInstanceController{})
	// beego.Router("/bs-v1/auditlogs", &controllers.BsAuditLogController{})

	beego.Router("*", &controllers.AppController{})

	// beego.Handler("/rpc", nil)
}
