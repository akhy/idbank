# idbank
Indonesian bank statement parser and exporter

Written in Go, but requires `camelot-py` python package.

I personally used this tool to convert bank statements to import to Firefly III. At least this project will try to support these banks:

- Bank BCA
- Bank Jago
- GoPay

Steps:

1. Parse PDF to JSON using camelot-py
2. Transform unstructured JSON to `idbank.Record` struct
3. Output the structs as CSV, along with JSON import config for Firefly III
