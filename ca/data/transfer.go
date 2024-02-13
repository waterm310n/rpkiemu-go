package data

import (
	"encoding/json"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// 一个handle有自己的资源，自己发布的roa，它的下属ca
type fileEntry struct {
	resource fs.DirEntry
	roas     fs.DirEntry
	children fs.DirEntry
}

type parentHandle struct {
	Ipv4         []string `json:"ipv4,omitempty"`
	Ipv6         []string `json:"ipv6,omitempty"`
	Asn          []string `json:"asn,omitempty"`
	PublishPoint string   `json:"publish_point,omitempty"`
}

// Handle结构体
type oldHandle struct {
	PublishPoint string                  `json:"publish_point,omitempty"`
	Parents      map[string]parentHandle `json:"parents,omitempty"`
	Layer        int                     `json:"-"`
}

type ca struct {
	Name  string `json:"name,omitempty"`
	RIR   bool   `json:"rir,omitempty"`
	RRDP  bool   `json:"rrdp,omitempty"`
	Rsync bool   `json:"rsync,omitempty"`
}

// 分割目录，资源证书文件，roas文件
func extract(dirEntries []fs.DirEntry) map[string]*fileEntry {
	entries := make(map[string]*fileEntry)
	for _, entry := range dirEntries {
		switch {
		case entry.IsDir():
			name := entry.Name()[:len(entry.Name())-9]
			if _, ok := entries[name]; ok {
				entries[name].children = entry
			} else {
				entries[name] = &fileEntry{
					children: entry,
				}
			}
		case strings.HasSuffix(entry.Name(), "roas"):
			name := entry.Name()[:len(entry.Name())-5]
			if _, ok := entries[name]; ok {
				entries[name].roas = entry
			} else {
				entries[name] = &fileEntry{
					roas: entry,
				}
			}
		default:
			name := entry.Name()
			if _, ok := entries[name]; ok {
				entries[name].resource = entry
			} else {
				entries[name] = &fileEntry{
					resource: entry,
				}
			}
		}
	}
	return entries
}

func getHandleFromPath(path string) Handle {
	var handle Handle
	if bytes, err := os.ReadFile(path); err != nil {
		slog.Error(err.Error())
	} else {
		json.Unmarshal(bytes, &handle)
	}
	return handle
}

func Transfer(sourceDataDir, destDataDir string) {
	mp := make(map[string]oldHandle, 0)
	var entries map[string]*fileEntry
	if dirEntries, err := os.ReadDir(sourceDataDir); err != nil {
		slog.Error(err.Error())
	} else {
		entries = extract(dirEntries)
	}
	layer := 1
	cas := make([]ca, 0)
	roasMp := make(map[string][]string, 0)
	for certName, v := range entries {
		handle := getHandleFromPath(filepath.Join(sourceDataDir, v.resource.Name()))
		publishPoint := handle.PublishPoint
		cas = append(cas, ca{Name: publishPoint, RIR: true, RRDP: true, Rsync: true})
		mp[certName] = oldHandle{PublishPoint: publishPoint, Parents: make(map[string]parentHandle), Layer: layer}
		mp[certName].Parents["testbed"] = parentHandle{handle.Ipv4, handle.Ipv6, handle.Asn, publishPoint}
		if v.children != nil {
			recursiveTransfer(publishPoint, certName, filepath.Join(sourceDataDir, v.children.Name()), destDataDir, layer+1, mp, roasMp)
		}
	}
	os.MkdirAll(destDataDir, 0744)
	writeOldHierachyFile(destDataDir, cas, mp)
	writeOldRoasFile(destDataDir, cas, roasMp)
}

func recursiveTransfer(parentPublishPoint string, parentCertName string, sourceDataDir string, destDataDir string, layer int, mp map[string]oldHandle, roasMp map[string][]string) {
	var entries map[string]*fileEntry
	if dirEntries, err := os.ReadDir(sourceDataDir); err != nil {
		slog.Error(err.Error())
	} else {
		entries = extract(dirEntries)
	}
	for certName, v := range entries {
		handle := getHandleFromPath(filepath.Join(sourceDataDir, v.resource.Name()))
		publishPoint := handle.PublishPoint
		mp[certName] = oldHandle{PublishPoint: publishPoint, Parents: make(map[string]parentHandle), Layer: layer}
		mp[certName].Parents[parentCertName] = parentHandle{handle.Ipv4, handle.Ipv6, handle.Asn, parentPublishPoint}
		if v.roas != nil {
			//发布roas
			sourcePath := filepath.Join(sourceDataDir, v.roas.Name())
			destPath := filepath.Join(destDataDir, "roas-delta", publishPoint)
			os.MkdirAll(destPath, 0744)
			destPath = filepath.Join(destPath, certName)
			source, err := os.Open(sourcePath)
			if err != nil {
				slog.Error(err.Error())
			}
			defer source.Close()
			destination, err := os.Create(destPath)
			if err != nil {
				slog.Error(err.Error())
			}
			defer destination.Close()
			if _, err := io.Copy(destination, source) ;err != nil{
				slog.Error(err.Error())
			}
			roasMp[publishPoint] = append(roasMp[publishPoint], certName)
		}
		if v.children != nil {
			recursiveTransfer(publishPoint, certName, filepath.Join(sourceDataDir, v.children.Name()),destDataDir, layer+1, mp, roasMp)
		}
	}
}

func writeOldHierachyFile(destDataDir string, cas []ca, mp map[string]oldHandle) {
	file, err := os.Create(filepath.Join(destDataDir, "data.json"))
	defer file.Close()
	if err != nil {
		slog.Error(err.Error())
	}
	configure := make([]byte, 0)
	configure = append(configure, '{')
	configure = append(configure, "\"cas\":"...)
	if content, err := json.Marshal(cas); err == nil {
		configure = append(configure, content...)
	}
	configure = append(configure, ",\"handles\":"...)
	if content, err := marshalForOrdered(mp); err == nil {
		configure = append(configure, content...)
	}
	configure = append(configure, '}')
	file.Write(configure)
}

func writeOldRoasFile(destDataDir string, cas []ca, roasMp map[string][]string) {
	destDataDir = filepath.Join(destDataDir, "roas-delta")
	os.MkdirAll(destDataDir, 0744)
	file, err := os.Create(destDataDir + "/roas.json")
	defer file.Close()
	if err != nil {
		slog.Error(err.Error())
	}
	configure := make([]byte, 0)
	configure = append(configure, '{')
	configure = append(configure, "\"cas\":"...)
	if content, err := json.Marshal(cas); err == nil {
		configure = append(configure, content...)
	}
	configure = append(configure, ",\"handles\":"...)
	if content, err := json.Marshal(roasMp); err == nil {
		configure = append(configure, content...)
	}
	configure = append(configure, '}')
	file.Write(configure)
}
func marshalForOrdered(mp map[string]oldHandle) ([]byte, error) {
	res := []byte{}
	res = append(res, '{')
	handles := make([]struct {
		Name   string
		Handle oldHandle
	}, 0)
	for name, v := range mp {
		handles = append(handles, struct {
			Name   string
			Handle oldHandle
		}{name, v})
	}
	sort.Slice(handles, func(i, j int) bool {
		return handles[i].Handle.Layer < handles[j].Handle.Layer
	})
	for i, handle := range handles {
		res = append(res, '"')
		res = append(res, []byte(handle.Name)...)
		res = append(res, '"')
		res = append(res, ':')
		if content, err := json.Marshal(handle.Handle); err == nil {
			res = append(res, content...)
		}
		if i != len(handles)-1 {
			res = append(res, ',')
		}
	}
	res = append(res, '}')
	return res, nil
}
