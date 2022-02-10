package cronutil

import (
	"context"
	"fmt"
	"time"

	"github.com/cn1211/goeva/core/exp"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type CronUtil struct {
	cron *cron.Cron

	cronJobs  []CronJob
	isRunning bool
}

type CronJob struct {
	spec     string
	isRunNow bool
	job      func(ctx context.Context)
}

// AddJob 定时任务必须在try包方法集中执行.避免任务导致的panic引起程序异常
func (c *CronUtil) AddJob(spec string, isRunNow bool, job func(ctx context.Context)) error {
	if c.isRunning {
		return errors.New("定时模块已启动，无法添加定时任务")
	}

	c.cronJobs = append(c.cronJobs, CronJob{
		spec:     spec,
		isRunNow: isRunNow,
		job:      job,
	})

	return nil
}

func (c *CronUtil) addCronFunc(ctx context.Context, cronJob CronJob) error {
	// 添加定时任务
	if _, err := c.cron.AddFunc(cronJob.spec, func() {
		cronJob.job(ctx)
	}); err != nil {
		return fmt.Errorf("添加cron定时 任务失败 内容 = %s  原因 = %s", cronJob.spec, err)
	}

	if cronJob.isRunNow {
		go exp.TryWithErr(func() {
			cronJob.job(ctx)
		}, func(ex error) {
			logrus.Error("初次执行定时任务失败 原因= ", ex)
		})
	}

	return nil
}

func (c *CronUtil) Init(ctx context.Context) error {
	// 设置为启动状态
	c.isRunning = true

	c.cron = cron.New(cron.WithLocation(time.Now().Location()))
	for _, cronJob := range c.cronJobs {
		if err := c.addCronFunc(ctx, cronJob); err != nil {
			return err
		}
	}

	// 协程启动.不会阻塞主协程
	c.cron.Start()
	return nil
}

func (c *CronUtil) Close() error {
	if !c.isRunning {
		// 未启动
		return nil
	}

	c.isRunning = false
	if c.cron != nil {
		c.cron.Stop()
	}

	return nil
}
