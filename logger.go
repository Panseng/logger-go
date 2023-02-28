package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Panseng/logger-go/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger

var loggerMutex sync.RWMutex

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// project name 格式: project_name 不包含 log 等后缀, 纯项目名称
func GetLogger(projectName string) *zap.SugaredLogger {
	loggerMutex.Lock()

	defer loggerMutex.Unlock()
	prestr := strings.ToUpper(projectName)
	if Logger == nil {
		logs := &Loggers{}

		// 获取 基础目录
		baseDir := os.Getenv(fmt.Sprintf("%s_PATH", prestr))
		fmt.Printf("Project(%s) base dir: %s.\n", projectName, baseDir)

		envConfigKeyLog = fmt.Sprintf("%s_log", projectName)

		if baseDir == "" {
			baseDir = utils.DefaultFolder
		}
		baseDir, _ = filepath.Abs(baseDir)
		logs.SetBaseDir(baseDir)
		cfgPath := filepath.Join(baseDir, utils.DefaultLoggerFile)
		cfgExists := false
		if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
			cfgExists = true
		}

		// 是否存在
		if cfgExists {
			// 存在
			logs.LoadJSONFileAndEnv(cfgPath)
		} else {
			// 不存在
			_ = logs.Default()
			err := logs.ApplyEnvVars()
			if err != nil {
				fmt.Printf("logger apply env vars get error: %s", err)
				os.Exit(1)
			}
		}
		initLogger(logs)
		logs.SaveLoggerToDisk(cfgPath)
	}
	return Logger
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// level: debug info warn error dpanic panic fatal
func initLogger(cfg *Loggers) {
	hook := lumberjack.Logger{
		Filename:   cfg.GetFilePath(),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	}

	var syncWrite zapcore.WriteSyncer
	if cfg.LogInConsole {
		syncWrite = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncWrite = zapcore.AddSync(&hook)
	}

	lvl := getLoggerLevel(cfg.Level)
	encoder := zap.NewProductionConfig()
	encoder.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder.EncoderConfig), syncWrite, zap.NewAtomicLevelAt(lvl))
	// log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	log := zap.New(core)
	Logger = log.Sugar()
}
