package logger

import (
	"os"
	"strconv"
	"sync"
	"testing"
)

func TestLogger(t *testing.T) {
	var wg sync.WaitGroup
	file := "../test123.log"
	loggerInstance := &Logger{File: file}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, logger *Logger, i int) {
			defer wg.Done()
			logger.Log("some info " + strconv.Itoa(i))
		}(&wg, loggerInstance, i)
	}

	wg.Wait()
	t.Cleanup(func() {
		os.Remove(file)
	})
}
