#
# Copyright IBM Corporation 2021 - 2023
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   IBM Corporation - implementation
#

FROM registry.redhat.io/openshift4/ose-operator-registry@sha256:ecef9cded6d99990770248529b074b127e2c69bca00e2cccaf09adffeb40c02b AS builder
#FROM registry.redhat.io/openshift4/ose-operator-registry:v4.13 AS builder

FROM registry.redhat.io/ubi8/ubi-minimal:latest

ARG PRODUCT_VERSION="4.0.0"
LABEL operators.operatorframework.io.index.database.v1=/database/index.db

COPY LICENSE /licenses/
COPY bundles.db /database/index.db
COPY --from=builder /bin/registry-server /bin/registry-server
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

RUN microdnf update -y && \
    microdnf clean all && \
    mkdir -p /registry && \
    chgrp -R 0 /registry && \
    chmod -R g+rwx /registry

WORKDIR /registry

EXPOSE 50051
USER 1001

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

ENTRYPOINT ["/bin/registry-server"]
CMD ["--database", "/database/index.db", "--termination-log=log.txt"]