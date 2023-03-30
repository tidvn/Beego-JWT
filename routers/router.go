// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"DPay/controllers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	ns := beego.NewNamespace("/v1",
		//beego.NSBefore(handlers.Jwt),
		beego.NSNamespace("/object",
			beego.NSRouter("/", &controllers.ObjectController{}, "Get:GetAll"),
			beego.NSRouter("/", &controllers.ObjectController{}, "Post:Post"),
			beego.NSRouter("/:objectId", &controllers.ObjectController{}, "Get:Get"),
			beego.NSRouter("/:objectId", &controllers.ObjectController{}, "Put:Put"),
			beego.NSRouter("/:objectId", &controllers.ObjectController{}, "Delete:Delete"),
		),
		beego.NSNamespace("/user",
			beego.NSRouter("/:uid", &controllers.UserController{}, "Put:Put"),
			beego.NSRouter("/:uid", &controllers.UserController{}, "Delete:Delete"),
			beego.NSRouter("/register", &controllers.UserController{}, "Post:Post"),
			beego.NSRouter("/login", &controllers.UserController{}, "Get:Login"),
			beego.NSRouter("/logout", &controllers.UserController{}, "Get:Logout"),
		),
		beego.NSNamespace("/mesh",
			beego.NSRouter("/generateNonce", &controllers.MeshController{}, "Post:Getnonce"),
		),
	)
	beego.AddNamespace(ns)
}
