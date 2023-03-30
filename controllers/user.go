package controllers

import (
	"DPay/models"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

func (u *UserController) Post() {

	hash, _ := models.HashPassword(u.GetString("password"))
	user := models.User{
		UserName:     u.GetString("username"),
		PasswordHash: hash,
		Email:        u.GetString("email"),
		FullName:     u.GetString("fullname"),
	}
	db, _ := models.Connect()
	uid, err := db.AddUser(&user)
	if err == nil {
		token := models.AddToken(user, u.Ctx.Input.Domain())
		dataResponse := models.APIResponse{
			Status:  200,
			Message: "Register Successfully",
			Errors:  "",
			Data:    map[string]string{"uid": uid, "token": token},
		}
		u.Data["json"] = dataResponse
	} else {
		dataResponse := models.APIResponse{
			Status:  200,
			Message: "Failed to Register",
			Errors:  err.Error(),
			Data:    map[string]string{},
		}
		u.Data["json"] = dataResponse
	}

	u.ServeJSON()
}

func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		db, _ := models.Connect()
		_, err := db.UpdateUser(&user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJSON()
}

func (u *UserController) Delete() {
	uid := u.GetString(":uid")
	db, _ := models.Connect()
	db.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")
	db, _ := models.Connect()
	user, err := db.Login(username, password)
	if err == nil {
		token := models.AddToken(user, u.Ctx.Input.Domain())
		dataResponse := models.APIResponse{
			Status:  200,
			Message: "Login Successfully",
			Errors:  "",
			Data:    map[string]string{"uid": user.ID.Hex(), "token": token},
		}
		u.Data["json"] = dataResponse
	} else {
		u.Data["json"] = models.APIResponse{
			Status:  200,
			Message: "Failed to login,username or password not match",
			Errors:  "Invalid Credential",
			Data:    map[string]string{},
		}
	}
	u.ServeJSON()
}
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}
