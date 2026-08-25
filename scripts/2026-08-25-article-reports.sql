-- 記事表示の不具合報告(アプリの記事画面⋮メニュー「表示の不具合を報告」)
-- 適用先: 本番 (matome DB)。冪等(IF NOT EXISTS)。MySQL 5.7 / 8.0 共通構文のみ(CHECK制約は使わない)
--
-- 設計:
--   - (device_hash, article_id) で1行。再報告は url/site_title/platform/app_version/rules_version/reason/updated_at
--     だけを上書きし、管理者が触る status/admin_note/status_updated_at と created_at は絶対に戻さない
--     (repository/articleReportMySQL.go の upsert 列に status 系を含めない)
--   - 対応後の再報告は updated_at > status_updated_at で検出できる
--   - articles へのFKは張らない(記事が削除されても報告は残す)。reason の検証は Go 側(entity)
--   - 運用手順・クエリ例: news-app-frontend/docs/ARTICLE-REPORTS.md

CREATE TABLE IF NOT EXISTS article_reports (
  id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  article_id        CHAR(36)      NOT NULL,                 -- アプリ内ID(articles.id)
  url               VARCHAR(2048) NOT NULL,                 -- 記事URL
  site_title        VARCHAR(100)  NOT NULL DEFAULT '',      -- 集計用に非正規化
  device_hash       VARCHAR(64)   NOT NULL,                 -- コメント/通知と同じ端末識別子
  platform          VARCHAR(10)   NOT NULL DEFAULT '',      -- ios / android
  app_version       VARCHAR(20)   NOT NULL DEFAULT '',      -- 1.52 など
  rules_version     INT           NOT NULL DEFAULT 0,       -- 報告時に有効だった scraper-rules の version
  reason            VARCHAR(32)   NOT NULL,                 -- missing_media / ad_remains / body_broken / other
  status            VARCHAR(16)   NOT NULL DEFAULT 'open',  -- open / in_progress / resolved / ignored(管理者がSQLで更新)
  admin_note        TEXT          NULL,                     -- 管理者メモ(原因・対応内容)
  status_updated_at DATETIME      NULL,                     -- status/admin_note を変えた日時
  updated_at        DATETIME      NOT NULL,                 -- 最終報告日時(再報告で更新)
  created_at        DATETIME      NOT NULL,                 -- 初回報告日時
  PRIMARY KEY (id),
  UNIQUE KEY uq_device_article (device_hash, article_id),
  KEY idx_status_updated (status, updated_at),
  KEY idx_article_id (article_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ---- 管理者向けクエリ例 ----
-- 未対応をサイト・理由別に集計
--   SELECT site_title, reason, COUNT(*) AS n, MAX(updated_at) AS latest
--   FROM article_reports WHERE status = 'open' GROUP BY site_title, reason ORDER BY n DESC;
-- 未対応の記事一覧(端末数・最終報告時刻)
--   SELECT article_id, site_title, url, COUNT(*) AS devices, GROUP_CONCAT(DISTINCT reason) AS reasons,
--          MAX(rules_version) AS rules_version, MAX(updated_at) AS last_reported_at
--   FROM article_reports WHERE status = 'open'
--   GROUP BY article_id, site_title, url ORDER BY last_reported_at DESC LIMIT 50;
-- 対応状態の変更(記事単位)
--   UPDATE article_reports SET status = 'resolved', admin_note = 'ルールv4でimgur復元', status_updated_at = NOW()
--   WHERE article_id = '<uuid>';
-- 対応後に再報告があったもの(修正が効いていない疑い)
--   SELECT id, article_id, reason, status, status_updated_at, updated_at
--   FROM article_reports WHERE status IN ('resolved', 'ignored') AND updated_at > status_updated_at;
