package cacmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/ca/attack"
)

var (
	AttackCmd = &cobra.Command{
		Use:   "attack",
		Short: "对rpkiemu网络中的CA方进行攻击",
		Run: func(cmd *cobra.Command, args []string) {
			if attackJson, err := cmd.Flags().GetString("attackJson"); err == nil {
				slog.Info(fmt.Sprintf("run rpkiemu ca attack -i %s",attackJson))
				attack.ExcuteAttack(attackJson)
			} else {
				slog.Error(err.Error())
			}

		},
	}
)
