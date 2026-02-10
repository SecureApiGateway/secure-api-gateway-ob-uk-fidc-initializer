# securebanking-openbanking-uk-fidc-initializer

**Note** this repository is still in active development. Please aim to check back often for updates. 

The *securebanking-openbanking-uk-fidc-initializer* configures an instance of the [ForgeRock identity platform](https://www.forgerock.com/identity-and-access-management-platform) for use by the [Secure Banking Access Toolkit (SBAT)](https://github.com/SecureBankingAccessToolkit/SecureBankingAccessToolkit/wiki). The combination of an SBAT deployment and the ForgeRock Identity Platform will address the requirements of;
- [UK OpenBanking Specification](https://www.openbanking.org.uk/)
- [Financial-grade API, Read and Write API Security Profile](https://standards.openbanking.org.uk/security-profiles/)

To understand how this initializer fits into the Secure Banking Access Toolkit please check out the SBAT documentation [here](https://github.com/SecureBankingAccessToolkit/SecureBankingAccessToolkit/wiki).

## Requirements

- [go 1.15](https://golang.org/doc/install)
- configure [gopath](https://golang.org/doc/gopath_code.html#GOPATH)
- [pact](https://github.com/pact-foundation/pact-go#installation-on-nix)

## Building the Initializer

To build the initializer you can use the make file;

```bash
$ make
go mod download
rm -f initialize
go build -o initialize
```

This will build an executable called `initialize` 

NOTE: With the recent split of the config from the initializer code to [fr-platform-config](https://github.com/SecureApiGateway/fr-platform-config) the build command in the make file will pull a version of the repo, the version is defined in the [.congif-git-ref](https://github.com/SecureApiGateway/securebanking-openbanking-uk-fidc-initializer/blob/master/.config-git-ref), so that it is possible to make changes to the master branch for testing without breaking anything. The SHA will need updating once changes have been made and are tested

### Build a docker image

A Dockerfile is provided that will build a docker image based on the popular Alpine linux image. This can then be deployed into your SBAT environment as a Job and used to initialize your ForgeRock platform instance. 
**Note:** The initializer is great for development environments, where it can be used to intialize a fresh ForgeRock platfrom as part of a CI/CD methodology. For production deployments, a GitOps approach to configuration may be more appropriate for your needs. 

### Build the docker file

The docker file can be built and pushed to your docker image registry (e.g. gcr, docker hub etc) using the docker build command;

```bash
$ export DOCKER_REPO=<your repo here>
$ env GOOS=linux GOARCH=amd64 go build -o initialize
$ docker build -t $DOCKER_REPO .
$ docker push $DOCKER_REPO
```

## Deploying the Initializer

The repository contains an example helm chart that can be used to deploy the intializer into kubernetes.

Currently there are two helm charts. The one to use is the chart found in `_infra/helm/securebanking-openbanking-uk-iam-initializer`

### Helm install

In the example command below the iam initializer is being deployed into a namespace called local-dev-sbat where IG is running on a domain `https://sbat.openbanking.bigbank.com`. It is being used to configure an instance of the ForgeRock Identity Platform that is available at `https://iam.openbanking.bigbank.com`.

```bash
$ helm upgrade iam-init ./ --install --namespace local-dev-sbat \
  --set-string environment.CloudType=CDK \
  --set-string environment.fr_platform.fqdn=iam.openbanking.bigbank.com
  --set-string environment.sbat.domain=sbat.openbanking.bigbank.com --wait
```

### Helm Chart Values

The following values may be supplied when installing/upgrading:

| variable                       | Required | Default value | Description                                                                                                                                                                                        |
|--------------------------------|----------|---------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| iam_initializer_image_location | yes      | None          | The location of the docker image                                                                                                                                                                   |
| environment.fr_platform.fqdn   | yes      | None          | The Fully Qualified domain name on which the ForgeRock Identity Platform may be reached, e.g. iam.openbanking.bigbank.com                                                                          |
| environment.fr_platform.type   | no       | FIDC          | The type of ForgeRock Platform Deployment. Value values are:<br>- **CDK:** deploying against a deployment of the forgeops CDK<br>- **FIDC:** Deploying against a ForgeRock Identity Cloud instance |
| environment.sbat.projectId     | yes      | None          | The GCP project in which the SBAT is deployed                                                                                                                                                      |
| environment.sbat.fqdn          | yes      | None          | The fqdn on which the SBAT components are deployed                                                                                                                                                 |


## Initializer Configuration

The initializer application provides a default configuration yaml file (properties values to run the application), the default configuration yaml file is loaded using the [viper library](https://github.com/spf13/viper),
the initializer application supports a personalized configuration file (as a profile) that it can be personalized for each required environment following the below rules:

- Path of environment file: `config/viper`
- Pattern environment file name: `viper-${environment-profile.viper_config}-configuration.yaml`
- Format configuration file (extension file): `yaml`

**Example:** `viper-default-configuration.yaml`
> You will find the provides default configuration yaml file in `config/viper` as an example.

> :warning: The initializer application only supports one configuration yaml file by application instance. It's recommended copy the provided default configuration yaml file and change its values.

> :memo: All the variables/properties values provides by the configuration file can be overwritten using environment variables or a kubernetes config map.
> ```shell
> go build -o initialize \
> env ENVIRONMENT.VERBOSE=true ./initialize
> ```

#### ConfigMap for variables mount example

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: initializer-config
data:
  HOSTS.IG_FQDN: obdemo.dev.forgerock.financial
  HOSTS.IDENTITY_PLATFORM_FQDN: iam.dev.forgerock.financial
  ...
          
```

**Check the variables/properties in [Configuration variables](#configuration-variables) section.**

### The application configuration file
The configuration file is loaded from the path `config/viper` following the pattern `viper- + ${environment-name} + -configuration`
where the environment/profile can be specified by environment variable, passing that environment variable to the program.
**Examples**
```shell
go build -o initialize
```
```shell
env ENVIRONMENT.VIPER_CONFIG=MY-ENVIRONMENT-PROFILE-VIPER_CONFIG ./initialize
```
> The application will attempt to load the configuration file `viper-MY-ENVIRONMENT-PROFILE-VIPER_CONFIG-configuration.yaml`

**Other example**
```shell
env ENVIRONMENT.VIPER_CONFIG=MY-ENVIRONMENT-PROFILE-VIPER_CONFIG ENVIRONMENT.VERBOSE=true ... ./initialize
```

#### Configuration variables
**Environment variables**
There are a variables used before load the configuration file and these variables can change the behaviour of the application.
- Behaviour variables:
  - `ENVIRONMENT.VERBOSE`
  - `ENVIRONMENT.VIPER_CONFIG`
  - `ENVIRONMENT.STRICT`
  - `ENVIRONMENT.ONLY_CONFIG`

<details>
<summary>Variables table</summary>
<!-- always an empty line before table -->

| variable                                     | Default value                        | Description                                                                                                                                                                                                       |
|----------------------------------------------|--------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ENVIRONMENT.VERBOSE`                        | false                                | Log level (verbose=true means debug mode)                                                                                                                                                                         |
| `ENVIRONMENT.VIPER_CONFIG`                   | default                              | The profile that contains the configuration to be overwritten from system env                                                                                                                                     |
| `ENVIRONMENT.ONLY_CONFIG`                    | false                                | Prints the configuration and exiting the program, to review the properties                                                                                                                                        |
| `ENVIRONMENT.STRICT`                         | false                                | true = strict mode on, otherwise off, will exit if go resty returns an error in STRICT mode enabled, be it client error, server error or other. Turning off STRICT mode will simply warn of client/server errors. |
| `ENVIRONMENT.CLOUDTYPE`                           | CDK                                  | values: CDK, CDM or FIDC,  to identify the kind of identity platform                                                                                                                                              |
| `ENVIRONMENT.PATHS.CONFIG_BASE_DIRECTORY`    | config/defaults/                     | Base configuration root path folder for data files and templates to populate them into identity platform                                                                                                          |
| `ENVIRONMENT.PATHS.CONFIG_SECURE_BANKING`    | config/defaults/secure-open-banking/ | Base configuration path folder for specific secure open banking data files and templates to populate them into identity platform                                                                                  |
| `ENVIRONMENT.PATHS.CONFIG_IDENTITY_PLATFORM` | config/defaults/identity-platform/   | Base configuration path folder for generic data files and templates to populate them into identity platform                                                                                                       |
</details>

**Host variables**
<details>
<summary>Table</summary>
<!-- always an empty line before table -->

| Environment variable           | default                        | description                                  |
|--------------------------------|--------------------------------|----------------------------------------------|
| `HOSTS.IDENTITY_PLATFORM_FQDN` | iam.dev.forgerock.financial    | Identity platform Full Qualified Domain Name |
| `HOSTS.IG_FQDN`                | obdemo.dev.forgerock.financial | Ig Full Qualified Domain Name                |
| `HOSTS.SCHEME`                 | https                          | URI scheme, Syntax part of a generic URI     |
</details>

**IG variables**
<details>
<summary>Table</summary>
<!-- always an empty line before table -->

| Environment variable   | default               | description                                                                                                                                                                                                                                            |
|------------------------|-----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `IG.IG_CLIENT_ID`            | ig-client             | The initializer creates an OAuth2 Client that the SBAT IG will use to authenticate to the FR Platform to ensure that Api Client's can't bypass IG and use the FR Identity Platform APIs directly. The OAuth2 client will be created with this id       |
| `IG.IG_CLIENT_SECRET`        | add-here-the-password | The initializer creates an OAuth2 Client that the SBAT IG will use to authenticate to the FR Platform to ensure that Api Client's can't bypass IG and use the FR Identity Platform APIs directly. The OAuth2 client will be created with this password |
| `IG.IG_RCS_SECRET`           | add-here-the-secret   | IG rcs secret for remote consent service |
| `IG.IG_SSA_SECRET`           | add-here-the-secret   | IG ssa secret for software publisher agent  |
| `IG.IG_AS_IDM_CLIENT_SECRET` |  add-here-the-secret    | IG secret for IDM Client used on AS  |
| `IG.IG_RS_IDM_CLIENT_SECRET` |  add-here-the-secret | IG secret for IDM Client used on RS |
| `IG.IG_AGENT_ID`             | ig-agent              | IG agent id for IG policy agent |
| `IG.IG_AGENT_PASSWORD`       | add-here-the-password | Ig agent password for IG policy agent |
</details>

**Identity variables**
<details>
<summary>Table</summary>
<!-- always an empty line before table -->

| Environment variable                          | default                 | description                                              |
|-----------------------------------------------|-------------------------|----------------------------------------------------------|
| `IDENTITY.AM_REALM`                           | alpha                   | The realm used for secure banking                        |
| `IDENTITY.IDM_CLIENT_ID`                      | policy-client           | Placeholder to create Open Banking Dynamic Policy script |
| `IDENTITY.IDM_CLIENT_SECRET`                  | password                | Placeholder to create Open Banking Dynamic Policy script |
| `IDENTITY.SERVICE_ACCOUNT_POLICY`             | service_account.policy  | Service account for Open banking policy                  |
| `IDENTITY.REMOTE_CONSENT_ID`                  | secure-open-banking-rcs | Identification of remote consent agent                   |
| `IDENTITY.OBRI_SOFTWARE_PUBLISHER_AGENT_NAME` | OBRI                    | software publisher agent name                            |
| `IDENTITY.TEST_SOFTWARE_PUBLISHER_AGENT_NAME` | test-publisher          | test software publisher agent                            |
</details>

**Users variables**
<details>
<summary>Table</summary>
<!-- always an empty line before table -->

| Environment variable               | default                                                    | description                                                               |
|------------------------------------|------------------------------------------------------------|---------------------------------------------------------------------------|
| `USERS.FR_PLATFORM_ADMIN_USERNAME` | amadmin (for cloud platform needs to be used a admin user) | Identity platform Username with admin grants (must exist previously)      |
| `USERS.FR_PLATFROM_ADMIN_PASSWORD` | replace-me                                                 | Identity platform User password with admin grants (must exist previously) |

</details>

## Json Identify platform configuration files
Identity Platform JSON files configuration can be added to the config/defaults/${type} directory under either the [additional](./config/defaults/managed-objects/additional) or [openbanking](./config/defaults/managed-objects/openbanking) path. The files must be json and the filenames must match the name of the managed object.
Eg: the managed object with name `apiClient` must be contained in a filename called `apiClient.json`
The initializer will attempt to match the filename (minus suffix) to an IDM managed object of the same name. If none are found within IDM then the initializer will create a new idm managed object.

## Kubernetes ConfigMap
You can override all identity platform configuration files with config predefined within a kubernetes config map.

> :warning: If a path variable as is set to the default relative path of `config/defaults/` then default pre-baked configuration json objects will be used and not your mounted ConfigMap

### ConfigMap for identity platform files mount example

```
spec:
  volumes:
  - name: ob-defaults-objects
    configMap:
      name: openbanking-objects
  containers:
  - name: init-container
    env:
    - name: ENVIRONMENT.PATHS.CONFIG_BASE_DIRECTORY
      value: /opt/config/
    volumeMounts:
    - mountPath: /opt/config/
      name: ob-managed-objects
      readOnly: true
    - name: ENVIRONMENT.PATHS.CONFIG_SECURE_BANKING
      value: /opt/config/secure-open-banking/
    volumeMounts:
    - mountPath: /opt/config/secure-open-banking/
      name: ob-managed-objects
      readOnly: true
    - name: ENVIRONMENT.PATHS.CONFIG_IDENTITY_PLATFORM
      value: /opt/config/identity-platform/
    volumeMounts:
    - mountPath: /opt/config/identity-platform/
      name: ob-managed-objects
      readOnly: true      
```

## Running tests
The tests run against a mockserver which is supplied by [Pact](https://docs.pact.io/). It is used specifically to test internal logic rather than to verify the provider contract.
running the `make test-ci` target will download the required binaries to be able to run the pact tests. this target is used for github actions but can work locally too (if you do not have the pact bonaries installed)

## Test locally the program against an environment
For default the deployment will pull the latest initializer image, the latest image is generated when a PR is merged into master.
When a PR is created a github action is triggered to build and push the initializer image tagged with the PR-NUMBER.
To test the initializer against an environment follow the below steps:
  - Build the program `go mod tidy` and `go build -o initialize`
  - Set your configuration, follow the section [initialize configuration](#initializer-configuration)
  - Run the program:
    - `./initialize` with the default configuration
    - `env ENVIRONMENT.VIPER_CONFIG=${MY-PROFILE} ./initialize` to set your configuration profile

### Commands
| Command                  | description                                                                                                          |
|--------------------------|----------------------------------------------------------------------------------------------------------------------|
| `go mod tidy`            | add missing and remove unused modules                                                                                |
| `go build -o initialize` | compiles the packages named by the import paths, along with their dependencies, but it does not install the results. |
| `./initialize`           | run the compiled program                                                                                             |
| `go run`                 | compiles and runs the named main Go package                                                                          |

> For more information about go command `go help`

