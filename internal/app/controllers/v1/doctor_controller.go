package v1

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"htst/internal/app/controllers"
	"htst/pkg/config"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"htst/internal/app/models"
	"htst/internal/app/services"
)

type DoctorInfoController struct {
	service services.IDoctorInfoService
}

func NewDoctorInfoController(service services.IDoctorInfoService) *DoctorInfoController {
	return &DoctorInfoController{service: service}
}

// CreateDoctorInfo 创建医生信息（支持上传头像）
// @Summary 创建医生信息（支持上传头像）
// @Description 创建新的医生信息记录，支持上传头像文件
// @Tags 医生信息
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "医生姓名"
// @Param title formData string true "职称"
// @Param department formData string true "科室"
// @Param degree formData string false "学历"
// @Param introduction formData string false "个人简介"
// @Param detailed_description formData string false "详细介绍"
// @Param specialties formData string false "专业领域"
// @Param research_achievements formData string false "科研成果"
// @Param avatar_file formData file false "头像文件"
// @Success 200 {object} models.DoctorInfo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /doctors [post]
func (c *DoctorInfoController) CreateDoctorInfo(ctx *gin.Context) {
	// 处理文件上传
	file, err := ctx.FormFile("file")
	if err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "文件上传失败", nil)
		return
	}

	dirPath := filepath.Join(config.GetAppConf().FilePath, "avatars")
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Errorf("本地目录:%s, 创建失败", dirPath)
		controllers.Response(ctx, http.StatusInternalServerError, "创建本地目录失败", nil)
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("avatar_%d%s", time.Now().Unix(), filepath.Ext(file.Filename))
	// 设置文件保存路径（可根据实际需求调整）
	savePath := filepath.Join(config.AppConf.FilePath, "avatars", filename)
	// 保存文件
	log.Infoln("保存文件到本地目录: %s", savePath)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "文件保存失败:"+err.Error(), nil)
		return
	}

	// 解析表单字段
	doctor := models.DoctorInfo{
		Name:                 ctx.PostForm("name"),
		Title:                ctx.PostForm("title"),
		Department:           ctx.PostForm("department"),
		AvatarURL:            "avatars/" + filename,
		Degree:               ctx.PostForm("degree"),
		Introduction:         ctx.PostForm("introduction"),
		DetailedDescription:  ctx.PostForm("detailed_description"),
		Specialties:          ctx.PostForm("specialties"),
		ResearchAchievements: ctx.PostForm("research_achievements"),
	}

	// 验证必填字段
	if doctor.Name == "" || doctor.Title == "" || doctor.Department == "" {
		controllers.Response(ctx, http.StatusOK, "姓名、职称、科室为必填项", nil)
		return
	}

	// 调用服务创建医生信息
	createdDoctor, err := c.service.CreateDoctorInfo(&doctor)
	if err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "创建医生信息失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "创建医生信息成功", createdDoctor)
}

// GetDoctorInfo 获取医生信息
// @Summary 获取医生信息
// @Description 根据ID获取医生详细信息
// @Tags 医生信息
// @Accept json
// @Produce json
// @Param id path int true "医生ID"
// @Success 200 {object} models.DoctorInfo
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /doctors/{id} [get]
func (c *DoctorInfoController) GetDoctorInfo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		controllers.Response(ctx, http.StatusOK, "无效的ID参数", nil)
		return
	}

	doctor, err := c.service.GetDoctorInfoByID(id)
	if err != nil {
		controllers.Response(ctx, http.StatusNotFound, "医生信息不存在", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "获取医生信息成功", doctor)
}

// UpdateDoctorInfo 更新医生信息
// @Summary 更新医生信息
// @Description 根据ID更新医生信息
// @Tags 医生信息
// @Accept json
// @Produce json
// @Param id path int true "医生ID"
// @Param doctor body models.DoctorInfo true "医生信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /doctors/{id} [put]
func (c *DoctorInfoController) UpdateDoctorInfo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		controllers.Response(ctx, http.StatusOK, "无效的ID参数", nil)
		return
	}

	var doctor models.DoctorInfo
	if err := ctx.ShouldBindJSON(&doctor); err != nil {
		controllers.Response(ctx, http.StatusOK, "参数错误"+err.Error(), nil)
		return
	}

	err = c.service.UpdateDoctorInfo(id, &doctor)
	if err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "更新医生信息失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "更新医生信息成功", nil)
}

// DeleteDoctorInfo 删除医生信息
// @Summary 删除医生信息
// @Description 根据ID删除医生信息
// @Tags 医生信息
// @Accept json
// @Produce json
// @Param id path int true "医生ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /doctors/{id} [delete]
func (c *DoctorInfoController) DeleteDoctorInfo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		controllers.Response(ctx, http.StatusOK, "无效的ID参数", nil)
		return
	}

	err = c.service.DeleteDoctorInfo(id)
	if err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "删除医生信息失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "删除医生信息成功", nil)
}

// ListDoctorInfos 列出医生信息
// @Summary 列出医生信息
// @Description 分页列出医生信息，支持查询条件
// @Tags 医生信息
// @Accept json
// @Produce json
// @Param name query string false "医生姓名"
// @Param department query string false "科室"
// @Param title query string false "职称"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /doctors [get]
func (c *DoctorInfoController) ListDoctorInfos(ctx *gin.Context) {
	var query models.DoctorInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		controllers.Response(ctx, http.StatusOK, "参数错误"+err.Error(), nil)
		return
	}

	_, doctors, err := c.service.ListDoctorInfos(&query)
	if err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "查询医生信息失败", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "查询医生信息成功", doctors)
}
