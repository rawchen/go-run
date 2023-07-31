package main

import (
	"crypto/tls"
	"fmt"
	"github.com/robfig/cron/v3"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// 是否发送过邮件
var IsSendFlag bool

func main() {
	// 记录日志
	//fileName := "./logs/Info_First_" + carbon.Now().Format("Ymd") + ".log"
	//logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err == nil {
	//	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
	//} else {
	//	logrus.SetOutput(os.Stdout)
	//	logrus.Error(err)
	//}

	// @every 10s
	// @midnight
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		logrus.Info("cron exec")
		CheckServer()
	})
	c.Start()
	//c.Stop()
	select {}
}

func CheckServer() {
	// 请求状态
	reqStatus := true
	client := &http.Client{}

	if !IsSendFlag {
		for i := 0; i < 3; i++ {
			// 创建http.Request对象
			req, err := http.NewRequest("GET", "http://rnproxy.zhangjiashu.tech", nil)
			if err != nil {
				logrus.Error("error: ", err)
				time.Sleep(2 * time.Second)
				reqStatus = false
				continue
			}

			// 发送请求并获取响应结果
			resp, err := client.Do(req)
			if err != nil {
				logrus.Error("error: ", err)
				time.Sleep(2 * time.Second)
				reqStatus = false
				continue
			}

			if resp.StatusCode == 200 {
				reqStatus = true
			}
		}

		if !reqStatus {
			if !IsSendFlag {
				fmt.Println("send")
				//err := Send(remoteAddr + " site unreachable,please check the server")
				//if err != nil {
				//	logrus.Error("send email error", err)
				//}
				IsSendFlag = true
			}

		} else {
			// 成功则把失败次数重置
			IsSendFlag = false
			fmt.Println("connect success")
			//time.Sleep(5 * time.Minute)
		}
	}

}

func Send(content string) error {
	message := `
    <p> Hello %s,</p>
   
      <p style="text-indent:2em">%s</p>
   `

	host := "邮箱的smtp地址"
	port := 25
	userName := "你的邮箱"
	password := "邮箱smtp的密码"

	m := gomail.NewMessage()
	m.SetHeader("From", userName)                        // 发件人
	m.SetHeader("To", userName)                          // 收件人
	m.SetHeader("Subject", "Warning! Site unreachable!") // 邮件主题

	// text/html 的意思是将文件的 content-type 设置为 text/html 的形式，浏览器在获取到这种文件时会自动调用html的解析器对文件进行相应的处理。
	// 可以通过 text/html 处理文本格式进行特殊处理，如换行、缩进、加粗等等
	m.SetBody("text/html", fmt.Sprintf(message, "QiuYiEr", content))

	d := gomail.NewDialer(
		host,
		port,
		userName,
		password,
	)
	// 关闭SSL协议认证
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	logrus.Info("email sent")
	return nil
}
