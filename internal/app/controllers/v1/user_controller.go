package v1

import (
	"htst/internal/app/controllers"
	"htst/pkg/util"
	"net/http"
	"strconv"

	"htst/internal/app/models"
	"htst/internal/app/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.IUserService
}

func NewUserController(userService services.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.SysUser
	if err := ctx.ShouldBindJSON(&user); err != nil {
		controllers.Response(ctx, http.StatusOK, "参数错误"+err.Error(), nil)
		return
	}

	// 对密码进行加密
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "密码加密失败", nil)
		return
	}
	user.Password = hashedPassword

	if err := c.userService.CreateUser(&user); err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "创建用户失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "创建成功", user)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	var user models.SysUser
	if err := ctx.ShouldBindJSON(&user); err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "参数错误", nil)
		return
	}

	if user.Password != "" {
		// 对密码进行加密
		hashedPassword, err := util.HashPassword(user.Password)
		if err != nil {
			controllers.Response(ctx, http.StatusInternalServerError, "密码加密失败", nil)
			return
		}
		user.Password = hashedPassword
	}

	if err := c.userService.UpdateUser(&user); err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "更新用户失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "更新成功", user)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "ID格式错误", nil)
		return
	}

	if err := c.userService.DeleteUserByID(id); err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "删除用户失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "删除成功", nil)
}

func (c *UserController) GetUserList(ctx *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := c.userService.GetUserList(page, pageSize)
	if err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "获取用户列表失败", nil)
		return
	}

	result := map[string]interface{}{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	controllers.Response(ctx, http.StatusOK, "获取成功", result)
}

func (c *UserController) AuthUser(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "登录失败，账号或者密码错误", nil)
		return
	}

	user, err := c.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		controllers.Response(ctx, http.StatusUnauthorized, "登录失败，账号或者密码错误", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "登录成功", user)
}
