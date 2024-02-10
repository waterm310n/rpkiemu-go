package cacmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/ca/setup"
)

var (
	AdaptCmd = &cobra.Command{
		Use:   "adapt",
		Short: "根据rpkiemu-go所需对",
		Run: func(cmd *cobra.Command, args []string) {
			topoYaml ,err := cmd.Flags().GetString("topoYaml")
			if err != nil {
				slog.Error(err.Error())
				return
			}
			publishPointsJson,err := cmd.Flags().GetString("publishPointsJson")
			if err != nil {
				slog.Error(err.Error())
				return
			}
			topoWithRPKIYaml,err := cmd.Flags().GetString("topoWithRPKIYaml")
			if err != nil {
				slog.Error(err.Error())
				return
			}
			setup.SequentialAdapt(topoYaml,publishPointsJson,topoWithRPKIYaml,1)
		},
	}
)
