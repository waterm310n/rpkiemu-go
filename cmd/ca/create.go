package ca

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

var (
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "创建ca层次结构",
		Run: func(cmd *cobra.Command, args []string) {
			if dataDir ,err := cmd.Flags().GetString("dataDir");err == nil{
				krillop.CreateHierarchy(dataDir)
			}else{
				slog.Error(err.Error())
			}
		},
	}
)

