package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/services"
	worker "github.com/sonnyvictok/miniapp_taptoearn/internal/workers"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type UserHandler struct {
	userServices    *services.UserService
	taskDistributor worker.TaskDistributor
}

func NewUserHandler(userServices *services.UserService, taskDistributor worker.TaskDistributor) *UserHandler {
	return &UserHandler{userServices: userServices, taskDistributor: taskDistributor}
}

type ScoreUserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Id        string `json:"id"`
	Username  string `json:"username"`
}

func (h *UserHandler) ClickToEarn(c *gin.Context) {
	user, exists := c.Get("tma")
	fmt.Println("user", user)
	if !exists {
		c.JSON(400, gin.H{"error": "user information not found in context"})
		return
	}
	userData, ok := user.(initdata.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data format"})
		return
	}
	payload := &worker.PayloadHandleTapToEarn{
		TelegramID: userData.ID,
		Score:      1,
	}

	err := h.taskDistributor.DistributeTaskHandleTapToEarn(c.Request.Context(), payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Click request received, processing...",
		"user_id":  userData.ID,
		"username": userData.Username,
	})

	// user, err := h.userServices.GetOrCreateUser((userData.ID), userData.Username)

	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	// score, err := h.userServices.ClickToEarn((userData.ID))
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"user_id":    userData.ID,
	// 	"last_name":  userData.LastName,
	// 	"first_name": userData.FirstName,
	// 	"username":   userData.Username,
	// 	"photo_url":  userData.PhotoURL,
	// 	"language":   userData.LanguageCode,
	// 	"score":      score,
	// })
}

func (h *UserHandler) GetUser(c *gin.Context) {
	user, exists := c.Get("tma")
	if !exists {
		c.JSON(400, gin.H{"error": "user information not found in context"})
		return
	}
	userData, ok := user.(initdata.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data format"})
		return
	}

	user, err := h.userServices.GetOrCreateUser(userData.ID, userData.Username)
	fmt.Println("user", user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	user, exists := c.Get("tma")
	if !exists {
		c.JSON(400, gin.H{"error": "user information not found in context"})
		return
	}
	userData, ok := user.(initdata.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data format"})
		return
	}

	user, err := h.userServices.GetOrCreateUser(userData.ID, userData.Username)
	fmt.Println("user", user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
