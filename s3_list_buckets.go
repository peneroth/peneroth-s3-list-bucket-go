/*
   This is a test program based on an example file from Amazon
   The program lists files in a S3 bucket
*/

package main

import (
	"fmt"
	"os"
	"time"

	// Peter: Documentation for AWX Go API: https://docs.aws.amazon.com/sdk-for-go/api/
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// Initialize a session in eu-north-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	//
	// Peter: Info about the credentials file at:
	// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/create-shared-credentials-file.html
	//
	// Peter: Region names at https://docs.aws.amazon.com/general/latest/gr/rande.html
	//
	start := time.Now()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1")},
	)
	elapsed := time.Since(start)
	fmt.Println("Time to create session: ", elapsed)

	// Create S3 service client
	start = time.Now()
	svc := s3.New(sess)
	elapsed = time.Since(start)
	fmt.Println("Time to create S3 service client: ", elapsed)

	start = time.Now()
	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}
	elapsed = time.Since(start)
	fmt.Println("Time to list buckets: ", elapsed)

	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	bucket := aws.StringValue(result.Buckets[0].Name)
	fmt.Println(bucket)

	start = time.Now()
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}
	elapsed = time.Since(start)
	fmt.Println("Time to list objects: ", elapsed)

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
