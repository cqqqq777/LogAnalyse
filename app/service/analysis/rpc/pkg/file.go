package pkg

import (
	"LogAnalyse/app/service/analysis/internal"
	"fmt"
	"io"
	"net/http"
	"os"
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

	path := GetFilePath(userId) + "task"
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil || info.Size() <= internal.FileSize {
		return []string{path}, nil
	}

	var Cap int64
	if info.Size()%internal.FileSize != 0 {
		Cap = info.Size()/internal.FileSize + 1
	} else {
		Cap = info.Size() / internal.FileSize
	}

	paths := make([]string, 0, Cap)

	b := make([]byte, internal.FileSize)
	var i int64 = 1
	for ; i <= Cap; i++ {
		out.Seek((i-1)*internal.FileSize, 0)
		if len(b) > int(info.Size()-(i-1)*internal.FileSize) {
			b = make([]byte, info.Size()-(i-1)*internal.FileSize)
		}

		out.Read(b)

		f, err := os.Create(path + fmt.Sprintf("%d", i))
		if err != nil {
			return []string{path}, nil
		}
		f.Write(b)
		f.Close()
		paths = append(paths, path+fmt.Sprintf("%d", i))
	}

	os.Remove(path)

	return paths, nil
}

func DeleteFile(paths []string) {

}

func CreateFile(path string, data []byte) error {
	return nil
}
