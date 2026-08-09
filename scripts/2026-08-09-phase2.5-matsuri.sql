-- Phase 2.5: 祭り速報通知のオプトアウト用カラム
-- 適用先: OCI本番 (matome DB)。冪等

SET @col_exists = (
  SELECT COUNT(*) FROM information_schema.columns
  WHERE table_schema = DATABASE() AND table_name = 'device_tokens' AND column_name = 'matsuri_enabled'
);
SET @ddl = IF(@col_exists = 0,
  'ALTER TABLE device_tokens ADD COLUMN matsuri_enabled TINYINT(1) NOT NULL DEFAULT 1 AFTER digest_enabled',
  'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
