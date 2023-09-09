package common

import "io"

type Collection[T any] struct {
	Stats              []T
	CsvHeadersProvider func() []string
	CsvRowTransformer  func(T) []string

	DisplayHeadersProvider func() []interface{}
	DisplayRowTransformer  func(T) []interface{}

	VerboseOutputTransformer func(io.Writer, []T)
}

func (c Collection[T]) ToCsvOutput() *CsvOutput {
	csvOutput := CsvOutput{
		Headers: c.CsvHeadersProvider(),
	}

	for _, s := range c.Stats {
		csvOutput.Rows = append(csvOutput.Rows, c.CsvRowTransformer(s))
	}

	return &csvOutput
}

func (c Collection[T]) ToDisplayOutput() *DisplayOutput {
	displayOutput := DisplayOutput{
		Headers: c.DisplayHeadersProvider(),
	}

	for _, s := range c.Stats {
		displayOutput.Rows = append(displayOutput.Rows, c.DisplayRowTransformer(s))
	}

	return &displayOutput
}

func (c Collection[T]) ToVerboseOutput(w io.Writer) {
	c.VerboseOutputTransformer(w, c.Stats)
}
