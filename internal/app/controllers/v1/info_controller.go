package v1

import (
	log "github.com/sirupsen/logrus"
	"htst/pkg/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"htst/internal/app/controllers"
	"htst/internal/app/models"
	"htst/internal/app/services"
)

type InfoController struct {
	infoService services.IInfoService
}

func NewInfoController(infoService services.IInfoService) *InfoController {
	return &InfoController{
		infoService: infoService,
	}
}

func (c *InfoController) CreateInfo(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "文件上传失败", nil)
		return
	}

	title := ctx.PostForm("title")
	infoType := ctx.PostForm("type")
	coreDescription := ctx.PostForm("core_description")
	permissionNote := ctx.PostForm("permission_note")
	remark := ctx.PostForm("remark")

	if title == "" || infoType == "" {
		controllers.Response(ctx, http.StatusBadRequest, "参数缺失", nil)
		return
	}

	// 保存文件
	log.Infoln("保存文件:", file.Filename)
	// 创建目录路径
	dirPath := filepath.Join("uploads")
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Errorf("本地目录:%s, 创建失败", dirPath)
		controllers.Response(ctx, http.StatusInternalServerError, "创建本地目录失败", nil)
		return
	}

	// 生成新的文件名：type值+yymmddhhMMss
	timestamp := time.Now().Format("060102150405") // yyMMddHHmmss格式
	newFileName := infoType + "-" + timestamp + filepath.Ext(file.Filename)

	// 保存文件到本地时间命名的目录
	filePath := filepath.Join(dirPath, newFileName)
	log.Infoln("保存文件:", filePath)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "文件保存失败", nil)
		return
	}

	md5, _ := util.CalculateFileMD5(filePath)

	// 创建信息记录
	info := &models.Info{
		Type:            infoType,
		Title:           title,
		Format:          filepath.Ext(file.Filename)[1:], // 移除点号
		FilePath:        filePath,
		Size:            uint64(file.Size),
		CoreDescription: coreDescription,
		PermissionNote:  permissionNote,
		Remark:          remark,
		Md5:             md5,
	}

	if err := c.infoService.CreateInfo(info); err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "创建失败,重复创建相同名字的内容", nil)
		return
	}

	controllers.Response(ctx, http.StatusOK, "创建成功", info)
}

func (c *InfoController) DeleteInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "ID格式错误", nil)
		return
	}

	info, err := c.infoService.GetInfoByID(id)
	if err != nil {
		controllers.Response(ctx, http.StatusNotFound, "信息不存在", nil)
		return
	}

	if err := c.infoService.DeleteInfoByID(id); err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "删除失败", nil)
		_ = os.Remove(info.FilePath)
		return
	}

	controllers.Response(ctx, http.StatusOK, "删除成功", nil)
}

func (c *InfoController) GetInfoList(ctx *gin.Context) {
	var query models.InfoQuery
	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&query); err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "参数错误", nil)
		return
	}

	// 设置默认值
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	total, infos, err := c.infoService.GetInfoListByQuery(query)
	if err != nil {
		controllers.Response(ctx, http.StatusInternalServerError, "获取列表失败", nil)
		return
	}

	result := map[string]interface{}{
		"list":      infos,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	}

	controllers.Response(ctx, http.StatusOK, "获取成功", result)
}

func (c *InfoController) Increment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		controllers.Response(ctx, http.StatusBadRequest, "ID格式错误", nil)
		return
	}

	//info, err := c.infoService.GetInfoByID(id)
	//if err != nil {
	//	controllers.Response(ctx, http.StatusNotFound, "信息不存在", nil)
	//	return
	//}

	// 累加访问次数
	if err := c.infoService.IncrementInfoCount(id); err != nil {
		// 记录日志，但不中断请求
	}

	controllers.Response(ctx, http.StatusOK, "获取成功", id)
}
