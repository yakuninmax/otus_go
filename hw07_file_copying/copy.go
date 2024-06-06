package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile         = errors.New("unsupported file")
	ErrOffsetExceedsFileSize   = errors.New("offset exceeds file size")
	ErrFileNotFound            = errors.New("source file not found")
	ErrOpenSourceFile          = errors.New("could not open source file")
	ErrGetOffset               = errors.New("could not get offset for source file")
	ErrCreatingDestinationFile = errors.New("creating destination file failed")
	ErrCopyFile                = errors.New("could not copy file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Get file info and size.
	sourceFileInfo, err := os.Stat(fromPath)
	// Check if file exists.
	if err != nil {
		return ErrFileNotFound
	}

	// Check if regular file.
	if !sourceFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// Check if offset greater file len.
	sourceFileSize := sourceFileInfo.Size()

	if offset > sourceFileSize {
		return ErrOffsetExceedsFileSize
	}

	// Create destination file.
	destinationFile, err := os.Create(toPath)
	if err != nil {
		return ErrCreatingDestinationFile
	}
	defer destinationFile.Close()

	// Open source file.
	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return ErrOpenSourceFile
	}

	// Set starting position using offset.
	_, err = sourceFile.Seek(offset, 0)
	if err != nil {
		return ErrGetOffset
	}

	// Get bytes count to copy except offset.
	bytesCount := sourceFileSize - offset
	if limit != 0 && limit < bytesCount {
		bytesCount = limit
	}

	// Create progress bar.
	bar := pb.New(int(bytesCount)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()

	// Create proxy reader.
	reader := bar.NewProxyReader(sourceFile)

	// Copy data from source to destination file.
	_, err = io.CopyN(destinationFile, reader, bytesCount)
	if err != nil {
		return ErrCopyFile
	}
	defer sourceFile.Close()

	// Finish progress bar.
	bar.Finish()

	return nil
}
