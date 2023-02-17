## Credhub Service Broker

A Cloud Foundry Service Broker that can store secrets in credhub and allow secure access to applications.

## Intro

The configuration for the broker consists of the following environment variables:
* **CSB_DEBUG** - Debugging on or off, default is false.
* **CSB_HTTP_TIMEOUT** - Timeout in seconds for connecting to either UAA or Credhub endpoint, default is 10.
* **CSB_CLIENT_ID** - The uaa client to use for logging in to credhub, should have credhub_admin scope.
* **CSB_CREDHUB_URL** - The URL where to reach credhub (i.e. https://credhub.service.cf.internal:8844). 
* **CSB_BROKER_USER** - The userid for the broker (should be specified issuing the _cf create-service-broker_ cmd).
* **CSB_CATALOG_DIR** - The directory where to find the cf catalog for the broker, the directory should contain a file called catalog.json.
* **CSB_LISTEN_PORT** - The port that the broker should listen on, default is 8080.
* **CSB_CFAPI_URL** - The URL of the cf api (i.e. https://api.sys.mydomain.com).
* **CSB_SKIP_SSL_VALIDATION** - Skip ssl validation or not, default is false.
* **CREDHUB_CREDS_PATH** - The path in credhub where csb should get it's credentials (i.e. /brokers/credhub-service-broker/credentials).

When you create a service instance for the service of this broker, you have to specify the -c option. The value of this option (usually a secret) will be stored in credhub under the key: **_/pcsb/<service-instance-guid>/credentials_**
When an application does a bind on the service-instance, a credhub permission will be created for this app, where the path of the permissions will be equal to the above, operation will be read and the actor will be **_mtls-app:<app-guid>_**, the net effect of this is that the app has permission to read the above credhub entry

## Deploying/installing the broker

First make sure the broker itself runs (as a cf app, since it needs access to credhub.service.cf.internal), and the URL is available to the Cloud Controller.
Then install the broker:
```
#  the user and password should match with the user/pass you use when starting the broker app
cf create-service-broker pcsb <broker-user> <broker-password> <https://url.where.the.broker.runs>
```
Give access to the service (all plans to all orgs):
```
cf enable-service-access pcsb-service
```

## Creating the credentials in the runtime credhub
The broker has the envvar CREDHUB_CREDS_PATH which points to an entry in the runtime credhub where the following 2 credentials should be stored:
* CSB_BROKER_PASSWORD - The password for the broker (should be specified when issuing the _cf create-service-broker_ cmd).
* CSB_CLIENT_SECRET - The password for CSB_CLIENT_ID

To create the proper credhub entry in the runtime credhub, use the following sample command: 
```
credhub set -n /brokers/credhub-service-broker/credentials --type json --value='{ "CSB_BROKER_PASSWORD": "pswd1", "CSB_CLIENT_SECRET": "pswd2" }'
```

## Checking the contents of credhub

To see what the broker is creating in credhub, you can do the following to see that:
* update your local hosts file and have 2 entries:
```
<IP of a cf credhub server> credhub.service.cf.internal
<IP of a cf uaa server> uaa.service.cf.internal
```
* login to credhub with : ``credhub login -s https://credhub.service.cf.internal:8844 --skip-tls-validation --client-name credhub-service-broker --client-secret <credhub admin secret>``
* find all entries created by broker: ``credhub find -n /pcsb``
* show the contents of a credential entry: ``credhub get -n /pcsb/<service-instance-guid>/credentials``
* show the permissions on a credential entry: ``credhub get-permission -p /pcsb/<service-instance-guid>/credentials -a mtls-app:<app-guid>``