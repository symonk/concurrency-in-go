package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// main demonstrates a much more advanced pipeline example.
//
// the example in the basic_pipeline package does not account
// for various scenarios such as error handling, or upstream
// channels that will not have all value(s) produced.
//
// For example, if a stage in the pipeline (non producer, non sink)
// decided it would only forward on messages that were divisible by
// 5 equally etc, this creates a scenario where the output data is
// significantly less than the input data.
//
// The example outlined below checks the checksum of all the files
// in this directory, a few .txt files.
func main() {
	done := make(chan struct{})
	defer close(done)
	root := "."
	m, err := md5All(root)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)

}

// result encapsulates the data for a single file
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// sumFilesStage walks the tree and digests each of the files in a
// goroutine, the results are sent to it's downstream channel.
// sumFilesStage will return on the first error.
func sumFilesStage(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	out := make(chan result)
	e := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			wg.Add(1)
			go func() {
				defer wg.Done()
				data, err := ioutil.ReadFile(path)
				select {
				case out <- result{path, md5.Sum(data), err}:
				case <-done:
				}
			}()
			select {
			case <-done:
				return errors.New("walking cancelled")
			default:
				return nil
			}
		})
		go func() {
			wg.Wait()
			close(out)
		}()
		e <- err
	}()

	return out, e
}

// md5All reads all the files in the current directory and returns
// a map for each file path (name) and an array of (16) bytes.
// if anything fails, an error is returned.
func md5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	done := make(chan struct{})
	defer close(done)

	in, errs := sumFilesStage(done, root)
	for reply := range in {
		if reply.err != nil {
			return nil, reply.err
		}
		m[reply.path] = reply.sum
	}
	if err := <-errs; err != nil {
		return nil, err
	}
	return m, nil
}

// merge is a generic fan in implementation.
// it launches N goroutines, dictated by the
// number of channels in the inbound slice.
func merge[T any](done <-chan struct{}, inbound ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	size := len(inbound)
	wg.Add(size)
	task := func(c <-chan T) {
		defer wg.Done()
		for {
			select {
			case out <- <-c:
			case <-done:
				return
			}
		}
	}

	for _, ch := range inbound {
		go task(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
