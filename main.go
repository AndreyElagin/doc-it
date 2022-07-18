package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	strings "strings"

	"gopkg.in/yaml.v3"
)

func main() {
	conf := Conf{
		includeFileTypes: []string{".yaml", ".yml"},
		metaMarker:       "@doc-it",
		outputDir:        "out",
	}

	yamls := readYamls("yamls", conf)
	meta := make([]Meta, 0)
	for _, y := range yamls {
		meta = append(meta, y.toMeta(conf))
	}

	log.Println(meta[0].path.fileName())
}

type Path string

func (p Path) fileName() string {
	l := len(p)
	name := make([]byte, 0)

	for i := l - 1; i > 0; i-- {
		if p[i] == '/' {
			break
		}
		name = append(name, p[i])
	}

	nameLength := len(name)
	for i := 0; i < nameLength/2; i++ {
		tmp := name[i]
		name[i] = name[nameLength-i-1]
		name[nameLength-i-1] = tmp
	}

	return string(name)
}

type Conf struct {
	includeFileTypes []string
	metaMarker       string
	outputDir        string
}

type Yaml struct {
	path    Path
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
	path     Path
	comments []string
}

func readYamls(path string, conf Conf) []Yaml {
	yamls := make([]Yaml, 0)
	err := filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
		if !info.IsDir() && matchAnySuffix(p, conf.includeFileTypes) {
			bytes, err := os.ReadFile(p)
			check(err)

			yamls = append(yamls, Yaml{Path(p), string(bytes)})
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
