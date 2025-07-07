package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
)

const (
	webhookURL = "asfafdasf" // this should be your real webhoolURL
)

func sendToDiscord(message string) {
	data := map[string]string{
		"content": "!!전설 출현 알림!!\n" + message,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON Parse Error:", err)
		return
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Creating Request Error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Transmit to WebHook Url Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		fmt.Println("WebHook fail status code:", resp.StatusCode)
	}

}

func checkMsg(log string) (string, error) {
	re := regexp.MustCompile(`\[Pixelmon\].*has spawned`)

	if re.MatchString(log) {
		parts := strings.SplitN(log, "]:", 2)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", nil

}

func followLog(filePath string) {
	t, err := tail.TailFile(filePath, tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
	})

	if err != nil {
		fmt.Println("Fail to Open File:", err)
		return
	}

	for line := range t.Lines {
		if line.Err != nil {
			fmt.Println("Fail to Read file:", line.Err)
			continue
		}

		text := line.Text
		// fmt.Println("읽은 로그:", text)

		parsing_log, err := checkMsg(text)
		if err != nil {
			fmt.Println("Parsing error:", err)
		}
		if parsing_log != "" {
			fmt.Println("find log!:", text)
			sendToDiscord(parsing_log)
		}
	}

}

// func parsingChat(msg string) string {
// 	idx := strings.Index(msg, "[CHAT]")
// 	if idx == -1 {
// 		return ""
// 	}
// 	return strings.TrimSpace(msg[idx+6:])
// }

func main() {
	logPath := "/app/logs/latest.log" //This should be your real log path

	fmt.Println("chat log service is start..")
	followLog(logPath)
}
