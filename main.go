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
	"os"
)

const BufferSize = 64 * 1024 // 64 KiB

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		writeFileName = ""
		regularFile   = true
		appendFile    = false

		buffer = make([]byte, BufferSize)

		writeFile   *os.File
		tmpFile     *os.File
		tmpFileName string
		err         error
	)

	flag.BoolVar(&appendFile, "a", false, "append")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "%s [-a] [file]: soak up all input from stdin and write it to [file] or stdout\n", os.Args[0])
	}
	flag.Parse()
	writeFileName = flag.Arg(0)

	tmpFile, err = ioutil.TempFile("", "sponge-*.tmp")
	CheckError(err)
	tmpFileName = tmpFile.Name()

	_, err = io.CopyBuffer(tmpFile, os.Stdin, buffer)
	if err != nil {
		tmpFile.Close()
		_ = os.Remove(tmpFileName)
		CheckError(err)
	}

	_, err = tmpFile.Seek(io.SeekStart, 0)
	if err != nil {
		tmpFile.Close()
		_ = os.Remove(tmpFileName)
		CheckError(err)
	}

	if writeFileName == "" {
		_, err = io.CopyBuffer(os.Stdout, tmpFile, buffer)

		tmpFile.Close()
		_ = os.Remove(tmpFileName)
		CheckError(err)

		return
	}

	if fi, err := os.Stat(writeFileName); err == nil {
		regularFile = fi.Mode().IsRegular()
	} else if !os.IsNotExist(err) {
		tmpFile.Close()
		_ = os.Remove(tmpFileName)
		CheckError(err)
	}

	if regularFile && !appendFile {
		tmpFile.Close()
		err = os.Rename(tmpFileName, writeFileName)
		CheckError(err)

		return
	}

	openFlags := os.O_WRONLY
	if regularFile {
		openFlags |= os.O_CREATE
	}
	if appendFile {
		openFlags |= os.O_APPEND
	}

	writeFile, err = os.OpenFile(writeFileName, openFlags, 0644)
	if err != nil {
		tmpFile.Close()
		_ = os.Remove(tmpFileName)
		CheckError(err)
	}

	_, err = io.CopyBuffer(writeFile, tmpFile, buffer)

	tmpFile.Close()
	writeFile.Close()
	_ = os.Remove(tmpFileName)
	CheckError(err)
}
