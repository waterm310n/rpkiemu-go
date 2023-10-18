package cmd

import (
	"os"
	"log/slog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/waterm310n/rpkiemu-go/cmd/cacmd"
	"github.com/waterm310n/rpkiemu-go/cmd/rpcmd"
)

var (
	// Used for flags.
	cfgFile     string
	rootCmd = &cobra.Command{
		Use:   "rpkiemu-go",
		Short: "一个基于go语言的用于RPKI网络模拟的命令行工具",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config","", "配置文件 (默认为 [$HOME|$CURRENT_WORKDIR]/.rpki-emu.json)")
	rootCmd.AddCommand(cacmd.CaCmd)
	rootCmd.AddCommand(rpcmd.RpCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile) //使用config字符串指向的路径上的配置文件
	} else {
		home, err := os.UserHomeDir() // 查找home目录
		cobra.CheckErr(err) //如果home目录不存在就报错
		viper.AddConfigPath(".") //添加当前目录到查找配置文件的路径列表中
		viper.AddConfigPath(home) //添加home目录到查找配置文件的路径列表中
		viper.SetConfigType("json") //配置文件的类型
		viper.SetConfigName(".rpkiemu-go") //配置文件的名称
	}
	// viper.AutomaticEnv() //暂时没有使用到环境变量
	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Using config file", "fileName", viper.ConfigFileUsed())
	}else if err != nil{
		slog.Error(err.Error())
	}
}
