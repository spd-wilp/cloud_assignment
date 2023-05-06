package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/spd-wilp/cloud_assignment/handlers"
	"github.com/spd-wilp/cloud_assignment/model"
	"github.com/spd-wilp/cloud_assignment/utils/timeutil"
)

var sesHandler *handlers.SESHandler
var s3Handler *handlers.S3Handler

func main() {
	sesHandler = handlers.InitSESHandler(model.Region, model.EmailSender, model.EmailReceiver)
	s3Handler = handlers.InitS3Handler(model.Region, model.S3SrcBucket, model.S3DestBucket, model.S3MetadataBucket, model.S3MetadataObjectKey)

	lambda.Start(handler)
}

/*
	Send emails with details displaying following information about objects uploaded the previous day
		- s3 uri
		- object name
		- object size
		- object type
		- thumbnail uri
*/
// todo: capture concrete type for event
func handler(ctx context.Context, event map[string]interface{}) error {
	st, et, _ := timeutil.FindTimeBoundOfPreviousDay(time.Now())

	metadata, err := s3Handler.GetMetadata(st, et)
	if err != nil {
		log.Printf("error while fetching metadata, err=%v", err.Error())
		return err
	}

	err = sesHandler.SendEmail(ctx, metadata)
	if err != nil {
		log.Printf("error while sending email, err=%v", err.Error())
	} else {
		log.Printf("successfully sent email")
	}

	return err
}
