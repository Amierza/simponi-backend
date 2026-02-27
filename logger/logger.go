package logger

import (
	"os"

	"github.com/Amierza/simponi-backend/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// New creates and configures a production-ready zap logger
func New() (*zap.Logger, error) {
	isDev := true
	if os.Getenv("APP_ENV") == constants.ENUM_RUN_PRODUCTION {
		isDev = false
	}

	// Pastikan folder logs ada
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", os.ModePerm); err != nil {
			return nil, err
		}
	}

	// Konfigurasi encoder
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	if isDev {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	// Writer untuk info log (rotate otomatis)
	infoFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/info.log",
		MaxSize:    10, // dalam MB
		MaxBackups: 5,
		MaxAge:     30, // dalam hari
		Compress:   true,
	})

	// Writer untuk error log (rotate otomatis)
	errorFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})

	// Writer ke stdout (terminal)
	consoleWriter := zapcore.Lock(os.Stdout)

	// Gabungkan semua core
	core := zapcore.NewTee(
		// INFO & WARN ke info.log + stdout
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriter, consoleWriter),
			zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l < zapcore.ErrorLevel
			}),
		),

		// ERROR & FATAL ke error.log + stdout
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriter, consoleWriter),
			zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l >= zapcore.ErrorLevel
			}),
		),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger, nil
}
