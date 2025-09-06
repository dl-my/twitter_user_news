package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"twitter_user_news/model"
	"twitter_user_news/service"
	"twitter_user_news/utils"
)

func Add(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.StartSearchTask(req.UserName, req.UserId, time.Duration(req.Interval)*time.Second)

	c.JSON(http.StatusOK, model.Response{Message: "success"})
}

func Del(c *gin.Context) {
	var req model.DelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.StopSearchTask(req.UserName)

	c.JSON(http.StatusOK, model.Response{Message: "success"})
}

func List(c *gin.Context) {
	result := service.GetAllTasks()
	c.JSON(http.StatusOK, gin.H{"taskList": result})
}

func DelAll(c *gin.Context) {
	service.StopAllTasks()

	c.JSON(http.StatusOK, model.Response{Message: "success"})
}

func ReloadCookie(c *gin.Context) {
	if err := utils.LoadTokens(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{Message: "success"})
}
