package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"store/pkg/errors"
)

func main() {
	//conf := sarama.NewConfig()
	//conf.Net.MaxOpenRequests = 2
	//setting := kafka.Setting{
	//	Addr: []string{"10.3.0.42:9092"},
	//	Conf: conf,
	//	SuccessHandler: func(msg *sarama.ProducerMessage) {
	//		fmt.Println("success:", msg.Value)
	//	},
	//	ErrorHandler: func(err error) {
	//		fmt.Println("err:", err)
	//	},
	//}
	//p, err := kafka.NewProducer(&setting)
	//if err != nil {
	//	fmt.Println("err:", err)
	//	return
	//}
	//
	//l := logger.NewLogger(p)
	//l.Info("id:12,action:access", "gateway")
	//time.Sleep(time.Second)
	//l.Error("sadfa", "aaaaa")
	//fmt.Println("success")
	//select {}

	//conf, err := config.NewGlobalConfig()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(conf)

	//token := "Bearer kjkjk"
	//fmt.Println(token[7:])

	//cipher, err := tools.Encrypt([]byte("role"), "0123456789012345")
	//if err != nil {
	//	errors.HandleError(err)
	//	return
	//}
	//fmt.Println(cipher)
	//
	//data, err := tools.Decrypt(cipher, "0123456789012345")
	//if err != nil {
	//	errors.HandleError(err)
	//	return
	//}
	//
	//fmt.Println(string(data))

	//conf, err := config.NewGlobalConfig()
	//if err != nil {
	//	fmt.Println("err:", err)
	//	return
	//}
	//rdb := redis.NewClient(conf)
	//data, err := rdb.Get(context.Background(), "1851894943077896192")
	//fmt.Println(data)
	//fmt.Println(err != nil)

	//srv, err := mock.NewServer()
	//if err != nil {
	//	errors.HandleError(err)
	//	return
	//}

	//if err := srv.Run(); err != nil {
	//	errors.HandleError(err)
	//	return
	//}

	//if err := srv.Clear(); err != nil {
	//	errors.HandleError(err)
	//	return
	//}

	r := gin.Default()
	r.POST("", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			errors.HandleError(err)
			return
		}

		fmt.Println(form.File["f"])
	})

	r.Run(":8080")
}
