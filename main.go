package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var url string
	var msg string
	var notifyStatus bool
	url = os.Getenv("SLACK_URL")
	flag.StringVar(&msg, "msg", "", "알림 제목 입력")
	flag.BoolVar(&notifyStatus, "notifyStatus", false, "빌드결과 알림 여부")
	flag.Parse()

	if len(url) == 0 {
		log.Fatal("Slack url을 환경변수로 설정해주세요")
	}
	CheckUrl(url)

	slackSender := NewSlackSender(url, notifyStatus)
	slackSender.SetPayload(msg)
	slackSender.Send()
}

type SlackSender struct {
	Url          string
	Payload      Payload
	NotifyStatus bool
}

type Payload struct {
	Blocks []Block `json:"blocks"`
}

type Block interface{}

type Section struct {
	Type string `json:"type"`
	Text struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"text"`
}

func NewSlackSender(url string, notifyStatus bool) *SlackSender {
	return &SlackSender{
		Url:          url,
		NotifyStatus: notifyStatus,
	}
}

// 메시지의 각 세션
func NewSection(text string) Section {
	return Section{
		Type: "section",
		Text: struct {
			Type string "json:\"type\""
			Text string "json:\"text\""
		}{
			Type: "mrkdwn",
			Text: text,
		},
	}
}

// slack 알림 발송
func (s *SlackSender) Send() int {
	payload, err := json.Marshal(s.Payload)
	errorHandler(err)

	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Post(s.Url, "application/json", bytes.NewBuffer(payload))
	errorHandler(err)

	fmt.Println("Slack 알림 (" + resp.Status + ")")
	return resp.StatusCode
}

// slack 메시지 내용 만들기
func (s *SlackSender) SetPayload(msg string) *SlackSender {
	s.Payload.Blocks = append(s.Payload.Blocks, NewSection(fmt.Sprintf("*%s*", msg)))

	s.Payload.Blocks = append(s.Payload.Blocks, struct {
		Type string "json:\"type\""
	}{
		Type: "divider",
	})

	t := fmt.Sprintf("- ARN : %s\n- 소스 버전 : %s\n- 시간 : %s",
		os.Getenv("CODEBUILD_BUILD_ARN"), os.Getenv("CODEBUILD_SOURCE_VERSION"), time.Now().UTC().Format(time.RFC3339))
	if s.NotifyStatus {
		if os.Getenv("CODEBUILD_BUILD_SUCCEEDING") == "1" {
			t = t + "\n- 빌드 결과 : ✅ Success"
		} else if os.Getenv("CODEBUILD_BUILD_SUCCEEDING") == "0" {
			t = t + "\n- 빌드 결과 : ❌ Fail"
		} else {
			t = t + "\n- 빌드 결과 : Unknown"
		}
	}
	s.Payload.Blocks = append(s.Payload.Blocks, NewSection(t))
	return s
}

// slack url이 유효한지 검증
func CheckUrl(url string) {
	resp, err := http.Get(url)
	errorHandler(err)

	if resp.StatusCode != 400 {
		log.Fatal("url이 유효하지 않습니다.")
	}
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
