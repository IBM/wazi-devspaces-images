# Default Image for IBM Developer for z/OS on Cloud

This folder contains the dockerfile for the default image provided for our samples that will be pull when starting a workspace from <https://github.com/IBM/zopeneditor-sample/tree/devfile>.

This image is based on [Red Hat OpenShift Dev Spaces - Universal Developer Image](https://catalog.redhat.com/software/containers/devspaces/udi-rhel9/673f8460bbf0c33aca0fe316) pulling from `registry.redhat.io/devspaces/udi-rhel9:latest`. This base image was release with Red Hat OpenShift Dev Spaces v3.20 which is now based on RHEL9.

The dockerfile in this folder adds the following components to this image:

- [IBM Semeruc Java 17 LTS runtime](https://developer.ibm.com/languages/java/semeru-runtimes/downloads/)
- [Zowe CLI](https://www.npmjs.com/package/@zowe/cli)
- [IBM CICS Plug-in for Zowe CLI](https://www.npmjs.com/package/@zowe/cics-for-zowe-cli)
- [IBM RSE API Plug-in for Zowe CLI](https://www.npmjs.com/package/@ibm/rse-api-for-zowe-cli)

Ansible is not included anymore as Red Hat now provides it own images with Dev Spaces that can be loaded as another sidecar iamge.
