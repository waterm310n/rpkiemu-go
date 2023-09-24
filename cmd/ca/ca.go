package ca

import (
	"github.com/spf13/cobra"
)

var (
	CaCmd = &cobra.Command{
		Use:   "ca",
		Short: "执行ca方的相关操作",
	}
)

func init(){
	CreateCmd.PersistentFlags().StringVarP(&inputDiR,"dataDir","d","examples", "数据目录")
	CaCmd.AddCommand(CreateCmd)
	GenerateCmd.PersistentFlags().StringVarP(&outputDir, "dataDir","d","examples", "数据生成目录")
	CaCmd.AddCommand(GenerateCmd)
}