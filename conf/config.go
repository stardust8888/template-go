package conf

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"

	"template-go/library/httpclient/kms"
)

var (
	config *Configuration
)

func InitConfig() error {
	var (
		vip     = viper.New()
		setting = map[string]any{}
		err     error
	)

	vip.SetConfigName("config")
	vip.SetConfigType("toml")
	vip.AddConfigPath(".")
	vip.AddConfigPath("./conf")
	vip.AddConfigPath(".././conf")
	vip.AddConfigPath("../.././conf")
	vip.AddConfigPath("../../.././conf")
	vip.AddConfigPath("../../../.././conf")

	// 读取配置文件
	if err = vip.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("配置文件未找到:%v", err))
		return err
	}

	setting = vip.AllSettings()
	err = mapstructure.Decode(setting, &config)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件失败:%v", err))
		return err
	}

	kms.InitClient()
	config.Mysql.Password, err = kms.GetClient().Decrypt(context.Background(), config.Mysql.Password)
	if err != nil {
		panic(fmt.Sprintf("解析mysql密码文件失败:%v", err))
		return err
	}
	config.MysqlAdmin.Password, err = kms.GetClient().Decrypt(context.Background(), config.MysqlAdmin.Password)
	if err != nil {
		panic(fmt.Sprintf("解析mysql-admin密码文件失败:%v", err))
		return err
	}
	config.Redis.Password, err = kms.GetClient().Decrypt(context.Background(), config.Redis.Password)
	if err != nil {
		panic(fmt.Sprintf("解析redis密码文件失败:%v", err))
		return err
	}
	config.MinIO.Password, err = kms.GetClient().Decrypt(context.Background(), config.MinIO.Password)
	if err != nil {
		panic(fmt.Sprintf("解析minio密码文件失败:%v", err))
		return err
	}

	return nil
}

func InitKmsConfig() error {
	var (
		vip     = viper.New()
		setting = map[string]any{}
		err     error
	)

	vip.SetConfigName("config")
	vip.SetConfigType("toml")
	vip.AddConfigPath(".")
	vip.AddConfigPath("./conf")
	vip.AddConfigPath(".././conf")
	vip.AddConfigPath("../.././conf")
	vip.AddConfigPath("../../.././conf")
	vip.AddConfigPath("../../../.././conf")

	// 读取配置文件
	if err = vip.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("配置文件未找到:%v", err))
		return err
	}

	setting = vip.AllSettings()
	err = mapstructure.Decode(setting, &config)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件失败:%v", err))
		return err
	}

	return nil
}

func GetConfiguration() *Configuration {
	return config
}

type Configuration struct {
	App           AppConfig     `mapstructure:"app"`
	Mysql         MysqlConfig   `mapstructure:"mysql"`
	MysqlAdmin    MysqlConfig   `mapstructure:"mysql_admin"`
	Redis         RedisConfig   `mapstructure:"redis"`
	Elasticsearch ESConfig      `mapstructure:"elasticsearch"`
	Log           LogConfig     `mapstructure:"log"`
	Wx            WxConfig      `mapstructure:"wx"`
	MinIO         MinIOConfig   `mapstructure:"minio"`
	Cron          CronConfig    `mapstructure:"cron"`
	Service       ServiceConfig `mapstructure:"service"`
	Kafka         KafkaConfig   `mapstructure:"kafka"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

type LogConfig struct {
	Output     string `mapstructure:"output"`
	Level      int    `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type MysqlConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Migrate  int    `mapstructure:"migrate"`
}

type RedisConfig struct {
	Endpoints []string `mapstructure:"endpoints"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
}

type ESConfig struct {
	Endpoints  []string `mapstructure:"endpoints"`
	Username   string   `mapstructure:"username"`
	Password   string   `mapstructure:"password"`
	CloudID    string   `mapstructure:"cloud_id"`
	APIKey     string   `mapstructure:"api_key"`
	CertPath   string   `mapstructure:"cert_path"`
	CertBase64 string   `mapstructure:"cert_base64"`
}

type WxConfig struct {
	Endpoint            string `mapstructure:"endpoint"`
	AppID               string `mapstructure:"app_id"`
	Secret              string `mapstructure:"secret"`
	PayEndpoint         string `mapstructure:"pay_endpoint"`
	PayEndpointCallback string `mapstructure:"pay_endpoint_callback"`
	MchID               string `mapstructure:"mch_id"`
}

type MinIOConfig struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	UseSSL   bool   `mapstructure:"use_ssl"`
	Bucket   string `mapstructure:"bucket"`
}

type CronConfig struct {
	SaveMessageCron string `mapstructure:"save_message"`
	DeleteUserCron  string `mapstructure:"delete_user"`
}

type ServiceConfig struct {
}

type KafkaConfig struct {
	Version             string   `mapstructure:"version"` // kafka 协议版本，如 "3.0.0"
	ProducerMessage     Producer `mapstructure:"producer_message"`
	ConsumerMessagePush Consumer `mapstructure:"consumer_message_push"`
	ConsumerMessageSave Consumer `mapstructure:"consumer_message_save"`
}

type Producer struct {
	Brokers       []string `mapstructure:"brokers"`
	Topic         string   `mapstructure:"topic"`
	Async         bool     `mapstructure:"async"`          //异步
	RequiredBacks int      `mapstructure:"required_backs"` // -1=all，1=leader，0=none
}

type Consumer struct {
	Brokers        []string `mapstructure:"brokers"`
	Topic          string   `mapstructure:"topic"`
	GroupID        string   `mapstructure:"group_id"`
	StartOffset    int64    `mapstructure:"start_offset"`
	CommitInterval int      `mapstructure:"commit_interval"` // 异步自动提交间隔；0 表示手动提交
}
