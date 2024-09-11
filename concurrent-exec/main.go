package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func CreateDataDir() {
	_, err := os.Stat("data")
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		err := os.Mkdir("data", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func TouchFile(count int) {
	fmt.Printf("TouchFile %d starting\n", count)
	file, err := os.Create(fmt.Sprintf("data/%d.txt", count))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileContents := []byte(fmt.Sprintf("This is file #%d", count))
	file.Write(fileContents)
	time.Sleep(time.Second)
	fmt.Printf("TouchFile %d done\n", count)

	stdout, err := os.Create(fmt.Sprintf("data/stdout_%d.log", count))
	if err != nil {
		panic(err)
	}
	defer stdout.Close()

	stderr, err := os.Create(fmt.Sprintf("data/stderr_%d.log", count))
	if err != nil {
		panic(err)
	}
	defer stderr.Close()

	cmd := exec.Command("cat", fmt.Sprintf("data/%d.txt", count*2))
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("error catting file", count*2, err.Error())
	}
}

func main() {
	CreateDataDir()
	fmt.Println("test")

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			TouchFile(i)
		}()
	}

	wg.Wait()
}
