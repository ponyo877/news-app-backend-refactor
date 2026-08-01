// matome-site-rss-list.md からサイト登録用のSQLとアイコンを生成するシードツール。
//
//	go run ./scripts/seed-sites                 # ドライラン(アイコン確認+SQL生成)
//	go run ./scripts/seed-sites -upload         # アイコンを POST /v1/static へアップロードしてからSQL生成
//
// 出力: scripts/seed-sites/insert-sites.sql(冪等: 同titleが存在する場合はINSERTしない)
// 既存8サイト・更新停止3サイトは除外する。手動手順(news-app-docker/MEMO.md)の置き換え。
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

const defaultAPI = "https://matome.folks-chat.com"

// 更新停止中のため登録しない(mdの調査結果より。ユーザ判断で除外)
var excluded = map[string]bool{
	"稲妻速報":         true,
	"【2ch】コピペ情報局": true,
	"じゃぱそく!":       true,
}

// 既存DBに登録済みのサイト(md上の名前 → DBのtitle)。INSERT対象から除外する
var existing = map[string]string{
	"暇人\\(^o^)/速報": "暇人速報",
	"痛いニュース(ノ∀`)":  "痛いニュース",
	"哲学ニュースnwk":     "哲学ニュース",
	"ニュー速クオリティ":    "ニュー速クオリティ",
	"VIPPERな俺":      "VIPPERな俺",
}

type seedSite struct {
	Name    string
	BlogURL string
	RSSURL  string
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

const ua = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"

func main() {
	upload := flag.Bool("upload", false, "アイコンをAPIへアップロードする(未指定はドライラン)")
	api := flag.String("api", defaultAPI, "アップロード先APIのベースURL")
	mdPath := flag.String("md", "matome-site-rss-list.md", "サイト一覧mdのパス")
	outPath := flag.String("out", "scripts/seed-sites/insert-sites.sql", "生成するSQLの出力先")
	flag.Parse()

	sites, err := parseSiteList(*mdPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サイト一覧の読み込みに失敗: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("登録対象: %d件(除外・既存を差し引き済み)\n", len(sites))

	var sql strings.Builder
	sql.WriteString("-- seed-sites が生成したサイト登録SQL(冪等: 同titleが存在すればスキップ)\n")
	sql.WriteString("-- 実行: mysql -h 127.0.0.1 -P 3306 -u root -p -D matome < insert-sites.sql\n\n")
	for _, site := range sites {
		iconURL := resolveIcon(site)
		if iconURL == "" {
			fmt.Printf("⚠ %s: アイコンが見つからないためフォールバック画像を使用\n", site.Name)
			iconURL = *api + "/static/myimage_1.png"
		} else if *upload {
			uploaded, err := uploadIcon(*api, site, iconURL)
			if err != nil {
				fmt.Printf("⚠ %s: アイコンのアップロードに失敗(%v)。元URLを直接使用\n", site.Name, err)
			} else {
				iconURL = uploaded
			}
		}
		id := uuid.New().String()
		sql.WriteString(fmt.Sprintf(
			"INSERT INTO sites (id, title, rss_url, image_url)\n"+
				"SELECT '%s', '%s', '%s', '%s' FROM DUAL\n"+
				"WHERE NOT EXISTS (SELECT 1 FROM sites WHERE title = '%s');\n\n",
			id, escapeSQL(site.Name), escapeSQL(site.RSSURL), escapeSQL(iconURL), escapeSQL(site.Name)))
		fmt.Printf("✓ %s\n    rss:  %s\n    icon: %s\n", site.Name, site.RSSURL, iconURL)
	}
	if err := os.WriteFile(*outPath, []byte(sql.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "SQL出力に失敗: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nSQLを %s に出力しました\n", *outPath)
	if !*upload {
		fmt.Println("(ドライラン: -upload を付けるとアイコンを /v1/static へアップロードします)")
	}
}

// mdの表(| サイト名 | ブログURL | RSS URL | 状態 |)を読む
func parseSiteList(path string) ([]seedSite, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rowRe := regexp.MustCompile(`(?m)^\| (.+?) \| (https?://\S+) \| (https?://\S+) \| (.+?) \|$`)
	var sites []seedSite
	for _, m := range rowRe.FindAllStringSubmatch(string(raw), -1) {
		name := strings.TrimSpace(m[1])
		if name == "サイト名" || excluded[name] {
			continue
		}
		if _, ok := existing[name]; ok {
			continue
		}
		sites = append(sites, seedSite{Name: name, BlogURL: m[2], RSSURL: m[3]})
	}
	return sites, nil
}

// ブログトップからapple-touch-icon → icon → /favicon.ico の順でアイコンURLを見つける
func resolveIcon(site seedSite) string {
	doc, base, err := fetchDoc(site.BlogURL)
	if err == nil {
		selectors := []string{
			`link[rel="apple-touch-icon"]`,
			`link[rel="apple-touch-icon-precomposed"]`,
			`link[rel="icon"]`,
			`link[rel="shortcut icon"]`,
		}
		for _, sel := range selectors {
			if href, ok := doc.Find(sel).First().Attr("href"); ok && href != "" {
				if abs := absoluteURL(base, href); abs != "" && urlIsImage(abs) {
					return abs
				}
			}
		}
	}
	fallback := absoluteURL(site.BlogURL, "/favicon.ico")
	if urlIsImage(fallback) {
		return fallback
	}
	return ""
}

func fetchDoc(pageURL string) (*goquery.Document, string, error) {
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", ua)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d", res.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", err
	}
	return doc, res.Request.URL.String(), nil
}

func absoluteURL(base, href string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return baseURL.ResolveReference(ref).String()
}

// 画像として取得できるか(Content-Type検査込み)
func urlIsImage(imageURL string) bool {
	data, contentType, err := download(imageURL)
	return err == nil && len(data) > 0 && strings.HasPrefix(contentType, "image/")
}

func download(imageURL string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", ua)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d", res.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(res.Body, 5<<20))
	if err != nil {
		return nil, "", err
	}
	return data, res.Header.Get("Content-Type"), nil
}

// アイコンをダウンロードして POST /v1/static へアップロードし、公開URLを返す
func uploadIcon(api string, site seedSite, iconURL string) (string, error) {
	data, contentType, err := download(iconURL)
	if err != nil {
		return "", err
	}
	filename := iconFilename(site, contentType)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(data); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", api+"/v1/static", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", res.StatusCode)
	}
	// ダウンロードはnginx直配信の/staticを使う(/v1/staticはappを経由するため)
	return api + "/static/" + filename, nil
}

// ホスト名ベースの一意なASCIIファイル名(例: site-icon-blog-livedoor-jp-itsoku.png)
func iconFilename(site seedSite, contentType string) string {
	parsed, err := url.Parse(site.RSSURL)
	slug := "unknown"
	if err == nil {
		slug = strings.ReplaceAll(strings.TrimPrefix(parsed.Host, "www."), ".", "-")
		if parsed.Host == "blog.livedoor.jp" {
			parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
			if len(parts) > 0 {
				slug += "-" + parts[0]
			}
		}
	}
	ext := ".png"
	switch {
	case strings.Contains(contentType, "jpeg"):
		ext = ".jpg"
	case strings.Contains(contentType, "gif"):
		ext = ".gif"
	case strings.Contains(contentType, "x-icon"), strings.Contains(contentType, "vnd.microsoft.icon"):
		ext = ".ico"
	}
	return "site-icon-" + slug + ext
}

func escapeSQL(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), "'", `''`)
}
