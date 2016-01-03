package main
import (
	"log"
	"io/ioutil"
	"bufio"
	"os"
)


func FilesInDirection(dir string) ([]string, error) {
	fileList, err := ioutil.ReadDir(dir) //要读取的目录地址DIR，得到列表
	if err != nil {
		log.Printf("read dir error")
		return nil, err
	}

	retFileList := []string{}
	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		retFileList = append(retFileList, file.Name())
	}

	return retFileList, nil
}

func readFile(fileName string) ([]byte, error) {
	srcDat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error:%v", err)
		return nil, err
	}

	return srcDat, nil
}


func writeFile(fileName, srcFileStr string) error {
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(srcFileStr)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	w.Flush()
	return nil
}
