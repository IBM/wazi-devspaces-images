# RHEL8 Image for IBM Developer for z/OS on Cloud

This folder contains the dockerfile of our previous default image that was based on RHEL8.

This image is based on [Red Hat OpenShift Dev Spaces - Universal Developer Image](https://catalog.redhat.com/software/containers/devspaces/udi-rhel8/622bce914a14c05796114be4) pulling from `registry.redhat.io/devspaces/udi-rhel8:latest`.

The dockerfile in this folder adds the following components to this image:

- [IBM Semeruc Java 17 LTS runtime](https://developer.ibm.com/languages/java/semeru-runtimes/downloads/)
- [Ansible](https://www.redhat.com/en/technologies/management/ansible)
- [Red Hat Ansible Certified Content for IBM Z](https://ibm.github.io/z_ansible_collections_doc/index.html)
- [Zowe CLI](https://www.npmjs.com/package/@zowe/cli)
- [IBM CICS Plug-in for Zowe CLI](https://www.npmjs.com/package/@zowe/cics-for-zowe-cli)
- [IBM RSE API Plug-in for Zowe CLI](https://www.npmjs.com/package/@ibm/rse-api-for-zowe-cli)
