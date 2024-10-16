package email

import (
	"article/pkg/config"
	"article/pkg/errors"
	"article/pkg/tools"
	"fmt"
	"gopkg.in/gomail.v2"
)

type Server struct {
	srv      *gomail.Message
	username string
	password string
	host     string
	port     int
}

func NewServer(conf *config.GlobalConfig) *Server {
	srv := new(Server)
	srv.init(conf)

	return srv
}

func (s *Server) init(conf *config.GlobalConfig) {
	s.srv = gomail.NewMessage()
	s.username = conf.Email.Username
	s.password = conf.Email.Password
	s.host = conf.Email.Host
	s.port = conf.Email.Port
}

func (s *Server) SendVerificationCode(target string) (string, error) {
	code := tools.RandomNumber(6)

	s.srv.SetHeader("From", fmt.Sprintf("xubo <%s>", s.username))
	s.srv.SetHeader("To", target)
	s.srv.SetHeader("Subject", "邮箱验证(测试)")
	s.srv.SetBody("text/plain", fmt.Sprintf("【注册验证】验证码：%s,有效期15分钟", code))
	d := gomail.NewDialer(s.host, s.port, s.username, s.password)

	if err := d.DialAndSend(s.srv); err != nil {
		return "", errors.EmailSendingFailed
	}

	return code, nil
}
