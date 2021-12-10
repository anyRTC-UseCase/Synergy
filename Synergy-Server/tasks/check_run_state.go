package tasks

import (
	"ARTeamViewService/global"
	"ARTeamViewService/services"
	"ARTeamViewService/utils"
	"github.com/robfig/cron/v3"
)

func ServeCheckRunLoop() {
	//开启执行任务
	startDayAndMonthTask()
}

/**
 * 执行天任务和月任务
 */
func startDayAndMonthTask() {
	//具体定时任务时间格式，参考https://godoc.org/github.com/robfig/cron
	c := cron.New()

	//每天0点50分执行任务删掉前一天的心跳包
	c.AddFunc("50 0 * * *", func() {
		dayDeleteUserOnlineInfo()
	})

	c.Start()
}

/**
 * 每天0点50分执行任务删掉前一天的心跳包
 */
func dayDeleteUserOnlineInfo() {
	global.GLogger.Info(" function task dayUpdateLicStatus ", utils.FormatFullDate(utils.FormatNowUnix()))

	services.GetInsUserSvr().DayDeleteUserOnlineInfo()
}
