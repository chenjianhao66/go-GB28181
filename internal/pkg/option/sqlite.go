package option

import "github.com/spf13/pflag"

// SqliteOptions 定义sqlite数据库的配置选项
type SqliteOptions struct {
	Path     string `json:"path,omitempty" mapstructure:"path"`
	File     string `json:"file" mapstructure:"file"`
	Username string `json:"username,omitempty" mapstructure:"username"`
	Password string `json:"password,omitempty" mapstructure:"password"`
	Database string `json:"database,omitempty" mapstructure:"database"`
}

func NewSqliteOptions() *SqliteOptions {
	return &SqliteOptions{
		Path:     "db",
		File:     "gbserver.db",
		Username: "root",
		Password: "root",
		Database: "go-gb28181",
	}
}

func (m *SqliteOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&m.Path, "sqlite.path", m.Path, "sqlite数据库的绝对路径")
	fss.StringVar(&m.File, "sqlite.file", m.File, "sqlite数据库文件名")
	fss.StringVar(&m.Username, "sqlite.username", m.Username, "sqlite数据库的用户名")
	fss.StringVar(&m.Password, "sqlite.password", m.Password, "sqlite数据库的密码")
	fss.StringVar(&m.Database, "sqlite.database", m.Database, "sqlite数据库名")
}
