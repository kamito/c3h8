package propane

import (
	"os"
	"path/filepath"
)

func CurDir() string {
	current, _ := filepath.Abs(".")
	return current
}

func IsDirectory(name string) (isDir bool, err error) {
	fInfo, err := os.Stat(name) // FileInfo型が返る。
	if err != nil {
		return false, err // もしエラーならエラー情報を返す
	}
	// ディレクトリかどうかチェック
	return fInfo.IsDir(), nil
}
