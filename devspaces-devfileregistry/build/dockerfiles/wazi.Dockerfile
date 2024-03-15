#
# Copyright (c) 2018-2022 Red Hat, Inc.
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

# Builder: check meta.yamls and create index.json
# https://registry.access.redhat.com/ubi8/python-311
FROM registry.access.redhat.com/ubi8/python-311:latest as devfile-builder
USER 0

################# 
# PHASE ONE: create ubi8 image with yq
################# 

ARG BOOTSTRAP=false
ENV BOOTSTRAP=${BOOTSTRAP}
ARG ALLOWED_REGISTRIES=""
ENV ALLOWED_REGISTRIES=${ALLOWED_REGISTRIES}
ARG ALLOWED_TAGS=""
ENV ALLOWED_TAGS=${ALLOWED_TAGS}

COPY ./build/dockerfiles/content_set*.repo /etc/yum.repos.d/
COPY ./build/dockerfiles/rhel.install.sh /tmp
RUN /tmp/rhel.install.sh && rm -f /tmp/rhel.install.sh

COPY ./build/scripts VERSION versions.json /build/
COPY ./devfiles /build/devfiles
COPY ./resources /build/resources

WORKDIR /build/

RUN \
    /build/generate_devworkspace_templates.sh && \
    /build/wazi_devfile.sh --devworkspaces

# Registry, organization, and tag to use for base images in dockerfiles. Devfiles
# will be rewritten during build to use these values for base images.
ARG PATCHED_IMAGES_REG="quay.io"
ARG PATCHED_IMAGES_ORG="eclipse"
ARG PATCHED_IMAGES_TAG="next"

# validate devfile content
RUN \
    ./check_referenced_images.sh devfiles --registries "${ALLOWED_REGISTRIES}" --tags "${ALLOWED_TAGS}" && \
    ./check_mandatory_fields.sh devfiles

# Cache projects in DS 
RUN \
    ./index.sh > /build/devfiles/index.json && \
    cp /build/devfiles/index.json /build/index && \
    ./list_referenced_images.sh devfiles > /build/devfiles/external_images.txt && \
    ./list_referenced_images_by_file.sh devfiles > /build/devfiles/external_images_by_devfile.txt && \
    chmod -R g+rwX /build/devfiles /build/resources

################# 
# PHASE TWO: configure registry image
################# 

# Build registry, copying meta.yamls and index.json from builder
# https://registry.access.redhat.com/ubi8/httpd-24
FROM registry.access.redhat.com/ubi8/httpd-24:latest AS devfile-registry
USER 0

# latest httpd container doesn't include ssl cert, so generate one
RUN \
    chmod +x /usr/share/container-scripts/httpd/pre-init/40-ssl-certs.sh && \
    /usr/share/container-scripts/httpd/pre-init/40-ssl-certs.sh
RUN \
    YUM_PKGS="jq" && \
    yum -y install --nodocs ${YUM_PKGS} && \
    yum -y -q update && \
    yum -y -q clean all && rm -rf /var/cache/yum && \
    echo "Installed Packages" && rpm -qa | sort -V && echo "End Of Installed Packages"

RUN \
    echo '<FilesMatch "^\.ht">' >> /etc/httpd/conf/httpd.conf && \
    echo "Require all denied" >> /etc/httpd/conf/httpd.conf && \
    echo "</FilesMatch>" >> /etc/httpd/conf/httpd.conf

RUN \
    sed -i /etc/httpd/conf.d/ssl.conf \
    -e "s,.*SSLProtocol.*,SSLProtocol all -SSLv3," \
    -e "s,.*SSLCipherSuite.*,SSLCipherSuite HIGH:!aNULL:!MD5,"

RUN \
    sed -i /etc/httpd/conf/httpd.conf \
    -e "s,Listen 80,Listen 8080," \
    -e "s,logs/error_log,/dev/stderr," \
    -e "s,logs/access_log,/dev/stdout," \
    -e "s,AllowOverride None,AllowOverride All," && \
    chmod a+rwX /etc/httpd/conf /run/httpd /etc/httpd/logs/

STOPSIGNAL SIGWINCH

ARG DS_BRANCH=devspaces-3.10-rhel-8
ENV DS_BRANCH=${DS_BRANCH}

WORKDIR /var/www/html

RUN mkdir -m 777 /var/www/html/devfiles
COPY README.md .htaccess /var/www/html/
COPY --from=devfile-builder /build/devfiles /var/www/html/devfiles
COPY --from=devfile-builder /build/resources /var/www/html/resources
COPY --from=devfile-builder /build/devfiles/index.json /var/www/html/index
COPY ./images /var/www/html/images
COPY ./LICENSE /licenses/
COPY ./build/dockerfiles/rhel.entrypoint.sh ./build/dockerfiles/entrypoint.sh /usr/local/bin/
RUN chmod g+rwX /usr/local/bin/entrypoint.sh /usr/local/bin/rhel.entrypoint.sh && \
    chgrp -R 0 /var/www/html && chmod -R g+rw /var/www/html

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["/usr/local/bin/rhel.entrypoint.sh"]

ARG PRODUCT_VERSION="4.0.0"
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

# Offline build
FROM devfile-builder AS offline-builder
RUN ./cache_projects.sh devfiles resources && \
    ./cache_images.sh devfiles resources && \
    chmod -R g+rwX /build

FROM devfile-registry AS offline-registry
COPY --from=offline-builder /build/devfiles /var/www/html/devfiles
COPY --from=offline-builder /build/devfiles/index.json /var/www/html/index
COPY --from=offline-builder /build/resources /var/www/html/resources

# append Brew metadata here
