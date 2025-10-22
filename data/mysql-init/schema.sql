/*
  =========================================================
  Ethereum Data Indexer & Analytics Dashboard
  Database Schema for MySQL 8.0
  Author: Alex 
  Created: 2025-10-22
  ===========================================================
*/

-- ==========================================================
-- 1. 区块表（Blocks）
-- ==========================================================
CREATE TABLE IF NOT EXISTS blocks (
  number BIGINT PRIMARY KEY COMMENT '区块号',
  hash VARCHAR(66) NOT NULL UNIQUE COMMENT '区块哈希',
  parent_hash VARCHAR(66) DEFAULT NULL COMMENT '父区块哈希',
  timestamp DATETIME NOT NULL COMMENT '出块时间',
  miner VARCHAR(66) DEFAULT NULL COMMENT '出块矿工地址',
  gas_used BIGINT DEFAULT 0 COMMENT '已使用 Gas',
  gas_limit BIGINT DEFAULT 0 COMMENT 'Gas 上限',
  tx_count INT DEFAULT 0 COMMENT '区块内交易数量',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_timestamp (timestamp),
  INDEX idx_miner (miner)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='以太坊区块数据表';


-- ==========================================================
-- 2. 交易表（Transactions）
-- ==========================================================
CREATE TABLE IF NOT EXISTS transactions (
  hash VARCHAR(66) PRIMARY KEY COMMENT '交易哈希',
  block_number BIGINT NOT NULL COMMENT '所属区块号',
  from_address VARCHAR(66) DEFAULT NULL COMMENT '发送方地址',
  to_address VARCHAR(66) DEFAULT NULL COMMENT '接收方地址',
  value DECIMAL(38,0) DEFAULT 0 COMMENT '交易金额（Wei）',
  gas_price DECIMAL(38,0) DEFAULT 0 COMMENT '交易 Gas 价格（Wei）',
  input_data MEDIUMTEXT COMMENT '交易输入数据（十六进制）',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_block_number (block_number),
  INDEX idx_to_address (to_address),
  INDEX idx_from_address (from_address),
  CONSTRAINT fk_tx_block FOREIGN KEY (block_number) REFERENCES blocks(number)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='以太坊交易数据表';


-- ==========================================================
-- 3. 回执表（Receipts）
-- ==========================================================
CREATE TABLE IF NOT EXISTS receipts (
  tx_hash VARCHAR(66) PRIMARY KEY COMMENT '交易哈希',
  status BOOLEAN DEFAULT TRUE COMMENT '执行状态（1=成功,0=失败）',
  gas_used BIGINT DEFAULT 0 COMMENT '实际使用的 Gas',
  logs JSON NULL COMMENT '事件日志数据',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_receipt_tx FOREIGN KEY (tx_hash) REFERENCES transactions(hash)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易执行结果（Receipt）表';


-- ==========================================================
-- 4. 每日指标汇总表（Metrics Daily）
-- ==========================================================
CREATE TABLE IF NOT EXISTS metrics_daily (
  day DATE PRIMARY KEY COMMENT '日期（按天汇总）',
  tx_count BIGINT DEFAULT 0 COMMENT '交易总数',
  active_addresses BIGINT DEFAULT 0 COMMENT '活跃地址数',
  new_addresses BIGINT DEFAULT 0 COMMENT '新增地址数',
  avg_gas_price DECIMAL(38,4) DEFAULT 0 COMMENT '平均 Gas 价格（Wei）',
  avg_block_time DECIMAL(10,4) DEFAULT 0 COMMENT '平均出块时间（秒）',
  fail_rate DECIMAL(5,2) DEFAULT 0 COMMENT '交易失败率（%）',
  congestion_index DECIMAL(5,2) DEFAULT 0 COMMENT '网络拥堵指数（%）',
  network_health_score DECIMAL(5,2) DEFAULT 0 COMMENT '网络健康评分（0-100）',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='每日分析指标汇总表';

