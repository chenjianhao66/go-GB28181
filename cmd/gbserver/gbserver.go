package main

import (
	_ "github.com/chenjianhao66/go-GB28181/api/swagger"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver"
)

// @title           Go-GB28181项目前端APi接口
// @version         1.0
// @description     Go-GB28181是一个基于GB28181-2016标准实现的网络视频平台，用 Go 语言实现，实现了 SIP 协议和信令服务器。
// @termsOfService  http://swagger.io/terms/

// @contact.name   github homepage
// @contact.url    https://github.com/chenjianhao66/go-GB28181
// @contact.email  jianhao_c@qq.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:18080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	gbserver.NewApp("gbserver").Run()
}
