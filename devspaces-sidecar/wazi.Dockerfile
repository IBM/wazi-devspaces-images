###############################################################################
# Licensed Materials - Property of IBM.
# Copyright IBM Corporation 2023. All Rights Reserved.
# U.S. Government Users Restricted Rights - Use, duplication or disclosure
# restricted by GSA ADP Schedule Contract with IBM Corp.
#
# Contributors:
#  IBM Corporation - initial API and implementation
###############################################################################

FROM registry.redhat.io/devspaces/udi-rhel8@sha256:d18f22ef1aa2e5d1da4e3356ee1fc8fa59f795cdc3ab9d54c666054fbcfecd8f AS core

###########################################
###
###   Core Instruction Set
###
###########################################

ARG PRODUCT_VERSION="3.0.0"
USER 0

ENV \
    JAVA_VERSION="17" \
    SEMERU_VERSION="17.0.7.7_0.38.0-1"

COPY LICENSE PRODUCT_LICENSE /licenses/
COPY scripts/* *.zip /tmp/

### *** General *** ###
RUN \
    echo $'alias ll=\'ls -l\'\nalias la=\'ls -la\'\nalias ld=\'ls -lad */\'\n' >> /home/user/.bashrc && \
    ln -sf /home/user/.bashrc /home/user/.profile

### *** Java (Semeru) *** ###
RUN \
    ARCH="$(uname -m)" && \
    SEMERU_JDK="jdk-17.0.7%2B7_openj9-0.38.0" && \
    SEMERU_RPM="https://github.com/ibmruntimes/semeru${JAVA_VERSION}-binaries/releases/download/${SEMERU_JDK}/ibm-semeru-open-${JAVA_VERSION}-jdk-${SEMERU_VERSION}.${ARCH}.rpm" && \
    DNF_PKGS="yum python39-wheel iputils libatomic libzip-tools cargo rust" && \
    YUM_PKGS="${SEMERU_RPM}" && \
    dnf -y update  --noplugins --nodocs --best --allowerasing && \
    dnf -y clean all --enablerepo='*' && dnf -y clean packages && \
    dnf -y clean all && rm -rf /var/cache/yum && \
    dnf -y install --noplugins --nodocs ${DNF_PKGS} && \
    yum -y install --nodocs ${YUM_PKGS} && \
    find /home/user/.java/current -maxdepth 1 -type l -delete && \
    ln -s /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/* /home/user/.java/current && \
    update-alternatives --set javac /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/bin/javac && \
    update-alternatives --set java /usr/lib/jvm/ibm-semeru-open-${JAVA_VERSION}-jdk/bin/java

### *** For CVE Remediation

RUN \
    ARCH="$(uname -m)" && \
    export HELM_INSTALL_DIR="/usr/bin" && \
    export ODO_INSTALL_DIR="$(which odo)" && \
    curl -o- -skL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash && \
    curl -skL https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/odo/latest/odo-linux-${ARCH} -o ${ODO_INSTALL_DIR}

FROM core AS code

###########################################
###
###   Code Instruction Set
###
###########################################

ENV \
    ZOWE_CLI_PLUGINS_DIR="/usr/local/lib/node_modules" \
    ZOWE_CLI_VERSION="7.16.2" \
    RSE_API_VERSION="latest" \
    ANSIBLE_LINT_VERSION="6.17.0" \
    ANSIBLE_CRYPT_VERSION="40.0.1" \
    NPM_VERSION="9.*"

### **** Update pip and npm *** ###

RUN \
    python -m pip install --upgrade pip --no-cache-dir && \
    python -m pip install responses "urllib3==1.26.16" --no-cache-dir && \
    npm install --prefix /usr/lib npm@${NPM_VERSION} --ignore-scripts --no-audit --no-fund --no-update-notifier

### *** Ansible *** ###
RUN \
    export CRYPTOGRAPHY_DONT_BUILD_RUST=1 CARGO_NET_GIT_FETCH_WITH_CLI=true && \
    ANSIBLE_COLLECTIONS="/usr/share/ansible/collections" && \
    PIP_PKGS="cryptography==$ANSIBLE_CRYPT_VERSION ansible-lint==$ANSIBLE_LINT_VERSION" && \
    ANSIBLE_PKGS="ibm.ibm_zos_core ibm.ibm_zosmf ibm.ibm_zos_cics ibm.ibm_zos_ims ibm.ibm_zhmc ibm.ibm_zos_sysauto" && \
    mkdir -pv /tmp/wheels && \
    python -m pip wheel --wheel-dir=/tmp/wheels cryptography==${ANSIBLE_CRYPT_VERSION} && \
    python -m pip wheel --wheel-dir=/tmp/wheels ansible-lint==${ANSIBLE_LINT_VERSION} && \
    python -m pip install --no-cache-dir --no-index --find-links=/tmp/wheels ${PIP_PKGS} && \
    python -m pip install --no-cache-dir ansible && \
    ansible-galaxy collection install --no-cache -p ${ANSIBLE_COLLECTIONS} ${ANSIBLE_PKGS}

### *** Install Zowe CLI, Zapp Core, RSE API *** ###
RUN \
    --mount=type=secret,id=docker_secret,dst=/run/secrets/docker_secret source /run/secrets/docker_secret && \
    /tmp/wazi_sidecar.sh --npmrc "/home/user/.npmrc" "${NPM_URI}" "${NPM_REG}" "${NPM_USER}" "${NPM_KEY}" && \
    NPM_PKGS=("@zowe/cli@${ZOWE_CLI_VERSION}" "@ibm/rse-api-for-zowe-cli@${RSE_API_VERSION}") && \
    NODE_PATH=/usr/local/lib/node_modules && \
    for NPM_PKG in "${NPM_PKGS[@]}"; do \
        echo "Installing ${NPM_PKG} ..."; \
        npm install -g ${NPM_PKG} --ignore-scripts --no-audit --no-fund --no-update-notifier; \
    done && \
    npm list -g --depth=0 && \
    zowe plugins install  "$NODE_PATH/@ibm/rse-api-for-zowe-cli" && \
    zowe plugins list && \
    rm -rfv "/home/user/.npmrc"

FROM code AS analyze

###########################################
###
###   Analyze Instruction Set
###
###########################################

ENV \
    WA="${HOME}/wazianalyze" \
    WADATA="${HOME}/wazianalyze/data" \
    WADATA_TEMPLATES="${HOME}/wazianalyze/templates" \
    AZN_SSL_LAX="true" \
    PATH="${HOME}/wazianalyze/script${PATH:+:${PATH}}"

EXPOSE 5000/tcp
EXPOSE 8001/tcp
EXPOSE 4680/tcp

RUN \
    ARCH="$(uname -m)" && \
    ANALYZE_BINARIES="analyze_binaries" && \
    ANALYZE_WORKDIR="/tmp/${ANALYZE_BINARIES}" && \
    mkdir -pv "${ANALYZE_WORKDIR}" "${WA}" && \
    echo "Unzipping Wazi Analyze binaries...(please wait)" && \
    unzip -o "/tmp/${ANALYZE_BINARIES}.zip" "wazianalyze/${ARCH}/"* -d "${ANALYZE_WORKDIR}" > /dev/null && \
    mv -fv "${ANALYZE_WORKDIR}/wazianalyze/${ARCH}/"* "${WA}"

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
