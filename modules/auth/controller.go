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
	e.GET("/me", Protected(ac.me))
}

type loginDto struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

func (ac *AuthController) login(c *gin.Context) {
	var dto loginDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err)
		return
	}

	user, jwt, err := ac.svc.Login(dto)

	if err != nil {
		c.Error(err)
		return
	}

	SetUserToContext(c, user.ID)

	c.JSON(200, gin.H{
		"user": user.View(),
		"jwt":  jwt,
	})

}

func (ac *AuthController) register(c *gin.Context) {
	var dto loginDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err)
		return
	}

	user, jwt, err := ac.svc.Register(dto)

	if err != nil {
		c.Error(err)
		return
	}

	SetUserToContext(c, user.ID)

	c.JSON(200, gin.H{
		"user": user.View(),
		"jwt":  jwt,
	})

}

func (ac *AuthController) me(c *gin.Context, uuid string) {
	user, err := ac.svc.userRepo.FindById(uuid)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, user.View())

}
