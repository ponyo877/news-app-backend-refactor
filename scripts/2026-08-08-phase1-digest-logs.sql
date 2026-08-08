-- Phase 1: ダイジェスト送信履歴
-- 適用先: OCI本番 (matome DB)。冪等(IF NOT EXISTS)
--
-- 用途:
--   1. 同一記事の連続送信防止(直近48時間に送った記事を次回の選定から除外)
--   2. CTR計測の分母(sent_count と GA4 の notification_open を突き合わせる)

CREATE TABLE IF NOT EXISTS digest_logs (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  article_id CHAR(36)        NOT NULL,
  sent_count INT             NOT NULL DEFAULT 0,
  created_at DATETIME        NOT NULL,
  PRIMARY KEY (id),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
