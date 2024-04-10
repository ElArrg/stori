package sources

import (
	"io"
	"os"
	"path"
)

type DiskOpener struct {
}

func NewDiskOpener() *DiskOpener {
	return &DiskOpener{}
}

func (c *DiskOpener) OpenFromSource(filePath string) (io.ReadCloser, error) {
	p := path.Clean(filePath)
	return os.Open(p)
}
