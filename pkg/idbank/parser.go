package idbank

import (
	"io"
	"os"
)

type StatementParser interface {
	Parse(r io.Reader) ([]Record, error)
}

func ParseFile(parser StatementParser, path string) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parser.Parse(f)
}
