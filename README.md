# eth-indexer

层级	     技术	      用途
数据源	Infura RPC	链上区块和交易抓取
抓取层	Go + go-ethereum	数据采集与序列化
存储层	MySQL 8.0 + GORM	交易数据与指标存储
缓存层	Redis	断点续传 / 状态管理
分析层	SQL + Grafana + Python	数据分析与可视化
部署层	Docker Compose	环境一键启动
