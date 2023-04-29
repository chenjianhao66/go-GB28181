package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

const configFlagName = "config"

var configFile string

func init() {
	pflag.StringVarP(&configFile, configFlagName, "c", configFile, "配置文件")
}

func addConfigFlag(basename string, fs *pflag.FlagSet) {

	fs.AddFlag(pflag.Lookup(configFlagName))

	viper.AutomaticEnv()

	cobra.OnInitialize(func() {
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			viper.AddConfigPath(".")

			viper.SetConfigName(basename)
		}

		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "错误，读取配置文件失败(%s): %v\n", configFile, err)
			os.Exit(1)
		}
	})
}
