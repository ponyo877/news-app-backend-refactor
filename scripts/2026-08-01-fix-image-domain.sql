-- 既存8サイトのアイコンと記事フォールバック画像を修復する
-- 実行方法: VM上で mysql -h 127.0.0.1 -P 3306 -u root -ppassword -D matome < このファイル
-- 前提調査(2026-08-01):
--   * 旧image_url(https://matome-kun.ga/v1/static/WaraNote.jpg 等)はドメイン失効かつ
--     ファイル自体が存在しない。実体は nginx直配信の /static/{DNT,HSK,...}.png(8件確認済み)
--   * /v1/static/<file> はアプリのDownload未実装により0バイトを返すため使わない
-- 実行前バックアップ: mysqldump -h127.0.0.1 -uroot -ppassword matome sites > /tmp/sites-backup.sql

UPDATE sites SET image_url = 'https://matome.folks-chat.com/static/WNT.png' WHERE id = '0e0f53b6-3d5c-f938-a65c-b5ef7dd89bf2' AND title = 'ワラノート';
UPDATE sites SET image_url = 'https://matome.folks-chat.com/static/HSK.png' WHERE id = '3d78d342-dfce-3ee8-7a11-12626c2852a8' AND title = '暇人速報';
UPDATE sites SET image_url = 'https://matome.folks-chat.com/static/VOR.png' WHERE id = '4f15c4b0-ee39-c322-df3e-6ca2eeef2027' AND title = 'VIPPERな俺';
UPDATE sites SET image_url = 'https://matome.folks-chat.com/static/DNT.png' WHERE id = '5fb5d46d-c44e-502e-ab36-350609900e1e' AND title = 'デジタルニューススレッド';
UPDATE sites SET image_url = 'https://matome.folks-chat.com/static/ITI.png' WHERE id = '73637763-53c5-7726-420b-0fe2055ecc19' AND title = '痛いニュース';
UPDATE sites SET image_url = 'https://matome.folks-chat.com/static/TGK.png' WHERE id = '76f8858b-45bc-08cb-3939-f89b85c08380' AND title = '哲学ニュース';
UPDATE sites SET image_url = 'https://matome.folks-chat.com/static/ISK.png' WHERE id = '7b9b58fb-4246-91b6-374b-0aff061b5c15' AND title = '稲妻速報';
UPDATE sites SET image_url = 'https://matome.folks-chat.com/static/NSK.png' WHERE id = '8bb0e13e-f574-8fb8-4d57-a7761c195c17' AND title = 'ニュー速クオリティ';

-- 独自ドメインへ移行済みサイトのRSS URLを新ドメインへ更新
UPDATE sites SET rss_url = 'https://itainews.com/index.rdf'
WHERE id = '73637763-53c5-7726-420b-0fe2055ecc19' AND title = '痛いニュース';
UPDATE sites SET rss_url = 'https://nwknews.jp/index.rdf'
WHERE id = '76f8858b-45bc-08cb-3939-f89b85c08380' AND title = '哲学ニュース';

-- 記事のフォールバック画像(https://matome-kun.ga/static/myimage_N.png)のドメイン置換
UPDATE articles SET image_url = REPLACE(image_url, 'https://matome-kun.ga/', 'https://matome.folks-chat.com/')
WHERE image_url LIKE 'https://matome-kun.ga/%';

-- 確認
SELECT id, title, rss_url, image_url FROM sites;
SELECT COUNT(*) AS remaining_dead_articles FROM articles WHERE image_url LIKE '%matome-kun.ga%';
