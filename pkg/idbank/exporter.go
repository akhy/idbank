package idbank

type Exporter interface {
	Export(records []Record) error
}
