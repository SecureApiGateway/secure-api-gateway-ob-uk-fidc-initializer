# Secure API Gateway - FIDC-Initializer


See [README](https://github.com/SecureApiGateway/secure-api-gateway-ob-uk-fidc-initializer/blob/master/README.md) for information on FIDC-Initializer

## Prerequisites

- Kubernetes v1.23 +
- Helm 3.0.0 +

To add the forgerock helm artifactory repository to your local machine to consume helm charts use the following;

```console
  helm repo add forgerock-helm https://maven.forgerock.org/artifactory/forgerock-helm-virtual/ --username [backstage_username]  --password [backstage_password]
  helm repo update
```

NOTE: You must have a valid [subscription](https://backstage.forgerock.com/knowledge/kb/article/a57648047#XAYQfS) to aquire the `backstage_username` and `backstage_password` values.

## Helm Charts
### Deployment
FIDC-Configurator is primarily used for development environments and not intended to run in production.

As part of a development deployment of the secure-api-gateway, you must build the java artifacts and built the docker image via the [Makefile](https://github.com/SecureApiGateway/secure-api-gateway-ob-uk-fidc-initializer/blob/master/Makefile). 

TODO: Add in commands to deploy the FIDC to a development environment without ArgoCD

### Example Manifest
This is an example manifest using the `values.yaml` file provided, there is no overlay values in this generated manifest hence why there is no repo URL in `spec.jobTemplate.spec.containers.0.image`.

```yaml
---
# Source: securebanking-openbanking-uk-iam-initializer/templates/cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: fidc-configurator
spec:
  schedule: "* * * * *"
  concurrencyPolicy: Forbid
  successfulJobsHistoryLimit: 1
  startingDeadlineSeconds: 180
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: fidc-configurator
              image: ":1.0.0"
              imagePullPolicy: Always
              env:
                - name: ENVIRONMENT.STRICT
                  value: "true"
                - name: ENVIRONMENT.CLOUDTYPE
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: CLOUD_TYPE
                - name: IDENTITY_PLATFORM_FQDN # variable to run the command shell, the shell doesn't support variables with dot.
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: IDENTITY_PLATFORM_FQDN
                - name: HOSTS.BASE_FQDN
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: BASE_FQDN
                - name: HOSTS.IG_FQDN
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: IG_FQDN
                - name: HOSTS.MTLS_FQDN
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: MTLS_FQDN
                - name: HOSTS.IDENTITY_PLATFORM_FQDN
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: IDENTITY_PLATFORM_FQDN
                - name: IDENTITY.DEFAULT_USER_AUTHENTICATION_SERVICE
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: IDENTITY_DEFAULT_USER_AUTHENTICATION_SERVICE
                      optional: true
                - name: IDENTITY.GOOGLE_SECRET_STORE_NAME
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: IDENTITY_GOOGLE_SECRET_STORE_NAME
                      optional: true
                - name: IDENTITY.GOOGLE_SECRET_STORE_OAUTH2_CA_CERTS_SECRET_NAME
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: IDENTITY_GOOGLE_SECRET_STORE_OAUTH2_CA_CERTS_SECRET_NAME
                      optional: true
                - name: USERS.FR_PLATFORM_ADMIN_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: initializer-secret
                      key: cdm-admin-password
                - name: USERS.FR_PLATFORM_ADMIN_USERNAME
                  valueFrom:
                    secretKeyRef:
                      name: initializer-secret
                      key: cdm-admin-user
                - name: IDENTITY.REMOTE_CONSENT_SIGNING_PUBLIC_KEY
                  valueFrom:
                    secretKeyRef:
                      name: rcs-signing
                      key: rcs-signing.pem
                - name: IDENTITY.REMOTE_CONSENT_SIGNING_KEY_ID
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: RCS_CONSENT_RESPONSE_JWT_SIGNINGKEYID
                - name: IDENTITY.REMOTE_CONSENT_ID
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: RCS_CONSENT_RESPONSE_JWT_ISSUER
                - name: IG.IG_CLIENT_ID
                  valueFrom:
                    secretKeyRef:
                      name: as-sapig-secrets
                      key: IG_CLIENT_ID
                - name: IG.IG_CLIENT_SECRET
                  valueFrom:
                    secretKeyRef:
                      name: as-sapig-secrets
                      key: IG_CLIENT_SECRET
                - name: IG.IG_IDM_USER
                  valueFrom:
                    secretKeyRef:
                      name: as-sapig-secrets
                      key: IG_IDM_USER
                - name: IG.IG_IDM_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: as-sapig-secrets
                      key: IG_IDM_PASSWORD
                - name: IG.IG_AGENT_ID
                  valueFrom:
                    secretKeyRef:
                      name: as-sapig-secrets
                      key: IG_AGENT_ID
                - name: IG.IG_AGENT_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: as-sapig-secrets
                      key: IG_AGENT_PASSWORD
              command: [ "/bin/sh", "-c" ]
              args:
                - |
                  echo "IDENTITY_PLATFORM_FQDN $IDENTITY_PLATFORM_FQDN"
                  until $(curl -X GET --output /dev/null --silent --head --fail -H "X-OpenIDM-Username: anonymous" \
                  -H "X-OpenIDM-Password: anonymous" -H "X-OpenIDM-NoSession: true" \
                  https://$IDENTITY_PLATFORM_FQDN/openidm/info/ping)
                  do
                  echo "IDM not ready"
                  sleep 10
                  done
                  ./initialize
          restartPolicy: OnFailure
```
### Environment Variables

These are the environment variables declared in the `cronjob.yaml`;
| Key | Default | Description | Source | Optional |
|-----|---------|-------------|--------|----------|
| ENVIRONMENT.STRICT | true | If true, any errors will cause the job to exit | cronjob.environment.strict |
| ENVIRONMENT.CLOUDTYPE | FIDC | Type of Cloud Instance being ran, depends on what environment you are running | as-sapig-deployment-config |
| IDENTITY_PLATFORM_FQDN | iam.forgerock.financial | Custom Domain created in Cloud Instance | as-sapig-deployment-config |
| HOSTS.BASE_FQDN | forgerock.financial | Base DNS to be used | as-sapig-deployment-config |
| HOSTS.IG_FQDN | sapig.forgerock.financial | IG DNS to be used | as-sapig-deployment-config |
| HOSTS.MTLS_FQDN | mtls.sapig.forgerock.financial | mtls DNS to be used | as-sapig-deployment-config |
| HOSTS.IDENTITY_PLATFORM_FQDN | iam.forgerock.financial | Custom Domain created in Cloud Instance | as-sapig-deployment-config |
| IDENTITY.DEFAULT_USER_AUTHENTICATION_SERVICE | | | as-sapig-deployment-config | X |
| IDENTITY.GOOGLE_SECRET_STORE_NAME | | | as-sapig-deployment-config | X |
| IDENTITY.GOOGLE_SECRET_STORE_OAUTH2_CA_CERTS_SECRET_NAME | | | as-sapig-deployment-config | X |
| USERS.FR_PLATFORM_ADMIN_PASSWORD | | Password for cloud instance. NOTE - This password can be used for `initializer-secret` or `am-env-secrets` depending on `ENVIRONMENT.CLOUDTYPE` set | If `ENVIRONMENT.CLOUDTYPE=FIDC` initializer-secret/cdm-admin-password else am-env-secrets/AM_PASSWORDS_AMADMIN_CLEAR |
| USERS.FR_PLATFORM_ADMIN_USERNAME | | Username for cloud instance, only populated if `ENVIRONMENT.CLOUDTYPE=FIDC` | initializer-secret/cdm-admin-user |
| IDENTITY.REMOTE_CONSENT_SIGNING_PUBLIC_KEY | PEM File | The pem file to be used for RCS signing | rcs-signing secret | 
| IDENTITY.REMOTE_CONSENT_SIGNING_KEY_ID | rcs-jwt-signer | | as-sapig-deployment-config/RCS_CONSENT_RESPONSE_JWT_SIGNINGKEYID 
IDENTITY.REMOTE_CONSENT_ID | secure-open-banking-rcs | | as-sapig-deployment-config/RCS_CONSENT_RESPONSE_JWT_SIGNINGKEYID | 
| IG.IG_CLIENT_ID | | | as-sapig-secrets |
| IG.IG_CLIENT_SECRET | | | as-sapig-secrets |
| IG.IG_IDM_USER | | | as-sapig-secrets |
| IG.IG_IDM_PASSWORD | | | as-sapig-secrets |
| IG.IG_AGENT_ID | | | as-sapig-secrets |
| IG.IG_AGENT_PASSWORD | | | as-sapig-secrets |

### Values
These are the values that are consumed in the `deployment.yaml` and `service.yaml`;
| Key | Type | Description | Default |
|-----|------|-------------|---------|
| cronjob.concurrencyPolicy | string |   # Policy of allowing concurrent pods to run | Forbid |
| cronjob.environment.cloudType | string | Type of Cloud Instance being ran, depends on what environment you are running | FIDC |
| cronjob.environment.strict | bool | If true, any errors will cause the job to exit | true |
| cronjob.image.repo | string | Repo to pull images from - Value should exist in values.yaml overlay in deployment repo | {} |
| cronjob.image.tag | string | Tag to deploy - Value should exist in values.yaml overlay in deployment repo | {} |
| cronjob.image.imagePullPolicy | string | Policy for pulling images | Always |
| cronjob.schedule | cron expression | What schedule the cronjob should run on | * * * * * (Every minute) |
| cronjob.seccessfulJobHistoryLimit | integer | How many successful jobs should be kept for histroy | 1 |
| cronjob.startingDeadlineSeconds | integer | Time in seconds to deplay starting the cronjob once deployed | 180 |
| cronjob.restartPolicy | string | When to restart the pod | OnFailure |

NOTE: There is no `deployment.image.repo` or `deployment.image.tag` specified in the `Values.yaml` - This needs to be done in a seperate 'deployments' repo using an additional `values.yaml` overlay. You may overwrite any of the other values in this additonal file if required.

Example of the RCS section of the additonal `values.yaml` file;
```yaml
fidc-configurator:
  deployment:  
    image:
      repo: [REPO_URL]
      # By default the AppVersion will be used so that users don't have to change this value, however you can override this by uncommenting the line and providing a valid verison.
      # tag: 1.0.1
```
## Support

For any issues or questions, please raise an issue within the [SecureApiGateway](https://github.com/SecureApiGateway/SecureApiGateway/issues) repository.