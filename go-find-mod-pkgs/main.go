package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SortDependencies struct {
	Key   string
	Value int
}

func main() {

	// 指定项目根目录
	rootDir := "your code path"

	// 创建一个 map 用于存储依赖路径和它们的计数
	dependencies := make(map[string]int)

	shouldSkipDir := func(dir string) bool {
		return false
	}

	// 遍历项目中的所有 go.mod 文件
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error walking the directory: %v\n", err)
			return nil
		}

		if info.IsDir() && shouldSkipDir(path) {
			// 跳过特定目录
			return filepath.SkipDir
		}

		if !info.IsDir() && info.Name() == "go.mod" {
			// 读取 go.mod 文件内容
			content, readErr := ioutil.ReadFile(path)
			if readErr != nil {
				fmt.Printf("Error reading file: %v\n", readErr)
				return nil
			}

			fmt.Printf("find: %s\n", path)

			// 解析 go.mod 文件并提取依赖路径
			dependenciesList := extractDependencies(string(content))

			// 更新依赖路径的计数
			for _, dep := range dependenciesList {
				dependencies[dep]++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
		return
	}

	// 重新存储
	var sortDependencies []*SortDependencies
	for dep, count := range dependencies {
		sortDependencies = append(sortDependencies, &SortDependencies{
			Key:   dep,
			Value: count,
		})
	}

	// 按照引用多的进行排序
	sort.Slice(sortDependencies, func(i, j int) bool {
		if sortDependencies[i].Value > sortDependencies[j].Value {
			return true
		}
		return false
	})

	for _, dependency := range sortDependencies {
		fmt.Printf("dependency: %s, count: %d\n", dependency.Key, dependency.Value)
	}
}

func extractDependencies(content string) []string {

	var dependencies []string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "\t") {
			reqPath := strings.ReplaceAll(line, "\t", "")
			reqPathSlice := strings.Split(reqPath, " ")
			if len(reqPathSlice) <= 0 {
				continue
			}

			dependencies = append(dependencies, reqPathSlice[0])
		}
	}

	return dependencies
}
