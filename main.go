// sponge - Soaks up all input from stdin and writes it to a file or stdout
// https://git.sr.ht/~tslocum/sponge
//
// Originally written in C by Tollef Fog Heen
// https://joeyh.name/code/moreutils
//
// Ported by Trevor Slocum
// https://rocketnine.space
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

const BufferSize = 64 * 1024 // 64 KiB

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("sponge: ")
}

func main() {
	appendFile := flag.Bool("a", false, "append")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "%s [-a] [file]: soak up all input from stdin and write it to [file] or stdout\n", os.Args[0])
	}
	flag.Parse()

	writeFileName := flag.Arg(0)

	tmpFile, err := ioutil.TempFile("", "sponge-*.tmp")
	CheckError(err)
	defer func(tmpFile *os.File) {
		if tmpFile == nil {
			return // File was closed elsewhere
		}

		tmpFileName := tmpFile.Name()
		_ = tmpFile.Close()
		_ = os.Remove(tmpFileName)
	}(tmpFile)

	buf := make([]byte, BufferSize)

	_, err = io.CopyBuffer(tmpFile, os.Stdin, buf)
	CheckError(err)

	_, err = tmpFile.Seek(io.SeekStart, 0)
	CheckError(err)

	var out io.Writer
	if writeFileName != "" {
		regularFile := true
		fileInfo, err := os.Stat(writeFileName)
		if err != nil {
			if !os.IsNotExist(err) {
				CheckError(err)
			}
		} else {
			regularFile = fileInfo.Mode().IsRegular()
		}

		if regularFile && !*appendFile {
			tmpFileName := tmpFile.Name()
			tmpFile.Close()
			tmpFile = nil

			err = os.Rename(tmpFileName, writeFileName)
			if err != nil {
				// Fall back to mv
				cmd := exec.Command("mv", tmpFileName, writeFileName)
				err = cmd.Run()
			}
			CheckError(err)

			return
		}

		openFlags := os.O_WRONLY
		if regularFile {
			openFlags |= os.O_CREATE
		}
		if *appendFile {
			openFlags |= os.O_APPEND
		}

		writeFile, err := os.OpenFile(writeFileName, openFlags, 0644)
		CheckError(err)
		defer writeFile.Close()

		out = writeFile
	} else {
		out = os.Stdout
	}

	_, err = io.CopyBuffer(out, tmpFile, buf)
	CheckError(err)
}
