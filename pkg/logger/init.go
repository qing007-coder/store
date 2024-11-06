package logger

import (
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"store/pkg/constant/rules"
	"store/pkg/kafka"
	"store/pkg/tools"
	"time"
)

type Logger struct {
	producer *kafka.Producer
	logger   *zap.Logger
}

func NewLogger(p *kafka.Producer) *Logger {
	l := new(Logger)
	l.producer = p
	l.init()
	return l
}

func (l *Logger) init() {
	logMode := zapcore.DebugLevel

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(l),
		logMode,
	)
	l.logger = zap.New(core)
}

func (l *Logger) Write(data []byte) (n int, err error) {
	fmt.Println(string(data))
	l.producer.Publish(&sarama.ProducerMessage{
		Topic:  rules.LOGTOPIC,
		Offset: 0,
		Key:    sarama.StringEncoder("log"),
		Value:  sarama.ByteEncoder(data),
	})

	return len(data), nil
}

func (l *Logger) Info(message, source string) {
	l.logger.Info(message,
		zap.Int64("time", time.Now().Unix()),
		zap.String("source", source),
		zap.String("id", tools.CreateID()),
	)
}

func (l *Logger) Error(message, source string) {
	l.logger.Error(message,
		zap.Int64("time", time.Now().Unix()),
		zap.String("source", source),
		zap.String("id", tools.CreateID()),
	)
}
