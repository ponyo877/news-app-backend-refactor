-- 失効ドメイン matome-kun.ga を現行ドメインへ置換する(既存8サイトのアイコン)
-- 実行方法: VM上で mysql -h 127.0.0.1 -P 3306 -u root -p -D matome < このファイル
-- 実行前確認: SELECT id, title, image_url FROM sites;

UPDATE sites
SET image_url = REPLACE(image_url, 'https://matome-kun.ga/', 'https://matome.folks-chat.com/')
WHERE image_url LIKE 'https://matome-kun.ga/%';

-- 独自ドメインへ移行済みサイトのRSS URLを新ドメインへ更新
-- (旧livedoor URLも現在は機能しているが、matome-site-rss-list.md の調査結果に従い正規URLへ)
UPDATE sites SET rss_url = 'https://itainews.com/index.rdf'
WHERE id = '73637763-53c5-7726-420b-0fe2055ecc19' AND title = '痛いニュース';

UPDATE sites SET rss_url = 'https://nwknews.jp/index.rdf'
WHERE id = '76f8858b-45bc-08cb-3939-f89b85c08380' AND title = '哲学ニュース';

-- 確認: 置換漏れゼロであること
SELECT id, title, rss_url, image_url FROM sites;
