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

type PieceOfRef struct {
	Comment         string
	ObjectReference string
}

func (ya Yaml) ToMeta(conf config.Conf) Meta {
	var n yaml.Node
	err := yaml.Unmarshal([]byte(ya.Content), &n)
	errorutils.Check(err)

	refPieces := make([]PieceOfRef, 0)
	collectMeta(n.Content[0], &refPieces, "", conf)

	return Meta{ya.Path, refPieces}
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

func collectMeta(n *yaml.Node, pieces *[]PieceOfRef, ref string, conf config.Conf) {
	if n == nil {
		return
	}

	for i := 0; i < len(n.Content); i += 2 {
		key := n.Content[i]
		var value *yaml.Node = nil
		if i+1 < len(n.Content) {
			value = n.Content[i+1]
		}

		updatedRef := ref + "." + key.Value

		if strings.Contains(key.HeadComment, conf.MetaMarker) {
			*pieces = append(*pieces, PieceOfRef{clearMetaComment(key.HeadComment), updatedRef})
		}

		collectMeta(value, pieces, updatedRef, conf)
		//switch value.Kind {
		//case yaml.ScalarNode:
		//  updatedRef := updatedRef + "." + n.Value
		//  continue
		//case yaml.MappingNode:
		//  collectMeta(n.Content[i], pieces, updatedRef, conf)
		//}
	}
}

func clearMetaComment(comment string) string {
	// First line in every meta comment is meta marker. We can delete it
	out := make([]byte, 0)

	dataStartPoint := 0
	for i, ch := range comment {
		if ch == '\n' {
			if i+3 > len(comment) {
				panic("empty meta section")
			}
			// skip linebreak, comment mark, and empty space
			dataStartPoint = i + 3
			break
		}
	}

	for i := dataStartPoint; i < (len(comment) - 2); i++ {
		if comment[i] == '\n' && comment[i+1] == '#' {
			// Skip for linebreak, comment mark and empty space
			i += 2
			out = append(out, '\n')
			continue
		}
		out = append(out, comment[i])
	}

	out = append(out, comment[len(comment)-2])
	out = append(out, comment[len(comment)-1])

	return string(out)
}
