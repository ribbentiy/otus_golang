package main

import (
	"errors"
	"io"
	"os"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotExist          = errors.New("file not exist")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExist
		}
		return err
	}
	defer srcFile.Close()
	fileInfo, _ := srcFile.Stat()
	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	fileLength := fileInfo.Size()
	if fileLength < offset {
		return ErrOffsetExceedsFileSize
	}
	copyLen := limit
	if limit >= fileLength-offset || limit == 0 {
		copyLen = fileLength - offset
	}
	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	p := mpb.New()
	bar := p.AddBar(copyLen,
		mpb.AppendDecorators(decor.Percentage()),
	)
	proxyReader := bar.ProxyReader(srcFile)
	defer proxyReader.Close()

	destFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	_, err = io.CopyN(destFile, proxyReader, copyLen)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}
	defer destFile.Close()
	p.Wait()
	return nil
}
