package ca

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "创建ca层次结构",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("create command is excuted")
		},
	}
)

