package sources

import "io"

type Opener interface {
	OpenFromSource(source string) (io.ReadCloser, error)
}
