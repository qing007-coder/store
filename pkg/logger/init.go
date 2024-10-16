package logger

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"store/pkg/constant"
	"store/pkg/kafka"
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
	l.producer.Publish(&sarama.ProducerMessage{
		Topic:  constant.LOGTOPIC,
		Offset: 0,
		Key:    sarama.StringEncoder("log"),
		Value:  sarama.ByteEncoder(data),
	})

	return len(data), nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}
