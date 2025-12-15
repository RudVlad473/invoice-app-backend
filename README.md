# Testing

All tests inside `/tests` folder are integrational.
<br/>
Here we're utilizing a local dynamodb setup using docker so we can run our tests by simply calling `go test ./...`,
no strings attached.
<br/>
Although you would need to have docker daemon running, in github actions this is a default,
no need for extra setup, but if you're running these tests on Windows, for example, you would need to figure out how to
run the docker daemon for corresponding platform (for Windows you would need to run `Docker Desktop` app)
