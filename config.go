package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"

	"github.com/Panseng/logger-go/utils"
)

const configKeyLog = "logger"
var envConfigKeyLog = "log"

const (
	DefaultLevel        string = "info"
	DefaultSubFolderLog string = "logs"
	DefaultFilename     string = "logs.log"
	DefaultMaxSize      int    = 30
	DefaultMaxBackups   int    = 30
	DefaultMaxAge       int    = 30
	DefaultLocalTime    bool   = true
	DefaultCompress     bool   = false
	defaultLogInConsole bool   = false // 日志文件默认不输出到 console
)

// level
// debug, info, warn, error, dpanic, panic, fatal
type Loggers struct {
	BaseDir      string
	Level        string
	Folder       string // 日志的文件夹 相对目录: 相对配置基础目录
	Filename     string
	MaxSize      int // mb
	MaxBackups   int // 文件数
	MaxAge       int // 日志保留市场 day
	LocalTime    bool
	Compress     bool
	LogInConsole bool // 是否输出到控制台
}

type jsonLogger struct {
	Level        string `json:"level,omitempty"`
	Folder       string `json:"folder,omitempty"`
	Filename     string `json:"filename,omitempty"`
	MaxSize      int    `json:"max_size,omitempty"`
	MaxBackups   int    `json:"max_backups,omitempty"`
	MaxAge       int    `json:"max_age,omitempty"`
	LocalTime    bool   `json:"local_time,omitempty"`
	Compress     bool   `json:"compress,omitempty"`
	LogInConsole bool   `json:"log_in_console,omitempty"`
}

func (cfg *Loggers) ConfigKey() string {
	return configKeyLog
}

func (cfg *Loggers) SaveLoggerToDisk(path string) error {
	f := filepath.Dir(path)
	err := utils.MakeAllDirIfNotExist(f)
	if err != nil {
		return nil
	}
	return cfg.SaveJSON(path)
}

func (cfg *Loggers) LoadJSONFileAndEnv(path string) error {
	if err := cfg.LoadJSONFromFile(path); err != nil {
		return err
	}
	return cfg.ApplyEnvVars()
}

func (cfg *Loggers) Default() error {
	// level: debug info warn error dpanic panic fatal
	cfg.Level = DefaultLevel
	cfg.Folder = DefaultSubFolderLog
	cfg.Filename = DefaultFilename
	cfg.MaxSize = DefaultMaxSize
	cfg.MaxAge = DefaultMaxAge
	cfg.MaxBackups = DefaultMaxBackups
	cfg.LocalTime = DefaultLocalTime
	cfg.Compress = DefaultCompress
	cfg.LogInConsole = defaultLogInConsole
	return nil
}

func (cfg *Loggers) ApplyEnvVars() error {
	jcfg := cfg.toJSONConfig()

	err := envconfig.Process(envConfigKeyLog, jcfg)
	if err != nil {
		return err
	}

	return cfg.applyJSONConfig(jcfg)
}

func (cfg *Loggers) toJSONConfig() *jsonLogger {
	return &jsonLogger{
		Level:        cfg.Level,
		Folder:       cfg.Folder,
		Filename:     cfg.Filename,
		MaxSize:      cfg.MaxSize,
		MaxBackups:   cfg.MaxBackups,
		MaxAge:       cfg.MaxAge,
		LocalTime:    cfg.LocalTime,
		Compress:     cfg.Compress,
		LogInConsole: cfg.LogInConsole,
	}
}

func (cfg *Loggers) applyJSONConfig(jcfg *jsonLogger) error {
	utils.SetIfNotDefault(jcfg.Level, &cfg.Level)
	utils.SetIfNotDefault(jcfg.Folder, &cfg.Folder)
	utils.SetIfNotDefault(jcfg.Filename, &cfg.Filename)
	utils.SetIfNotDefault(jcfg.MaxSize, &cfg.MaxSize)
	utils.SetIfNotDefault(jcfg.MaxBackups, &cfg.MaxBackups)
	utils.SetIfNotDefault(jcfg.MaxAge, &cfg.MaxAge)
	utils.SetIfNotDefault(jcfg.LocalTime, &cfg.LocalTime)
	utils.SetIfNotDefault(jcfg.Compress, &cfg.Compress)
	utils.SetIfNotDefault(jcfg.LogInConsole, &cfg.LogInConsole)

	return nil
}

func (cfg *Loggers) LoadJSON(raw []byte) error {
	jcfg := &jsonLogger{}
	err := json.Unmarshal(raw, jcfg)
	if err != nil {
		fmt.Println("Error unmarshaling logger config")
		return err
	}

	cfg.Default()

	return cfg.applyJSONConfig(jcfg)
}

func (cfg *Loggers) LoadJSONFromFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return cfg.LoadJSON(file)
}

func (cfg *Loggers) ToJSON() ([]byte, error) {
	jcfg := cfg.toJSONConfig()
	return utils.DefaultJSONMarshal(jcfg)
}

func (cfg *Loggers) SaveJSON(path string) error {
	bs, err := cfg.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, bs, 0600)
}

func (cfg *Loggers) Validate() error {
	if cfg.Filename == "" {
		return errors.New("logger file name is empty")
	}
	return nil
}

func (cfg *Loggers) ToDisplayJSON() ([]byte, error) {
	return utils.DisplayJSON(cfg.toJSONConfig())
}

func (cfg *Loggers) SetBaseDir(dir string) {
	cfg.BaseDir = dir
}

// 获取 文件夹路径
func (cfg *Loggers) GetFolder() string {
	if filepath.IsAbs(cfg.Folder) {
		return cfg.Folder
	}

	return filepath.Join(cfg.BaseDir, cfg.Folder)
}

// 获取 文件路径
func (cfg *Loggers) GetFilePath() string {
	return filepath.Join(cfg.GetFolder(), cfg.Filename)
}
