package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	tempDir     = "/tmp"
	tempPattern = "env_reader_tmp."
)

func TestReadDir(t *testing.T) {
	t.Run("empty dir name should return error", func(t *testing.T) {
		_, err := ReadDir("")
		require.ErrorIs(t, err, ErrUnsupportedInputPath)
	})
	t.Run("not existing directory return path error", func(t *testing.T) {
		_, err := ReadDir("./someRandomText")
		require.Error(t, err)
	})
	t.Run("not dir should return error", func(t *testing.T) {
		file, err := os.CreateTemp(tempDir, tempPattern)
		defer func(n string) {
			err := os.Remove(n)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(file.Name())
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		n := file.Name()
		_, err = ReadDir(n)
		require.ErrorIs(t, err, ErrUnsupportedInputPath)
	})
	t.Run("empty dir should returns empty map", func(t *testing.T) {
		emptyDir, err := os.MkdirTemp(tempDir, tempPattern)
		if err != nil {
			fmt.Println(err)
		}
		defer func(d string) {
			err := os.RemoveAll(d)
			if err != nil {
				fmt.Println(err)
			}
		}(emptyDir)
		env, err := ReadDir(emptyDir)
		require.NoError(t, err)
		require.Empty(t, env)
	})
	t.Run("non alphanumeric symbols in file name should return error", func(t *testing.T) {
		dir, err := os.MkdirTemp(tempDir, tempPattern)
		if err != nil {
			fmt.Println(err)
		}
		defer func(d string) {
			err := os.RemoveAll(d)
			if err != nil {
				fmt.Println(err)
			}
		}(dir)
		_, err = os.Create(dir + "/" + "some-filename-with-unexpected-symbols!")
		if err != nil {
			fmt.Println(err)
		}
		_, err = ReadDir(dir)
		require.ErrorIs(t, err, ErrUnsupportedSymbolsInFilename)
	})
	t.Run("positive test", func(t *testing.T) {
		dir, err := os.MkdirTemp(tempDir, tempPattern)
		if err != nil {
			fmt.Println(err)
		}
		defer func(d string) {
			err := os.RemoveAll(d)
			if err != nil {
				fmt.Println(err)
			}
		}(dir)
		key1 := "ValidWithVal"
		val1 := "ValidVal"
		key2 := "ValidWithoutVal"
		f1, err := os.Create(dir + "/" + key1)
		if err != nil {
			fmt.Println(err)
		}
		_, err = f1.WriteString(val1)
		if err != nil {
			fmt.Println(err)
		}
		_, err = os.Create(dir + "/" + key2)
		if err != nil {
			fmt.Println(err)
		}
		env, err := ReadDir(dir)
		require.NoError(t, err)
		env1, ok1 := env[key1]
		require.True(t, ok1)
		require.False(t, env1.NeedRemove)
		require.Equal(t, val1, env1.Value)
		env2, ok2 := env[key2]
		require.True(t, ok2)
		require.True(t, env2.NeedRemove)
		require.Empty(t, env2.Value)
	})
}
