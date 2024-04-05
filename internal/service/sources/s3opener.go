package sources

import "io"

type S3Opener struct {
}

func NewS3Opener() *S3Opener {
	return &S3Opener{}
}

func (s *S3Opener) OpenFromSource(source string) (io.ReadCloser, error) {
	return nil, nil
}
