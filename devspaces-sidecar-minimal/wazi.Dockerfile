###############################################################################
# Licensed Materials - Property of IBM.
# Copyright IBM Corporation 2024. All Rights Reserved.
# U.S. Government Users Restricted Rights - Use, duplication or disclosure
# restricted by GSA ADP Schedule Contract with IBM Corp.
#
# Contributors:
#  IBM Corporation - initial API and implementation
###############################################################################

FROM registry.redhat.io/ubi8/nodejs-20:latest AS core

###########################################
###
###   Core Instruction Set
###
###########################################

ARG PRODUCT_VERSION="5.1.0"
USER 0

COPY LICENSE /licenses/
COPY entrypoint.sh /entrypoint.sh
COPY *.sh /tmp/

### *** Creating user and general updates *** ###

ENV \
    HOME=/home/user
RUN \
    useradd -u 1000 -G wheel,root -d /home/user --shell /bin/bash -m user && \
    touch /etc/profile.d/udi_prompt.sh && \
    echo "export PS1='\W \`git branch --show-current 2>/dev/null | sed -r -e \"s@^(.+)@\(\1\) @\"\`$ '" >> /etc/profile.d/udi_prompt.sh && \
    # Change permissions to let any arbitrary user
    mkdir -p /projects && \
    for f in "${HOME}" "/etc/passwd" "/etc/group" "/projects"; do \
        echo "Changing permissions on ${f}" && chgrp -R 0 ${f} && \
        chmod -R g+rwX ${f}; \
    done && \
    # Generate passwd.template
    cat /etc/passwd | \
    sed s#user:x.*#user:x:\${USER_ID}:\${GROUP_ID}::\${HOME}:/bin/bash#g \
    > ${HOME}/passwd.template && \
    cat /etc/group | \
    sed s#root:x:0:#root:x:0:0,\${USER_ID}:#g \
    > ${HOME}/group.template && \
    echo $'alias ll=\'ls -l\'\nalias la=\'ls -la\'\nalias ld=\'ls -lad */\'\n' >> /home/user/.bashrc && \
    ln -sf /home/user/.bashrc /home/user/.profile && \
    DNF_PKGS="yum iputils libzip-tools curl findutils git git-lfs vim" && \
    npm install npm@latest --ignore-scripts --no-audit --no-fund --no-update-notifier \
    dnf -y update  --noplugins --nodocs --nobest && \
    dnf -y clean all --enablerepo='*' && dnf -y clean packages && \
    dnf -y clean all && rm -rf /var/cache/yum && \
    dnf -y install --noplugins --nodocs ${DNF_PKGS}

### *** For CVE Remediation

RUN \
    dnf -y update python3-libs --noplugins --nodocs && \
    dnf -y update platform-python --noplugins --nodocs

###########################################
###
###   Install Zowe CLI and RSE API plugin
###
###########################################

ENV \
    ZOWE_CLI_VERSION="zowe-v3-lts" \
    RSE_API_VERSION="latest" \
    ZOWE_CLI_HOME=${HOME}

RUN \
    --mount=type=secret,id=docker_secret,dst=/run/secrets/docker_secret source /run/secrets/docker_secret && \
    if [[ -n "${NPM_REG}" ]] ; then \
      echo "Fetching RSE API Plugin from ${NPM_REG}" ; \
      /tmp/wazi_sidecar.sh --npmrc "/home/user/.npmrc" "${NPM_URI}" "${NPM_REG}" "${NPM_USER}" "${NPM_KEY}" ; \
    fi && \
    NPM_PKGS=("@zowe/cli@${ZOWE_CLI_VERSION}" "@ibm/rse-api-for-zowe-cli@${RSE_API_VERSION}") && \
    for NPM_PKG in "${NPM_PKGS[@]}"; do \
        echo "Installing ${NPM_PKG} ..."; \
        npm install -g ${NPM_PKG} --ignore-scripts --no-audit --no-fund --no-update-notifier; \
    done && \
    npm list -g --depth=0 && \
    zowe plugins install  "@ibm/rse-api-for-zowe-cli" && \
    zowe plugins list && \
    rm -rfv "/home/user/.npmrc"

###########################################
###
###   Permissions / Clean-up / Startup
###
###########################################

RUN \
    /tmp/wazi_sidecar.sh --permissions && \
    /tmp/wazi_sidecar.sh --cleanup && \
    rm -rf /tmp/*

USER 10001
ENV HOME=/home/user
WORKDIR /projects

ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["tail", "-f", "/dev/null"]

ENV \
    SUMMARY="IBM Developer for z/OS on Red Hat OpenShift Dev Spaces - Zowe Terminal" \
    DESCRIPTION="Developer image for IBM Developer for z/OS on Red Hat OpenShift Dev Spaces that provides Zowe CLI with the IBM RSE API Plugin. Use in combination with the UDI." \
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
      cloudpakMetric="$CLOUDPAK_METRIC" \
      productName="$PRODNAME" \
      productID="$PRODUCT_ID" \
      productMetric="$PRODUCT_METRIC" \
      productChargedContainers="$PRODUCT_CHARGED_CONTAINERS" \
      productCloudpakRatio="$PRODUCT_CLOUDPAK_RATIO"
