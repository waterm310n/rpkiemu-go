package ca

import (
	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

var (
	inputDiR string
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "创建ca层次结构",
		Run: func(cmd *cobra.Command, args []string) {
			krillop.Create(inputDiR)
		},
	}
)

