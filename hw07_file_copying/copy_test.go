package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const tmpPattern = "copy_test."

func TestCopy(t *testing.T) {
	tmpFile, err := os.CreateTemp("/tmp", tmpPattern)
	if err != nil {
		log.Fatal(err)
	}
	fromFileName := tmpFile.Name()
	defer os.Remove(fromFileName)
	content := []byte("temporary file content\n")
	contentLen := int64(len(content))
	if _, err := tmpFile.Write(content); err != nil {
		fmt.Println(err)
		return
	}

	if err := tmpFile.Close(); err != nil {
		fmt.Println(err)
		return
	}

	t.Run("offset more than files length returns error", func(t *testing.T) {
		err := Copy(fromFileName, "", contentLen+1, 0)
		require.Error(t, err, ErrOffsetExceedsFileSize)
	})
	t.Run("copy from /dev/urandom returns error", func(t *testing.T) {
		err := Copy("/dev/urandom", "", 0, 0)
		require.Error(t, err, ErrUnsupportedFile)
	})
	t.Run("copy from not existing file returns error", func(t *testing.T) {
		notExistingFileName := fromFileName + "randomstring"
		err := Copy(notExistingFileName, "", 0, 0)
		require.Error(t, err, ErrFileNotExist)
	})
}
