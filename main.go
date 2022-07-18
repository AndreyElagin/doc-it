package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	strings "strings"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	includeFileTypes []string
	metaMarker       string
}

type Yaml struct {
	path    string
	content string
}

func (ya Yaml) toMeta(conf Conf) Meta {
	var n yaml.Node
	err := yaml.Unmarshal([]byte(ya.content), &n)
	check(err)

	comments := make([]string, 0)
	collectMeta(&n, &comments, conf)

	return Meta{ya.path, comments}
}

type Meta struct {
	path     string
	comments []string
}

func main() {
	conf := Conf{
		includeFileTypes: []string{".yaml", ".yml"},
		metaMarker:       "@doc-it",
	}

	yamls := readYamls("yamls", conf)
	meta := make([]Meta, 0)
	for _, y := range yamls {
		meta = append(meta, y.toMeta(conf))
	}

	log.Println(meta)
}

func readYamls(path string, conf Conf) []Yaml {
	yamls := make([]Yaml, 0)
	err := filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
		if !info.IsDir() && matchAnySuffix(p, conf.includeFileTypes) {
			bytes, err := os.ReadFile(p)
			check(err)

			yamls = append(yamls, Yaml{p, string(bytes)})
		}
		return nil
	})
	check(err)

	return yamls
}

func matchAnySuffix(path string, suffixes []string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
}

func collectMeta(n *yaml.Node, comments *[]string, conf Conf) {
	if n == nil {
		return
	}
	if strings.Contains(n.HeadComment, conf.metaMarker) {
		*comments = append(*comments, n.HeadComment)
	}

	for _, node := range n.Content {
		collectMeta(node, comments, conf)
	}
}

func check(e error) {
	if e != nil {
		panic("Wow wow")
	}
}
