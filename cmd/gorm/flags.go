package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
)

type Flags struct {
	cmd.GlobalFlags
	configPath string
	forceGen   bool

	username string
	password string
	host     string
	port     int32
	database string
	params   []string
}

var flags Flags

func (f *Flags) addFlags(c *cobra.Command) {
	f.GlobalFlags = cmd.GetGlobalFlags()
	c.PersistentFlags().StringVarP(&f.configPath, "config", "c", "", "config file")
	c.PersistentFlags().BoolVarP(&f.forceGen, "force-gen", "f", false, "force generate code, if the code already exists, it will be overwitten")
	c.PersistentFlags().StringVar(&f.username, "username", "root", "mysql username")
	c.PersistentFlags().StringVar(&f.password, "password", "123456", "mysql password")
	c.PersistentFlags().StringVar(&f.host, "host", "localhost", "mysql host")
	c.PersistentFlags().Int32Var(&f.port, "port", 3306, "mysql port")
	c.PersistentFlags().StringVar(&f.database, "database", "rabbit", "mysql database")
	c.PersistentFlags().StringSliceVar(&f.params, "params", []string{"charset=utf8mb4", "parseTime=true", "loc=Asia/Shanghai"}, "mysql params")
}

func (f *Flags) applyToBootstrap(bc *conf.Bootstrap) {
	if pointer.IsNil(bc) || pointer.IsNil(bc.GetMain()) {
		return
	}
	mysqlConf := bc.GetMain()
	if username := mysqlConf.GetUsername(); strutil.IsNotEmpty(username) {
		f.username = username
	}
	if password := mysqlConf.GetPassword(); strutil.IsNotEmpty(password) {
		f.password = password
	}
	if host := mysqlConf.GetHost(); strutil.IsNotEmpty(host) {
		f.host = host
	}
	if port := mysqlConf.GetPort(); port > 0 {
		f.port = port
	}
	if database := mysqlConf.GetDatabase(); strutil.IsNotEmpty(database) {
		f.database = database
	}
}

func (f *Flags) databaseDSN() string {
	urlParams := url.Values{}
	for _, param := range f.params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) == 2 {
			urlParams.Add(parts[0], parts[1])
		}
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", f.username, f.password, f.host, f.port, f.database, urlParams.Encode())
}

func (f *Flags) connectDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/", f.username, f.password, f.host, f.port)
}
