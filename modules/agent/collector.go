package agent

import (
	"encoding/json"
	"time"

	"dudu/models"
	"dudu/modules/agent/collector"
	_ "dudu/modules/agent/collector/collect"
)

// 初始化采集管理器
func (app *AgentNode) initCollect() error {
	app.collectorMag = collector.NewCollectorManager(app.logger, app.cfg.Agent.Collects)
	return nil
}

// 开始采集
func (app *AgentNode) startCollect() {
	collectResultChan := app.collectorMag.Run()
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		app.asyncPush(collectResultChan)
	}()
}

// 关闭采集
func (app *AgentNode) stopCollect() {
	app.collectorMag.Stop()
}

// 推送采集信息
func (app *AgentNode) asyncPush(collectResultChan <-chan *collector.CollectResult) {
	// 过期时长
	batchDuration := time.Second * time.Duration(app.cfg.Agent.BatchDuration)
	if batchDuration <= 0 {
		batchDuration = time.Second * 5 // 默认最长间隔5秒
	}

	batchLength := app.cfg.Agent.BatchLength
	if batchLength <= 0 {
		batchLength = 100 // 默认最大100条
	}

	timer := time.NewTimer(batchDuration)
	collectResults := make([]*collector.CollectResult, 0, batchLength)
	for {
		select {
		case <-timer.C:
			if len(collectResults) > 0 {
				if err := app.push(collectResults); err != nil {
					app.logger.Warnf("push collect results err:%s", err.Error())
				}

				collectResults = collectResults[0:0] // 重置
			}

			timer.Reset(batchDuration) // 重置
		case collectResult, ok := <-collectResultChan:
			if !ok {
				if len(collectResults) > 0 {
					if err := app.push(collectResults); err != nil {
						app.logger.Warnf("push collect results err:%s", err.Error())
					}

					collectResults = collectResults[0:0] // 重置
				}
				app.logger.Info("清理完成剩余日志")
				timer.Stop()
				return
			}

			collectResults = append(collectResults, collectResult)
			if len(collectResults) >= batchLength {
				if err := app.push(collectResults); err != nil {
					app.logger.Warnf("push collect results err:%s", err.Error())
				}
				collectResults = collectResults[0:0] // 重置
				timer.Reset(batchDuration)           // 重置
			}
		}
	}
}

func (app *AgentNode) push(collectResults []*collector.CollectResult) (err error) {
	value, err := json.Marshal(collectResults)
	if err != nil {
		return
	}

	var compactorName string
	if app.compactor != nil {
		// comparess
		value, err = app.compactor.Encode(value)
		compactorName = app.compactor.Name()
	}

	if err != nil {
		// compactor err
		return
	}

	metric := &models.MetricValue{
		Endpoint:  app.cfg.IP,
		HostName:  app.cfg.HostName,
		Compactor: compactorName,
		Value:     value,
		Tags:      "",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	data, err := json.Marshal(metric)
	if err != nil {
		return
	}

	return app.pipe.Push(data)
}
