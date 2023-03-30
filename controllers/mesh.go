package controllers

import (
	"DPay/models"

	beego "github.com/beego/beego/v2/server/web"
)

// MeshController operations for Mesh
type MeshController struct {
	beego.Controller
}

func (c *MeshController) Getnonce() {
	nonce, _ := models.GenerateNonce(32)
	c.Data["json"] = nonce
	c.ServeJSON()
}
