package user

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc UserService
}

type UserResponse struct {
	Address     string    `json:"address"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewUserHandler(userSvc UserService) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
	}
}

func (h *UserHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if len(query) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid username/address",
		})
		return
	}

	users, err := h.userSvc.Search(c.Request.Context(), query, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		log.Printf("search failed: %v", err)
		return
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{
			Address:     user.Address,
			Username:    user.Username.String,
			DisplayName: user.DisplayName.String,
			CreatedAt:   user.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Me(c *gin.Context) {
	address := c.GetString("address")

	user, err := h.userSvc.GetUser(c.Request.Context(), address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	response := UserResponse{
		Address:     user.Address,
		Username:    user.Username.String,
		DisplayName: user.DisplayName.String,
		CreatedAt:   user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) SetUsername(c *gin.Context) {
	address := c.GetString("address")

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.userSvc.UpdateUsername(c.Request.Context(), address, req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "username updated", "username": req.Username})
}
