package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Reference file name.
	inputFile := path.Join("testdata", "input.txt")

	// Test case structure.
	type testCase struct {
		name            string
		inputFile       string
		destinationFile string
		referenceFile   string
		offset          int64
		limit           int64
		expectedError   error
	}

	// Tests list.
	tests := []testCase{
		{
			name:          "offset 0 limit 0",
			inputFile:     inputFile,
			referenceFile: path.Join("testdata", "out_offset0_limit0.txt"),
			offset:        0,
			limit:         0,
		},
		{
			name:          "offset 0 limit 10",
			inputFile:     inputFile,
			referenceFile: path.Join("testdata", "out_offset0_limit10.txt"),
			offset:        0,
			limit:         10,
		},
		{
			name:          "offset 0 limit 1000",
			inputFile:     inputFile,
			referenceFile: path.Join("testdata", "out_offset0_limit1000.txt"),
			offset:        0,
			limit:         1000,
		},
		{
			name:          "offset 0 limit 10000",
			inputFile:     inputFile,
			referenceFile: path.Join("testdata", "out_offset0_limit10000.txt"),
			offset:        0,
			limit:         10000,
		},
		{
			name:          "offset 100 limit 1000",
			inputFile:     inputFile,
			referenceFile: path.Join("testdata", "out_offset100_limit1000.txt"),
			offset:        100,
			limit:         1000,
		},
		{
			name:          "offset 6000 limit 1000",
			inputFile:     inputFile,
			referenceFile: path.Join("testdata", "out_offset6000_limit1000.txt"),
			offset:        6000,
			limit:         1000,
		},
		{
			name:          "irregular source file",
			inputFile:     "/dev/urandom",
			expectedError: ErrUnsupportedFile,
		},
		{
			name:          "offset exceeds file size",
			inputFile:     inputFile,
			offset:        100000,
			expectedError: ErrOffsetExceedsFileSize,
		},
		{
			name:          "source file not found",
			inputFile:     inputFile + "hhgjh",
			expectedError: ErrFileNotFound,
		},
		{
			name:          "invalid offset value",
			inputFile:     inputFile + "hhgjh",
			expectedError: ErrInvalidOffsetValue,
			offset:        -5,
		},
		{
			name:          "invalid limit value",
			inputFile:     inputFile,
			expectedError: ErrInvalidLimitValue,
			limit:         -5,
		},
		{
			name:            "same file",
			inputFile:       path.Join("testdata", "samefile"),
			destinationFile: path.Join("testdata", "samefile"),
			expectedError:   ErrSameFile,
		},
	}

	// Run tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var destinationFile string

			if test.destinationFile != "" {
				destinationFile = test.destinationFile
			} else {
				// Create temp file in temp dir.
				tempFile, err := os.CreateTemp("", "ddtest*")
				if err != nil {
					log.Println(err)
				}
				destinationFile = tempFile.Name()
				// Cleanup.
				defer os.Remove(destinationFile)
			}

			// Run copy function.
			err := Copy(test.inputFile, destinationFile, test.offset, test.limit)

			// Check results.
			if err != nil {
				// Check for errors.
				require.EqualError(t, err, test.expectedError.Error())
			} else {
				// Compare reference & temp copied file.
				_, err := compareFiles(t, inputFile, destinationFile)
				require.Nil(t, err)
			}
		})
	}
}

// Compare files.
func compareFiles(t *testing.T, sourceFileName, destinationFileName string) (bool, error) {
	// Define as testing helper.
	t.Helper()

	// Read source file.
	sourceFile, err := os.ReadFile(sourceFileName)
	if err != nil {
		return false, err
	}

	// Read destination file.
	destinationFile, err := os.ReadFile(destinationFileName)
	if err != nil {
		return false, err
	}

	// Compare files.
	return bytes.Equal(sourceFile, destinationFile), nil
}
