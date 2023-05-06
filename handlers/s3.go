package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spd-wilp/cloud_assignment/model"
	"github.com/spd-wilp/cloud_assignment/utils/image"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Handler struct {
	sess              *session.Session
	svc               *s3.S3
	region            string
	srcBucket         string
	destBucket        string
	metadataBucket    string
	metadataObjectKey string
}

func InitS3Handler(region, srcBucket, destBucket, metadataBucket, metadataObjectKey string) *S3Handler {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region: aws.String(region),
	})

	return &S3Handler{
		sess:              sess,
		svc:               svc,
		region:            region,
		srcBucket:         srcBucket,
		destBucket:        destBucket,
		metadataBucket:    metadataBucket,
		metadataObjectKey: metadataObjectKey,
	}
}

func (handler S3Handler) GenerateAndStoreThumbnailForObjects(objects []events.S3Object) ([]model.ObjectMetadata, error) {
	var objectType, thumbnailURI string
	objectDetails := make([]model.ObjectMetadata, 0)

	for _, object := range objects {
		objectType = "normal"
		thumbnailURI = "n/a"
		if image.IsImage(object.Key) {
			objectType = "image"
			thumbnailData, err := handler.createThumbnail(object.Key)
			if err != nil {
				log.Printf("error while generating thumbnail, err=%v", err.Error())
				continue
			}
			err = handler.writeThumbnail(thumbnailData, object.Key)
			if err != nil {
				log.Printf("error while creating and storing thumbnail, err=%v", err.Error())
				continue
			}
			thumbnailURI = handler.generateResourceURI(handler.destBucket, object.Key)
			log.Printf("successfully created and stored thumbnail, src_object_key=%s thumbnail_uri=%s", object.Key, thumbnailURI)
		} else {
			log.Printf("skipping thumbnail generation as object is not an image, src_object_key=%s", object.Key)
		}
		t := time.Now()
		objectDetails = append(objectDetails, model.ObjectMetadata{
			Name:            object.Key,
			SourceURI:       handler.generateResourceURI(handler.srcBucket, object.Key),
			LastModified:    t.Unix(),
			LastModifiedStr: t.Format("3:04:05 PM"),
			Size:            object.Size,
			Type:            objectType,
			ThumbnailURI:    thumbnailURI,
		})
	}

	return objectDetails, nil
}

func (handler S3Handler) createThumbnail(objectKey string) ([]byte, error) {
	out, err := handler.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(handler.srcBucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("error while reading object, bucket=%s, key=%s, err=%v", handler.srcBucket, objectKey, err.Error())
		return nil, err
	}
	defer out.Body.Close()

	objectContent, err := io.ReadAll(out.Body)
	if err != nil {
		log.Printf("error while reading object content, bucket=%s, key=%s, err=%v", handler.srcBucket, objectKey, err.Error())
		return nil, err
	}
	return objectContent, nil
}

func (handler S3Handler) writeThumbnail(data []byte, objectKey string) error {
	thumbnail, err := image.CreateThumbnail(objectKey, data)
	if err != nil {
		log.Printf("error while creating thumbnail, bucket=%s, key=%s, err=%v", handler.srcBucket, objectKey, err.Error())
		return err
	}

	_, err = handler.svc.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(thumbnail),
		Bucket: aws.String(handler.destBucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("error while writing thumbnail, src_bucket=%s, key=%s, dest_bucket=%s, err=%v", handler.srcBucket, objectKey, handler.destBucket, err.Error())
		return err
	}

	return nil
}

func (handler S3Handler) generateResourceURI(bucket string, key string) string {
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucket, handler.region, key)
}

func (handler S3Handler) GetMetadata(st, et int64) ([]model.ObjectMetadata, error) {
	out, err := handler.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(handler.metadataBucket),
		Key:    aws.String(handler.metadataObjectKey),
	})
	if err != nil {
		log.Printf("error while reading object, bucket=%s, key=%s, err=%v", handler.metadataBucket, handler.metadataObjectKey, err.Error())
		return nil, err
	}
	defer out.Body.Close()

	var curMetadata []model.ObjectMetadata
	metaKeyMap := make(map[string]model.ObjectMetadata)
	var uniqMetadata []model.ObjectMetadata

	err = json.NewDecoder(out.Body).Decode(&curMetadata)
	if err != nil {
		log.Printf("error while decoding current metadata, err=%s", err.Error())
		return nil, err
	}

	// fing uniq metadata entries, there might be duplicate ones due to multiple lamdba runs / if same file is deleted and re-uploaded
	for _, metadata := range curMetadata {
		data, exists := metaKeyMap[metadata.Name]
		if !exists {
			metaKeyMap[metadata.Name] = metadata
		} else {
			// keep the latest modification time
			if data.LastModified > metadata.LastModified {
				metaKeyMap[metadata.Name] = data
			}
		}
	}

	for _, data := range metaKeyMap {
		uniqMetadata = append(uniqMetadata, data)
	}

	// return all
	if st < 0 || et < 0 {
		return uniqMetadata, nil
	}

	// filter based on time
	filteredMetadata := make([]model.ObjectMetadata, 0)
	for _, metadata := range uniqMetadata {
		if metadata.LastModified >= st && metadata.LastModified <= et {
			filteredMetadata = append(filteredMetadata, metadata)
		}
	}

	return filteredMetadata, nil
}

func (handler S3Handler) UpdateMetadata(metadata []model.ObjectMetadata) error {
	var curMetadata []model.ObjectMetadata
	var updatedMetadata []model.ObjectMetadata
	var err error

	curMetadata, err = handler.GetMetadata(-1, -1)
	if err != nil {
		log.Printf("error while reading object, bucket=%s, key=%s, err=%v", handler.metadataBucket, handler.metadataObjectKey, err.Error())
		return err
	}

	updatedMetadata = append(curMetadata, metadata...)
	update, err := json.Marshal(updatedMetadata)
	if err != nil {
		log.Printf("error while marshaling updated metadata, err=%s", err.Error())
		return err
	}

	_, err = handler.svc.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(update),
		Bucket: aws.String(handler.metadataBucket),
		Key:    aws.String(handler.metadataObjectKey),
	})
	if err != nil {
		log.Printf("error while writing metadata, bucket=%s, key=%s, err=%v", handler.metadataBucket, handler.metadataObjectKey, err.Error())
		return err
	}

	return nil
}
