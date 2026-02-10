service := secureopenbanking-uk-iam-initializer
name := securebanking-openbanking-uk-iam-initializer
repo := europe-west4-docker.pkg.dev/sbat-gcr-develop/sapig-docker-artifact
helm_repo := forgerock-helm/secure-api-gateway/${name}/
binary-name := initialize
latesttagversion := latest

.PHONY: all
all: mod build
	
mod:
	go mod download

build: clean mod
	go build -o ${binary-name}

test:
	go test ./...

test-ci: mod
	$(eval localPath=$(shell pwd))
	curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash
	PATH=$(PATH):${localPath}/pact/bin go test ./...

clean:
	go clean
	rm -f ${binary-name}

docker: clean mod
ifndef tag
	$(warning no tag supplied; latest assumed)
	$(eval TAG=latest)
else
	$(eval TAG=$(shell echo $(tag) | tr A-Z a-z))
endif
ifndef setlatest
	$(warning no setlatest true|false supplied; false assumed)
	$(eval setlatest=false)
endif
	env GOOS=linux GOARCH=amd64 go build -o initialize
	@if [ "${setlatest}" = "true" ]; then \
		docker buildx build --platform linux/amd64 -t ${repo}/securebanking/${service}:${TAG} -t ${repo}/securebanking/${service}:${latesttagversion} . ; \
		docker push ${repo}/securebanking/${service} --all-tags; \
    else \
   		docker buildx build --platform linux/amd64 -t ${repo}/securebanking/${service}:${TAG} . ; \
   		docker push ${repo}/securebanking/${service}:${TAG}; \
   	fi;

package_helm:
ifndef version
	$(error A version must be supplied, Eg. make helm version=1.0.0)
endif
	helm dependency update _infra/helm/${name}
	helm template _infra/helm/${name}
	helm package _infra/helm/${name} --version ${version} --app-version ${version}

publish_helm:
ifndef version
	$(error A version must be supplied, Eg. make helm version=1.0.0)
endif
	jf rt upload  ./*-${version}.tgz ${helm_repo}
