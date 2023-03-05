# qqbot-chatgpt
使用chatgpt-3.5turbo回复的qq机器人. 

详细请看https://ywt56xqv2i.feishu.cn/docx/ToIeds8H0oxPogxJzNAc8tDAnEf. 

云服务器部署
服务器
个人用的阿里云服务器，目前有活动可以白嫖
https://developer.aliyun.com/plan/student?spm=5176.21103406.J_3207526240.19.44c152cbUyHX9u&scm=20140722.M_4259405._.V_1
os
Alibaba Cloud Linux  3.2104 LTS 64位
宝塔
https://blog.csdn.net/gaojingsong/article/details/124153335
辅助文件上传、环境安装、文件编辑、解压缩、端口管理
yum install -y wget && wget -O install.sh http://download.bt.cn/install/install_6.0.sh && sh install.sh
修改面板用户名、密码
虚拟内存
https://blog.csdn.net/qq_40371220/article/details/125223115
cd /usr
mkdir swap
dd if=/dev/zero of=/usr/swap/swapfile bs=1M count=16384
swapon /usr/swap/swapfile
free -m
vim /etc/fstab
#/usr/swap/swapfile swap swap defaults 0 0
Go
go1.20.1.linux-arm64.tar.gz
https://studygolang.com/dl
配置环境变量
修改go proxy
go env -w GOPROXY=https://goproxy.cn,direct
ChatGPT
账号
登陆后在https://platform.openai.com/account/api-keys 获取api-key
Http接口（Go）
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

   // 通过梯子代理请求
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
   req.Header.Add("Authorization", "Bearer {api key}")

   // 请求过程日志
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

// 解析gpt响应
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
clash
参考http://www.zztongyun.com/article/clash%E6%80%8E%E4%B9%88%E9%85%8D%E7%BD%AE%E8%A7%84%E5%88%99
https://ikuuu.eu/user/tutorial安装好本地可用的梯子
配置文件
把配置文件导出iKuuu_V2.yaml
[图片]

clash安装
上传到服务器
暂时无法在飞书文档外展示此内容
暂时无法在飞书文档外展示此内容
gunzip clash-linux-amd64-v1.10.6.gz
解压后不是可执行文件则
chmod +x 解压后的文件名
iKuuu_V2.yaml放在解压后的相同目录
Web ui
tar xvJf yacd.tar.xz
解压后改名为ui
暂时无法在飞书文档外展示此内容
把iKuuu_V2.yaml、Country.mmdb和ui复制到~/.config/clash/
修改配置
#配置外网访问
external-controller: 0.0.0.0:9090
#ui登陆密码
secret: "123qwe" 
#ui文件夹路径
external-ui: "/home/clash/ui"
[图片]
启动
./clash -f iKuuu_V2.yaml
验证
curl --proxy localhost:7890 http://www.google.com
端口放开
服务器防火墙、安全组配置放开用到的端口。
QQbot
https://github.com/mamoe/mirai
mirai
暂时无法在飞书文档外展示此内容
根据提示安装JDK17
安装好后启动安装插件 ./mcl
https://github.com/mamoe/mirai/blob/dev/docs/ConsoleTerminal.md
./mcl --update-package net.mamoe:mirai-api-http --type plugin --channel stable-v2
./mcl --update-package net.mamoe:chat-command --type plugin --channel stable
登陆
准备qq小号，启动miral
mirai中登陆qq，可能需要验证，把验证连接打开然后F12获取ticket，填入mcl
/login qq号 密码 ANDROID_WATCH
[图片]
插件
基于https://mirai.mamoe.net/topic/1193/%E4%B8%80%E4%B8%AA%E7%AE%80%E6%98%93%E7%9A%84%E8%87%AA%E5%8A%A8%E5%9B%9E%E5%A4%8D%E6%8F%92%E4%BB%B6魔改
https://github.com/soundofautumn/mirai-autoreply
git clone https://github.com/soundofautumn/mirai-autoreply.git
[图片]
魔改AutoReply和ReplyManager
魔改（kotlin）
魔改onEnable()
    override fun onEnable() {
        Config.reload()
        println("启动auto replay")
        registerCommands()
        println("注册auto replay")
        //添加群组监听
        GlobalEventChannel
            .filterIsInstance<GroupMessageEvent>()
//            .filter { enabledGroups.contains(it.group.id) }
            .subscribeAlways<GroupMessageEvent> {
                println("监听私聊auto replay")
                val content = it.message.contentToString()
                if (content.length >= 5 && content.startsWith("/chat", false)) {
                    it.group.sendMessage(getResponse(it))
                }
            }
        //添加私聊监听
        GlobalEventChannel
            .subscribeAlways<UserMessageEvent> {
                println("监听私聊auto replay")
                val content = it.message.contentToString()
                if (content.length >= 5 && content.startsWith("/chat", false)) {
                    it.subject.sendMessage(getResponse(it))
                }
            }
    }

魔改getResponse
fun getResponse(event: MessageEvent): Message {
    val msg = event.message
    val content = msg.contentToString()

    val okHttpClient = OkHttpClient.Builder()
        .connectTimeout(60, TimeUnit.SECONDS)
        .writeTimeout(60, TimeUnit.SECONDS)
        .readTimeout(60, TimeUnit.SECONDS)
        .build()
    val url = "http://go获取chatgptapi的响应接口/chat?msg=" + content.subSequence(5, content.length);
    var respstr = url?.let { Request.Builder().url(it).get() }
        ?.let { it.header("Content-type", "text/plain") }
        ?.let { it.build() }
        ?.let { okHttpClient.newCall(it).execute() }.body?.string()

    println(respstr)
    val jsonObject = JSONObject(respstr)
    respstr = jsonObject["Msg"] as String?

    val chainBuilder = MessageChainBuilder()
    val response = respstr?.let { Response(it) }
    response?.let { addExtraMessage(event, it, chainBuilder) }
    respstr?.let { chainBuilder.add(it) }

    return chainBuilder.asMessageChain()
}
打包（gradle）
buildPlugin
auto-reply-1.3.2.mriai2.jar

[图片]
[图片]
使用
把auto-reply-1.3.2.mriai2.jar上传到mirai安装目录的plugin文件夹下，mcl所在的目录下./plugin
然后启动./mcl，插件会自动加载
部署
Clash 
nohup ./clash -f iKuuu_V2.yaml > nohup.log &
Http接口
把本地项目放到服务器上
然后到代码目录下go build main.go
nohup ./main > nohup.log &
Mirai
mcl是交互进程所以不能直接nohup挂起
sudo yum install screen
#开启新screen
screen -S mcl
#然后启动mirai
#退出screen ctrl + a ，d
#查看screen
screen -ls
#进入screen
screen -r id
完成后台执行mirai
其他
https://github.com/lss233/chatgpt-mirai-qq-bot
