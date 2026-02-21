package challenge

import (
	"net/http"

	"github.com/ArchDevs/radix/internal/service"
	"github.com/ArchDevs/radix/internal/validation"
	"github.com/gin-gonic/gin"
)

type ChallengeHandler struct {
	challengeSvc *ChallengeService
	jwtSvc       *service.JWTService
}

func NewChallengeHandler(challengeSvc *ChallengeService, jwtSvc *service.JWTService) *ChallengeHandler {
	return &ChallengeHandler{
		challengeSvc: challengeSvc,
		jwtSvc:       jwtSvc,
	}
}

func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {
	address := c.Query("address")
	if err := validation.ValidateAddress(address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	nonce, timestamp, err := h.challengeSvc.CreateChallenge(c.Request.Context(), address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create challenge",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"nonce":     nonce,
		"timestamp": timestamp,
	})
}

func (h *ChallengeHandler) Verify(c *gin.Context) {
	var req struct {
		Address   string `json:"address"`
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
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

	if err := validation.ValidateNonce(req.Nonce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validation.ValidateSignature(req.Signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	verified, err := h.challengeSvc.Verify(c.Request.Context(), req.Address, req.Nonce, req.Signature)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "verification failed",
			"message": err.Error(),
		})
		return
	}

	if !verified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "signature verification failed",
		})
		return
	}

	token, err := h.jwtSvc.Generate(req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
