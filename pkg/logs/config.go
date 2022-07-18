package logs

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type OutputType string

const (
	OutputConsole OutputType = "console"
	OutputFile    OutputType = "file"
)

type Config struct {
	Level      string             `json:"level" yaml:"level"`
	Output     OutputType         `json:"output" yaml:"output"`
	Lumberjack *lumberjack.Logger `json:"lumberjack" yaml:"lumberjack"`
}
