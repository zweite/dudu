package models

type MetricValue struct {
	Endpoint  string `json:"endpoint"`  // 节点IP
	HostName  string `json:"hostName"`  // 节点主机名
	Compactor string `json:"compactor"` // 压缩算法
	Value     []byte `json:"value"`     // 采集数据
	Tags      string `json:"tags"`      // 标签
	Timestamp int64  `json:"timestamp"` // 上传时间戳
}
