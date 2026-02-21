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
