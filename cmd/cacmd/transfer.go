package cacmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/ca/data"
)

var (
	TransferCmd = &cobra.Command{
		Use:   "transfer",
		Short: "为兼容rpkiemu,将rpkiemu-go的CA配置文件转化为rpkiemu对应的CA配置文件",
		Run: func(cmd *cobra.Command, args []string) {
			var sourceDataDir,destDataDir string 
			if dataDir ,err := cmd.Flags().GetString("sourceDataDir");err == nil{
				sourceDataDir = dataDir
			}else{
				slog.Error(err.Error())
			}
			if dataDir,err := cmd.Flags().GetString("destDataDir");err == nil{
				destDataDir = dataDir
			}else{
				slog.Error(err.Error())
			}
			data.Transfer(sourceDataDir,destDataDir)
		},
	}
)


