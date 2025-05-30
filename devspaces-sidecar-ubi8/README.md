# Minimal Image for IBM Developer for z/OS on Cloud

This folder contains the dockerfile that shows how you can build a minimal image for IBM Developer for z/OS on Cloud. It is not based on the very large Red Hat OpenShift Dev Spaces - Universal Developer Image, but rather based on the [Red Hat Universal Base Image 8](https://catalog.redhat.com/software/containers/ubi8/5c647760bed8bd28d0e38f9f) pulling from `registry.redhat.io/ubi8/nodejs-20`.

The dockerfile in this folder adds the following components to this image:

- [Zowe CLI](https://www.npmjs.com/package/@zowe/cli)
- [IBM RSE API Plug-in for Zowe CLI](https://www.npmjs.com/package/@ibm/rse-api-for-zowe-cli)
