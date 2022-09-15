# Copyright (c) 2019-2022 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   Red Hat, Inc. - initial API and implementation
#   IBM Corporation - implementation
#

# https://access.redhat.com/containers/?tab=tags#/registry.access.redhat.com/rhel8/go-toolset
# FROM registry.redhat.io/rhel8/go-toolset:1.16.12-4 as builder
FROM registry.redhat.io/rhel8/go-toolset@sha256:2d56a94f1e83fe0e874f099dd8d03b876186985a7dbeea118e8c725b1b07b9f0 as builder
ENV GOPATH=/go/
ARG DEV_WORKSPACE_CONTROLLER_VERSION="v0.9.0"
ARG DEV_HEADER_REWRITE_TRAEFIK_PLUGIN="v0.1.2"
ARG TESTS="true"
USER root

COPY asset-* /tmp/

WORKDIR /che-operator

RUN unzip /tmp/asset-devworkspace-operator.zip */deploy/deployment/* -d /tmp && \
    mkdir -p /tmp/devworkspace-operator/templates/ && \
    mv /tmp/devfile-devworkspace-operator-*/deploy/* /tmp/devworkspace-operator/templates/

RUN unzip /tmp/asset-header-rewrite-traefik-plugin.zip -d /tmp && \
    mkdir -p /tmp/header-rewrite-traefik-plugin && \
    mv /tmp/*-header-rewrite-traefik-plugin-*/headerRewrite.go /tmp/*-header-rewrite-traefik-plugin-*/.traefik.yml /tmp/header-rewrite-traefik-plugin

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Copy the go source
COPY main.go main.go
COPY vendor/ vendor/
COPY mocks/ mocks/
COPY api/ api/
COPY templates/ templates/
COPY config/ config/
COPY controllers/ controllers/
COPY pkg/ pkg/

# build operator
RUN export ARCH="$(uname -m)" && if [[ ${ARCH} == "x86_64" ]]; then export ARCH="amd64"; elif [[ ${ARCH} == "aarch64" ]]; then export ARCH="arm64"; elif [[ ${ARCH} == "s390x" ]]; then export ARCH="s390x"; fi && \
    if [[ ${TESTS} == "true" && ${ARCH} == "x86_64" ]]; then export MOCK_API=true && go test -mod=vendor -v ./...; fi && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} GO111MODULE=on go build -mod=vendor -a -o che-operator main.go

# https://access.redhat.com/containers/?tab=tags#/registry.access.redhat.com/ubi8-minimal
# FROM registry.redhat.io/ubi8/ubi-minimal:8.5
FROM registry.redhat.io/ubi8/ubi-minimal@sha256:3f32ebba0cbf3849a48372d4fc3a4ce70816f248d39eb50da7ea5f15c7f9d120

# Build Arguments
ARG PRODUCT_VERSION=2.2.0

ENV PRODUCT="IBM Wazi for Dev Spaces" \
    COMPANY="IBM" \
    VERSION=$PRODUCT_VERSION \
    RELEASE="1" \
    SUMMARY="IBM Wazi for Dev Spaces" \
    DESCRIPTION="IBM Wazi for Dev Spaces - Operator" \
    PRODTAG="ibm-wazi-for-devspaces" \
    CLOUDPAK_ID="9d41d2d8126f4200b62ba1acc0dffa2e" \
    CLOUDPAK_METRIC="VIRTUAL_PROCESSOR_CORE" \
    PRODUCT_NAME="Wazi Code" \
    PRODUCT_ID="0e775d0d3f354a5ca074a6a4398045f3" \
    PRODUCT_METRIC="AUTHORIZED_USER" \
    PRODUCT_CHARGED_CONTAINERS="All" \
    PRODUCT_CLOUDPAK_RATIO="5:1"
        
LABEL name="$PRODUCT" \
      vendor="$COMPANY" \
      version="$VERSION" \
      release="$RELEASE" \
      license="EPLv2" \
      summary="$SUMMARY" \
      description="$DESCRIPTION" \
      io.k8s.description="$DESCRIPTION" \
      io.k8s.display-name="$SUMMARY" \
      io.openshift.tags="$PRODTAG,$COMPANY" \
      com.redhat.component="$PRODTAG" \
      io.openshift.expose-services="" \
      cloudpakName="$PRODUCT" \
      cloudpakId="$CLOUDPAK_ID" \
      cloudpakMetric="$CLOUDPAK_METRIC" \
      productName="$PRODUCT_NAME" \
      productVersion="$VERSION" \
      productID="$PRODUCT_ID" \
      productMetric="$PRODUCT_METRIC" \
      productChargedContainers="$PRODUCT_CHARGED_CONTAINERS" \
      productCloudpakRatio="$PRODUCT_CLOUDPAK_RATIO"

# install httpd-tools for /usr/bin/htpasswd
RUN microdnf install -y httpd-tools && microdnf -y update && microdnf -y clean all && rm -rf /var/cache/yum && echo "Installed Packages" && rpm -qa | sort -V && echo "End Of Installed Packages" && \
    mkdir ~/.ssh && chmod 0766  ~/.ssh

COPY ./LICENSE /licenses/
COPY --from=builder /tmp/devworkspace-operator/templates /tmp/devworkspace-operator/templates
COPY --from=builder /tmp/header-rewrite-traefik-plugin /tmp/header-rewrite-traefik-plugin
COPY --from=builder /che-operator/templates/*.sh /tmp/
COPY --from=builder /che-operator/che-operator /manager

WORKDIR /
USER 65532:65532

ENTRYPOINT ["/manager"]

# append Brew metadata here - see https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/crw-operator_2.x.jenkinsfile
