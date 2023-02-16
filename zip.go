package xutils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	SkipDir = fs.SkipDir
)

// ZipDir 将dir整个目录打包到为zip文件
func ZipDir(dir string, filename string, noWrap ...bool) error {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	zipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	var baseDir string
	noWarpDir := len(noWrap) > 0 && noWrap[0]
	if !noWarpDir {
		baseDir = filepath.Base(dir)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if noWarpDir && path == dir {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimLeft(strings.TrimPrefix(path, dir), "/")
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, header.Name)
		}
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		return err
	})

	return err
}

// Unzip 解压 zip 文件到指定目录
func Unzip(filename string, dir string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	return UnzipReader(f, fi.Size(), dir)
}

// UnzipReader 解压 zip 文件到指定目录
func UnzipReader(rd io.ReaderAt, size int64, dir string) error {
	r, err := zip.NewReader(rd, size)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(dir, PrivateDirMode); err != nil {
		return err
	}
	for _, file := range r.File {
		path := filepath.Join(dir, file.Name)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(path, PrivateDirMode)
			continue
		}
		if filepath.Base(file.Name) == ".DS_Store" {
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		if d := filepath.Dir(path); !IsDir(d) {
			if err := os.MkdirAll(d, PrivateDirMode); err != nil {
				fileReader.Close()
				return err
			}
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			return err
		}
		if _, err = io.Copy(targetFile, fileReader); err != nil {
			fileReader.Close()
			return err
		}
		fileReader.Close()
	}
	return nil
}

// ZipWalk 遍历zip中的文件，执行回调函数，如果回调函数返回error，则终止迭代
func ZipWalk(filename string, fun func(zf *zip.File) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	return ZipWalkReader(f, fi.Size(), fun)
}

// ZipWalkReader 遍历zip reader中的文件，执行回调函数，如果回调函数返回error，则终止迭代
func ZipWalkReader(rd io.ReaderAt, fsize int64, fun func(zf *zip.File) error) (err error) {
	zr, err := zip.NewReader(rd, fsize)
	if err != nil {
		return
	}
	for _, v := range zr.File {
		if err = fun(v); err != nil {
			if err == SkipDir {
				return nil
			}
			return
		}
	}
	return
}
