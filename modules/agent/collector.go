package agent

import (
	"encoding/json"
	"time"

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
				app.push(collectResults)
				collectResults = collectResults[0:0] // 重置
			}

			timer.Reset(batchDuration) // 重置
		case collectResult, ok := <-collectResultChan:
			if !ok {
				if len(collectResults) > 0 {
					app.push(collectResults)
					collectResults = collectResults[0:0] // 重置
				}
				app.logger.Info("清理完成剩余日志")
				timer.Stop()
				return
			}

			collectResults = append(collectResults, collectResult)
			if len(collectResults) >= batchLength {
				app.push(collectResults)
				collectResults = collectResults[0:0] // 重置
				timer.Reset(batchDuration)           // 重置
			}
		}
	}
}

func (app *AgentNode) push(collectResults []*collector.CollectResult) {
	value, err := json.Marshal(collectResults)
	rawLength := len(value)
	if err == nil {
		if app.compactor != nil {
			// comparess
			value, err = app.compactor.Encode(value)
		}
	}

	if err != nil {
		// compactor err
		app.logger.Warnf("push encode err:%s", err.Error())
	}

	compressLength := len(value)
	app.logger.Infof("rawLength:%d compressLength:%d", rawLength, compressLength)
}
