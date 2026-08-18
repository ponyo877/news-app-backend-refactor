-- 差別的表現が中心になりやすい情報源を配信対象から外す。
--
-- 背景: Play のコンテンツレーティングを再申告するにあたり、IARC の
-- 「アプリに差別的な言葉は含まれていますか?」に「いいえ」と答えられる状態にするため、
-- 先に実態を変える。申告だけ「いいえ」にすると虚偽申告になる。
-- (news-app-frontend/docs/REMAINING_TASKS.md のレーティング再申告のセクションを参照)
--
-- sites だけ消すと articles の LEFT JOIN が NULL になり一覧が壊れるため、記事も一緒に削除する。
--
-- ⚠️ ID ではなく title で指定すること。
-- scripts/seed-sites/insert-sites.sql は「同titleが存在すればスキップ」する冪等INSERTなので、
-- ファイルに書かれたUUIDと実DBのidは一致しない(実測: キムチ速報はファイル
-- e1967c35-... に対しDBは 028c7d07-bde4-4fe2-8a0f-b99a813c684f だった)。
--
-- ⚠️ 実行前に必ずバックアップを取ること:
--   mysqldump -h127.0.0.1 -uroot -ppassword matome sites articles \
--     > /tmp/backup-before-site-removal-$(date +%Y%m%d).sql
--
-- 実行:
--   mysql -h 127.0.0.1 -P 3306 -u root -p -D matome < 2026-08-16-remove-discriminatory-sites.sql

START TRANSACTION;

-- 対象サイト(必要に応じて増減する)
--   キムチ速報 (kimsoku.com)
--
-- 検討候補(名称からは判断できないため、実際の記事を確認してから足すこと):
--   なんJ政治ネタまとめ (j-seiji.blog.jp)
--   モナニュース (mona-news.com)

DELETE FROM articles
WHERE site_id IN (SELECT id FROM sites WHERE title IN ('キムチ速報'));

DELETE FROM sites
WHERE title IN ('キムチ速報');

COMMIT;

-- 実行後の確認:
--   SELECT COUNT(*) FROM sites;                          -- 69 → 68
--   SELECT COUNT(*) FROM sites WHERE title = 'キムチ速報'; -- 0
