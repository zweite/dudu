package proxy

import (
	"encoding/json"

	"dudu/commons/log"
	"dudu/models"
	"dudu/modules/collector"
	_ "dudu/modules/collector/collect"
)

type Parser struct {
	logger log.Logger
}

func NewParser(logger log.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

func (p *Parser) Parser(metric *models.MetricValue) (successParseCollectResults []*models.CollectResult, err error) {
	collectResults := make([]*models.CollectResult, 0, 100)
	if err = json.Unmarshal(metric.Value, &collectResults); err != nil {
		return
	}

	successParseCollectResults = make([]*models.CollectResult, 0, len(collectResults))
	for _, collectResult := range collectResults {
		if collectResult.Err == "" {
			if err := p.parser(collectResult); err != nil {
				p.logger.Warnf("parse %s err:%s", collectResult.Metric, err.Error())
				continue
			}
		}

		successParseCollectResults = append(successParseCollectResults, collectResult)
	}
	return
}

func (p *Parser) parser(result *models.CollectResult) (err error) {
	result.RelValue, err = collector.UnmarshalResult(result.Metric, result.Value)
	return
}
