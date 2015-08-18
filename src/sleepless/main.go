// sleepless - keep my pesky Toshiba HDD awake
// Copyright (c) 2015, Matteo Panella <morpheus@level28.org>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//    1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
//
//    2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"crypto/rand"
	"github.com/ncw/directio"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

// Lifted lock, stock and barrel from FiloSottile's whosthere
func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Perform post-exit cleanup
func cleanup(fd *os.File, path string) {
	// Don't bother with error checking
	fd.Close()
	os.Remove(path)
}

func main() {
	// Get hold of the current working dir
	cwd, err := os.Getwd()
	fatalIfErr(err)

	// Create a temporary file
	fd, err := ioutil.TempFile(cwd, "sleepless")
	fatalIfErr(err)
	// We just need the file name, we're going to reopen it with the right flags
	path := fd.Name()
	fd.Close()

	// Allocate a block which will be filled with garbage on each iteration
	garbage := directio.AlignedBlock(directio.BlockSize)

	// Get ready
	out, err := directio.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0600)
	fatalIfErr(err)

	// Setup signal handling
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		<-sigint
		cleanup(out, path)
		os.Exit(0)
	}()

	// Start the main loop
	for {
		// Garbage in...
		if _, err = rand.Read(garbage); err != nil {
			cleanup(out, path)
			log.Fatal(err)
		}
		// ... garbage out
		if _, err = out.Write(garbage); err != nil {
			cleanup(out, path)
			log.Fatal(err)
		}
		// Rewind the file to the start so that it does not grow
		// indefintely
		if _, err = out.Seek(0, os.SEEK_SET); err != nil {
			cleanup(out, path)
			log.Fatal(err)
		}
		// Take a nap
		time.Sleep(10 * time.Second)
	}
}
