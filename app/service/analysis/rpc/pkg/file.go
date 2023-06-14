package pkg

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"LogAnalyse/app/service/analysis/internal"
	"LogAnalyse/app/shared/log"
)

func GetFilePath(userId int64) string {
	return fmt.Sprintf("%s%d/", internal.FilePrefix, userId)
}

func DownloadFile(url string, userId int64) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	path := GetFilePath(userId) + "task.log"
	tmp, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer os.Remove(path)
	defer tmp.Close()

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(tmp)
	linesCh := make(chan string)
	fileCount := 1
	paths := make([]string, 0)

	path = fmt.Sprintf("%stask%d", GetFilePath(userId), fileCount)
	out, err := os.Create(path)
	if err != nil {
		return []string{path}, nil
	}
	paths = append(paths, path)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range linesCh {
			go possessFile(line, out, &wg)
		}
	}()

	for i := 1; scanner.Scan(); {
		linesCh <- scanner.Text()
		i++
		if i >= internal.ChunkLine {
			out.Close()
			fileCount++
			path = fmt.Sprintf("./task%d.log", fileCount)
			paths = append(paths, path)
			out, err = os.Create(path)
			if err != nil {
				log.Zlogger.Warn("create task file failed err:" + err.Error())
				return paths, nil
			}
			i = 1
		}
	}

	close(linesCh)
	out.Close()
	wg.Wait()
	return paths, nil
}

func possessFile(line string, file *os.File, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Fprintln(file, line)
}

func DeleteFile(paths []string) {
	for _, path := range paths {
		os.Remove(path)
	}
}

func CreateFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}
