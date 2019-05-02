package files

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateProjectFile(curFolder string, projectName string) (err error) {
	targetFile := filepath.Join(curFolder, "templates", "content", fmt.Sprintf("projects-%s.html", projectName))
	f, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("can't create file %v, with error: %v", projectName, err)
	}
	f.Sync()

	defer f.Close()

	var tmplBytes []byte
	tmplFile := filepath.Join(curFolder, "templates", "layout", "projects_content.tmpl.html")

	tmplBytes, err = ioutil.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("can't read file %s error: %v", tmplFile, err)
	}

	_, err = io.Copy(f, bytes.NewBuffer(tmplBytes))
	if err != nil {
		return fmt.Errorf("can't copy content from %s to %s  error: %v", tmplFile, targetFile, err)
	}

	return
}
