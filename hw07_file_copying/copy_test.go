package main

import (
	"crypto/md5"
	"io"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Reference file name.
	inputFile := path.Join("testdata", "input.txt")

	// Test case structure.
	type testCase struct {
		name          string
		inputFile     string
		referenceFile string
		offset        int64
		limit         int64
		expectedError error
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
	}

	// Run tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create temp file in temp dir.
			tempFile, err := os.CreateTemp("", "ddtest*")
			if err != nil {
				log.Println(err)
			}
			// Cleanup.
			defer os.Remove(tempFile.Name())

			// Run copy function.
			err = Copy(test.inputFile, tempFile.Name(), test.offset, test.limit)

			// Check results.
			if err != nil {
				// Check for errors.
				require.EqualError(t, err, test.expectedError.Error())
			} else {
				// Compare reference & temp copied file.
				assert.Equal(t, md5checksum(t, test.referenceFile), md5checksum(t, tempFile.Name()))
			}
		})
	}
}

// Calculate MD5 checksum.
func md5checksum(t *testing.T, filename string) []byte {
	// Define as testing helper.
	t.Helper()

	// Open file.
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// Create md5 hash.
	md5Hash := md5.New()

	// Get file md5 hash.
	_, err = io.Copy(md5Hash, file)
	if err != nil {
		log.Println(err)
	}

	// Return md5 checksum.
	return md5Hash.Sum(nil)
}
