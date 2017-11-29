package cruncy

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// Inmemory database over processed files. Mutex protected

// ParlyDb handles prosessed files
type ParlyDb struct {
	sync.RWMutex
	Name  string
	Lines []string
}

// NewParlyDb constructs a new struct to handle file proessesing
func NewParlyDb(name string) *ParlyDb {
	db := ParlyDb{Name: name}

	db.Lines = make([]string, 1)
	db.ReadLines()
	return &db
}

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
func (db *ParlyDb) ReadLines() error {
	db.Lock()
	defer db.Unlock()

	file, err := os.Open(db.Name)
	if err != nil {
		return err
	}
	defer file.Close()

	db.Lines = make([]string, 1)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		db.Lines = append(db.Lines, line)
	}
	return scanner.Err()
}

// WriteLines writes the lines to the given file.
func (db *ParlyDb) WriteLines() error {

	file, err := os.Create(db.Name)
	if err != nil {
		return err
	}
	defer file.Close()

	db.Lock()
	defer db.Unlock()

	w := bufio.NewWriter(file)
	for _, line := range db.Lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

// Add - Adding file to is Prossesed list
func (db *ParlyDb) Add(fileName string) {
	db.Lock()
	defer db.Unlock()
	db.Lines = append(db.Lines, fileName)
}

// IsProsessed - Checks if a file is already prosessed
func (db *ParlyDb) IsProsessed(fileName string) bool {
	db.RLock()
	db.RUnlock()
	// fmt.Println(db)
	if len(db.Lines) == 0 {
		return false
	}

	for _, line := range db.Lines {
		if line == fileName {
			return true
		}
	}
	return false
}
