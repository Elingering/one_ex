package base

import (
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

var formatter *prefixed.TextFormatter
var lfh *utils.LineNumLogrusHook

func InitLogrus() {
	// 定义日志格式
	formatter = &prefixed.TextFormatter{}
	//设置高亮显示的色彩样式
	formatter.ForceColors = true
	formatter.DisableColors = false
	formatter.ForceFormatting = true
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "41",
		PanicLevelStyle: "41",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})
	//开启完整时间戳输出和时间戳格式
	formatter.FullTimestamp = true
	//设置时间格式
	formatter.TimestampFormat = "2006-01-02.15:04:05.000000"
	//设置日志formatter
	log.SetFormatter(formatter)
	log.SetOutput(colorable.NewColorableStdout())
	//日志级别，通过环境变量来设置
	// 后期可以变更到配置中来设置
	if os.Getenv("log.debug") == "true" {
		log.SetLevel(log.DebugLevel)
	}
	//开启调用函数、文件、代码行信息的输出
	log.SetReportCaller(true)
	//设置函数、文件、代码行信息的输出的hook
	SetLineNumLogrusHook()
}

func SetLineNumLogrusHook() {
	lfh = utils.NewLineNumLogrusHook()
	lfh.EnableFileNameLog = true
	lfh.EnableFuncNameLog = true
	log.AddHook(lfh)
}
