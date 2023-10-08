package krillop

import (
	_ "os"
	"fmt"
	"log/slog"
)

//发布roa
func PublishRoas(dataDir string){
	slog.Debug(fmt.Sprintf("func PublishRoas Run with %s",dataDir))
}