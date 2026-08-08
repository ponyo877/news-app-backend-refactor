-- Phase 1: プッシュ通知の送信先テーブル
-- 適用先: OCI本番 (matome DB)。冪等(IF NOT EXISTS)
--
-- expo_token を主キーとする自然キー設計:
--   - 同一端末のトークンローテーションで古い行が残っても、送信時の
--     DeviceNotRegistered 検出でアプリ側が消すため実害がない
--   - device_hash はユーザー突合・デバッグ用(コメント機能と同じ端末識別子)

CREATE TABLE IF NOT EXISTS device_tokens (
  expo_token     VARCHAR(128) NOT NULL,
  device_hash    VARCHAR(64)  NOT NULL,
  platform       VARCHAR(10)  NOT NULL DEFAULT '',
  digest_enabled TINYINT(1)   NOT NULL DEFAULT 1,
  updated_at     DATETIME     NOT NULL,
  created_at     DATETIME     NOT NULL,
  PRIMARY KEY (expo_token),
  KEY idx_device_hash (device_hash),
  KEY idx_digest_enabled (digest_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
