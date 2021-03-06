package cruncy

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"syscall"

	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// EnsureFileSave creates directory unless exists for a given file
func EnsureFileSave(fileName string) {
	pt := filepath.Dir(fileName)
	CreateDirUnlessExists(pt)
}

// Exists returns true if file/path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// CreateDirUnlessExists creates a directory if to do not exist
func CreateDirUnlessExists(path string) {
	if !Exists(path) {
		os.MkdirAll(path, os.ModeDir|0755)
	}
}

// Decompress decompresses a file using os utility gzip
func Decompress(fileName string) (string, error) {
	err := DoCmd("gzip", "-d", fileName)
	if err != nil {
		return "", err
	}
	return strings.Replace(fileName, ".gz", "", -1), nil
}

// DoCmd Runs a os command
func DoCmd(command string, arg ...string) error {
	cmd := exec.Command(command, arg...) // no need to call Output method here
	cmd.Stdout = os.Stdout               // instead use Stdout
	cmd.Stderr = os.Stderr               // attach Stderr as well

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Unable to execute %s with error %s", command, err)
	}
	return nil
}

// FileToReader gets a buffered reader from a fileName
func FileToReader(fileName string) (io.Reader, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("File %s do not exist", fileName)
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %s with error: %s", fileName, err)
	}

	return bufio.NewReader(f), nil
}

// DiskUsage disk usage of path/disk
func DiskUsage(path string, checkSize uint64) float64 {

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return 0.0
	}
	all := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)

	// Adding download size to the calculation.
	used := all - free - checkSize

	return (float64(used) / float64(all)) * 100.0
}

// IsDirectory is path a directory
func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// DeleteFile deletes a file
func DeleteFile(fileName string) error {
	err := os.Remove(fileName)
	if err != nil {
		return fmt.Errorf("unable to delete %s with error: %s", fileName, err)
	}

	return nil
}

// IsTextFile tells if a file is a text file
func IsTextFile(fileName string) (bool, error) {
	peekLen := 64

	f, err := os.Open(fileName)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("isTextFile Unable to open  %s with error %s", fileName, err))
	}
	defer f.Close()
	b := make([]byte, peekLen)
	_, err = f.Read(b)
	if err != nil && err != io.EOF {
		return false, errors.Wrap(err, fmt.Sprintf("isTextFile Unable to open  %s with error %s", fileName, err))
	}
	return IsTextData(b), nil
}

// IsTextData checks if a byte stream is text
func IsTextData(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	return strings.Contains(http.DetectContentType(data), "text/")
}
