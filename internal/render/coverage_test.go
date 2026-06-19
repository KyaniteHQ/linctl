package render

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type failingWriter struct{}

func (writer failingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failed")
}

func Test_RenderScenarios_write_text_json_and_report_writer_errors(t *testing.T) {
	output := bytes.Buffer{}

	require.NoError(t, WriteLine(&output, "hello %s", "Omer"))
	require.Equal(t, "hello Omer\n", output.String())

	err := WriteLine(failingWriter{}, "hello")
	require.Error(t, err)
	require.Contains(t, err.Error(), "write line")

	err = WriteJSON(failingWriter{}, map[string]string{"hello": "Omer"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "write json")
}
