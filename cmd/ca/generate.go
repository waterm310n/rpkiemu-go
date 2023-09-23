package ca

import (
	"github.com/spf13/cobra"
	"github.com/waterm310n/rpkiemu-go/ca/data"
)

var dataDir string

var (
	GenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "根据数据库内容生成ca侧数据",
		Run: func(cmd *cobra.Command, args []string) {
			data.GenerateData(dataDir)
		},
	}
)
