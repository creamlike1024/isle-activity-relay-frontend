package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	RedisHost            string
	RedisPort            int
	RedisPassword        string
	OutputDir            string
	OutputFilename       string
	TemplateInfoFilename string
	TemplateHeadFilename string
	templateLogFilename  string
	TimeFormat           string
	Timezone             string
	logPath              string
	offset               int
}

var config Config

func main() {
	// 读取配置文件
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Failed to read config file")
		panic(err)
	}
	config = Config{
		RedisHost:            viper.GetString("redis.host"),
		RedisPort:            viper.GetInt("redis.port"),
		RedisPassword:        viper.GetString("redis.password"),
		OutputDir:            viper.GetString("output.dir"),
		OutputFilename:       viper.GetString("output.file"),
		TemplateInfoFilename: viper.GetString("template.info"),
		templateLogFilename:  viper.GetString("template.log"),
		TemplateHeadFilename: viper.GetString("template.html_head"),
		TimeFormat:           viper.GetString("time.format"),
		Timezone:             viper.GetString("time.timezone"),
		logPath:              viper.GetString("relay_log.path"),
		offset:               viper.GetInt("relay_log.offset"),
	}
	// 建立 redis 连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})
	ctx := context.Background()
	// 获取订阅者和关注者列表
	subscribeAndFollowers := GetSubcribesAndFollowers(rdb, ctx)
	// 生成 HTML 文件
	htmlBytes := GenHtml(&subscribeAndFollowers)
	// 输出HTML文件
	// 检查输出目录是否存在，如果不存在则创建
	if _, err := os.Stat(config.OutputDir); os.IsNotExist(err) {
		err = os.Mkdir(config.OutputDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create output dir: %s\n", config.OutputDir)
			panic(err)
		}
	}
	// 输出 HTML 文件
	err = os.WriteFile(config.OutputDir+"/"+config.OutputFilename, htmlBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write HTML file: %s", config.OutputDir+"/"+config.OutputFilename)
		panic(err)
	}
	fmt.Printf("Successfully generated HTML file: %s\n", config.OutputDir+"/"+config.OutputFilename)
}
