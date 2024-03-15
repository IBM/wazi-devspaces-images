#
# Copyright (c) 2023 Red Hat, Inc.
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

# OpenVSX
FROM ghcr.io/eclipse/openvsx-server:v0.14.2 AS openvsx-server

### Open VSX Server - Builder
FROM registry.access.redhat.com/ubi8/ubi:latest as ovsx-server-builder

ENV \
    CURRENT_BRANCH="devspaces-3.10-rhel-8"

RUN \
    YUM_PKGS="java-17-openjdk-devel git jq unzip curl" && \
    yum install -y --nodocs ${YUM_PKGS} && \
    yum update -q -y && \
    mkdir -pv /openvsx-server

WORKDIR /openvsx-server
COPY --from=openvsx-server --chown=0:0 /home/openvsx/server /openvsx-server/
COPY ./openvsx-sync.json /openvsx-server/
COPY ./build/scripts/download_vsix.sh /tmp
RUN \
    /tmp/download_vsix.sh -b ${CURRENT_BRANCH} && \
    mv -v /tmp/vsix /openvsx-server

WORKDIR /
COPY ./build/dockerfiles/application.yaml /openvsx-server/config/
RUN \
    tar -czf openvsx-server.tar.gz /openvsx-server

### Open VSX Modules - Builder
FROM registry.access.redhat.com/ubi8/nodejs-18:latest as ovsx-lib-builder
USER 0

ENV \
    ovsx_version=0.8.3 \
    npm_config_cache=/tmp/opt/cache

# Install pre-requisites for ovsx (multi-arch support)
RUN \
    YUM_PKGS="libsecret" && \
    yum install -y --nodocs "${YUM_PKGS}" && \
    { \
    ARCH=$(uname -m); \
    if [[ ${ARCH} != "x86_64" ]]; then \
        FEDORA_REPO_FILE="/etc/yum.repos.d/fedora.repo"; \
        touch ${FEDORA_REPO_FILE}; \
        echo "[fedora]" >> ${FEDORA_REPO_FILE}; \
        echo "name=Fedora - Secondary \$releasever on \$basearch" >> ${FEDORA_REPO_FILE}; \
        echo "baseurl=https://rpmfind.net/linux/fedora-secondary/releases/38/Everything/\$basearch/os" >> ${FEDORA_REPO_FILE}; \
        echo "enabled=1" >> ${FEDORA_REPO_FILE}; \
        echo "gpgcheck=0" >> ${FEDORA_REPO_FILE}; \
        echo "skip_if_unavailable=True" >> ${FEDORA_REPO_FILE}; \
        YUM_PKGS="libsecret"; \
        yum install -y --nodocs "${YUM_PKGS}"; \
    fi; \ 
    }

USER 1001
RUN \
    npm install --location=global ovsx@${ovsx_version} --prefix /tmp/opt/ovsx --cache ${npm_config_cache} --no-audit --no-fund --no-update-notifier && \
    chmod -R g+rwX /tmp/opt/ovsx && \
    tar -czf ovsx.tar.gz /tmp/opt/ovsx

### Plugin Generator - Builder
FROM registry.access.redhat.com/ubi8/nodejs-18:latest as plugin-builder
USER 0

ENV \
    JOB_CONFIG="job-config.json"

COPY che-*.yaml ${JOB_CONFIG} /tmp/

WORKDIR /tmp
RUN \
    YUM_PKGS="jq" && \
    yum install -y --nodocs ${YUM_PKGS} && \
    npm install --prefix /usr/lib npm@latest --ignore-scripts --no-audit --no-fund --no-update-notifier && \
    REGISTRY_VERSION=$(jq -r '.Version' "${JOB_CONFIG}") && \
    REGISTRY_GENERATOR_VERSION=$(jq -r --arg REGISTRY_VERSION "${REGISTRY_VERSION}" '.Other["@eclipse-che/plugin-registry-generator"][$REGISTRY_VERSION]' "${JOB_CONFIG}") && \
    npm install @eclipse-che/plugin-registry-generator@$REGISTRY_GENERATOR_VERSION --ignore-scripts --no-audit --no-fund --no-update-notifier && \
    npx @eclipse-che/plugin-registry-generator@"${REGISTRY_GENERATOR_VERSION}" --root-folder:"$(pwd)" --output-folder:"$(pwd)/output" --embed-vsix:true --skip-digest-generation:true && \
    tar -czf resources.tar.gz ./output/v3/

### Plugin Registry Image
FROM registry.redhat.io/rhel8/postgresql-15:latest
USER 0
WORKDIR /

