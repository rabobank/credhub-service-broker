module github.com/rabobank/credhub-service-broker

go 1.23

replace (
	golang.org/x/text => golang.org/x/text v0.19.0
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20220930021109-9c4e6c59ccf1
	github.com/cloudfoundry-community/go-cfenv v1.18.0
	github.com/cloudfoundry-community/go-uaa v0.3.3
	github.com/gorilla/mux v1.8.1
	golang.org/x/oauth2 v0.23.0
)

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.30.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
