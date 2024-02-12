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
	TransferCmd.PersistentFlags().StringP("sourceDataDir","i","examples", "数据目录")
	TransferCmd.PersistentFlags().StringP("destDataDir","o","transfered_examples", "临时转换的数据目录")
	CaCmd.AddCommand(TransferCmd)
	AdaptCmd.PersistentFlags().StringP("topoYaml","t","topo.yaml","bgpemu使用的topo.yaml配置")
	AdaptCmd.PersistentFlags().StringP("publishPointsJson","i","tmp_publishPoints.json","写有发布点")
	AdaptCmd.PersistentFlags().StringP("topoWithRPKIYaml","o","topoWithRPKI.yaml","运行CA方容器和依赖方容器的topo.yaml配置")
	CaCmd.AddCommand(AdaptCmd)
	AttackCmd.PersistentFlags().StringP("attackJson","i","attak.json","攻击配置文件")
	CaCmd.AddCommand(AttackCmd)
}