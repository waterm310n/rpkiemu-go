package rpcmd

import (
	"github.com/spf13/cobra"
)

var (
	RpCmd = &cobra.Command{
		Use:   "rp",
		Short: "执行依赖方的相关操作",

	}
)

func init(){
	RpCmd.AddCommand(SetUpCmd)
}