package flow

import (
	"online_fashion_shop/initializers/storage"
	"os"
	"testing"
)

func TestAzureStoreBlob(t *testing.T) {
	asb := storage.AzureStorageBlob{
		AccountName: "huuvuongassert",
		AccountKey:  "zxM15v2/reZN/38o4zj8xa5Hhbk4KZHBON+G/n9cWKA/+jeGNU18Sf7LXKthJiE6EJ2uJU710TIU+AStdi3JHg==",
	}

	files, err := os.ReadDir("./assert")
	if err != nil {
		panic(err)
	}

	var filePointers []*os.File

	for _, file := range files {
		if !file.IsDir() {
			filePointer, err := os.Open("./assert/" + file.Name())
			if err != nil {
				panic(err)
			}
			filePointers = append(filePointers, filePointer)
		}
	}

	t.Run("delete", func(t *testing.T) {
		photoUrl := "https://huuvuongassert.blob.core.windows.net/photo/0a3fe296ecbfb0226c5602ed44faca96.jpg"
		err := asb.Delete(photoUrl)
		if err != nil {
			panic(err)
		}
	})
}
