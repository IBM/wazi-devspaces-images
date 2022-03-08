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

FROM scratch

LABEL operators.operatorframework.io.bundle.mediatype.v1=registry+v1
LABEL operators.operatorframework.io.bundle.manifests.v1=manifests/
LABEL operators.operatorframework.io.bundle.metadata.v1=metadata/
LABEL operators.operatorframework.io.bundle.package.v1=ibm-developer-for-zos-enterprise-edition
LABEL operators.operatorframework.io.bundle.channels.v1=v2.0
LABEL operators.operatorframework.io.bundle.channel.default.v1=v2.0

COPY manifests /manifests/
COPY metadata /metadata/
