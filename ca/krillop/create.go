package krillop

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/waterm310n/rpkiemu-go/ca/data"
	"github.com/waterm310n/rpkiemu-go/k8sexec"
)

// 一个handle有自己的资源，自己发布的roa，它的下属ca
type fileEntry struct {
	resource fs.DirEntry
	roas     fs.DirEntry
	children fs.DirEntry
}

/*一个CA可以操作的所有接口*/
type CA interface {
	createHandle(handle string)                                              //在当前CA中创建指定handle
	deleteHandle(handle string)                                              //在当前CA中删除指定handle
	getHandles() []string                                                    //获取当前CA管理的所有handle名
	existHandle(handle string) bool                                          //检查一个handle是否在当前CA
	getChildren(handle string) []string                                      //获取当前handle的所有子handle名
	getParent(handle string) []string                                        //获取当前handle的所有父亲handle名
	getRoaName(handle string, asn int) (map[string]interface{}, error)       //获取当前handle指定asn的所有Roa证书名
	getCertName(handle, parentHandle string) (map[string]interface{}, error) //获取当前handle从指定父handle获取的资源证书名

	//下面是发布点处理方法

	getRepoRequest(handle string) error   //在当前CA中为指定handle创建repo请求
	setRepoConfigure(handle string) error //在当前发布点中为指定handle配置发布点信息
	setPubserver(handle string) error     //在当前CA中为指定handle配置它的发布点位置信息

	//下面是上下级CA资源处理方法

	getParentRequest(handle string) error                                //在当前CA中为指定handle创建request请求
	setChild(handle, childHandle string, ipv4, ipv6, asn []string) error //在当前CA从指定handle为下属handle分配资源
	setParent(handle, parentHandle string) error                         //在当前CA为指定handle配置其上级handle

	//下面是roa发布处理方法

	AddAsnIpPair(handle, ip, asn string)    //在当前CA为指定handle添加一条ASN-IP对
	RemoveAsnIpPair(handle, ip, asn string) //在当前CA为指定handle删除一条ASN-IP对
	addDeltaRoa(handle, file string)        //在当前CA中,为指定handle插入一个文件包含的ASN-IP对

	//攻击时使用的操作
	DeleteRoa(handle string, asn int)                                    //在handle中为破坏包含asn的所有roa证书
	CorruptRoa(handle string, asn int)                                   //在handle中为破坏包含asn的所有roa证书
	DeleteCert(handle string, parentHandle string, parentCa KrillK8sCA)  //在parentHandle中为删除handle的资源证书
	CorruptCert(handle string, parentHandle string, parentCa KrillK8sCA) //在parentHandle中为删除handle的资源证书

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
		case strings.HasSuffix(entry.Name(), ".conf"):
			continue
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

// 创建CA操作接口，从配置文件的publishPoints中读取
func CreateCAOp() map[string]*KrillK8sCA {
	var publishPoints map[string]struct {
		Namespace       string `json:"namespace,omitempty" mapstructure:"namespace,omitempty"`
		PodName         string `json:"pod_name,omitempty" mapstructure:"pod_name,omitempty"`
		CAContainerName string `json:"ca_container_name,omitempty" mapstructure:"ca_container_name,omitempty"`
		IsRIR           bool   `json:"is_rir,omitempty" mapstructure:"is_rir,omitempty"`
	}
	caOps := make(map[string]*KrillK8sCA)
	viper.Sub("publish_points").Unmarshal(&publishPoints)
	for name, v := range publishPoints {
		if execOptions, err := k8sexec.NewExecOptions(v.Namespace, v.PodName, v.CAContainerName); err == nil {
			kCA := NewKrillK8sCA(execOptions, v.IsRIR)
			caOps[name] = kCA
		}
	}
	return caOps
}

func getHandleFromPath(path string) data.Handle {
	var handle data.Handle
	if bytes, err := os.ReadFile(path); err != nil {
		slog.Error(err.Error())
	} else {
		json.Unmarshal(bytes, &handle)
	}
	return handle
}

// 创建ca层次结构中的rir部分
func CreateHierarchy(dataDir string) {
	slog.Debug(fmt.Sprintf("func CreateHierarchy Run with %s", dataDir))
	caOps := CreateCAOp()
	var entries map[string]*fileEntry
	if dirEntries, err := os.ReadDir(dataDir); err != nil {
		slog.Error(err.Error())
	} else {
		entries = extract(dirEntries)
	}
	for certName, v := range entries {
		handle := getHandleFromPath(filepath.Join(dataDir, v.resource.Name()))
		publishPoint := handle.PublishPoint
		if _, ok := caOps[publishPoint]; !ok {
			slog.Error(publishPoint + "is not exist")
			continue
		}
		caOps[publishPoint].createHandle(certName)
		if err := setRepo(publishPoint, publishPoint, certName, caOps); err != nil {
			slog.Error(err.Error())
			continue
		}
		if err := setParentChildrenRel(publishPoint, publishPoint, certName, "testbed", handle, caOps); err != nil {
			slog.Error(err.Error())
			continue
		}
		if v.children != nil {
			recursiveCreateHierarchy(publishPoint, certName, filepath.Join(dataDir, v.children.Name()), caOps)
		}
	}
}

// 创建ca层次结构rir之下的部分
func recursiveCreateHierarchy(parentPublishPoint, parentCertName, dataDir string, caOps map[string]*KrillK8sCA) {
	var entries map[string]*fileEntry
	if dirEntries, err := os.ReadDir(dataDir); err != nil {
		slog.Error(err.Error())
	} else {
		entries = extract(dirEntries)
	}
	var wg sync.WaitGroup
	for certName, v := range entries {
		certName := certName
		children := v.children
		handle := getHandleFromPath(filepath.Join(dataDir, v.resource.Name()))
		publishPoint := handle.PublishPoint
		wg.Add(1)
		go func() {
			if caOps[publishPoint] == nil {
				return
			}
			if publishPoint == parentPublishPoint {
				caOps[publishPoint].createHandle(certName)
				if err := setRepo(publishPoint, publishPoint, certName, caOps); err != nil {
					slog.Error(err.Error())
				}
				if err := setParentChildrenRel(publishPoint, publishPoint, certName, parentCertName, handle, caOps); err != nil {
					slog.Error(err.Error())
				}
				if children != nil {
					recursiveCreateHierarchy(publishPoint, certName, filepath.Join(dataDir, children.Name()), caOps)
				}
			} else {
				slog.Debug(publishPoint)
				caOps[publishPoint].createHandle(certName)
				if err := setRepo(publishPoint, parentPublishPoint, certName, caOps); err != nil {
					slog.Error(err.Error())
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func setRepo(publishPoint, parentPublishPoint, certName string, caOps map[string]*KrillK8sCA) error {
	var err error
	if publishPoint == parentPublishPoint {
		cnt := 0
		for ; cnt < 5; cnt++ {
			if err := caOps[publishPoint].getRepoRequest(certName); err != nil {
				cnt++
				time.Sleep(1 * time.Second)
				continue
			}
			if err := caOps[publishPoint].setPubserver(certName); err != nil {
				cnt++
				time.Sleep(1 * time.Second)
				continue
			}
			if err := caOps[publishPoint].setRepoConfigure(certName); err != nil {
				cnt++
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

	} else {
		cnt := 0
		for ; cnt < 5; cnt++ {
			if err := caOps[publishPoint].getRepoRequest(certName); err != nil {
				cnt++
				time.Sleep(1 * time.Second)
				continue
			}

		}
	}
	return err
}

func setParentChildrenRel(publishPoint, parentPublishPoint, certName, parentCertName string, handle data.Handle, caOps map[string]*KrillK8sCA) error {
	var err error
	if publishPoint == parentPublishPoint {
		cnt := 0
		for ; cnt < 5; cnt++ {
			if err = caOps[publishPoint].getParentRequest(certName); err != nil {
				cnt++
				time.Sleep(1 * time.Second)
				continue
			}
			if err = caOps[publishPoint].setChild(parentCertName, certName, handle.Ipv4, handle.Ipv6, handle.Asn); err != nil {
				cnt++
				time.Sleep(1 * time.Second)
				continue
			}
			if err = caOps[publishPoint].setParent(certName, parentCertName); err != nil {
				cnt++
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
	}
	return err
}
