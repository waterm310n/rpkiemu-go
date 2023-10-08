package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/waterm310n/rpkiemu-go/cmd"
)

func setLog() (*os.File,error) {
	f, err := os.Create("execute.log")
    if err != nil {
        return nil,fmt.Errorf("rpkiemu-go can not create execute.log,so use stderr as ErrOut")
    }
    logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelDebug,
	}))
    slog.SetDefault(logger)
	return f,nil
}

func main() {
	if f,err := setLog() ; err != nil{
		slog.Error(err.Error())
	}else{
		defer f.Close()
	}
	cmd.Execute()
}
