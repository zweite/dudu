### dudu 监控平台

主要分为几大模块

- agent 节点日志采集器、接收文件、执行命令
- proxy 数据中心代理节点，收集agent日志并发送给中央数据中心，接收命令并分发给管理的agent执行，接收文件并分发给管理的agent
- dashboard 监控面板，包含机器管理，命令分发，文件分发
