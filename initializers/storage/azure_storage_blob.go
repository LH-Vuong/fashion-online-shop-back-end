package storage

import (
	"fmt"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	uuid2 "github.com/google/uuid"
	"io"
	"net/url"
	"online_fashion_shop/initializers"
	"time"
)

type PhotoStorage interface {
	Upload(file io.Reader) (string, error)
	MultiUpload(files []io.Reader) ([]string, error)
	Delete(photoUrl string) error
	DeleteMany(photoUrls []string) error
}

func NewAzureStorageBlob(accountName, key2 string) PhotoStorage {
	return &AzureStorageBlob{
		AccountName: accountName,
		AccountKey:  key2,
	}
}

type AzureStorageBlob struct {
	AccountName string
	AccountKey  string
}

func (asb *AzureStorageBlob) DeleteMany(photoUrls []string) error {
	for _, photoUrl := range photoUrls {
		err := asb.Delete(photoUrl)
		if err != nil {
			return err
		}

	}
	return nil
}

var uploadOption = azblob.UploadStreamToBlockBlobOptions{
	BufferSize:      8 * 1024 * 1024,
	MaxBuffers:      16,
	BlobHTTPHeaders: azblob.BlobHTTPHeaders{ContentType: "image"},
}

//func (asb *AzureStorageBlob) containerUrl() (*azblob.BlockBlobURL, error) {
//	// Create a connection string to your Azure Storage account
//	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net", asb.AccountName, asb.AccountKey)
//	// Create a BlobServiceClient object
//	credential, err := azblob.NewSharedKeyCredential(connectionString, asb.AccountKey)
//	if err != nil {
//		return nil, err
//	}
//
//	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
//	cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/photo", asb.AccountName))
//	containerURL := azblob.NewContainerURL(*cURL, p)
//	uuid, err := uuid2.NewUUID()
//	blobURL := containerURL.NewBlockBlobURL(uuid.String())
//	return &blobURL, err
//}

func (asb *AzureStorageBlob) createPipeline() (*pipeline.Pipeline, error) {
	credential, err := azblob.NewSharedKeyCredential(asb.AccountName, asb.AccountKey)
	if err != nil {
		return nil, err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	return &p, err
}

func (asb *AzureStorageBlob) Upload(file io.Reader) (string, error) {

	// From the Azure portal, get your Storage account blob service URL endpoint.
	cURL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/photo", asb.AccountName))
	if err != nil {
		return "", err
	}
	pipeline, err := asb.createPipeline()
	if err != nil {
		return "", err
	}
	// Create an ServiceURL object that wraps the service URL and a request pipeline to making requests.
	containerURL := azblob.NewContainerURL(*cURL, *pipeline)
	// Here's how to create a blob with HTTP headers and metadata (I'm using the same metadata that was put on the container):
	uuid, _ := uuid2.NewUUID()
	blobURL := containerURL.NewBlockBlobURL(uuid.String() + "_" + time.Now().Format("02_01_2006_15h_04m_05s"))

	ctx, cancel := initializers.InitContext()
	defer cancel()
	_, err = azblob.UploadStreamToBlockBlob(ctx, file, blobURL, uploadOption)
	if err != nil {
		return "", err
	}

	blobURLWithSAS := azblob.NewBlobURL(blobURL.URL(), azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})).URL()
	return blobURLWithSAS.String(), nil
}

func (asb *AzureStorageBlob) MultiUpload(files []io.Reader) ([]string, error) {
	photos := make([]string, len(files))
	for index, file := range files {
		if photo, err := asb.Upload(file); err == nil {
			photos[index] = photo
		} else {
			return nil, err
		}

	}
	return photos, nil
}

func (asb *AzureStorageBlob) Delete(path string) error {

	pipeline, err := asb.createPipeline()

	if err != nil {
		return err
	}
	cURL, err := url.Parse(path)
	if err != nil {
		return err
	}
	blodURL := azblob.NewBlobURL(*cURL, *pipeline)
	ctx, cancel := initializers.InitContext()
	defer cancel()
	_, err = blodURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})

	return err

}
