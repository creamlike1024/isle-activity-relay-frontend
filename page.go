package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func MdToHTML(md []byte) []byte {
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

func GenHtml(subList *[]string) []byte {
	var htmlStr string
	// 添加 html 头
	htmlHead, err := os.ReadFile(config.TemplateHeadFilename)
	if err != nil {
		fmt.Printf("Failed to read template file: %s\n", config.TemplateHeadFilename)
		panic(err)
	}
	htmlStr = string(htmlHead) + "\n<body>\n" + "<section class=\"main\">\n"
	// 添加中继信息
	infoBody, err := os.ReadFile(config.TemplateInfoFilename)
	if err != nil {
		fmt.Printf("Failed to read template file: %s\n", config.TemplateInfoFilename)
		panic(err)
	}
	htmlStr += "\n" + string(infoBody)
	// 添加订阅列表
	subListHtml := MdToHTML(GenSubMdList(*subList))
	htmlStr += "\n" + string(subListHtml)
	htmlStr += fmt.Sprintf("\n<p>感谢大家的支持</p>\n")
	count, err := getAcceptedNotesCount(config.logPath)
	if err != nil {
		panic(err)
	}
	htmlStr += fmt.Sprintf("\n<p>累计转发 %d 条帖子</p>\n", count+config.offset)
	loc, _ := time.LoadLocation(config.Timezone)
	htmlStr += fmt.Sprintf("\n<p>数据最后更新于 %s</p>\n", time.Now().In(loc).Format(config.TimeFormat))
	// 添加实时日志
	logBody, err := os.ReadFile(config.templateLogFilename)
	if err != nil {
		fmt.Printf("Failed to read template file: %s\n", config.templateLogFilename)
		panic(err)
	}
	htmlStr += "\n" + string(logBody)
	// 添加 html 尾
	htmlStr += "\n</section>\n</body>\n</html>"
	return []byte(htmlStr)
}

func GenSubMdList(domainList []string) []byte {
	var list string
	count := 0
	var wg sync.WaitGroup
	for _, d := range domainList {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			nodeInfo, err1 := GetNodeInfo(d)
			nodeName, err2 := GetNodeName(d)
			if err1 == nil && err2 == nil {
				list += fmt.Sprintf("- %s | [%s](%s) | 👥 %d 📝 %d 🎯 %s %s\n", nodeName, d, "https://"+d, nodeInfo.Usage.Users.Total, nodeInfo.Usage.LocalPosts, GetSoftwareName(nodeInfo.Software.Name), nodeInfo.Software.Version)
			} else if err1 != nil && err2 == nil {
				list += fmt.Sprintf("- %s | [%s](%s)\n", nodeName, d, "https://"+d)
			} else if err1 == nil && err2 != nil {
				list += fmt.Sprintf("- [%s](%s) | 👥 %d 📝 %d 🎯 %s %s\n", d, "https://"+d, nodeInfo.Usage.Users.Total, nodeInfo.Usage.LocalPosts, GetSoftwareName(nodeInfo.Software.Name), nodeInfo.Software.Version)
			} else {
				list += fmt.Sprintf("- [%s](%s)\n", d, "https://"+d)
			}
			count++
		}(d)
	}
	wg.Wait()
	var md []byte
	md = append(md, []byte(fmt.Sprintf("\n共 **%d** 个站点\n", count))...)
	return append(md, []byte(list)...)
}

func getAcceptedNotesCount(path string) (int, error) {
	out, err := exec.Command("grep", "-c", "Accepted", path).Output()
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, err
	}
	return count, nil
}
