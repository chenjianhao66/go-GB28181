package parser

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/sirupsen/logrus"
	"runtime"
)

// GetCmdTypeFromXML 根据body获取XML配置文件中的根元素
func GetCmdTypeFromXML(body string) (key string, err error) {
	document := etree.NewDocument()

	err = document.ReadFromString(body)
	if err != nil {
		logrus.Debugf("解析XML配置失败，%s", err)
		return "", err
	}

	command := document.Root().Tag
	path := etree.MustCompilePath(command)
	cmdType := document.FindElementPath(path).FindElementPath(etree.MustCompilePath("CmdType")).Text()
	key = command + ":" + cmdType
	defer func() {
		// 捕获因为解析xml所抛出的panic
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("runtime error:", err)
		case nil:
			return
		default: // 非运行时错误
			fmt.Println("error:", err)
		}
	}()
	return
}
