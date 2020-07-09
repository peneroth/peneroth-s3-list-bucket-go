// snippet-comment:[These are tags for the AWS doc team's sample catalog. Do not remove.]
// snippet-sourceauthor:[Doug-AWS]
// snippet-sourcedescription:[Lists your S3 buckets.]
// snippet-keyword:[Amazon Simple Storage Service]
// snippet-keyword:[Amazon S3]
// snippet-keyword:[ListBuckets function]
// snippet-keyword:[Go]
// snippet-sourcesyntax:[go]
// snippet-service:[s3]
// snippet-keyword:[Code Sample]
// snippet-sourcetype:[full-example]
// snippet-sourcedate:[2018-03-16]
/*
   Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at
    http://aws.amazon.com/apache2.0/
   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
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
	// Initialize a session in us-west-2 that the SDK will use to load
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
