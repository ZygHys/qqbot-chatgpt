package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func main() {
	r := gin.Default()
	r.GET("/chat", chat)
	r.Run(":801") // port
}

var limit = 0
var lastUpdate = time.Now()

// TODO 实现连续对话
func chat(c *gin.Context) {
	// 简陋限流器
	if time.Since(lastUpdate) > 60*time.Second {
		limit = 0
	}
	if limit > 30 {
		c.JSON(http.StatusOK, "等一会......限流了")
	}

	limit++
	c.JSON(http.StatusOK, chatSend(c.Query("msg")))
}

type MyResp struct {
	Msg string
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func chatSend(msg string) MyResp {
	myresp := MyResp{}
	id := r.Intn(int(^uint32(0) >> 1))
	if strings.Compare(msg, "") == 0 {
		myresp.Msg = "null"
		return myresp
	}

	fmt.Println(strconv.Itoa(id) + msg)

	urlstr := "https://api.openai.com/v1/chat/completions"
	method := "POST"

	payload := strings.NewReader(`{
    "model": "gpt-3.5-turbo",
    "messages":[
        {"role": "system", "content": "你是无所不知的高情商的私人智能助手。"},
        {"role": "user", "content": "` + msg + `"}
    ]
}`)

	proxy := "http://127.0.0.1:7890/"
	proxyAddress, _ := url.Parse(proxy)

	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyAddress),
		},
	}
	req, err := http.NewRequest(method, urlstr, payload)

	if err != nil {
		fmt.Println(err)
		myresp.Msg = "build request err"
		return myresp
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Chrome/75.0.3770.142")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Authorization", "Bearer {your api-key}")

	send := true
	go func() {
		i := 0
		for send {
			time.Sleep(5 * time.Second)
			i += 5
			println(strconv.Itoa(id) + " wait " + strconv.Itoa(i) + "s")
		}
	}()
	res, err := client.Do(req)
	send = false

	if err != nil {
		fmt.Println(err)
		myresp.Msg = "send err"
		return myresp
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		myresp.Msg = "io read err"
		return myresp
	}

	resp := Resp{}
	if err := json.Unmarshal(body, &resp); err != nil {
		myresp.Msg = "json err"
		return myresp
	}

	content := ""
	for i := range resp.Choices {
		i := i
		content += resp.Choices[i].Message.Content
	}
	fmt.Println(strconv.Itoa(id) + ":" + content)
	myresp.Msg = content
	return myresp
}

type Resp struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Usage   Usage     `json:"usage"`
	Choices []Choices `json:"choices"`
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
type Choices struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}
