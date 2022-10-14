package wasm

import "github.com/gin-gonic/gin"

const (
	CONTROLLER_ROUTE string = "/wasm"
)

type WasmController struct {
	svc *WasmService
}

func NewController() *WasmController {
	svc := NewService()

	return &WasmController{
		svc: svc,
	}

}

func (wc *WasmController) Routes(e *gin.RouterGroup) {
	e.POST("compile", wc.compile)
}

type compileDTO struct {
	Code string `json:"code"`
}

func (wc *WasmController) compile(c *gin.Context) {
	var dto compileDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err)
		return
	}

	wasm, err := wc.svc.Compile(dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"path": wasm})
}
