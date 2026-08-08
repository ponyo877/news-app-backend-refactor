package entity

// PushMessage プッシュ通知1件分の内容
type PushMessage struct {
	To    string
	Title string
	Body  string
	Data  map[string]string
}
