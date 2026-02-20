[![Release](https://img.shields.io/github/release/IBM/wazi-devspaces-images.svg)](../../releases/latest)
[![License](https://img.shields.io/github/license/IBM/wazi-devspaces-images)](./LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-blue?color=1f618d)](https://ibm.github.io/zopeneditor-about/Docs/cloud_overview.html)

## What's inside?

This repository contains the source sources for container images that can be used with IBM Developer for z/OS on VS Code extension in Cloud-based development environments such as Red Hat OpenShift Dev Spaces. You can use these file as a starting point for creating your own customized variants of these images that can also contain your specific prerequisites and development tools that were not included in the IBM image.

To learn more about running Z Open Editor, Zowe ClI and other tools in Red Hat Dev Spaces review our [documentation](https://ibm.github.io/zopeneditor-about/Docs/cloud_overview.html), in particular the section [Create and use custom images](https://ibm.github.io/zopeneditor-about/Docs/cloud_custom_images.html).

## Dockerfile examples

Currently, there are several examples in separate folders:

### Dockerfiles based on RHEL9 base images (Dev Spaces 3.20 or newer)

- devspaces-sidecar-udi-rhel9: Contains a dockerfile that represents our Dev Spaces solution that takes the Red Hat Universal Developer Image based on RHEL9 as the base and adds our tools to it. At the moment we are adding Zowe CLI with the CICS plugin and our RSE API Plugin.
- devspaces-sidecar-udibase-rhel9: Similar to the image above, but instead of being based on the full Red Hat UDI (Universal Developer Image) it based on the UDI-Base image, which does not have any languages and development tools installed, so it is half the size. We added Nodejs for Zowe CLI and IBM Semeru Java for the Z Open Editor to this image. This is the dockerfile used for the image currently available as `icr.io/wazi-code/ibm-wazi-for-devspaces-sidecar:latest`.

### Dockerfiles based on the now deprecated RHEL8 base images (Dev Spaces 3.19 or older)

- devspaces-sidecar-rhel8: Contains a dockerfile that represents another Dev Spaces solution that takes the Red Hat Universal Developer Image based on RHEL8 as the base and adds our tools to it. At the moment we are adding Zowe CLI with our RSE API Plugins as well as Ansible, Ansible development tools, as well as our Red Hat Ansible Certified Content for IBM Z collections so that you can write playbooks for z/OS right out of the box.
- devspaces-sidecar-ubi8: Contains a dockerfile that follows and alternative approach of providing a separate image that only contains tools not available in Red Hat Universal Developer Image. It only contains nodejs, Zowe CLI, and the RSE API Plugin. The idea is to use it in combination with other images, such as the UDI and Red Hat's image for Ansible development. So instead of having one big image, you would create a devfile that loads several smaller images.

## How to build the images

The images build without many prerequisites. However, as they are all based on Red Hat base images you need a Red Hat account that is entitled to pull these images.

Once you logged in with `docker login registry.redhat.io` you can build the image simply by running a command such as such using the directory name of the variant that you want to build

```bash
docker build -f devspaces-sidecar-udibase-rhel9/wazi.Dockerfile -t idzee-devspaces-sidecar:6.4.0 ./devspaces-sidecar-udibase-rhel9
```

Then you can push the image to your private image registry and use them from your devfiles as described in the [OpenShift Dev Spaces documentation](https://docs.redhat.com/en/documentation/red_hat_openshift_dev_spaces/3.20/html/user_guide/getting-started-with-devspaces). We recommend that you run security scans on these images yourself to assess if you can live with any open security issues that are found. Red Hat and IBM are committed to keeping images up-to-date, but as new issues are discovered every day you want to scan and evaluate the risks as well.

To build a multi-architecture image that can be used with s390x as well as x86-based OpenShift clusters you can use docker's buildx commands such as

```bash
cd devspaces-sidecar-udibase-rhel9
docker buildx create --use --name devspaces-builder
docker buildx build --platform linux/amd64,linux/s390x --tag idzee-devspaces:6.4.0 --builder devspaces-builder -f wazi.Dockerfile .
```

## Images available

In addition to the source files for building images, we also provide one image in the IBM Cloud image registry that you can use to try out Z Open Editor, Zowe Explorer, and Zowe CLI from any OpenShift cluster that has access to the internet. See our [tutorial](https://ibm.github.io/zopeneditor-about/Docs/cloud_developer_sandbox.html) for using the Red Hat Developer Sandbox as one way of doing a quick trial. We publish frequent updates to this image to include the latest security updates published by Red Hat and to update our own components such as IBM RSE API Plugin for Zowe CLI that are included on this image.

The most recent image is available as

```txt
stg.icr.io/ibm/wazi-code/ibm-wazi-for-devspaces-sidecar:latest
```

A simple devfile to use it could look like this:

```yaml
schemaVersion: 2.3.0
metadata:
  name: idzee-devspaces

components:
  - name: zowe
    volume:
      size: 100Mi
  - name: idzee-terminal
    container:
      image: stg.icr.io/ibm/wazi-code/ibm-wazi-for-devspaces-sidecar:latest
      memoryLimit: 3072Mi
      mountSources: true
      volumeMounts:
        - name: zowe
          path: /home/user/.zowe
```

Instead of `latest` you can use the version number such as `6.4.0` or the SHA listed in the table below. Here is a list of released images.

| <sub>Registry and Image</sub> | <sub>Version</sub> |
| :--- | --- |
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-sidecar@sha256:TBD</sub>|<sub><b>6.4.0</b></sub>|
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-sidecar@sha256:2ef78deed87dd21d1f39cfd3c52e035b9b2d74c3a4c2af6a6b615eb3b20b374a</sub>|<sub><b>5.3.0</b></sub>|
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-sidecar@sha256:925aee1f34ee72c65cfa1fe1b03b71dd2872c71aea24a1c82e86c61fa772eeb5</sub>|<sub><b>5.1.0</b></sub>|

## Feedback

We would love to hear feedback from you about IBM Wazi for Dev Spaces.
File an issue or provide feedback here: [IBM Wazi for Dev Spaces Issues](https://github.com/IBM/wazi-devspaces-images/issues)
