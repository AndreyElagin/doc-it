package main

import (
	"os"
	strings "strings"

	"doc-it/pkg/config"
	"doc-it/pkg/docit"
	"doc-it/pkg/errorutils"
	"doc-it/pkg/fsutils"
)

func main() {
	conf := config.Conf{
		IncludeFileTypes: []string{".yaml", ".yml"},
		MetaMarker:       "@doc-it",
		OutputDir:        "out",
	}

	yamls := docit.ReadYamls("yamls", conf)
	meta := make([]docit.Meta, 0)
	for _, y := range yamls {
		meta = append(meta, y.ToMeta(conf))
	}

	err := fsutils.CreateDirIfNotExist(conf.OutputDir)
	errorutils.Check(err)

	for _, m := range meta {
		err := os.WriteFile(
			conf.OutputDir+"/"+m.Path.FileName()+".md",
			[]byte(strings.Join(m.Comments, "\n\n")),
			0666,
		)
		errorutils.Check(err)
	}
}
