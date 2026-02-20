package handler

import (
	"net/http"

	"github.com/ArchDevs/radix/internal/service"
	"github.com/ArchDevs/radix/internal/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Address   string `json:"address"`
		PublicKey string `json:"public_key"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := validation.ValidateAddress(req.Address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validation.ValidatePublicKey(req.PublicKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := h.authSvc.Register(c.Request.Context(), req.Address, req.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "created",
	})
}
