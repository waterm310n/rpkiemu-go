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
	GenerateCmd.PersistentFlags().StringVarP(&dataDir, "dataDir","d","examples", "数据生成目录")
	CaCmd.AddCommand(CreateCmd)
	CaCmd.AddCommand(GenerateCmd)
}