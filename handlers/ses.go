package handlers

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/spd-wilp/cloud_assignment/model"
)

type SESHandler struct {
	sess     *session.Session
	svc      *ses.SES
	sender   string
	receiver string
}

const charset = "UTF-8"

func InitSESHandler(region, sender, receiver string) *SESHandler {
	sess := session.Must(session.NewSession())
	svc := ses.New(sess, &aws.Config{
		Region: aws.String(region),
	})

	return &SESHandler{
		sess:     sess,
		svc:      svc,
		sender:   sender,
		receiver: receiver,
	}
}

func (handler SESHandler) SendEmail(ctx context.Context, metadata []model.ObjectMetadata) error {
	payload := &ses.SendEmailInput{
		Source: aws.String(handler.sender),
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(handler.receiver),
			},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String("test email"),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String("<p>this is a test email</p>"),
				},
			},
		},
	}

	if out, err := handler.svc.SendEmail(payload); err != nil {
		log.Printf("error while sending email, err=%v", err.Error())
		return err
	} else {
		log.Printf("successfully sent email")
		log.Printf("out: %+v", out)
		return nil
	}
}
