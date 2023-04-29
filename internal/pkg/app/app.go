package app

import (
	"flag"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type App struct {
	// 编译出来的程序文件名称
	basename string
	// 应用名称
	name string
	// 详细描述
	description string
	// 程序横幅logo
	banner string
	// 应用真正要执行的函数
	runFunc RunFunc
	// app的选项
	options option.GbOption
	cmd     *cobra.Command
}

type Option func(*App)

type RunFunc func(basename string) error

func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runFunc = run
	}
}

func WithDescription(description string) Option {
	return func(app *App) {
		app.description = description
	}
}

func WithBanner(banner string) Option {
	return func(app *App) {
		app.banner = banner
	}
}

func WithOption(opt option.GbOption) Option {
	return func(app *App) {
		app.options = opt
	}
}

func NewApp(name, basename string, opts ...Option) *App {
	a := &App{
		basename: basename,
		name:     name,
	}

	for _, opt := range opts {
		opt(a)
	}

	a.buildCommand()

	return a
}

func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func (a *App) buildCommand() {
	cmd := &cobra.Command{
		Use:   strings.ToLower(a.basename),
		Short: a.name,
		Long:  a.description,
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	if a.options != nil {
		addConfigFlag(a.basename, cmd.Flags())
		fss := a.options.Flags()
		cmd.Flags().AddFlagSet(fss)
	}

	if a.runFunc != nil {
		cmd.RunE = a.launch
	}

	a.cmd = cmd
}

func (a *App) launch(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		_, _ = fmt.Fprint(os.Stderr, "解析标志失败")
		os.Exit(1)
	}

	if a.options != nil {
		if err := viper.Unmarshal(&a.options); err != nil {
			color.RedString("err: %v", err)
			_, _ = fmt.Fprint(os.Stderr, "解析配置文件失败")
			os.Exit(1)
		}
	}

	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}
	return nil
}
