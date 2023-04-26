package main

import (
	"github.com/678go/xcos/cmd"
	"golang.org/x/exp/slog"
	"os"
)

func main() {
	opt := slog.HandlerOptions{
		AddSource: true,
	}
	slog.SetDefault(slog.New(opt.NewJSONHandler(os.Stdout)))
	if err := cmd.RootCmd.Execute(); err != nil {
		return
	}
}
