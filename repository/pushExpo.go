package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
)

// PushExpo Expo Push Service(exp.host)経由でプッシュ通知を送る。
// APNs/FCMの鍵管理はEASに委譲しているため、バックエンドはHTTP POSTのみ
type PushExpo struct {
	endpoint string
	client   *http.Client
}

// NewPushExpo create new repository
func NewPushExpo() *PushExpo {
	return &PushExpo{
		endpoint: "https://exp.host/--/api/v2/push/send",
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

type expoPushRequest struct {
	To    string            `json:"to"`
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
	Sound string            `json:"sound"`
}

type expoPushTicket struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Details struct {
		Error string `json:"error"`
	} `json:"details"`
}

type expoPushResponse struct {
	Data []expoPushTicket `json:"data"`
}

// Expo Push APIは1リクエスト最大100件
const expoPushChunkSize = 100

// Push 通知を送信し、失効していたトークン(DeviceNotRegistered)の一覧を返す
func (p *PushExpo) Push(messages []entity.PushMessage) ([]string, error) {
	var invalidTokens []string
	for start := 0; start < len(messages); start += expoPushChunkSize {
		end := min(start+expoPushChunkSize, len(messages))
		chunk := messages[start:end]
		tickets, err := p.pushChunk(chunk)
		if err != nil {
			// 一部チャンクの失敗で全体を止めない(次回のダイジェストで再送される)
			log.Warnf("Expo Pushのチャンク送信に失敗しました(%d-%d件目): %v", start, end, err)
			continue
		}
		// チケットはリクエストの並び順に対応する
		for i, ticket := range tickets {
			if ticket.Status == "ok" {
				continue
			}
			if ticket.Details.Error == "DeviceNotRegistered" {
				invalidTokens = append(invalidTokens, chunk[i].To)
				continue
			}
			log.Warnf("Expo Pushチケットエラー(token=%s): %s %s", chunk[i].To, ticket.Details.Error, ticket.Message)
		}
	}
	return invalidTokens, nil
}

func (p *PushExpo) pushChunk(messages []entity.PushMessage) ([]expoPushTicket, error) {
	requests := make([]expoPushRequest, 0, len(messages))
	for _, message := range messages {
		requests = append(requests, expoPushRequest{
			To:    message.To,
			Title: message.Title,
			Body:  message.Body,
			Data:  message.Data,
			Sound: "default",
		})
	}
	requestBody, err := json.Marshal(requests)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, p.endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	response, err := p.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expo PushがHTTP %dを返しました", response.StatusCode)
	}
	var pushResponse expoPushResponse
	if err := json.NewDecoder(response.Body).Decode(&pushResponse); err != nil {
		return nil, err
	}
	if len(pushResponse.Data) != len(messages) {
		return nil, fmt.Errorf("Expo Pushのチケット数(%d)が送信数(%d)と一致しません", len(pushResponse.Data), len(messages))
	}
	return pushResponse.Data, nil
}
