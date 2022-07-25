package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("running without command should returns 1", func(t *testing.T) {
		code := RunCmd([]string{}, nil)
		require.Equal(t, failCode, code)
	})
	t.Run("positive test", func(t *testing.T) {
		env := make(Environment)
		envName := "Bar"
		envVal := "Baz"
		env[envName] = EnvValue{Value: envVal}
		com := []string{"ls", "-l"}
		code := RunCmd(com, env)
		require.Equal(t, okCode, code)
	})
	t.Run("test with error", func(t *testing.T) {
		env := make(Environment)
		com := []string{"ls", "/tmp/notExist/definitelyNotExist"}
		code := RunCmd(com, env)
		require.Equal(t, failCode, code)
	})
}
