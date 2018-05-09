package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../models"
	"../module"
	"../utils"
)

type AccountController struct {
	*BaseController
}

func (self AccountController) GetAdd(c *gin.Context) {
	themisView.AccountView{self.BaseView}.GetAdd(c)
}

func (self AccountController) PostAdd(c *gin.Context) {
	var addRequest models.AccountAddRequestJson
	c.ShouldBindJSON(&addRequest)

	addResult := &models.AccountAddResultJson{}

	if len(addRequest.Name) < 1 {
		addResult.Message = "id is not allowed empty"
		themisView.AccountView{self.BaseView}.PostAdd(c, addResult)
		return
	}

	if len(addRequest.Name) > 128 {
		addResult.Message = "maximum name length is 128 characters"
		themisView.AccountView{self.BaseView}.PostAdd(c, addResult)
		return
	}

	accountModule := module.NewAccountModule(self.DB)

	if accountModule.Get(addRequest.Name) > 0 {
		addResult.Message = "this id is exist"
		themisView.AccountView{self.BaseView}.PostAdd(c, addResult)
		return
	}

	// TODO: ロールによる権限管理の追加

	password := utils.RandomString(24)

	accountModule.Add(addRequest.Name, password)

	addResult.Name = addRequest.Name
	addResult.Success = true
	addResult.Password = password
	themisView.AccountView{self.BaseView}.PostAdd(c, addResult)
}