package proxy

import (
	"dudu/commons/log"
	"dudu/models"
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/cpu"
)

type Parser struct {
	logger log.Logger
}

func NewParser(logger log.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

func (p *Parser) Parser(metric models.MetricValue) []*models.CollectResult {
	collectResults := make([]*models.CollectResult, 0, 100)
	for _, collectResult := range collectResults {
		if err := p.parser(collectResult); err != nil {
			p.logger.Warnf("parse %s err:%s", collectResult.Metric, err.Error())
		}
	}
	return collectResults
}

func (p *Parser) parser(result *models.CollectResult) (err error) {
	if result.Err != "" {
		return fmt.Errorf(result.Err)
	}

	switch result.Metric {
	case "CPUCount":
		var cpuCoount int
		err = json.Unmarshal(result.Value, &cpuCoount)
		result.RelValue = cpuCoount
		return
	case "CPUInfo":
		infos := make([]cpu.InfoStat, 0, 10)
		err = json.Unmarshal(result.Value, &infos)
		result.RelValue = infos
		return
	case "CPUTimes":
		stats := make([]cpu.TimesStat, 0, 10)
		err = json.Unmarshal(result.Value, &stats)
		result.RelValue = stats
		return
	}
	return fmt.Errorf("%s type err, data:%+v", result.Metric, string(result.Value))
}
