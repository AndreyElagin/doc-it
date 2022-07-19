package docit

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"doc-it/pkg/config"
	"doc-it/pkg/errorutils"
	"gopkg.in/yaml.v3"
)

type Yaml struct {
	Path    Path
	Content string
}

func (ya Yaml) ToMeta(conf config.Conf) Meta {
	var n yaml.Node
	err := yaml.Unmarshal([]byte(ya.Content), &n)
	errorutils.Check(err)

	comments := make([]string, 0)
	collectMeta(&n, &comments, conf)

	return Meta{ya.Path, comments}
}

func ReadYamls(path string, conf config.Conf) []Yaml {
	yamls := make([]Yaml, 0)
	err := filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
		if !info.IsDir() && matchAnySuffix(p, conf.IncludeFileTypes) {
			bytes, err := os.ReadFile(p)
			errorutils.Check(err)

			yamls = append(yamls, Yaml{Path(p), string(bytes)})
		}
		return nil
	})
	errorutils.Check(err)

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

func collectMeta(n *yaml.Node, comments *[]string, conf config.Conf) {
	if n == nil {
		return
	}
	if strings.Contains(n.HeadComment, conf.MetaMarker) {
		*comments = append(*comments, n.HeadComment)
	}

	for _, node := range n.Content {
		collectMeta(node, comments, conf)
	}
}
