package kong

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	var cli struct{}
	p, err := New(&cli, Name("name"), Description("description"), Writers(nil, nil), ExitFunction(nil))
	require.NoError(t, err)
	require.Equal(t, "name", p.Model.Name)
	require.Equal(t, "description", p.Model.Help)
	require.Nil(t, p.Stdout)
	require.Nil(t, p.Stderr)
	require.Nil(t, p.Exit)
}

func TestConfigLoading(t *testing.T) {
	first, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer first.Close()
	defer os.Remove(first.Name())
	second, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer second.Close()
	defer os.Remove(second.Name())

	var cli struct {
		Flag string `json:"flag,omitempty"`
	}

	cli.Flag = "first"
	err = json.NewEncoder(first).Encode(&cli)
	require.NoError(t, err)

	cli.Flag = ""
	err = json.NewEncoder(second).Encode(&cli)
	require.NoError(t, err)

	p := mustNew(t, &cli, Configuration(JSON, first.Name(), second.Name()))
	_, err = p.Parse(nil)
	require.NoError(t, err)
	require.Equal(t, "first", cli.Flag)
}
