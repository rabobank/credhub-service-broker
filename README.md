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

## Deploying the broker

### Preparing the app

First make sure the broker itself runs (as a cf app, since it needs access to credhub.service.cf.internal), and the URL is available to the Cloud Controller.

To do this, get the latest binary of the `credhub-service-broker` and create an empty folder on your machine.

Put the `credhub-service-broker` binary in the folder.

Create a directory called `config` with a file called `catalog.json` in it.

An example `catalog.json` could look as follows:

```
{
  "services": [
    {
      "name": "pcsb-service",
      "id": "6d45a057-dd05-44d6-8ed7-732b8cc8abe4", # replace GUIDs
      "description": "Credhub Service",
      "bindable": true,
      "plans": [
        {
          "name": "default",
          "id": "39bdd3f6-c4a1-4083-807d-26d13ddb35d1",  # replace GUIDs
          "description": "The default plan."
        }
      ]
    }
  ]
}
```
You can create a `manifest.yml` as well if you want to.

> **_NOTE:_**  Make sure to use `command` with the following property `command: chmod 755 credhub-service-broker && ./credhub-service-broker`

Afterwards your structure of your directory should look like this:

```
.
├── config
│   └── catalog.json
├── credhub-service-broker
└── manifest.yml
```

Push your app with `--no-start`.

```
cf push -f manifest.yml --no-start
```

### Creating the credentials in the CF runtime credhub
To be able to talk to the CF runtime Credhub make sure to prepare the following things.

* update your local hosts file and have 2 entries:
```
<IP of a cf credhub server> credhub.service.cf.internal
<IP of a cf uaa server> uaa.service.cf.internal
```
* login to credhub with : ``credhub login -s https://credhub.service.cf.internal:8844 --skip-tls-validation --client-name credhub-service-broker --client-secret <credhub admin secret>``

The broker has the envvar `CREDHUB_CREDS_PATH` which points to an entry in the CF runtime credhub where the following 2 credentials should be stored:
* `CSB_BROKER_PASSWORD` - The password for the broker (should be specified when issuing the _cf create-service-broker_ cmd).
* `CSB_CLIENT_SECRET` - The password for `CSB_CLIENT_ID´

To create the proper credhub entry in the CF runtime credhub, use the following sample command: 
```
credhub set -n /brokers/credhub-service-broker/credentials --type json --value='{ "CSB_BROKER_PASSWORD": "pswd1", "CSB_CLIENT_SECRET": "pswd2" }'
```

Afterwards make sure to allow the app to access the credentials. This can be done by extracting your app guid via

```
cf app <app_name> --guid
```

Then use the following command

```
credhub set-permission -p /brokers/credhub-service-broker/credentials -a mtls-app:<app_guid> -o read
```

You can double-check if the app has the right permission now, if you run the `get-permission' command.

```
credhub get-permission -p /brokers/credhub-service-broker/credentials -a mtls-app:<app_guid>
```

which should give you a similar output to the one below

```
credhub get-permission -a mtls-app:<app_guid> -p /brokers/credhub-service-broker/credentials
actor: mtls-app:<app_guid>
operations:
    - read
path: /brokers/credhub-service-broker/credentials
uuid: 6fa4db98-972d-4e30-a2af-2ebb7620a37f
```

if the permissions are wrong you will receive the following message

```
The request could not be completed because the permission does not exist or you do not have sufficient authorization.
```

Now that the necessary credentials and permissions are in place we can go ahead and start the app.

```
cf start <app_name>
```

## Installing the broker

Then install the broker:
```
#  the user and password should match with the user/pass you use when starting the broker app
cf create-service-broker pcsb <broker-user> <broker-password> <https://url.where.the.broker.runs>
```
Give access to the service (all plans to all orgs):
```
cf enable-service-access pcsb-service
```

## Creating a service instance

When creating a service instance you need provide a configuration via the `-c` flag. This can either be inline or the path to a json formatted file. 

For example you can prepare a json file similar to this:

```
example_creds_config.json
{
  "superimportantcred": "test1234"
}
```

Then you create your service instance via 

```
cf create-service pcsb-service default test-instance -c example_creds_config.json
```

This causes a credential to be created within CF credhub. The path of the credential created follows the pattern 
`/pcsb/<service-instance-guid>/credentials`.

## Using a service instance

To make the credential available to your app you need to bind your app to the service instance.

```
cf bind-service <app_name> test-instance
```

As soon as you bind an app to the instance the credhub-service-broker will set the respective permissions for the secret in the CF credhub allowing the app to access the secret.

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