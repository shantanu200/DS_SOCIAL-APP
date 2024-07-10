package s3

import (
	"fmt"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AWSRegion  = "ap-south-1"
	BucketName = "skyboxstudiobucketv2"
)

// Function which will upload file on s3 and give all links
func FileUploader(files []*multipart.FileHeader) ([]string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("ACCESSKEY_ID"),
			os.Getenv("ACCESSKEY_SECRET"),
			"",
		),
	})

	svc := s3.New(sess)

	if err != nil {
		return nil, err
	}
	var urls []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var uploadErr error
	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			fileKey := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)

			srcFile, err := file.Open()

			if err != nil {
				fmt.Println("Error on upload!!", err)
			}
			defer srcFile.Close()

			_, err = svc.PutObject(&s3.PutObjectInput{
				Bucket:             aws.String(BucketName),
				Key:                aws.String(fileKey),
				Body:               srcFile,
				ContentDisposition: aws.String("inline"),
			})

			if err != nil {
				uploadErr = err
				fmt.Println("Error on upload!!", err)
			}

			url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", BucketName, AWSRegion, fileKey)

			fmt.Println(url)

			mu.Lock()

			urls = append(urls, url)

			mu.Unlock()
		}(file)
	}

	wg.Wait()

	if uploadErr != nil {
		return nil, uploadErr
	}

	return urls, nil

}
