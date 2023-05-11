package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"os"
	"sync"
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

func FillHtml(html []byte) []byte {
	htmlStr := string(html)
	// 添加 html 头
	htmlHead, err := os.ReadFile(config.TemplateHeadFilename)
	if err != nil {
		fmt.Printf("Failed to read template file: %s\n", config.TemplateHeadFilename)
		panic(err)
	}
	htmlStr = string(htmlHead) + "\n<body>\n" + "<section class=\"main\">\n" + htmlStr
	// 添加 html 尾
	htmlStr += "\n</section>\n</body>"
	return []byte(htmlStr)
}

func FillMarkdownTemplate(md []byte, domainList []string) []byte {
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
	md = append(md, []byte(fmt.Sprintf("\n共 **%d** 个站点\n", count))...)
	return append(md, []byte(list)...)
}
