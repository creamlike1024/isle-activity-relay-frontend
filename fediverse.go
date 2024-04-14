package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type NodeInfo struct {
	Software struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"software"`
	Usage struct {
		Users struct {
			Total int `json:"total"`
		} `json:"users"`
		LocalPosts int `json:"localPosts"`
	} `json:"usage"`
	OpenRegistrations bool `json:"openRegistrations"`
}

type nodeInfoWellKnown struct {
	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
}

func GetNodeInfo(domain string) (NodeInfo, error) {
	var nodeInfo NodeInfo
	url, err := getNodeInfoLink(domain)
	if err != nil || url == "" {
		return nodeInfo, err
	}
	// 发送 get 请求
	resp, err := http.Get(url)
	if err != nil {
		return nodeInfo, err
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nodeInfo, err
	}
	err = json.Unmarshal(body, &nodeInfo)
	if err != nil {
		return nodeInfo, err
	}
	return nodeInfo, nil
}

func getNodeInfoLink(domain string) (string, error) {
	url := "https://" + domain + "/.well-known/nodeinfo"
	// 发送 get 请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var nodeInfoLinks nodeInfoWellKnown
	err = json.Unmarshal(body, &nodeInfoLinks)
	if err != nil {
		return "", err
	}
	for _, link := range nodeInfoLinks.Links {
		if link.Rel == "http://nodeinfo.diaspora.software/ns/schema/2.0" {
			return link.Href, nil
		}
	}
	return "", nil
}

func GetNodeName(domain string) (string, error) {
	url := "https://" + domain + "/manifest.json"
	// 发送 get 请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var manifest struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		return "", err
	}
	return manifest.Name, nil
}

func GetSoftwareName(name string) string {
	// 将 name 转换为 byte 数组
	nameBytes := []byte(name)
	if len(nameBytes) == 0 {
		return ""
	}
	// 将 nameBytes 中的首字母转换为大写
	nameBytes[0] -= 32
	// 将 nameBytes 转换为 string
	name = string(nameBytes)
	return name
}
