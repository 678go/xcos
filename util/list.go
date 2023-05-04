package util

import (
	"fmt"
	"golang.org/x/exp/slog"
	"io/fs"
	"os"
	"path/filepath"
)

func GetLocalFilesListRecursive(localPath string) (files []string) {
	wd, _ := os.Getwd()
	localPath = filepath.Join(wd, localPath)

	// bfs遍历文件夹
	var dirs []string
	dirs = append(dirs, localPath)
	for len(dirs) > 0 {
		dirName := dirs[0]
		dirs = dirs[1:]

		fileInfos, err := os.ReadDir(dirName)
		if err != nil {
			slog.Error("上传文件夹失败,", "err", err)
			os.Exit(1)
		}

		for _, f := range fileInfos {
			fileName := dirName + "/" + f.Name()
			if f.Type().IsRegular() { // 普通文件，直接添加
				fileName = fileName[len(localPath)+1:]
				files = append(files, fileName)
			} else if f.IsDir() { // 普通目录，添加到继续迭代
				dirs = append(dirs, fileName)
			} else if f.Type()&os.ModeSymlink == fs.ModeSymlink { // 软链接
				slog.Info(fmt.Sprintf("List %s file is Symlink, will be excluded, please list or upload it from realpath", fileName))
				continue
			} else {
				slog.Info(fmt.Sprintf("List %s file is not regular file, will be excluded", fileName))
				continue
			}
		}
	}
	fmt.Println(files)
	return files
}
