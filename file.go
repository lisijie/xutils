package xutils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	PrivateFileMode = 0644
	PrivateDirMode  = 0755
)

var (
	DefaultTempDir = os.TempDir()
)

// TempDir 创建临时目录
// 返回临时目录地址和清理函数
func TempDir(prefix string, dir ...string) (string, func()) {
	tmpDir := DefaultTempDir
	if len(dir) > 0 {
		tmpDir = dir[0]
	}
	p, err := os.MkdirTemp(tmpDir, prefix)
	if err != nil {
		panic("TempDir: " + err.Error())
	}
	p, err = filepath.Abs(p) // 使用绝对路径
	if err != nil {
		panic("TempDir: " + err.Error())
	}
	return p, func() {
		if err := os.RemoveAll(p); err != nil {
			log.Printf("os.RemoveAll %s: %s\n", p, err)
		}
	}
}

// TempFile 创建临时文件，调用方需要负责删除
func TempFile(pattern string, dir ...string) (*os.File, error) {
	tmpDir := DefaultTempDir
	if len(dir) > 0 {
		tmpDir = dir[0]
	}
	return os.CreateTemp(tmpDir, pattern)
}

// FileExt 返回文件小写扩展名，如 foo.PNG 返回 .png
func FileExt(path string) string {
	return strings.ToLower(filepath.Ext(path))
}

// FileSize 返回文件大小
func FileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fi.Size()
}

// FileMD5 计算文件的MD5
func FileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))
	return hashStr, nil
}

// FormatSize 格式化文件大小
func FormatSize(size int64) string {
	units := []string{"Byte", "KB", "MB", "GB", "TB", "PB"}
	n := 0
	fs := float64(size)
	for fs >= 1024 {
		fs /= 1024
		n += 1
	}
	return fmt.Sprintf("%.2f %s", fs, units[n])
}

// IsDirWritable 检查目录是否可写
func IsDirWritable(dir string) error {
	tmpFile := filepath.Join(dir, ".tmp")
	if err := os.WriteFile(tmpFile, []byte(""), PrivateFileMode); err != nil {
		return err
	}
	return os.Remove(tmpFile)
}

// IsFile 检查是否文件
func IsFile(filename string) bool {
	info, err := os.Stat(filename)
	if err == nil && !info.IsDir() {
		return true
	}
	return false
}

// IsDir 检查是否目录
func IsDir(dir string) bool {
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		return true
	}
	return false
}

// ReadDir 读取目录下的所有文件名
func ReadDir(dir string) ([]string, error) {
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

// ReadDirAll 返回 dir 目录及其子目录下的所文件列表
func ReadDirAll(dir string) ([]string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path[len(dir)+1:])
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// MakeDirAll 创建目录
func MakeDirAll(path string) error {
	return os.MkdirAll(path, PrivateDirMode)
}

// CopyDir 拷贝指定目录下的所有文件到另一个目录
// 如果目标文件存在，则会覆盖
func CopyDir(src string, dst string) error {
	src, _ = filepath.Abs(src)
	dst, _ = filepath.Abs(dst)
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		dstPath := filepath.Join(dst, path[len(src):])
		if info.IsDir() {
			return MakeDirAll(dstPath)
		} else {
			fw, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, PrivateFileMode)
			if err != nil {
				return err
			}
			defer fw.Close()
			fr, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fr.Close()
			_, err = io.Copy(fw, fr)
			return err
		}
	})
}

// CopyFile 拷贝文件到指定目录，如果目标文件已存在将进行覆盖
func CopyFile(src string, dst string) error {
	if IsDir(src) {
		return errors.New(src + " is a directory")
	}
	if IsDir(dst) {
		dst = filepath.Join(dst, filepath.Base(src))
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, PrivateFileMode)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, info.Mode())
}

// MoveFile 移动文件到指定目录
func MoveFile(src string, dst string) error {
	if IsDir(src) {
		return errors.New(src + " is a directory")
	}
	if IsDir(dst) {
		dst = filepath.Join(dst, filepath.Base(src))
	}
	return os.Rename(src, dst)
}

// AppendFile 追加写入文件
func AppendFile(filename string, data []byte) (err error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, PrivateFileMode)
	if err != nil {
		return
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		return io.ErrShortWrite
	}
	return f.Close()
}

// WriteFile 写文件
func WriteFile(fn string, data []byte) error {
	dir := filepath.Dir(fn)
	if !IsDir(dir) {
		if err := MakeDirAll(dir); err != nil {
			return err
		}
	}
	return os.WriteFile(fn, data, PrivateFileMode)
}

// CleanDir 清理目录，只保留指定文件
func CleanDir(dir string, excludeFiles []string) error {
	if !IsDir(dir) {
		return fmt.Errorf("%s is not a directory", dir)
	}
	if len(excludeFiles) == 0 {
		return os.RemoveAll(dir)
	}
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	fis, err := d.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		fn := filepath.Join(dir, fi.Name())
		if fi.IsDir() {
			_ = CleanDir(fn, excludeFiles)
			fis2, _ := os.ReadDir(fn)
			if len(fis2) == 0 {
				os.Remove(fn)
			}
		} else {
			if !InArray(fn, excludeFiles) {
				os.Remove(fn)
			}
		}
	}
	return nil
}

// MoveUploadFile 将HTTP上传的文件保存到指定目录
func MoveUploadFile(upfile *multipart.FileHeader, toPath string) (string, error) {
	if IsDir(toPath) {
		toPath = filepath.Join(toPath, filepath.Base(upfile.Filename))
	}
	if IsFile(toPath) {
		return "", errors.New("目标文件已存在")
	}
	f1, err := upfile.Open()
	if err != nil {
		return "", err
	}
	f2, err := os.OpenFile(toPath, os.O_CREATE|os.O_RDWR, PrivateFileMode)
	if err != nil {
		f1.Close()
		return "", err
	}
	_, err = io.Copy(f2, f1)
	f1.Close()
	f2.Close()
	if err != nil {
		return "", err
	}
	return toPath, nil
}
