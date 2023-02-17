module github.com/rabobank/credhub-service-broker

go 1.20

require (
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20211117203709-9b81b3940cc7
	github.com/cloudfoundry-community/go-uaa v0.3.1
	github.com/gorilla/mux v1.8.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

require (
	code.cloudfoundry.org/gofileutils v0.0.0-20170111115228-4d0c80011a0f // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	golang.org/x/net v0.0.0-20190611141213-3f473d35a33a // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

exclude (
	golang.org/x/text v0.3.0
	golang.org/x/text v0.3.2
	gopkg.in/yaml.v2 v2.2.1
	gopkg.in/yaml.v2 v2.2.2
)
