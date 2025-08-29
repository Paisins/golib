package utils

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"sync"
)

func Stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	for i := skip; ; i++ {   // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "exception occurred in %s[%s:%d]\n", runtime.FuncForPC(pc).Name(), file, line)
	}
	return buf.Bytes()
}

func GoAndWait(handlers ...func()) error {
	var (
		wg   sync.WaitGroup
		once sync.Once
		err  error
	)
	for _, f := range handlers {
		wg.Add(1)
		go func(handler func()) {
			defer func() {
				if e := recover(); e != nil {
					stack := Stack(1)
					fmt.Printf("go exception err:%v, %s\n", e, stack)
					once.Do(func() {
						err = errors.New("exception found in call handlers")
					})
				}
				wg.Done()
			}()
			handler()
		}(f)
	}
	wg.Wait()
	return err
}
