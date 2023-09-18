package rows

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type Writer interface {
	Write([]string) error
}

func NewTabRowWriter(w *tabwriter.Writer) Writer {
	if w == nil {
		w = tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)
	}
	return &TabRowWriter{*w}
}

type TabRowWriter struct {
	tabwriter.Writer
}

func (w *TabRowWriter) Write(record []string) error {
	fmt.Fprintln(&w.Writer, strings.Join(record, "\t"))
	return nil
}
