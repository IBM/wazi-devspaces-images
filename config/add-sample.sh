#!/bin/bash

###############################################################################
# Licensed Materials - Property of IBM.
# Copyright IBM Corporation 2025. All Rights Reserved.
# U.S. Government Users Restricted Rights - Use, duplication or disclosure
# restricted by GSA ADP Schedule Contract with IBM Corp.
###############################################################################

# See <https://docs.redhat.com/en/documentation/red_hat_openshift_dev_spaces/3.19/html/administration_guide/configuring-devspaces#configuring-getting-started-samples> for details
# Check where Dev Spaces is deployed on your cluster: openshift-operators or openshift-devspaces
# and adjust the oc command parameter.

sh ./dashboard-sample.sh
oc project openshift-operators
oc create configmap getting-started-samples --from-file=wazi-samples.json -n openshift-operators
oc label configmap getting-started-samples app.kubernetes.io/part-of=che.eclipse.org app.kubernetes.io/component=getting-started-samples -n openshift-operators
