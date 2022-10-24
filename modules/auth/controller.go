package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/modules/user"
)

const SESSION_STORE_NAME = "session"
const SESSION_USER_ID_KEY = "user_id"

type AuthController struct {
	svc *AuthService
}

func NewController() *AuthController {
	return &AuthController{
		svc: NewService(
			user.NewRepository(),
		),
	}
}

func (ac *AuthController) Routes(e *gin.RouterGroup) {
	e.POST("/login", ac.login)
	e.POST("/register", ac.register)
	e.GET("/me", ac.me)
}

type loginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ac *AuthController) login(c *gin.Context) {
	var dto loginDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err)
		return
	}

	user, jwt, err := ac.svc.Login(dto.Username, dto.Password)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"user":  user.View(),
		"token": jwt,
	})

}

func (ac *AuthController) register(c *gin.Context) {
	var dto loginDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err)
		return
	}

	user, jwt, err := ac.svc.Register(dto.Username, dto.Password)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"user":  user.View(),
		"token": jwt,
	})

}

func (ac *AuthController) me(c *gin.Context) {
	bearerHeader := c.GetHeader("Authorization")
	if bearerHeader == "" {
		c.Error(errors.New("no authorization header provided"))
		return
	}

	tok := strings.Replace(bearerHeader, "Bearer ", "", 1)

	claims, err := VerifyJWT(tok)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if user_id, ok := claims["user_id"].(string); ok {
		user, err := ac.svc.userRepo.FindById(user_id)

		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, user.View())
	}

}
