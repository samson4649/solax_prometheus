package solax_prometheus

import (
  // "os"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
  "strings"
)

func InitLogger( f, l string ) ( *zap.Logger, error ) {
  var ll zapcore.Level
  switch strings.ToLower(l) {
  case "verbose":
  case "debug":
  case "dev":
    ll = zapcore.DebugLevel
  default:
    ll = zapcore.InfoLevel
  }
  cfg :=  zap.Config{
    Encoding: "json",
    Level: zap.NewAtomicLevelAt(ll),
    OutputPaths: []string{f},
    EncoderConfig: zapcore.EncoderConfig{
      MessageKey: "message",
      LevelKey:    "level",
      EncodeLevel: zapcore.CapitalLevelEncoder,
      TimeKey:    "time",
      EncodeTime: zapcore.ISO8601TimeEncoder,
      CallerKey:  "caller",
      EncodeCaller: zapcore.ShortCallerEncoder,
    },
  }

  return cfg.Build()
}
