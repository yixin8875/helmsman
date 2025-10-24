-- 用户表：存储系统用户信息
CREATE TABLE users (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,        -- 用户唯一ID
                       username TEXT NOT NULL UNIQUE,               -- 用户名（唯一）
                       password_hash TEXT NOT NULL,                 -- 密码哈希值
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 用户创建时间
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- 用户最后更新时间
);

-- 交易账户表：用户关联的交易账户
CREATE TABLE accounts (
                          id INTEGER PRIMARY KEY AUTOINCREMENT,        -- 账户唯一ID
                          user_id INTEGER NOT NULL,                    -- 关联的用户ID
                          name TEXT NOT NULL,                          -- 账户名称
                          initial_balance REAL NOT NULL,               -- 初始余额
                          currency TEXT DEFAULT 'USD',                 -- 账户货币类型
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 账户创建时间
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- 账户最后更新时间
);

-- 交易策略表：用户定义的交易策略
CREATE TABLE strategies (
                            id INTEGER PRIMARY KEY AUTOINCREMENT,        -- 策略唯一ID
                            user_id INTEGER NOT NULL,                    -- 关联的用户ID
                            name TEXT NOT NULL,                          -- 策略名称
                            description TEXT,                            -- 策略详细描述
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 策略创建时间
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 策略最后更新时间
                            UNIQUE(user_id, name)                       -- 确保用户下策略名称唯一
);

-- 交易日志表：记录完整的交易生命周期
CREATE TABLE trades (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,        -- 交易唯一ID
                        account_id INTEGER NOT NULL,                -- 关联的账户ID
                        strategy_id INTEGER,                        -- 关联的策略ID

    -- 计划阶段字段
                        status TEXT NOT NULL DEFAULT 'planned',     -- 交易状态：planned/active/closed
                        symbol TEXT NOT NULL,                       -- 交易品种（如BTC/USD）
                        direction TEXT NOT NULL,                    -- 交易方向：long/short
                        planned_entry_price REAL,                   -- 计划入场价格
                        planned_stop_loss REAL,                     -- 计划止损价格
                        planned_take_profit REAL,                   -- 计划止盈价格
                        position_size REAL,                         -- 计划持仓大小
                        planned_risk_amount REAL,                   -- 计划风险金额
                        plan_notes TEXT,                            -- 交易计划备注

    -- 执行阶段字段
                        actual_entry_time TIMESTAMP,                -- 实际入场时间
                        actual_entry_price REAL,                    -- 实际入场价格
                        actual_exit_time TIMESTAMP,                 -- 实际出场时间
                        actual_exit_price REAL,                     -- 实际出场价格
                        commission REAL,                            -- 交易佣金费用

    -- 结果与复盘字段
                        pnl REAL,                                   -- 盈亏金额（Profit and Loss）
                        r_multiple REAL,                            -- 风险回报倍数
                        exit_reason TEXT,                           -- 出场原因
                        execution_score INTEGER,                    -- 执行评分（1-5分）
                        reflection_notes TEXT,                      -- 交易反思笔记

                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 交易创建时间
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- 交易最后更新时间
);

-- 交易快照表：存储交易相关的图表截图
CREATE TABLE snapshots (
                           id INTEGER PRIMARY KEY AUTOINCREMENT,        -- 快照唯一ID
                           trade_id INTEGER NOT NULL,                   -- 关联的交易ID
                           type TEXT NOT NULL,                          -- 快照类型：pre_trade/post_trade
                           image_url TEXT NOT NULL,                     -- 图片存储路径或URL
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 快照创建时间
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- 快照最后更新时间
);

-- 标签表：用于交易分类的标签系统
CREATE TABLE tags (
                      id INTEGER PRIMARY KEY AUTOINCREMENT,        -- 标签唯一ID
                      user_id INTEGER NOT NULL,                    -- 关联的用户ID
                      name TEXT NOT NULL,                          -- 标签名称
                      color TEXT,                                  -- 标签颜色（HEX格式）
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 标签创建时间
                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 标签最后更新时间
                      UNIQUE(user_id, name)                       -- 确保用户下标签名称唯一
);

-- 交易标签关联表：建立交易与标签的多对多关系
CREATE TABLE trade_tags (
                            trade_id INTEGER NOT NULL,                   -- 关联的交易ID
                            tag_id INTEGER NOT NULL,                     -- 关联的标签ID
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 关联创建时间
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 关联最后更新时间
                            PRIMARY KEY (trade_id, tag_id)               -- 复合主键确保唯一关联
);