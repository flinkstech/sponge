// sponge - Soaks up all input from stdin and writes it to a file or stdout
// https://gitlab.com/tslocum/sponge
//
// Originally written in C by Tollef Fog Heen
// https://joeyh.name/code/moreutils
//
// Ported by Trevor Slocum
// https://rocketnine.space

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// BufferSize determines the size of the input/output buffer. The buffer is
// used directly when soaking initially. If the input is larger than the buffer
// it is written to a temporary file instead. The buffer is reused when
// writing to and reading from the temporary file.
const BufferSize = 64 * 1024 // 64 KiB

var (
	buf            = make([]byte, BufferSize)
	bufInitialRead int
	appendFile     bool
	tmpFile        *os.File
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("sponge: ")
}

func main() {
	flag.BoolVar(&appendFile, "a", false, "append")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s [-a] [file]: soak up all input from stdin and write it to [file] or stdout\n", os.Args[0])
	}
	flag.Parse()

	err := soak()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanUp()

	err = squeeze(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
}

func soak() error {
	var err error
	bufInitialRead, err = io.ReadFull(os.Stdin, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return err
	} else if bufInitialRead == BufferSize && err != io.ErrUnexpectedEOF {
		tmpFile, err = ioutil.TempFile("", "sponge-*.tmp")
		if err != nil {
			return err
		}

		_, err = tmpFile.Write(buf[:bufInitialRead])
		if err != nil {
			tmpFile.Close()
			return err
		}

		_, err = io.CopyBuffer(tmpFile, os.Stdin, buf)
		if err != nil {
			tmpFile.Close()
			return err
		}

		_, err = tmpFile.Seek(io.SeekStart, 0)
		if err != nil {
			tmpFile.Close()
			return err
		}
	}

	return nil
}

func squeeze(filename string) error {
	var out io.Writer
	if filename != "" {
		regularFile := true
		fileInfo, err := os.Stat(filename)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			regularFile = fileInfo.Mode().IsRegular()
		}

		if tmpFile != nil && regularFile && !appendFile {
			tmpFile.Close()
			// Rename temporary file
			err = os.Rename(tmpFile.Name(), filename)
			if err != nil {
				// Fall back to mv
				cmd := exec.Command("mv", tmpFile.Name(), filename)
				err = cmd.Run()
			}
			tmpFile = nil
			return err
		}

		openFlags := os.O_WRONLY
		if regularFile {
			openFlags |= os.O_CREATE
		}
		if appendFile {
			openFlags |= os.O_APPEND
		} else if regularFile {
			openFlags |= os.O_TRUNC
		}

		writeFile, err := os.OpenFile(filename, openFlags, 0644)
		if err != nil {
			return err
		}
		defer writeFile.Close()

		out = writeFile
	} else {
		out = os.Stdout
	}

	var err error
	if tmpFile != nil {
		_, err = io.CopyBuffer(out, tmpFile, buf)
	} else {
		_, err = out.Write(buf[:bufInitialRead])
	}
	return err
}

func cleanUp() {
	if tmpFile == nil {
		return
	}

	tmpFile.Close()
	err := os.Remove(tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}
}
