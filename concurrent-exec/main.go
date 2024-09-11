package main

import (
	"fmt"
	"os"
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
