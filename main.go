package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/bytedance/sonic"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

const (
	button        = `<a href="">点我试试</a>`
	buttonWithUrl = `<a href="%v">点我试试</a>`
	htmlTemp      = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>猜猜我在哪</title>
    <style type="text/css">
        body {
            text-align: center;
        }
        a {
            background-color: #5496ce; /* 是Vidar蓝! */
            border: none;
            color: white;
            padding: 15px 32px;
            text-align: center;
            text-decoration: none;
            display: inline-block;
            font-size: 16px;
        }
    </style>
</head>
<body>
    <h1>%s</h1>
    <p>红豆泥私密马赛，我忘记我把flag丢在哪一关了，下面有个按钮让你前往下一关，慢慢找叭~XD</p>
    %s
</body>
</html>
`
)

type Info struct {
	TeamId string `json:"team_id"`
	Level  int64  `json:"level"`
	Noise  int64  `json:"noise"`
}

func hello(ctx context.Context, req events.APIGatewayRequest) (resp events.APIGatewayResponse, err error) {
	var info Info
	if key, ok := req.QueryString["key"]; ok {
		Cipher, err := base64.StdEncoding.DecodeString(key[0])
		if err != nil {
			return Resp(genHtml("你这key有毒啊！", "https://hgame.vidar.club/", 1))
		}
		plaintext := AESDecrypt(Cipher, aseKey)

		err = sonic.UnmarshalString(string(plaintext), &info)
		if err != nil {
			return Resp(genHtml("你这key有毒啊！", "https://hgame.vidar.club/", 1))
		}
	} else {
		p := strings.Split(req.Path, "/")
		if len(p) > 1 && p[1] != "" {
			info = Info{
				TeamId: p[1],
				Level:  0,
				Noise:  time.Now().Unix(),
			}
		} else {
			return Resp(genHtml("请从比赛平台进来哦，不然我不晓得你所属的团队欸~", "https://hgame.vidar.club/", 1))
		}
	}

	info.Level++
	if info.Level >= 100 {
		s := sha256.Sum256([]byte(info.TeamId + flag))
		hash2 := hex.EncodeToString(s[:])
		return events.APIGatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      200,
			Headers: map[string]string{
				"Content-Type":     "text/html; charset=utf-8",
				"fI4g":             "hgame{" + hash2 + "}",
				"auth0r":           "asjdf",
				"Welcome-To-HGame": "See you next week!",
			},
			Body: genHtml("我好像在就是把flag落在这里了欸~ 快帮我找找x", "", 1),
		}, nil
	}
	info.Noise = time.Now().Unix()
	infoByte, _ := sonic.Marshal(&info)
	Cipher := AESEncrypt(infoByte, aseKey)
	nextUrl := "?key=" + url.QueryEscape(base64.StdEncoding.EncodeToString(Cipher))
	return Resp(genHtml(fmt.Sprintf("你现在在第%v关", info.Level), nextUrl, int(info.Level)))
}

func main() {
	cloudfunction.Start(hello)
}

func genHtml(title string, nextUrl string, buttonNum int) string {
	rand.Seed(time.Now().UnixNano())
	trueButtonPos := rand.Intn(buttonNum)

	var buttons string
	for i := 0; i < buttonNum; i++ {
		if i == trueButtonPos {
			buttons += fmt.Sprintf(buttonWithUrl, nextUrl)
		} else {
			buttons += button
		}
	}

	return fmt.Sprintf(htmlTemp, title, buttons)
}

func Resp(body string) (events.APIGatewayResponse, error) {
	return events.APIGatewayResponse{
		IsBase64Encoded: false,
		StatusCode:      200,
		Headers: map[string]string{
			"Content-Type": "text/html; charset=utf-8",
			"hint":         "赶紧爬吧！蛛蛛...嘿嘿...我的蛛蛛...",
		},
		Body: body,
	}, nil
}
