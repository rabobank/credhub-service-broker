module github.com/rabobank/credhub-service-broker

go 1.26.0

replace (
	golang.org/x/net => golang.org/x/net v0.59.0
	golang.org/x/text => golang.org/x/text v0.42.0
)

require (
	github.com/cloudfoundry-community/go-cfenv v1.24.4
	github.com/cloudfoundry-community/go-uaa v0.5.0
	github.com/cloudfoundry/go-cfclient/v3 v3.0.0-beta.1
	github.com/gorilla/mux v1.8.1
	golang.org/x/oauth2 v0.37.0
)

require (
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/go-martini/martini v0.0.0-20170121215854-22fa46961aab // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/martini-contrib/render v0.0.0-20150707142108-ec18f8345a11 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.12.1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
)
