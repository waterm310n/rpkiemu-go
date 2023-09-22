package rp

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	RpCmd = &cobra.Command{
		Use:   "rp",
		Short: "执行依赖方的相关操作",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("rp is excuted")
		},
	}
)

func init(){
	
}