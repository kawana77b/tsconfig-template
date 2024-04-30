package basesrepo

import (
	"embed"
	"io/fs"
	"path"
	"path/filepath"
)

const (
	DEFAULT_JSON_FILENAME string = "recommended.json"
	baseDir               string = "bases/bases"
)

type BasesRepo struct {
	fs embed.FS
}

func NewBasesRepo(fs embed.FS) *BasesRepo {
	return &BasesRepo{
		fs: fs,
	}
}

func (b *BasesRepo) TemplateFiles() ([]TemplateFile, error) {
	files := []TemplateFile{}
	d, err := fs.ReadDir(b.fs, baseDir)
	if err != nil {
		return files, err
	}

	for _, v := range d {
		if v.IsDir() {
			continue
		}
		if filepath.Ext(v.Name()) != ".json" {
			continue
		}

		files = append(files, TemplateFile{
			Name: v.Name(),
			Path: path.Join(baseDir, v.Name()),
		})
	}

	return files, nil
}

// ReadFile reads a file from the file system.
func (b *BasesRepo) ReadFile(path string) ([]byte, error) {
	return b.fs.ReadFile(path)
}

type TemplateFile struct {
	Name string
	Path string
}

// GetFileMap returns a map of file names to paths and a slice of file names.
func GetFileMap(files []TemplateFile) (map[string]string, []string) {
	fileMap := make(map[string]string)
	var names []string

	for _, f := range files {
		fileMap[f.Name] = f.Path
		names = append(names, f.Name)
	}

	return fileMap, names
}
