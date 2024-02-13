package rpcmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/rp/routinatorop"
)

var (
	SetUpCmd = &cobra.Command{
		Use:   "setup",
		Short: "对CA方中的krill,rsyncd,nginx等容器进行配置",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("run rp setup")
			routinatorop.Setup()
		},
	}
)
