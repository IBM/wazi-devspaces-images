# Image for IBM Developer for z/OS on Cloud with minimal tools

This folder contains a dockerfile for an image that can be used with our sample files when starting a workspace from <https://github.com/IBM/zopeneditor-sample/tree/devfile>.

This image is based on [Red Hat OpenShift Dev Spaces - Universal Developer Base Image](https://catalog.redhat.com/software/containers/devspaces/udi-base-rhel9/67e45c79c28cc5d1205c529c) pulling from `registry.redhat.io/devspaces/udi-base-rhel9:latest`. This base image was released with Red Hat OpenShift Dev Spaces v3.20 which is based on RHEL9. It is a simpler version than the full UDI image that does not yet have any development tools installed. In this dockerfile we are adding only the minimal set of tools, which is nodejs for Zowe CLI and Java for Z Open Editor.

The dockerfile in this folder adds the following components to this image:

- [IBM Semeruc Java 17 LTS runtime](https://developer.ibm.com/languages/java/semeru-runtimes/downloads/)
- [Zowe CLI](https://www.npmjs.com/package/@zowe/cli)
- [IBM CICS Plug-in for Zowe CLI](https://www.npmjs.com/package/@zowe/cics-for-zowe-cli)
- [IBM RSE API Plug-in for Zowe CLI](https://www.npmjs.com/package/@ibm/rse-api-for-zowe-cli)

Ansible is not included anymore as Red Hat now provides it own images with Dev Spaces that can be loaded as another sidecar image.
