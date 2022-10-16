package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/modules/user"
)

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

	ac.svc.Login(dto.Username, dto.Password)
}

func (ac *AuthController) register(c *gin.Context) {
	var dto loginDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err)
		return
	}

	ac.svc.Register(dto.Username, dto.Password)

}
