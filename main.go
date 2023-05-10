package main

import (
	"context"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	RedisHost            string
	RedisPort            int
	RedisPassword        string
	OutputDir            string
	OutputFilename       string
	TemplateFilename     string
	TemplateHeadFilename string
}

var config Config

// Subscriber : Manage for Mastodon Traditional Style Relay Subscriber
type Subscriber struct {
	Domain     string `json:"domain,omitempty"`
	InboxURL   string `json:"inbox_url,omitempty"`
	ActivityID string `json:"activity_id,omitempty"`
	ActorID    string `json:"actor_id,omitempty"`
}

// Follower : Manage for LitePub Style Relay Follower
type Follower struct {
	Domain         string `json:"domain,omitempty"`
	InboxURL       string `json:"inbox_url,omitempty"`
	ActivityID     string `json:"activity_id,omitempty"`
	ActorID        string `json:"actor_id,omitempty"`
	MutuallyFollow bool   `json:"mutually_follow,omitempty"`
}

func mdToHTML(md []byte) []byte {
	// create Markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)
	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func fillHtml(html []byte) []byte {
	htmlStr := string(html)
	// 添加 html 头
	htmlHead, err := os.ReadFile(config.TemplateHeadFilename)
	if err != nil {
		fmt.Printf("Failed to read template head file: %s\n", config.TemplateHeadFilename)
		panic(err)
	}
	htmlStr = string(htmlHead) + "\n<body>\n" + "<section class=\"main\">\n" + htmlStr
	// 添加 html 尾
	htmlStr += "\n</section>\n</body>"
	return []byte(htmlStr)
}

func fillMarkdownTemplate(md []byte, sub []Subscriber, fo []Follower) []byte {
	var list string
	count := 0
	for _, s := range sub {
		list += fmt.Sprintf("- [%s](%s)\n", s.Domain, "https://"+s.Domain)
		count++
	}
	for _, f := range fo {
		list += fmt.Sprintf("- [%s](%s)\n", f.Domain, "https://"+f.Domain)
		count++
	}
	md = append(md, []byte(fmt.Sprintf("\n共 **%d** 个站点\n", count))...)
	return append(md, []byte(list)...)
}

func main() {
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
		TemplateFilename:     viper.GetString("template.markdown"),
		TemplateHeadFilename: viper.GetString("template.html_head"),
	}
	// 建立 redis 连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})
	ctx := context.Background()

	var subscribers []Subscriber
	var followers []Follower
	// 获取订阅者列表
	domains, _ := rdb.Keys(ctx, "relay:subscription:*").Result()
	for _, domain := range domains {
		domainName := strings.Replace(domain, "relay:subscription:", "", 1)
		inboxURL, _ := rdb.HGet(ctx, domain, "inbox_url").Result()
		activityID, err := rdb.HGet(ctx, domain, "activity_id").Result()
		if err != nil {
			activityID = ""
		}
		actorID, err := rdb.HGet(ctx, domain, "actor_id").Result()
		if err != nil {
			actorID = ""
		}
		subscribers = append(subscribers, Subscriber{domainName, inboxURL, activityID, actorID})
		// subscribersAndFollowers = append(subscribersAndFollowers, Subscriber{domainName, inboxURL, activityID, actorID})
	}
	// 获取关注者列表
	domains, _ = rdb.Keys(ctx, "relay:follower:*").Result()
	for _, domain := range domains {
		domainName := strings.Replace(domain, "relay:follower:", "", 1)
		inboxURL, _ := rdb.HGet(ctx, domain, "inbox_url").Result()
		activityID, err := rdb.HGet(ctx, domain, "activity_id").Result()
		if err != nil {
			activityID = ""
		}
		actorID, err := rdb.HGet(ctx, domain, "actor_id").Result()
		if err != nil {
			actorID = ""
		}
		mutuallyFollow, err := rdb.HGet(ctx, domain, "mutually_follow").Result()
		if err != nil {
			mutuallyFollow = "0"
		}
		followers = append(followers, Follower{domainName, inboxURL, activityID, actorID, mutuallyFollow == "1"})
		// subscribersAndFollowers = append(subscribersAndFollowers, Subscriber{domainName, inboxURL, activityID, actorID})
	}
	// 读取模板文件
	templateFile, err := os.ReadFile(config.TemplateFilename)
	if err != nil {
		fmt.Printf("Failed to read template file: %s\n", config.TemplateFilename)
		panic(err)
	}
	// 生成 HTML 文件
	bodyHtml := mdToHTML(fillMarkdownTemplate(templateFile, subscribers, followers))
	htmlBytes := fillHtml(bodyHtml)
	// 输出HTML文件
	// 检查输出目录是否存在，如果不存在则创建
	if _, err := os.Stat(config.OutputDir); os.IsNotExist(err) {
		err = os.Mkdir(config.OutputDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create output dir: %s\n", config.OutputDir)
			panic(err)
		}
	}
	err = os.WriteFile(config.OutputDir+"/"+config.OutputFilename, htmlBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write HTML file: %s", config.OutputDir+"/"+config.OutputFilename)
		panic(err)
	}
	fmt.Printf("Successfully generated HTML file: %s\n", config.OutputDir+"/"+config.OutputFilename)
}
