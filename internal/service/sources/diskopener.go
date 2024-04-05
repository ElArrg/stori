package sources

import (
	"io"
	"os"
)

type DiskOpener struct {
}

func NewDiskOpener() *DiskOpener {
	return &DiskOpener{}
}

func (c *DiskOpener) OpenFromSource(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}
