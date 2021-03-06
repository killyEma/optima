package controllers

import (
	"log"

	"encoding/hex"
	"github.com/astaxie/beegae"
	"github.com/golang/crypto/bcrypt"
	"github.com/jackgris/optima/models"
	"math/rand"
)

type RegisterAuth struct {
	beegae.Controller
}

func (this *RegisterAuth) Post() {

	user, err := DecodeUserData(this.Ctx.Input.Request.Body)
	if err != nil {
		log.Println("RegisterAuth: Error decode user data on register", err)
		this.Data["json"] = &models.Token{}
		return
	}

	token, _ := models.GenerateToken(user.Email)
	user.Token = token
	exist, err := models.CheckExist(user, this.AppEngineCtx)
	if err != nil {
		log.Println("RegisterAuth: Error at verify user", err)
		this.Data["json"] = &models.Token{}
		return
	}
	if exist {
		log.Println("RegisterAuth: The user already exist")
		this.Data["json"] = &models.Token{}
		return
	} else {
		cost := rand.Intn(27) + 4 //cost min 4 and max 31
		hashHex, err := bcrypt.GenerateFromPassword([]byte(user.Pass), cost)
		hashString := hex.EncodeToString(hashHex)
		user.Pass = hashString
		if err != nil {
			log.Println("bcrypt broke")
			return
		}
		models.AddUser(user, this.AppEngineCtx)
	}
	// It's all ok, return the user data on json format
	this.Data["json"] = &user.Token
}

func (this *RegisterAuth) Get() {
	// not implemmented
}

func (this *RegisterAuth) Render() error {
	if _, ok := this.Data["json"].(error); ok {
		this.AppEngineCtx.Errorf("register auth error: %v", this.Data["json"])
	}
	this.ServeJson()
	return nil
}
