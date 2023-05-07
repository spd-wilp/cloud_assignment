1. Make sure go v1.20 is installed and configured.

2. To setup codebase run the following

```sh
go mod tidy
go mod vendor
mkdir deploy
```

3. To build generate-metadata lambda zip, run the follwing.This will generate `deploy/main.zip`, which can be deployed in aws lambda.

```sh
make build-generate-metadata-zip
``` 
   
4. To build send-email lambda zip, run the following. This will generate `deploy/main.zip`, which can be deployed in aws lambda. 

```sh
make build-send-email-zip
```