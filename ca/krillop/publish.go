package krillop

import (
	"os"
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"
)

//发布roa
func PublishRoas(dataDir string){
	slog.Debug(fmt.Sprintf("func PublishRoas Run with %s",dataDir))
	caOps := createCAOp()
	var entries map[string]*fileEntry
	if dirEntries, err := os.ReadDir(dataDir); err != nil {
		slog.Error(err.Error())
	} else {
		entries = extract(dirEntries)
	}
	for certName, v := range entries {
		handle := getHandleFromPath(filepath.Join(dataDir, v.resource.Name()))
		publishPoint := handle.PublishPoint
		children := v.children
		if v.roas != nil{
			//发布roas
			path := filepath.Join(dataDir, v.roas.Name())
			caOps[publishPoint].addDeltaRoa(certName,path)
		}
		if children != nil {
			recursivePublishRoas(filepath.Join(dataDir, children.Name()),caOps)
		}
	}
}

func recursivePublishRoas(dataDir string,caOps map[string]CA){
	var entries map[string]*fileEntry
	if dirEntries, err := os.ReadDir(dataDir); err != nil {
		slog.Error(err.Error())
	} else {
		entries = extract(dirEntries)
	}
	var wg sync.WaitGroup
	for certName,v := range entries{
		certName := certName
		children := v.children
		roas := v.roas
		handle := getHandleFromPath(filepath.Join(dataDir, v.resource.Name()))
		publishPoint := handle.PublishPoint
		wg.Add(1)
		go func(){
			if roas != nil{
				roasPath := filepath.Join(dataDir, roas.Name())
				caOps[publishPoint].addDeltaRoa(certName,roasPath)
			}
			if children != nil{
				recursivePublishRoas(filepath.Join(dataDir, children.Name()),caOps)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}