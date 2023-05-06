package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"html/template"

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

	htmlBodyTemplate *template.Template
}

const charset = "UTF-8"

func InitSESHandler(region, sender, receiver string) *SESHandler {
	sess := session.Must(session.NewSession())
	svc := ses.New(sess, &aws.Config{
		Region: aws.String(region),
	})

	htmlBodyTemplate, _ := template.New("html_email_body").Parse(`
	<html>
		<body>
			<p>This email captures all the objects modified for the given date.</p>
			<br />
			<table style="border-collapse: collapse;font-family: Tahoma, Geneva, sans-serif;">
				<tr>
					<th style="background-color: #54585d;color: #ffffff;font-weight: bold;font-size: 13px;border: 1px solid #54585d;">Name</th>
					<th style="background-color: #54585d;color: #ffffff;font-weight: bold;font-size: 13px;border: 1px solid #54585d;">URI</th>
					<th style="background-color: #54585d;color: #ffffff;font-weight: bold;font-size: 13px;border: 1px solid #54585d;">Type</th>
					<th style="background-color: #54585d;color: #ffffff;font-weight: bold;font-size: 13px;border: 1px solid #54585d;">Size (Byte)</th>
					<th style="background-color: #54585d;color: #ffffff;font-weight: bold;font-size: 13px;border: 1px solid #54585d;">Last Modified Time (GMT)</th>
					<th style="background-color: #54585d;color: #ffffff;font-weight: bold;font-size: 13px;border: 1px solid #54585d;">Thumbnail</th>
				</tr>
				{{ range .}}
					<tr style="background-color: #f9fafb;">
						<td style="color: #636363;border: 1px solid #dddfe1; padding: 15px;">{{ .Name }}</td>
						<td style="color: #636363;border: 1px solid #dddfe1; padding: 15px;">{{ .SourceURI }}</td>
						<td style="color: #636363;border: 1px solid #dddfe1; padding: 15px;">{{ .Type }}</td>
						<td style="color: #636363;border: 1px solid #dddfe1; padding: 15px;">{{ .Size }}</td>
						<td style="color: #636363;border: 1px solid #dddfe1; padding: 15px;">{{ .LastModifiedStr }}</td>
						<td style="color: #636363;border: 1px solid #dddfe1; padding: 15px;">{{ .ThumbnailURI }}</td>
					</tr>
				{{ end}}
		  </table> 
		  <br />
		  <br />
		  <p>Sent from Report Automation Job, please reach out to 2022mt93539@wilp.bits-pilani.ac.in in case of queries</p>
		  <p></p>
		</body>
	</html>
	`)

	return &SESHandler{
		sess:             sess,
		svc:              svc,
		sender:           sender,
		receiver:         receiver,
		htmlBodyTemplate: htmlBodyTemplate,
	}
}

func (handler SESHandler) SendEmail(ctx context.Context, metadata []model.ObjectMetadata) error {
	log.Printf("metadata: %+v", metadata)

	var emailBodyBuffer []byte
	emailBody := bytes.NewBuffer(emailBodyBuffer)
	err := handler.htmlBodyTemplate.Execute(emailBody, metadata)
	if err != nil {
		log.Printf("error while compiling email body template, err=%v", err.Error())
		return err
	}

	prevDayStr := time.Now().Add(-24 * time.Hour).Format("January 2, 2006")
	emailSubject := fmt.Sprintf("Summary Email for %s", prevDayStr)

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
				Data:    aws.String(emailSubject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(emailBody.String()),
				},
			},
		},
	}

	if _, err := handler.svc.SendEmail(payload); err != nil {
		log.Printf("error while sending email, err=%v", err.Error())
		return err
	}
	return nil
}