COPY ./LICENSE /licenses/
COPY ./build/dockerfiles/content_sets*.repo /etc/yum.repos.d/
COPY ./build/dockerfiles/rhel.install.sh /tmp
COPY --chown=0:0 --chmod=775 ./build/scripts/import_vsix.sh ./build/scripts/start_services.sh ./build/dockerfiles/entrypoint.sh /usr/local/bin/
COPY --chown=0:0 --chmod=664 ./build/dockerfiles/openvsx.conf /etc/httpd/conf.d/
COPY --chown=0:0 --chmod=664 README.md .htaccess /var/www/html/
COPY --chown=0:0 --chmod=775 ./build/scripts/*.sh che-*.yaml /build/

COPY --chown=0:0 v3/plugins/ ./build/dockerfiles/wazi.external_images.yaml /output/v3/plugins/
COPY --chown=0:0 v3/images/*.png /output/v3/images/

COPY --chown=0:0 --from=ovsx-server-builder /openvsx-server.tar.gz .
COPY --chown=0:0 --from=ovsx-lib-builder /opt/app-root/src/ovsx.tar.gz .
COPY --chown=0:0 --from=plugin-builder /tmp/resources.tar.gz .

RUN \
    tar -xzf openvsx-server.tar.gz && \
    tar -xzf ovsx.tar.gz && \
    tar -xzf resources.tar.gz && \
    mv -v /output /build && \
    rm -rvf /build/output/v3/che-editors.yaml && \
    /tmp/rhel.install.sh && \
    /build/wazi_plugin.sh /build/output/v3 && \
    /build/list_referenced_images.sh /build/output/v3 --use-generated-content > /build/output/v3/external_images.txt && \
    rm -rvf /build/output/v3/plugins/wazi.external_images.yaml && \
    cat /build/output/v3/external_images.txt && \
    mv /build/output/v3 /var/www/html/ && \
    cat /etc/passwd | sed s#root:x.*#root:x:\${USER_ID}:\${GROUP_ID}::\${HOME}:/bin/bash#g > /.passwd.template && \
    cat /etc/group  | sed s#root:x:0:#root:x:0:0,\${USER_ID}:#g > /.group.template && \
    for f in "/etc/passwd" "/etc/group" "/var/log/httpd" "/run/httpd" "/usr/local/bin/*.sh" "/var/www/html/v3" "openvsx-server"; do \
           chgrp -R 0 ${f} && \
           chmod -R g+rwX ${f}; \
    done && \
    localedef -f UTF-8 -i en_US en_US.UTF-8 && \
    usermod -a -G apache,root,postgres postgres && \
    sed -i /etc/httpd/conf/httpd.conf \
    -e "s,Listen 80,Listen 8080," \
    -e "s,logs/error_log,/dev/stderr," \
    -e "/<IfModule log_config_module>/a SetEnvIf User-Agent \"^kube-probe/\" dontlog" \
    -e 's,CustomLog "logs/access_log" combined,CustomLog /dev/stdout combined env=!dontlog,' \
    -e "s,AllowOverride None,AllowOverride All," && \
    chmod a+rwX /etc/httpd/conf /etc/httpd/conf.d /run/httpd /etc/httpd/logs/ && \
    rm -rf /tmp/rhel.install.sh /openvsx-server.tar.gz /ovsx.tar.gz /resources.tar.gz /build

STOPSIGNAL SIGWINCH

USER postgres
ARG DS_BRANCH=devspaces-3.10-rhel-8
ENV \
    LC_ALL=en_US.UTF-8 \
    LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    PGDATA=/var/lib/pgsql/15/data/database \
    PATH="/tmp/opt/ovsx/bin:$PATH" \
    JVM_ARGS="-DSPDXParser.OnlyUseLocalLicenses=true -Xmx2048m" \
    DS_BRANCH=${DS_BRANCH}

RUN \
    echo "======================" && \
    echo -n "node:  "; node --version && \
    echo "======================" \
    echo -n "ovsx:  "; /tmp/opt/ovsx/bin/ovsx --version && \
    echo "======================" && \
    chmod 777 /var/run/postgresql && \
    initdb && \
    /usr/local/bin/import_vsix.sh && \
    chmod -R 777 /tmp/file && \
    rm /var/lib/pgsql/15/data/database/postmaster.pid && \
    rm /var/run/postgresql/.s.PGSQL* && \
    rm /tmp/.s.PGSQL* && \
    rm /tmp/.lock && \
    chmod -R g+rwx /var/lib/pgsql/15 /var/lib/pgsql/data /var/lib/pgsql/backups && \
    chgrp -R 0 /var/lib/pgsql/15 /var/lib/pgsql/data /var/lib/pgsql/backups && \
    mv /var/lib/pgsql/15/data/database /var/lib/pgsql/15/data/old

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

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
