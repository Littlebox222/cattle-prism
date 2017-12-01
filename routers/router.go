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

	// beego.Router("/bsgroups", &controllers.GroupController{})
	// beego.Router("/bsidcs", &controllers.IdcController{})
	// beego.Router("/bsareas", &controllers.AreaController{})
	// beego.Router("/bscarrieroperators", &controllers.CarrierOperatorController{})
	// beego.Router("/bsgroupidcmaps", &controllers.GroupIdcMapController{})
	beego.Router("/bs-v1/resource/total", &controllers.BsUserResourceTotalController{})

	beego.Router("*", &controllers.AppController{})

	// beego.Handler("/rpc", nil)
}
