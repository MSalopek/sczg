package config

import (
	"errors"
	"io/ioutil"
	"sczg/dbutil"

	"github.com/fatih/structs"
	"gopkg.in/yaml.v2"
)

// Env keeps current running environment params
type Env struct {
	CfgPath string
	DbPath  string
	Cfg     *Config
	DB      *dbutil.Storage
	Err     error
}

func SetupEnv(e *Env) {
	if e.CfgPath == "" || e.DbPath == "" {
		e.Err = errors.New("config and db path not specified")
	}
	cfg, err := InitCfg(e.CfgPath)
	if err != nil {
		e.Err = err
	}
	db, err := dbutil.InitStorage(e.DbPath)
	if err != nil {
		e.Err = err
	}
	e.Cfg = cfg
	e.DB = db
	e.Err = nil
}

type Urls struct {
	Index       string
	Marketing   string
	Ducani      string
	Ugostitelji string
	Ciscenje    string
	Proizvodnja string
	Turizam     string
	Fizicki     string
	Razno       string
	Admin       string
	Skladiste   string
}

type Config struct {
	Base     string
	Urls     Urls
	Agent    string
	Interval int
	Timeout  int
	Port     string
	Archive  int
}

// InitCfg initializes new config
func InitCfg(path string) (*Config, error) {
	var y Config
	YFile, err := ioutil.ReadFile(path)
	if err != nil {
		return &y, err
	}
	err = yaml.Unmarshal(YFile, &y)
	if err != nil {
		return &y, err
	}
	return &y, nil
}

// MapURLs returns urls as map
func (c *Config) MapURLs() map[string]string {
	urlMapInterface := structs.Map(c.Urls)
	urlMapStr := make(map[string]string)
	for k, v := range urlMapInterface {
		urlMapStr[k] = v.(string)
	}
	return urlMapStr
}
