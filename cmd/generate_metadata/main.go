package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/spd-wilp/cloud_assignment/handlers"
	"github.com/spd-wilp/cloud_assignment/model"

	"github.com/aws/aws-lambda-go/events"
)

var s3Handler *handlers.S3Handler

func main() {
	s3Handler = handlers.InitS3Handler(model.Region, model.S3SrcBucket, model.S3DestBucket, model.S3MetadataBucket, model.S3MetadataObjectKey)

	lambda.Start(handler)
}

/*
Whenever a new object is uploaded in source s3 bucket, this will check
  - if uploaded object is an image (.jpg, .jpeg, .png)
  - if object is an image, then this will create a thumbnail and store that into the destination bucket
  - adds metadata about all objects in metadata.json in metadata bucket
*/
func handler(ctx context.Context, event events.S3Event) error {
	objects := make([]events.S3Object, 0)
	for _, record := range event.Records {
		objects = append(objects, record.S3.Object)
	}

	if len(objects) == 0 {
		log.Printf("no object information received in incoming event")
		return fmt.Errorf("no object detail received")
	}

	metadata, err := s3Handler.GenerateAndStoreThumbnailForObjects(objects)
	if err != nil {
		log.Printf("error while performing the thumbnail backup process, err=%v", err.Error())
	} else {
		log.Printf("successfully performed thumbnail creation and backup")
	}

	err = s3Handler.UpdateMetadata(metadata)
	if err != nil {
		log.Printf("error while adding metadata in db, err=%v", err.Error())
	} else {
		log.Printf("successfully stored metadata")
	}

	return err
}
