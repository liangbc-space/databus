package utils

import (
	"github.com/liangbc-space/databus/system"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"reflect"
)

const (
	MAX_SIZE       = 128 //日志切割：日志最大大小（M）
	MAX_BACKUPS    = 10  //日志切割：日志最大备份数
	LOG_VAILD_DAYS = 10  //日志切割：日志最大保留多少天
)

type LoggerCfg struct {
	Level           zapcore.Level
	Hook            lumberjack.Logger //lumberjack日志切割器
	WithCaller      bool              //是否开启堆栈跟踪
	OutputToConsole bool              //输出到控制台
}

func NewDefaultLogger() *zap.Logger {
	defaultCfg := system.ApplicationCfg.LoggerConfig

	cfg := new(LoggerCfg)
	switch defaultCfg.Level {
	case "debug":
		cfg.Level = zap.DebugLevel
	case "info":
		cfg.Level = zap.InfoLevel
	case "warn":
		cfg.Level = zap.WarnLevel
	case "error":
		cfg.Level = zap.ErrorLevel
	default:
		if system.ApplicationCfg.Debug {
			cfg.Level = zap.InfoLevel
		} else {
			cfg.Level = zap.WarnLevel
		}
	}

	cfg.Hook = lumberjack.Logger{
		Filename: defaultCfg.LogPath,
		MaxAge:   int(defaultCfg.LogValidDays),
	}
	if system.ApplicationCfg.Debug {
		cfg.WithCaller = true
	}

	return cfg.NewLogger()
}

func (cfg LoggerCfg) NewLogger() *zap.Logger {
	cfg = cfg.initCfg()

	writer := make([]zapcore.WriteSyncer, 0)
	if cfg.OutputToConsole {
		writer = append(writer, zapcore.AddSync(os.Stdout))
	}

	if !reflect.DeepEqual(cfg.Hook, lumberjack.Logger{}) {
		writer = append(writer, zapcore.AddSync(&cfg.Hook))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(logEncoder()),
		zapcore.NewMultiWriteSyncer(writer...),
		zap.NewAtomicLevelAt(cfg.Level),
	)

	options := []zap.Option{
		// 设置初始化字段
		zap.Fields(zap.String("service_name", system.ApplicationCfg.AppName)),
		zap.AddStacktrace(zap.ErrorLevel),
	}

	if cfg.WithCaller {
		options = append(options, zap.AddCaller())
		options = append(options, zap.Development())
	}

	return zap.New(core, options...)
}

func (cfg LoggerCfg) initCfg() LoggerCfg {
	if !reflect.DeepEqual(cfg.Hook, lumberjack.Logger{}) {
		if cfg.Hook.MaxSize <= 0 {
			cfg.Hook.MaxSize = MAX_SIZE
		}

		if cfg.Hook.MaxBackups <= 0 {
			cfg.Hook.MaxBackups = MAX_BACKUPS
		}

		if cfg.Hook.MaxAge <= 0 {
			if system.ApplicationCfg.LoggerConfig.LogValidDays > 0 {
				cfg.Hook.MaxAge = int(system.ApplicationCfg.LoggerConfig.LogValidDays)
			} else {
				cfg.Hook.MaxAge = LOG_VAILD_DAYS
			}
		}
	}

	return cfg
}

func logEncoder() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	return encoderConfig
}
