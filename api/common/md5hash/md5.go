package md5hash

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func GetFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

func GetFolderHash(folderPath string) (string, error) {
	hash := md5.New()

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(hash, file); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}
