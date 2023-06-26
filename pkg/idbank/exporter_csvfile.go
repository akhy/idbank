package idbank

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type csvFileExporter struct {
	filePath string
}

// Export implements Exporter.
func (c *csvFileExporter) Export(records []Record) error {
	f, err := os.Create(c.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)

	for _, record := range records {
		amount := record.Amount

		if record.Debit {
			amount = -amount
		}

		writer.Write([]string{
			record.Date.String(),
			fmt.Sprintf("%f", amount),
			fmt.Sprintf(`%s`, strings.ReplaceAll(record.Notes, "\n", " ")),
		})
	}

	writer.Flush()

	return nil
}

type CSVFileExporterOpts struct {
	FilePath string
}

func NewCSVFileExporter(opts CSVFileExporterOpts) Exporter {
	return &csvFileExporter{
		filePath: opts.FilePath,
	}
}
