package ca

import (
	"log/slog"
	"github.com/waterm310n/rpkiemu-go/ca/setup"
	"github.com/spf13/cobra"
)

var (
	SetUpCmd = &cobra.Command{
		Use:   "setup",
		Short: "对CA方中的krill,rsyncd,nginx等容器进行配置",
		Run: func(cmd *cobra.Command, args []string) {
			if dataDir ,err := cmd.Flags().GetString("dataDir");err == nil{
				setup.SetUp(dataDir)
			}else{
				slog.Error(err.Error())
			}
		},
	}
)
