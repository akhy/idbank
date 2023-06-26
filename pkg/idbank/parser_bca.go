package idbank

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
)

type bcaRawRecord struct {
	date    string
	debit   bool
	notes   []string
	numbers []string
}

type bcaParser struct {
	year     int
	currency string
}

func (c *bcaParser) cellsFromFile(r io.Reader) (cells []string, err error) {
	var content []map[string]string
	if err := json.NewDecoder(r).Decode(&content); err != nil {
		return nil, err
	}

outer:
	for _, row := range content {
		if row["0"] == "TANGGAL" {
			continue
		}

	inner:
		for i := 0; i < len(row); i++ {
			cell := row[fmt.Sprintf("%d", i)]
			if cell == "" {
				continue inner
			}

			if cell == "SALDO AWAL" {
				continue outer
			}

			cells = append(cells, cell)
		}
	}

	return cells, err
}

func (c *bcaParser) rawRecordsFromCells(cells []string) (rawRecords []bcaRawRecord) {
	dateRegex := regexp.MustCompile(`^\d\d\/\d\d$`)
	moneyRegex := regexp.MustCompile(`^[\d,]+\.\d\d$`)

	begIdx := map[int]bool{}
	for i, cell := range cells {
		if dateRegex.MatchString(cell) {
			begIdx[i] = true
		}
	}

	// pre-allocate as we already know exactly how many records
	rawRecords = make([]bcaRawRecord, len(begIdx))

	cur := -1
	for i, cell := range cells {
		if begIdx[i] {
			cur += 1
			rawRecords[cur].date = cell
		}

		if cur < 0 {
			continue
		}

		if !begIdx[i] {
			if cell == "DB" {
				rawRecords[cur].debit = true
			} else if moneyRegex.MatchString(cell) {
				rawRecords[cur].numbers = append(rawRecords[cur].numbers, strings.ReplaceAll(cell, ",", ""))
			} else {
				rawRecords[cur].notes = append(rawRecords[cur].notes, strings.Split(cell, "\n")...)
			}
		}
	}

	return rawRecords
}

func (c *bcaParser) createRecord(raw bcaRawRecord) (*Record, error) {
	daymonth := strings.Split(raw.date, "/")
	if len(daymonth) != 2 {
		return nil, fmt.Errorf("invalid date: %s", raw.date)
	}

	day, err := strconv.ParseInt(daymonth[0], 10, 32)
	if err != nil {
		return nil, err
	}

	month, err := strconv.ParseInt(daymonth[1], 10, 32)
	if err != nil {
		return nil, err
	}

	if len(raw.numbers) == 0 {
		return nil, nil
	}

	amountFloat, err := strconv.ParseFloat(raw.numbers[0], 64)
	if err != nil {
		return nil, err
	}

	notes := append([]string{}, raw.notes...)
	for _, number := range raw.numbers[1:] {
		notes = append(notes, number)
	}

	return &Record{
		Date: civil.Date{
			Year:  c.year,
			Day:   int(day),
			Month: time.Month(month),
		},
		Amount: amountFloat,
		Unit:   c.currency,
		Debit:  raw.debit,
		Notes:  strings.Join(notes, "\n"),
	}, nil
}

// Parse implements StatementParser.
func (c *bcaParser) Parse(r io.Reader) ([]Record, error) {
	cells, err := c.cellsFromFile(r)
	if err != nil {
		return nil, err
	}

	records := []Record{}

	for _, raw := range c.rawRecordsFromCells(cells) {
		record, err := c.createRecord(raw)
		if err != nil {
			return nil, err
		}

		if record != nil {
			records = append(records, *record)
		}
	}

	return records, nil
}

type BcaParserOpts struct {
	Year     int
	Currency string
}

func NewBcaParser(opts BcaParserOpts) (StatementParser, error) {
	return &bcaParser{
		year:     opts.Year,
		currency: opts.Currency,
	}, nil
}
