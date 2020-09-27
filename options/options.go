package options

import (
	"flag"



	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CommandLine содержит параметры командной строки
type CommandLine struct {
	CfgFileName string
}

// getCommandLine считывает параметры командной строки
func getCommandLine(appName string) (c CommandLine) {
	flag.StringVar(&c.CfgFileName, "config", appName, "configuration file name")
	flag.Parse()

	return
}

// Options содержит опции приложения
type Options struct {
	ApplicationName string
	CmdLine         CommandLine
	Log             application.LogOptions
	GRPC            grpc.Options
	HTTP            http.Options

	SecretAPI secret_api.Options
}

// Load считывает опции приложения
func Load(appName string) (o *Options) {
	o = &Options{}

	o.ApplicationName = appName

	// считываем командную строку
	o.CmdLine = getCommandLine(appName)

	// считываем опции
	viper.SetConfigName(o.CmdLine.CfgFileName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath("/etc")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warnf("Configuration file %s not found", o.CmdLine.CfgFileName)
		} else {
			log.Fatalf("Unable to load configuration file %s: %s", o.CmdLine.CfgFileName, err)
		}
	}

	o.Log = application.LogOptions{
		LogLevel:     viper.GetString("log.level"),
		Syslog:       viper.GetString("log.syslog"),
		GraylogHost:  viper.GetString("log.graylog_host"),
		GraylogLevel: viper.GetString("log.graylog_level"),
	}

	o.SecretAPI = secret_api.Options{
		Timeout:        viper.GetDuration("secret_api.timeout"),
		PolligInterval: viper.GetDuration("secret_api.polling_interval"),
		PollingDelay:   viper.GetDuration("secret_api.polling_delay"),
	}

	o.GRPC = grpc.Options{
		Addr: viper.GetString("api.grpc.addr"),
	}

	o.HTTP = http.Options{
		Addr: viper.GetString("api.http.addr"),
	}

	return
}

// setDefaults устанавливает дефолтные значения опций
func setDefaults() {
	viper.SetDefault("log.level", "DEBUG")
	viper.SetDefault("log.graylog_level", "INFO")

	viper.SetDefault("secret_api.timeout", "5s")
	viper.SetDefault("secret_api.polling_interval", "100ms")
	viper.SetDefault("secret_api.delay", "0s")

	viper.SetDefault("api.grpc.addr", ":9002")
	viper.SetDefault("api.http.addr", ":8003")
}
