package cacmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

var (
	PublishCmd = &cobra.Command{
		Use:   "publish",
		Short: "发布Roas",
		Run: func(cmd *cobra.Command, args []string) {
			if dataDir ,err := cmd.Flags().GetString("dataDir");err == nil{
				krillop.PublishRoas(dataDir)
			}else{
				slog.Error(err.Error())
			}
		},
	}
)
