package source

import "go/ast"

func NewFileCrawlerByAF(file *ast.File) *FileCrawler {
	return &FileCrawler{
		file: file,
	}
}

type FileCrawler struct {
	file *ast.File
}

func (c *FileCrawler) SetPackageIfNotDefined(packageName string) {
	if nil == c.file.Name {
		c.file.Name = ast.NewIdent(packageName)
	}
}

func (c *FileCrawler) GetStruct(name string) *StructCrawler {
	return nil // TODO: implement me
}
