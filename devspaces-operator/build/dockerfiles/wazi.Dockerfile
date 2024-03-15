# Copyright (c) 2019-2023 Red Hat, Inc.
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

FROM registry.redhat.io/rhel8/go-toolset:latest as operator-builder
USER 0

ARG SKIP_TESTS="true"

ENV \
    GOPATH=/go/ \
    DEV_WORKSPACE_VERSION="v0.23.0" \
    TRAEFIK_VERSION="v0.1.2"

WORKDIR /devspaces-operator

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Copy the go source
COPY main.go main.go
COPY vendor/ vendor/
COPY mocks/ mocks/
COPY api/ api/
COPY config/ config/
COPY controllers/ controllers/
COPY pkg/ pkg/

RUN \
    DEV_WORKSPACE_URL="https://api.github.com/repos/devfile/devworkspace-operator/zipball/${DEV_WORKSPACE_VERSION}" && \
    DEV_WORKSPACE_ASSET="devworkspace-operator" && DEV_WORKSPACE_ZIP="asset-${DEV_WORKSPACE_ASSET}.zip" && \
    TRAEFIK_URL="https://api.github.com/repos/che-incubator/header-rewrite-traefik-plugin/zipball/${TRAEFIK_VERSION}" && \
    TRAEFIK_ASSET="header-rewrite-traefik-plugin" && TRAEFIK_ZIP="asset-${TRAEFIK_ASSET}.zip" && \
    mkdir -pv /tmp/{"${DEV_WORKSPACE_ASSET}"/templates,"${TRAEFIK_ASSET}"} && \
    curl -skL "${DEV_WORKSPACE_URL}" -o "/tmp/${DEV_WORKSPACE_ZIP}" && unzip -q -d "/tmp/${DEV_WORKSPACE_ASSET}" "/tmp/${DEV_WORKSPACE_ZIP}" && \
    curl -skL "${TRAEFIK_URL}" -o "/tmp/${TRAEFIK_ZIP}" && unzip -jq -d "/tmp/${TRAEFIK_ASSET}" "/tmp/${TRAEFIK_ZIP}" "*/headerRewrite.go" "*/.traefik.yml" && \
    mv /tmp/${DEV_WORKSPACE_ASSET}/*${DEV_WORKSPACE_ASSET}*/deploy/deployment/* /tmp/${DEV_WORKSPACE_ASSET}/templates/ &&  \
    rm -rf /tmp/${DEV_WORKSPACE_ASSET}/*${DEV_WORKSPACE_ASSET}* && \
    export ARCH="$(uname -m)" && \
    if [[ ${ARCH} == "x86_64" ]]; then export ARCH="amd64"; elif [[ ${ARCH} == "aarch64" ]]; then export ARCH="arm64"; fi && \
    if [[ ${SKIP_TESTS} == "false" ]]; then export MOCK_API=true && go test -mod=vendor -v ./...; fi && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} GO111MODULE=on go build -mod=vendor -v -a -o che-operator main.go


FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ARG PRODUCT_VERSION="4.0.0"

COPY --from=operator-builder /tmp/devworkspace-operator/templates /tmp/devworkspace-operator/templates
COPY --from=operator-builder /tmp/header-rewrite-traefik-plugin /tmp/header-rewrite-traefik-plugin
COPY --from=operator-builder /devspaces-operator/che-operator /manager

ENTRYPOINT ["/manager"]

ENV \
    SUMMARY="IBM Wazi for Dev Spaces" \
    DESCRIPTION="IBM Wazi for Dev Spaces" \
    PRODNAME="Wazi Code" \
    COMPNAME="Wazi" \
    CLOUDPAK_ID="9d41d2d8126f4200b62ba1acc0dffa2e" \
    PRODUCT_ID="0e775d0d3f354a5ca074a6a4398045f3" \
    PRODUCT_METRIC="AUTHORIZED_USER" \
    PRODUCT_CHARGED_CONTAINERS="All" \
    PRODUCT_CLOUDPAK_RATIO="5:1"

LABEL \
      version="${PRODUCT_VERSION}" \
      productVersion="${PRODUCT_VERSION}" \
      maintainer="IBM Corporation" \
      vendor="IBM Corporation" \
      license="EPLv2" \
      name="$SUMMARY" \
      summary="$SUMMARY" \
      description="$DESCRIPTION" \
      io.k8s.description="$DESCRIPTION" \
      io.k8s.display-name="$DESCRIPTION" \
      io.openshift.tags="$PRODNAME,$COMPNAME" \
      com.redhat.component="$PRODNAME-$COMPNAME-container" \
      io.openshift.expose-services="" \
      usage="" \
      cloudpakName="$SUMMARY" \
      cloudpakId="$CLOUDPAK_ID" \
      cloudpakMetric="$CLOUDPAK_METRIC" \
      productName="$PRODNAME" \
      productID="$PRODUCT_ID" \
      productMetric="$PRODUCT_METRIC" \
      productChargedContainers="$PRODUCT_CHARGED_CONTAINERS" \
      productCloudpakRatio="$PRODUCT_CLOUDPAK_RATIO"
