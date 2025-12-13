package servers

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"
	"htst/internal/app/routers"
	"htst/pkg/config"
)

var (
	apiImpl *api
	apiOnce sync.Once
)

type api struct {
	server *http.Server
	cron   *cron.Cron
}

func ApiServer() IServer {
	apiOnce.Do(func() {
		apiImpl = &api{}
		apiImpl.cron = cron.New(cron.WithSeconds())
		apiImpl.server = &http.Server{
			Addr:    fmt.Sprintf(":%d", config.AppConf.Port),
			Handler: routers.SetUp(),
		}
		log.Infoln("api init finish")
	})
	return apiImpl
}

func (r *api) Start() error {
	log.Infoln("api sever start")
	// 添加文件清理任务，每天凌晨1点执行
	_, err := r.cron.AddFunc("0 14 18 * * *", func() {
		log.Infoln("开始清理文件定时任务：", config.AppConf.FilePath)
		//services.IInfo.CleanFile(config.AppConf.FilePath)
	})
	if err != nil {
		log.Errorf("添加文件清理任务失败: %v", err)
	} else {
		log.Info("定时任务调度器启动成功")
		r.cron.Start()
	}
	return r.server.ListenAndServe()
}

func (r *api) Stop() error {
	log.Infoln("api server stop")
	return r.server.Shutdown(context.Background())
}
