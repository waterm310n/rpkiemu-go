package krillop

import (
	"fmt"
	"log/slog"
	"os"
)

/*
为ca容器，创建krill.conf文件，完成一些前置准备
*/

const KRILL_TEMPLATE = `admin_token = "krillTestBed"
data_dir = "/var/krill/data/"
log_type = "stderr"
ip = "0.0.0.0"
service_uri = "https://nginx.%s.publication/"
bgp_risdumps_enabled = false
`

const KRILL_TESTBED_TEMPLATE = `admin_token = "krillTestBed"
data_dir = "/var/krill/data/"
log_type = "stderr"
ip = "0.0.0.0"
service_uri = "https://nginx.%s.publication/" 
bgp_risdumps_enabled = false

[testbed]
rrdp_base_uri = "https://nginx.%s.publication/rrdp/"
rsync_jail = "rsync://rsyncd.%s.publication/repo/"
ta_uri = "https://nginx.%s.publication/ta/ta.cer"
ta_aia = "rsync://rsyncd.%s.publication/ta/ta.cer"
`

// 在tmp目录下创建容器名.conf
func createKrillConfig(dataDir, containerName string, isRIR bool) error {
	var err error
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil  {
		slog.Error(err.Error())
		os.Exit(1)
	}
	file, err := os.OpenFile(dataDir+"/"+containerName+".conf",os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if isRIR {
		_, err = file.WriteString(fmt.Sprintf(KRILL_TESTBED_TEMPLATE, containerName, containerName, containerName, containerName, containerName))
	} else {
		_, err = file.WriteString(fmt.Sprintf(KRILL_TEMPLATE, containerName))
	}
	return err
}
