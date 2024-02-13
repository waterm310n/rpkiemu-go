package cacmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/ca/data"
)

var (
	GenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "根据数据库内容生成ca侧数据",
		Run: func(cmd *cobra.Command, args []string) {
			if dataDir ,err := cmd.Flags().GetString("dataDir");err == nil{
				data.GenerateData(dataDir)
			}else{
				slog.Error(err.Error())
			}
		},
	}
)
