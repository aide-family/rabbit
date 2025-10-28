package gorm

import (
	"fmt"
	"net/url"

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

	username  string
	password  string
	host      string
	port      int32
	database  string
	charset   string
	parseTime string
	loc       string
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
	c.PersistentFlags().StringVar(&f.charset, "charset", "utf8mb4", "mysql charset")
	c.PersistentFlags().StringVar(&f.parseTime, "parseTime", "true", "mysql parseTime")
	c.PersistentFlags().StringVar(&f.loc, "loc", "Asia/Shanghai", "mysql location")
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
	if charset := mysqlConf.GetCharset(); strutil.IsNotEmpty(charset) {
		f.charset = charset
	}
	if parseTime := mysqlConf.GetParseTime(); strutil.IsNotEmpty(parseTime) {
		f.parseTime = parseTime
	}
	if loc := mysqlConf.GetLoc(); strutil.IsNotEmpty(loc) {
		f.loc = loc
	}
}

func (f *Flags) databaseDSN() string {
	loc := url.QueryEscape(f.loc)
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s", f.username, f.password, f.host, f.port, f.database, f.charset, f.parseTime, loc)
}

func (f *Flags) connectDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/", f.username, f.password, f.host, f.port)
}
