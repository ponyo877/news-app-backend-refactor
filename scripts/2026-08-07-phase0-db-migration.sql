-- Phase 0 DB修理(docs/GROWTH-PLAN.md)— 2026-08-07
-- 実行方法: VM上で mysql -h 127.0.0.1 -P 3306 -u root -ppassword -D matome < このファイル
-- 実行前バックアップ(必須):
--   mysqldump -h127.0.0.1 -uroot -ppassword matome articles users comments sites > /tmp/phase0-backup.sql
--
-- 内容:
--   1. カラム長の拡張(title/url/image_url が VARCHAR(100) で、WordPress系サイトの
--      長いURL・長文タイトルがINSERT失敗 → クロールサイクル中断の主因になっていた)
--   2. 記事の重複排除(検索結果ベースで28%が重複。OnConflict句がUNIQUEキー不在で空振り)
--   3. UNIQUE(site_id, title) 追加(以後の重複を根絶し、OnConflictを機能させる)
--   4. users/comments のアバターURL修復(失効ドメイン matome-kun.ga が1,770件残存。
--      2026-08-01のスクリプトは sites/articles のみで users/comments が漏れていた)
--
-- 注意: 手順2で消えた記事IDがRedisのランキングZSETに残るが、
--       バックエンド側で欠損IDをスキップする修正済み(repository/articleMySQL.go List)。
--       本SQLは新バイナリのデプロイ後に実行すること。
-- 再実行について: DELETE/UPDATEは再実行可。ALTERは2回目に
--       「Duplicate key name」等のエラーになるが、その場合は該当行を飛ばしてよい。

-- ============================================================
-- 1. カラム長の拡張
-- ============================================================

ALTER TABLE articles
    MODIFY title     VARCHAR(255) NOT NULL,
    MODIFY url       VARCHAR(500) NOT NULL,
    MODIFY image_url VARCHAR(500) NOT NULL;

ALTER TABLE sites
    MODIFY rss_url   VARCHAR(500) NOT NULL,
    MODIFY image_url VARCHAR(500) NOT NULL;

ALTER TABLE users
    MODIFY image_url VARCHAR(255) NOT NULL;

ALTER TABLE comments
    MODIFY image_url VARCHAR(255) NOT NULL;

-- ============================================================
-- 2. 記事の重複排除
-- ============================================================
-- 同一(site_id, title)のグループごとに「created_atが最新(同値ならidが大きい)」の1件を残す。
-- 先に作業用インデックスを張らないと自己結合が総当たりになり1GB VMでは終わらない。

ALTER TABLE articles ADD INDEX idx_tmp_dedupe (site_id, title);

DELETE a FROM articles a
JOIN articles b
  ON  a.site_id = b.site_id
  AND a.title   = b.title
  AND (a.created_at < b.created_at OR (a.created_at = b.created_at AND a.id < b.id));

ALTER TABLE articles DROP INDEX idx_tmp_dedupe;

-- ============================================================
-- 3. UNIQUE制約の追加
-- ============================================================
-- 以後、再クロールされた同一記事は Create の OnConflict(updated_at/created_atの更新)に落ちる

ALTER TABLE articles ADD UNIQUE INDEX uq_site_id_title (site_id, title);

-- ============================================================
-- 4. アバターURLの修復
-- ============================================================
-- 旧ドメインを現行ドメインへ置換。/v1/static/ はnginx直配信の /static/ に寄せる
-- (GET /v1/static は長らく0バイトを返すバグがあり、直配信の方が確実 + CDNに載る)

UPDATE users
   SET image_url = REPLACE(image_url, 'matome-kun.ga/v1/static/', 'matome.folks-chat.com/static/')
 WHERE image_url LIKE '%matome-kun.ga/v1/static/%';

UPDATE users
   SET image_url = REPLACE(image_url, 'matome-kun.ga/', 'matome.folks-chat.com/')
 WHERE image_url LIKE '%matome-kun.ga/%';

UPDATE comments
   SET image_url = REPLACE(image_url, 'matome-kun.ga/v1/static/', 'matome.folks-chat.com/static/')
 WHERE image_url LIKE '%matome-kun.ga/v1/static/%';

UPDATE comments
   SET image_url = REPLACE(image_url, 'matome-kun.ga/', 'matome.folks-chat.com/')
 WHERE image_url LIKE '%matome-kun.ga/%';

-- ============================================================
-- 確認クエリ
-- ============================================================

SELECT COUNT(*) AS dup_groups_should_be_0
  FROM (SELECT site_id, title FROM articles GROUP BY site_id, title HAVING COUNT(*) > 1) t;

SELECT COUNT(*) AS dead_domain_users_should_be_0    FROM users    WHERE image_url LIKE '%matome-kun.ga%';
SELECT COUNT(*) AS dead_domain_comments_should_be_0 FROM comments WHERE image_url LIKE '%matome-kun.ga%';

SHOW INDEX FROM articles WHERE Key_name = 'uq_site_id_title';
