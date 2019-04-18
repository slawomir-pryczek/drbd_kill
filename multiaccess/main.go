package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"multiaccess/fslock"
	"os"
	"runtime"
	"time"
)

var buff []byte = nil

func init() {
	buff = make([]byte, 1000000)
	for i := 0; i < len(buff); i++ {
		buff[i] = byte(rand.Int() % 255)
	}

	runtime.GOMAXPROCS(80)
}

func main() {
	finish := false

	// start 30 threads
	for i := 0; i < 30; i++ {
		go func() {
			for !finish {
				randOperation()
			}
		}()
	}

	// run for 5 seconds, then close the app
	start := time.Now().Unix()
	for time.Now().Unix()-start < 5 {
		time.Sleep(time.Second)
		fmt.Println(".")
	}
	finish = true
}

func randOperation() {

	op := rand.Int() % 10
	fname := fmt.Sprintf("out-%d.tmp", rand.Int()%30)

	_start := rand.Int() % 10000
	_end := _start + rand.Int()%800000 + 1
	_b := buff[_start:_end]

	if op == 0 {
		ioutil.WriteFile(fname, _b, 0777)
	}

	if op == 1 {
		f, err := os.Create(fname)
		defer f.Close()
		if err != nil {
			return
		}

		f.Write(_b)
		f.Seek(int64(len(_b)/2), os.SEEK_SET)
		f.Read(_b)

		f.Sync()
		f.Write(_b)
		f.Read(_b)

		f.Seek(int64(len(_b)/2), os.SEEK_SET)
		f.Stat()
		f.Truncate(int64(rand.Int() % 100))
		f.Write(_b)
		f.Seek(int64(len(_b)/2), os.SEEK_SET)
		return
	}

	// other operations
	mode := os.O_RDWR | os.O_CREATE
	if rand.Int()%2 == 0 {
		mode = mode | os.O_APPEND
	}

	f, err := os.OpenFile(fname, mode, 0755)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if op < 5 {
		fsl := fslock.New(fname)
		fsl.LockWithTimeout(time.Microsecond * 500)
		for i := 0; i < 10; i++ {
			f.Seek(int64(rand.Int()%10000), os.SEEK_SET)
			f.Read(_b)

		}
		fsl.Unlock()
		return
	}

	f.Close()
	fsl := fslock.New(fname)
	if rand.Int()%10 > 5 {
		fsl.Lock()
	} else {
		// CRUCIAL PART - THIS CODE WILL CAUSE KERNEL TO PANIC!
		// IT SEEMS THE ONLY ISSUE IS RELATED TO EXCLUSIVE FILE LOCK(s)
		// IT USUALLY TAKES 5-10s TO KILL DRBD/NFSv4
		// All other tasks beside this one are safe to run
		fsl.LockEx()
	}
	time.Sleep(time.Duration(rand.Int()%10) * time.Millisecond)
	fsl.Unlock()

	if op == 9 {
		if rand.Int()%10 > 5 {
			fname_new := fmt.Sprintf("out-%d.tmp", rand.Int()%30)
			os.Rename(fname, fname_new)
		}
	}

	if op == 8 {
		if rand.Int()%10 > 5 {
			os.Remove(fname)
		}

	}
}
