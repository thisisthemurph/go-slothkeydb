package repository

import (
	"errors"
	"log"
	"os"
	"slothkeydb/utils"
)

type LogReader interface {
	Read(k string) (string, error)
}

type LogWriter interface {
	Write(k, v string, entryType logEntryType)
}

type LogReaderWriter interface {
	LogReader
	LogWriter
}

type logFile struct {
	path string
}

// Writes to the log file, the entryType indicates the type of entry to be written
func (lf logFile) Write(k, v string, entryType logEntryType) {
	switch entryType {
	case LET_Deleted:
		lf.writeDeletedEntry(k)
		break
	case LET_Live:
		lf.writeLiveEntry(k, v)
		break
	}
}

// Writes a deletion to the log file
func (lf logFile) writeDeletedEntry(k string) {
	file, err := os.OpenFile(lf.path, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		log.Println(err)
		log.Fatalln("Could not open log file.")
	}

	defer file.Close()

	keyBytes := []byte(k)
	keyLen := byte(len(keyBytes))
	valType := byte(LET_Deleted)
	valLen := byte(1)

	var entry []byte
	entry = append(entry, keyLen)
	entry = append(entry, keyBytes...)
	entry = append(entry, valLen)
	entry = append(entry, valType)

	if _, err := file.Write(entry); err != nil {
		log.Fatalln(err)
	}
}

// Writes a live entry to the log file
func (lf logFile) writeLiveEntry(k, v string) {
	file, err := os.OpenFile(lf.path, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		log.Println(err)
		log.Fatalln("Could not open log file.")
	}

	defer file.Close()

	keyBytes := []byte(k)
	valType := byte(LET_Live)
	valBytes := []byte(v)

	keyLen := byte(len(keyBytes))
	valLen := byte(len(valBytes) + 1)
	
	var entry []byte
	entry = append(entry, keyLen) 
	entry = append(entry, keyBytes...) 
	entry = append(entry, valLen) 
	entry = append(entry, valType) 
	entry = append(entry, valBytes...)

	if _, err := file.Write(entry); err != nil {
		log.Fatalln(err)
	}
}

// Read a specific key from the log
func (lf logFile) Read(k string) (string, error) {
	entries := lf.getEntries()
	entry, ok := entries[k]

	// If the key does not exist or has been deleted, return unknown key
	if !ok || entry.Type == LET_Deleted {
		return "", makeUnknownKeyError(k)
	} 

	return entry.Value, nil
}

// Retrieves the log entries from the log file
func (lf logFile) getEntries() map[string]logEntry {
	data := lf.readLogData()
	return processLogData(data)
}

// Retrieves the data (bytes) from the log file
func (lf logFile) readLogData() []byte {
	file, err := os.OpenFile(lf.path, os.O_RDONLY, 0664)
	if err != nil {
		log.Println(err)
		log.Fatalln("Could not open log file.")
	}

	defer file.Close()

	data, err := os.ReadFile(lf.path)
	if err != nil {
		log.Println(err)
		log.Fatalln("Could not read log file.")
	}

	return data
}

// Builds the log file struct and creates the log file on disk if required
func MakeLogFile(path string) (LogReaderWriter, error) {
	if !utils.PathExists(path) {
		_, err := os.Create(path)
		if err != nil {
			return nil, errors.New("Could not create store log in that location.")
		}
	}

	return logFile{path}, nil
}

// Processes the bytes from the log data
func processLogData(data []byte) map[string]logEntry {
	results := make(map[string]logEntry)

	cursor := 0
	for cursor < len(data) {
		keyLen := int8(data[cursor])
		base := cursor + int(keyLen)
		valLen := int64(data[base + 1])

		entry := logEntry{
			KeyLength:   keyLen,
			Key:         string(data[cursor + 1:cursor + int(keyLen) + 1]), 
			ValueLength: valLen, 
			Type:        logEntryType(int8(data[base + 2])),
			Value:       string(data[base + 3:base + 2 + int(valLen)]),
		}
		
		results[entry.Key] = entry
		cursor += entry.Length()
	}

	return results
}