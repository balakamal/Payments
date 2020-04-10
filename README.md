# go-rest-api

## Folder structure

- client - contains client libraries for external services to consume
- cmd - entry point to the application
- implementation - service layer containing the business logic
- middleware - performs tasks such as logging, monitoring etc
- pkg - includes external dependencies
- repository - db access code
- service - interface for teh implementation
- transport - contains the transport protocols eg: http and protobuf

## How to compile the applicationm
`cd cmd && go build -o main .`

## How to dockerize the applciation
We’re disabling cgo which gives us a static binary. We’re also setting the OS to Linux (in case someone builds this on a Mac or Windows) and the -a flag means to rebuild all the packages we’re using, which means all the imports will be rebuilt with cgo disabled. 

`CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .`

Run this from the root folder
`docker run -it example-scratch`

`