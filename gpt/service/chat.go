package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	GPT3_5 = `"gpt-3.5-turbo"`
	GPT4   = "gpt-4"
)

type Resp struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Usage   Usage     `json:"usage"`
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MyResp struct {
	Msg string
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func ChatSend(msg, model, apiKey string) MyResp {
	var myResp MyResp
	logId := r.Intn(int(^uint32(0) >> 1))
	if strings.Compare(msg, "") == 0 {
		myResp.Msg = "null, logId: " + strconv.Itoa(logId)
		return myResp
	}

	fmt.Println(strconv.Itoa(logId) + msg)

	urlstr := "https://api.openai.com/v1/chat/completions"
	method := "POST"

	bd := fmt.Sprint(`{
    "model": `, model, `,
    "messages":[
        {"role": "system", "content": "你是无所不知的高情商的私人智能助手。"},
        {"role": "user", "content": "`+msg+`"}
    ]
}`)

	payload := strings.NewReader(bd)

	proxy := "http://127.0.0.1:7890/"
	proxyAddress, _ := url.Parse(proxy)

	client := &http.Client{
		Timeout: 120 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyAddress),
		},
	}
	req, err := http.NewRequest(method, urlstr, payload)

	if err != nil {
		fmt.Println(err)
		myResp.Msg = "build request err"
		return myResp
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", apiKey))

	send := true
	go func() {
		i := 0
		for send {
			time.Sleep(5 * time.Second)
			i += 5
			println(strconv.Itoa(logId) + " wait " + strconv.Itoa(i) + "s")
		}
	}()
	res, err := client.Do(req)
	send = false

	if err != nil {
		fmt.Println(err)
		myResp.Msg = "send err, logId: " + strconv.Itoa(logId)
		return myResp
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		myResp.Msg = "io read err, logId: " + strconv.Itoa(logId)
		return myResp
	}

	resp := Resp{}
	if err := json.Unmarshal(body, &resp); err != nil {
		bodyString := ""
		json.Unmarshal(body, &bodyString)
		myResp.Msg = "json err, logId: " + strconv.Itoa(logId) + " " + string(body)
		return myResp
	}

	content := ""
	for i := range resp.Choices {
		i := i
		content += resp.Choices[i].Message.Content
	}
	fmt.Println(strconv.Itoa(logId) + ":" + content)
	myResp.Msg = content
	return myResp
}
