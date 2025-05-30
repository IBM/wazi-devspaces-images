###############################################################################
# Licensed Materials - Property of IBM.
# Copyright IBM Corporation 2023, 2025. All Rights Reserved.
# U.S. Government Users Restricted Rights - Use, duplication or disclosure
# restricted by GSA ADP Schedule Contract with IBM Corp.
#
# Contributors:
#  IBM Corporation - initial API and implementation
###############################################################################

FROM registry.redhat.io/devspaces/udi-base-rhel9:latest AS core

###########################################
###
###   Core Instruction Set
###
###########################################

ARG PRODUCT_VERSION="5.3.0"
USER 0

ENV \
    JAVA_VERSION="17" \
    SEMERU_JDK="jdk-17.0.15%2B5_openj9-0.51.0-m2" \
    SEMERU_VERSION="17.0.15.5_0.51.0-1" \
    NODEJS_VERSION="20" \
    MAVEN_VERSION="3.8"

COPY LICENSE /licenses/
COPY *.sh /tmp/

### **** Prepare installation based on codeready-builder repos *** ###
RUN \
    dnf -y module enable nodejs:${NODEJS_VERSION} maven:${MAVEN_VERSION}

### *** Install and upgrade core packages *** ###
RUN \
    DNF_PKGS="yum iputils nodejs npm maven" && \
    dnf -y update  --noplugins --nodocs --nobest && \
    dnf -y clean all --enablerepo='*' && dnf -y clean packages && \
    dnf -y clean all && rm -rf /var/cache/yum && \
    dnf -y install --noplugins --nodocs ${DNF_PKGS}

### *** For CVE Remediation
RUN \
    dnf -y remove podman

### *** Install IBM Semeru Java *** ###
RUN \
    ARCH="$(uname -m)" && \
    SEMERU_RPM="https://github.com/ibmruntimes/semeru${JAVA_VERSION}-binaries/releases/download/${SEMERU_JDK}/ibm-semeru-open-${JAVA_VERSION}-jdk-${SEMERU_VERSION}.${ARCH}.rpm" && \
    YUM_PKGS="${SEMERU_RPM}" && \
    yum -y install --nodocs ${YUM_PKGS} && \
    mkdir -p /home/user/.java/current && \
    rm -f /usr/bin/java && \
    find /home/user/.java/current -maxdepth 1 -type l -delete && \
    ln -s /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/* /home/user/.java/current && \
    update-alternatives --set javac /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/bin/javac && \
    update-alternatives --set java /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/bin/java

FROM core AS code

### *** Install Zowe CLI, RSE API Plugin, CICS Plugin *** ###
ENV \
    ZOWE_CLI_PLUGINS_DIR="/usr/local/lib/node_modules" \
    ZOWE_CLI_VERSION="zowe-v3-lts" \
    RSE_API_VERSION="latest" \
    CICS_CLI_VERSION="latest" \
    JAVA_HOME="/home/user/.java/current" \
    M2_HOME="/usr/share/maven" \
    PATH="/home/user/.java/current/bin:/usr/share/maven/bin:/home/user/node_modules/.bin/:${PATH:-/bin:/usr/bin}"

RUN \
    NPM_PKGS=("@zowe/cli@${ZOWE_CLI_VERSION}" "@zowe/cics-for-zowe-cli@${CICS_CLI_VERSION}" "@ibm/rse-api-for-zowe-cli@${RSE_API_VERSION}") && \
    NODE_PATH=/usr/lib/node_modules && \
    for NPM_PKG in "${NPM_PKGS[@]}"; do \
    echo "Installing ${NPM_PKG} ..."; \
    npm install -g ${NPM_PKG} --prefix /usr --ignore-scripts --no-audit --no-fund --no-update-notifier; \
    done && \
    npm list -g --depth=0 && \
    zowe plugins install  "$NODE_PATH/@ibm/rse-api-for-zowe-cli" && \
    zowe plugins install  "$NODE_PATH/@zowe/cics-for-zowe-cli" && \
    zowe plugins list

###########################################

RUN \
    echo $'alias ll=\'ls -l\'\nalias la=\'ls -la\'\nalias ld=\'ls -lad */\'\n' >> /home/user/.bashrc && \
    ln -sf /home/user/.bashrc /home/user/.profile

### *** Permissions / Clean-up *** ###
RUN \
    /tmp/wazi_sidecar.sh --permissions && \
    /tmp/wazi_sidecar.sh --cleanup && \
    rm -rf /tmp/*

USER 10001
WORKDIR /projects

ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["tail", "-f", "/dev/null"]

ENV \
    SUMMARY="IBM Developer for z/OS on Red Hat OpenShift Dev Spaces" \
    DESCRIPTION="Extended developer image for using enterprise appliation development with IBM Developer for z/OS on VS Code tools." \
    PRODNAME="IBM Developer for z/OS on Red Hat OpenShift Dev Spaces" \
    COMPNAME="IDzEE" \
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
    productName="$PRODNAME" \
    productID="$PRODUCT_ID" \
    productMetric="$PRODUCT_METRIC" \
    productChargedContainers="$PRODUCT_CHARGED_CONTAINERS" \
    productCloudpakRatio="$PRODUCT_CLOUDPAK_RATIO"
