package cacmd

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
	GenerateCmd.PersistentFlags().StringP("dataDir","d","examples", "数据生成目录")
	CaCmd.AddCommand(GenerateCmd)
	CreateCmd.PersistentFlags().StringP("dataDir","d","examples", "数据目录")
	CaCmd.AddCommand(CreateCmd)
	PublishCmd.PersistentFlags().StringP("dataDir","d","examples", "数据目录")
	CaCmd.AddCommand(PublishCmd)
	SetUpCmd.PersistentFlags().StringP("dataDir","d","tmp", "临时数据目录")
	CaCmd.AddCommand(SetUpCmd)
}