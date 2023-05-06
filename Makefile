build-generate-metadata-zip:
	rm -rf /tmp/aws_deploy
	mkdir /tmp/aws_deploy
	GOOS=linux GOARCH=amd64 go build -o bin/generate_metadata_linux_amd64 cmd/generate_metadata/main.go
	chmod +x bin/generate_metadata_linux_amd64
	cp bin/generate_metadata_linux_amd64 /tmp/aws_deploy/main
	cd /tmp/aws_deploy/ && zip main.zip main
	cp /tmp/aws_deploy/main.zip deploy/
	rm -rf /tmp/aws_deploy

build-send-email-zip:
	rm -rf /tmp/aws_deploy
	mkdir /tmp/aws_deploy
	GOOS=linux GOARCH=amd64 go build -o bin/send_email_linux_amd64 cmd/send_email/main.go
	chmod +x bin/send_email_linux_amd64
	cp bin/send_email_linux_amd64 /tmp/aws_deploy/main
	cd /tmp/aws_deploy/ && zip main.zip main
	cp /tmp/aws_deploy/main.zip deploy/
	rm -rf /tmp/aws_deploy

