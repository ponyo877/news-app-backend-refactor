package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ponyo877/news-app-backend-refactor/entity"
)

type stubReportUseCase struct {
	err          error
	articleID    entity.ID
	url          string
	siteTitle    string
	deviceHash   string
	platform     string
	appVersion   string
	rulesVersion int
	reason       string
	calls        int
}

func (s *stubReportUseCase) ReportArticle(articleID entity.ID, url, siteTitle, deviceHash, platform, appVersion string, rulesVersion int, reason string) error {
	s.calls++
	s.articleID, s.url, s.siteTitle, s.deviceHash = articleID, url, siteTitle, deviceHash
	s.platform, s.appVersion, s.rulesVersion, s.reason = platform, appVersion, rulesVersion, reason
	return s.err
}

func performReport(service *stubReportUseCase, articleID, form string) *httptest.ResponseRecorder {
	e := echo.New()
	MakeReportHandlers(e, service)
	req := httptest.NewRequest(http.MethodPost, "/v1/article/report/"+articleID, strings.NewReader(form))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

const reportForm = "url=https%3A%2F%2Fexample.com%2Farchives%2F1.html&sitetitle=%E3%83%86%E3%82%B9%E3%83%88&devicehash=hash-1&platform=ios&appversion=1.52&rulesversion=4&reason=missing_media"

func TestReportArticle(t *testing.T) {
	articleID := entity.NewID()

	t.Run("フォーム値がそのままユースケースに渡り200", func(t *testing.T) {
		stub := &stubReportUseCase{}
		rec := performReport(stub, articleID.String(), reportForm)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d", rec.Code)
		}
		if stub.calls != 1 || stub.articleID.String() != articleID.String() {
			t.Fatalf("ユースケースが正しく呼ばれていません: %+v", stub)
		}
		if stub.url != "https://example.com/archives/1.html" || stub.siteTitle != "テスト" || stub.deviceHash != "hash-1" ||
			stub.platform != "ios" || stub.appVersion != "1.52" || stub.rulesVersion != 4 || stub.reason != "missing_media" {
			t.Fatalf("フォーム値が渡っていません: %+v", stub)
		}
	})

	t.Run("rulesversionが数値でなければ0", func(t *testing.T) {
		stub := &stubReportUseCase{}
		performReport(stub, articleID.String(), strings.Replace(reportForm, "rulesversion=4", "rulesversion=x", 1))
		if stub.rulesVersion != 0 {
			t.Fatalf("rulesVersion = %d", stub.rulesVersion)
		}
	})

	t.Run("article_idがUUIDでなければ400でユースケースを呼ばない", func(t *testing.T) {
		stub := &stubReportUseCase{}
		rec := performReport(stub, "not-a-uuid", reportForm)
		if rec.Code != http.StatusBadRequest || stub.calls != 0 {
			t.Fatalf("status = %d, calls = %d", rec.Code, stub.calls)
		}
	})

	t.Run("内容が不正なら400", func(t *testing.T) {
		rec := performReport(&stubReportUseCase{err: entity.ErrInvalidEntity}, articleID.String(), reportForm)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", rec.Code)
		}
	})

	t.Run("保存に失敗したら500(アプリ側で再試行を促す)", func(t *testing.T) {
		rec := performReport(&stubReportUseCase{err: errors.New("db down")}, articleID.String(), reportForm)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", rec.Code)
		}
	})
}
