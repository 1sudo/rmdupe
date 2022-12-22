package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	checkPath := true

	if len(os.Args) < 4 {
		fmt.Println("Not enough args! Usage:")
		fmt.Printf("arg1: path to remove dupes from\narg2: path to compare with\narg3: number of threads\narg4: directory path must match (true/false)\n\n")
		fmt.Println("Example: \nrmdupe folder1/data folder2/data 10 true")
		return
	}

	if len(os.Args) > 5 {
		if os.Args[4] == "false" {
			checkPath = false
		} else if os.Args[4] != "true" {
			fmt.Println("Invalid boolean input, checkPath must be true or false, or not supplied (defaults to true)")
			return
		}
	}

	threads, err := strconv.Atoi(os.Args[3])

	if err != nil {
		fmt.Println("Invalid integer input, threads must be an integer value")
		return
	}

	RemoveDuplicates(os.Args[1], os.Args[2], threads, checkPath)
}

func GetFiles(path string) []string {
	var s []string

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !f.IsDir() {
			s = append(s, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return s
}

func RemoveDuplicates(sourcePath string, comparisonPath string, threads int, checkPath bool) {
	var files []string = GetFiles(sourcePath)
	var cFiles []string = GetFiles(comparisonPath)

	ch := make(chan int, threads)

	if files != nil && cFiles != nil {
		for _, file := range files {

			ch <- 1

			go func(sFile string) {
				res1 := strings.Split(sFile, sourcePath)[1]

				for _, cFile := range cFiles {
					res2 := strings.Split(cFile, comparisonPath)[1]

					if checkPath {
						if res1 == res2 {
							CompareSums(sFile, cFile)
						}
					} else {
						idx := strings.LastIndex(sFile, "/")
						idx2 := strings.LastIndex(cFile, "/")

						fileName := sFile[idx+1:]
						cFileName := cFile[idx2+1:]

						if fileName == cFileName {
							CompareSums(sFile, cFile)
						}
					}
				}
				<-ch
			}(file)
		}
	}
}

func CompareSums(file string, cFile string) {
	var bytes, err = ioutil.ReadFile(file)

	if err != nil {
		return
	}

	var bytes2, err2 = ioutil.ReadFile(cFile)

	if err2 != nil {
		return
	}

	sum := sha256.Sum256(bytes)
	sum2 := sha256.Sum256(bytes2)

	if sum == sum2 {
		os.Remove(file)
		fmt.Printf("Removing match: %s, %x, %x\n", file, sum, sum2)
	}
}
