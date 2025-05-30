###############################################################################
# Licensed Materials - Property of IBM.
# Copyright IBM Corporation 2023, 2025. All Rights Reserved.
# U.S. Government Users Restricted Rights - Use, duplication or disclosure
# restricted by GSA ADP Schedule Contract with IBM Corp.
#
# Contributors:
#  IBM Corporation - initial API and implementation
###############################################################################

FROM registry.redhat.io/devspaces/udi-rhel9:latest AS core

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
    SEMERU_VERSION="17.0.15.5_0.51.0-1"

COPY LICENSE /licenses/
COPY *.sh /tmp/

### *** General *** ###
RUN \
    echo $'alias ll=\'ls -l\'\nalias la=\'ls -la\'\nalias ld=\'ls -lad */\'\n' >> /home/user/.bashrc && \
    ln -sf /home/user/.bashrc /home/user/.profile && \
    DNF_PKGS="yum iputils" && \
    dnf -y update  --noplugins --nodocs --nobest && \
    dnf -y clean all --enablerepo='*' && dnf -y clean packages && \
    dnf -y clean all && rm -rf /var/cache/yum && \
    dnf -y install --noplugins --nodocs ${DNF_PKGS}

### *** Java (Semeru) *** ###
RUN \
    ARCH="$(uname -m)" && \
    SEMERU_RPM="https://github.com/ibmruntimes/semeru${JAVA_VERSION}-binaries/releases/download/${SEMERU_JDK}/ibm-semeru-open-${JAVA_VERSION}-jdk-${SEMERU_VERSION}.${ARCH}.rpm" && \
    YUM_PKGS="${SEMERU_RPM}" && \
    yum -y install --nodocs ${YUM_PKGS} && \
    find /home/user/.java/current -maxdepth 1 -type l -delete && \
    ln -s /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/* /home/user/.java/current && \
    update-alternatives --set javac /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/bin/javac && \
    update-alternatives --set java /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/bin/java

### *** For CVE Remediation

RUN \
    ARCH="$(uname -m)" && \
    YUM_PKGS="git-lfs" && \
    export HELM_INSTALL_DIR="/usr/bin" && \
    export ODO_INSTALL_DIR="$(which odo)" && \
    curl -o- -skL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash && \
    curl -o- -skL https://packagecloud.io/install/repositories/github/git-lfs/script.rpm.sh | bash && \
    curl -skL https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/odo/latest/odo-linux-${ARCH} -o ${ODO_INSTALL_DIR} && \
    yum -y install --nodocs ${YUM_PKGS}

FROM core AS code

###########################################
###
###   Code Instruction Set
###
###########################################

ENV \
    ZOWE_CLI_PLUGINS_DIR="/usr/local/lib/node_modules" \
    ZOWE_CLI_VERSION="zowe-v3-lts" \
    RSE_API_VERSION="latest" \
    CICS_CLI_VERSION="latest"

### *** Install Zowe CLI, RSE API *** ###
# for building from an inhouse npm registry with Zowe add
# --mount=type=secret,id=docker_secret,dst=/run/secrets/docker_secret source /run/secrets/docker_secret && \
# if [[ -n "${NPM_REG}" ]] ; then \
#   /tmp/wazi_sidecar.sh --npmrc "/home/user/.npmrc" "${NPM_URI}" "${NPM_REG}" "${NPM_USER}" "${NPM_KEY}" ; \
# fi && \
# then build with
# docker build -f devspaces-sidecar/wazi.Dockerfile -t idzee-devspaces-sidecar:5.0.0 --secret id=docker_secret,src=.docker_secret ./devspaces-sidecar
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
# rm -rfv "/home/user/.npmrc"

###########################################

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
