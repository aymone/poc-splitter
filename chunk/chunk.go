package chunk

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

// Chunkprefix is to fix chunk name
// appended by underline to part number
const Chunkprefix = "temp_"

// Split dir file to chunks
func Split(dirpath, filename string, size int) error {
	filetobechunked := fmt.Sprintf("%s/%s", dirpath, filename)
	file, err := os.Open(filetobechunked)
	if err != nil {
		return err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(size)))

	for i := uint64(0); i < totalPartsNum; i++ {
		chunksize := int(math.Min(float64(size), float64(fileSize-int64(i*uint64(size)))))
		chunkbuffer := make([]byte, chunksize)

		file.Read(chunkbuffer)

		chunkID := fmt.Sprintf("%s%d", Chunkprefix, i)
		chunkfullname := fmt.Sprintf("%s/%s", dirpath, chunkID)
		_, err := os.Create(chunkfullname)
		if err != nil {
			return err
		}

		ioutil.WriteFile(chunkfullname, chunkbuffer, os.ModeAppend)
	}

	return nil
}

// Join dir chunks to file
func Join(dirpath, filename string) error {
	// get chunks info
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}

	// store chunk counter
	totalPartsNum := uint64(len(files) - 1)

	// create output file to append chunks
	outputfilename := fmt.Sprintf("%s/%s", dirpath, filename)
	_, err = os.Create(outputfilename)
	if err != nil {
		return err
	}

	// open output to append chunks in append mode
	file, err := os.OpenFile(outputfilename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	var writePosition int64
	for i := uint64(0); i < totalPartsNum; i++ {
		//read a chunk
		currentChunkFileName := fmt.Sprintf("%s%d", Chunkprefix, i)
		newFileChunk, err := os.Open(fmt.Sprintf("%s/%s", dirpath, currentChunkFileName))
		if err != nil {
			return err
		}

		defer newFileChunk.Close()
		chunkInfo, err := newFileChunk.Stat()
		if err != nil {
			return err
		}

		// calculate the bytes size of each chunk
		chunkSize := chunkInfo.Size()
		chunkBufferBytes := make([]byte, chunkSize)
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(newFileChunk)
		_, err = reader.Read(chunkBufferBytes)
		if err != nil {
			return err
		}

		// write/save buffer to disk
		_, err = file.Write(chunkBufferBytes)
		if err != nil {
			return err
		}

		//flush direct to disk
		file.Sync()

		// reset buffer
		chunkBufferBytes = nil
	}

	// now, we close the filename
	file.Close()

	return nil
}
