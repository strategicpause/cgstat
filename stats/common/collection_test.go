package common

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestCollection_ToCsvOutput(t *testing.T) {
	// Given
	headers := []string{"header"}
	collection := Collection[string]{
		Stats: []string{"a", "b", "c"},
		CsvHeadersProvider: func() []string {
			return headers
		},
		CsvRowTransformer: func(s string) []string {
			return []string{s}
		},
	}

	// When
	csvOutput := collection.ToCsvOutput()

	// Then
	assert.Equal(t, headers, csvOutput.Headers)
	assert.Equal(t, [][]string{
		{"a"}, {"b"}, {"c"},
	}, csvOutput.Rows)
}

func TestCollection_ToDisplayOutput(t *testing.T) {
	// Given
	headers := []interface{}{"header"}
	collection := Collection[string]{
		Stats: []string{"a", "b", "c"},
		DisplayHeadersProvider: func() []interface{} {
			return headers
		},
		DisplayRowTransformer: func(s string) []interface{} {
			return []interface{}{s}
		},
	}

	// When
	displayOutput := collection.ToDisplayOutput()

	// Then
	assert.Equal(t, headers, displayOutput.Headers)
	assert.Equal(t, [][]interface{}{
		{"a"}, {"b"}, {"c"},
	}, displayOutput.Rows)
}

type FakeWriter struct {
	writtenData [][]byte
}

func (f *FakeWriter) Write(p []byte) (n int, err error) {
	f.writtenData = append(f.writtenData, p)

	return 0, nil
}

func (f *FakeWriter) GetWrittenData() []string {
	var writtenStrings []string

	for _, data := range f.writtenData {
		writtenStrings = append(writtenStrings, string(data))
	}

	return writtenStrings
}

func TestCollection_ToVerboseOutput(t *testing.T) {
	data := []string{"a", "b", "c"}
	// Given
	collection := Collection[string]{
		Stats: data,
		VerboseOutputTransformer: func(writer io.Writer, strings []string) {
			for _, str := range strings {
				writer.Write([]byte(str))
			}
		},
	}
	fakeWriter := &FakeWriter{}

	// When
	collection.ToVerboseOutput(fakeWriter)

	// Then
	assert.Equal(t, data, fakeWriter.GetWrittenData())
}
