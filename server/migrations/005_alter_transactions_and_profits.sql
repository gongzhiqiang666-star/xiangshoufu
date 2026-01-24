-- 享收付 - 扩展现有交易表和分润表

-- 交易表增加退款状态和扩展字段
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS refund_status SMALLINT DEFAULT 0;
COMMENT ON COLUMN transactions.refund_status IS '退款状态: 0-正常, 1-已退款';

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS channel_code VARCHAR(32);
COMMENT ON COLUMN transactions.channel_code IS '通道编码';

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS brand_code VARCHAR(32);
COMMENT ON COLUMN transactions.brand_code IS '品牌编号';

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS d0_fee BIGINT DEFAULT 0;
COMMENT ON COLUMN transactions.d0_fee IS 'D0手续费（分）';

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS high_rate DECIMAL(10,4);
COMMENT ON COLUMN transactions.high_rate IS '调价费率（%）';

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS card_no VARCHAR(32);
COMMENT ON COLUMN transactions.card_no IS '卡号（脱敏）';

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS ext_data JSONB;
COMMENT ON COLUMN transactions.ext_data IS '通道特有扩展字段';

-- 分润记录表增加撤销相关字段
ALTER TABLE profit_records ADD COLUMN IF NOT EXISTS is_revoked BOOLEAN DEFAULT FALSE;
COMMENT ON COLUMN profit_records.is_revoked IS '是否已撤销';

ALTER TABLE profit_records ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMPTZ;
COMMENT ON COLUMN profit_records.revoked_at IS '撤销时间';

ALTER TABLE profit_records ADD COLUMN IF NOT EXISTS revoke_reason VARCHAR(255);
COMMENT ON COLUMN profit_records.revoke_reason IS '撤销原因';

-- 创建退款状态索引
CREATE INDEX IF NOT EXISTS idx_transactions_refund ON transactions(refund_status) WHERE refund_status = 1;

-- 创建撤销状态索引
CREATE INDEX IF NOT EXISTS idx_profit_records_revoked ON profit_records(is_revoked) WHERE is_revoked = TRUE;
