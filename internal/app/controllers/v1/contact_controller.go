package v1

import (
	"htst/internal/app/controllers"
	"htst/internal/app/models"
	"htst/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactController struct {
	service services.IContactService
}

func NewContactController(service services.IContactService) *ContactController {
	return &ContactController{
		service: service,
	}
}

// CreateContact 创建联系人
func (c *ContactController) CreateContact(ctx *gin.Context) {
	var req models.Contact
	if err := ctx.ShouldBindJSON(&req); err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "参数错误: %v", err)
		return
	}

	contact, err := c.service.CreateContact(&req)
	if err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "创建联系人失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "创建联系人成功", contact)
}

// UpdateContact 更新联系人
func (c *ContactController) UpdateContact(ctx *gin.Context) {
	var req models.Contact
	if err := ctx.ShouldBindJSON(&req); err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "参数错误: %v", err)
		return
	}

	if err := c.service.UpdateContact(req.ID, &req); err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "更新联系人失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "更新联系人成功", nil)
}

// FirstContact 获取第一个联系人
func (c *ContactController) FirstContact(ctx *gin.Context) {
	contact, err := c.service.FirstContact()
	if err != nil {
		controllers.Response(ctx, http.StatusOK, "请添加联系人", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "获取第一个联系人成功", contact)
}
