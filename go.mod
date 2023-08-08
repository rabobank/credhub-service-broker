module github.com/rabobank/credhub-service-broker

go 1.20

require (
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20211117203709-9b81b3940cc7
	github.com/cloudfoundry-community/go-uaa v0.3.1
	github.com/gorilla/mux v1.8.0
	golang.org/x/oauth2 v0.11.0
)

require (
	code.cloudfoundry.org/gofileutils v0.0.0-20170111115228-4d0c80011a0f // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/cloudfoundry-community/go-cfenv v1.18.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.14.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

exclude (
	golang.org/x/text v0.3.0
	golang.org/x/text v0.3.2
	gopkg.in/yaml.v2 v2.2.1
	gopkg.in/yaml.v2 v2.2.2
)
