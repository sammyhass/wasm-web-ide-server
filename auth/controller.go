package auth

import (
	"github.com/gin-gonic/gin"
)

type controller struct {
	svc *Service
}

func NewController() *controller {
	return &controller{
		svc: NewService(),
	}
}

func (ac *controller) Routes(e *gin.RouterGroup) {
	e.POST("/login", ac.login)
	e.POST("/register", ac.register)
	e.GET("/me", Protected(ac.me))
}

type loginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ac *controller) login(c *gin.Context) {
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

func (ac *controller) register(c *gin.Context) {
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

func (ac *controller) me(c *gin.Context, uuid string) {
	user, err := ac.svc.userRepo.findById(uuid)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, user.View())

}
