package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg       *ini.File
	RunMode   string
	HttpPort  int
	ReadTime  time.Duration
	WriteTime time.Duration
	PageSize  int
	JwtSecret string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")

	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini' : %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	RunMode = Cfg.Section("").Key("RUN_MODE_").MustString("debug")
	HttpPort = sec.Key("HTTP_PORT").MustInt(9000)
	ReadTime = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTime = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app' : %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("_this&$_-$@--=_gin+")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}
