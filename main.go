package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/akhy/idbank/pkg/idbank"
)

func main() {
	sourceFiles := os.Args[1:]

	for _, file := range sourceFiles {
		export(file)
	}
}

func export(sourceFile string) {
	baseName := strings.TrimSuffix(sourceFile, ".pdf")
	year := 2023
	currency := "IDR"

	cmd := exec.Command("camelot",
		"--format", "json",
		"--output", baseName+".json",
		"--pages", "all",
		"stream",
		sourceFile,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	files, err := filepath.Glob(fmt.Sprintf("%s-*.json", baseName))
	if err != nil {
		panic(err)
	}

	bcaParser, err := idbank.NewBcaParser(idbank.BcaParserOpts{
		Year:     year,
		Currency: currency,
	})
	if err != nil {
		panic(err)
	}

	records := []idbank.Record{}

	for _, file := range files {
		result, err := idbank.ParseFile(bcaParser, file)
		if err != nil {
			panic(err)
		}

		records = append(records, result...)
	}

	exporter := idbank.NewCSVFileExporter(idbank.CSVFileExporterOpts{
		FilePath: baseName + ".csv",
	})

	if err := exporter.Export(records); err != nil {
		panic(err)
	}
}
