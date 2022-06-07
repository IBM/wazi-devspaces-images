#
# Copyright IBM Corporation 2021 - 2022
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   IBM Corporation - implementation
#

FROM registry.redhat.io/openshift4/ose-operator-registry:v4.10 AS builder

FROM registry.redhat.io/ubi8/ubi-minimal:8.5

# Build Arguments
ARG PRODUCT_VERSION=2.1.0

ENV PRODUCT="IBM Wazi for Dev Spaces" \
    COMPANY="IBM" \
    VERSION=$PRODUCT_VERSION \
    RELEASE="1" \
    SUMMARY="IBM Wazi for Dev Spaces" \
    DESCRIPTION="IBM Wazi for Dev Spaces - Catalog" \
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
      productCloudpakRatio="$PRODUCT_CLOUDPAK_RATIO" \
      operators.operatorframework.io.index.database.v1=/database/index.db

RUN microdnf update -y && \
    microdnf clean all && \
    mkdir -p /registry && \
    chgrp -R 0 /registry && \
    chmod -R g+rwx /registry

COPY LICENSE /licenses/
COPY bundle/bundles.db /database/index.db
COPY --from=builder /bin/registry-server /bin/registry-server
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

WORKDIR /registry

EXPOSE 50051
USER 1001

ENTRYPOINT ["/bin/registry-server"]
CMD ["--database", "/database/index.db", "--termination-log=log.txt"]
