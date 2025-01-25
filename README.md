[![Build](https://app.travis-ci.com/IBM/wazi-devspaces-images.svg?branch=main)](https://app.travis-ci.com/IBM/wazi-devspaces-images)
[![Release](https://img.shields.io/github/release/IBM/wazi-devspaces-images.svg)](../../releases/latest)
[![License](https://img.shields.io/github/license/IBM/wazi-devspaces-images)](./LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-blue?color=1f618d)](https://ibm.github.io/zopeneditor-about/Docs/cloud_overview.html)

## What's inside?

This repository contains the source sources for container images that can be used with IBM Developer for z/OS on VS Code VS Code extension in Cloud-based development environments such as Red Hat OpenShift Dev Spaces. You can use these file as a starting point to create your own customized variants of these images that can also contain your specific prerequistes and development tools that were not included in the IBM image.

To learn more about running Z Open Editor, Zowe ClI and other tools in Red Hat Dev Spaces review our [documentation](https://ibm.github.io/zopeneditor-about/Docs/cloud_overview.html).

## Dockerfile examples

This repository provides examples for building container images that can be used with OpenShift Dev Spaces and the IBM Developer for z/OS on VS Code development tools. You can build these dockerfiles and use them as is or you can modify them to add or replace tools that you want to use in your projects. Currently, there are two examples in two separate folders:

- devspaces-sidecar: Contains a dockerfile that represents our classic Dev Spaces solution that takes the Red Hat Universal Developer Image as the base and adds our tools to it. At the moment we are adding Zowe CLI with our RSE API Plugins as well as Ansible, Ansible development tools, as well as our Red Hat Ansible Certified Content for IBM Z collections so that you can write playbooks for z/OS right out of the box. So you have one big image that loads all the tools you need for development.
- devspaces-sidecar-minimal: Contains a dockerfile that follows and alternative approch of providing a separate image that only contains tools not available in Red Hat Universal Developer Image. It only contains nodejs, Zowe CLI, and the RSE API Plugin. The idea is to use it in combination with other images, such as the UDI and Red Hat's image for Ansible development. So indead of having one big image, you would create a devfile that loads several smaller images.

Both approaches have pros and cons. Find more details about these and how to use such images in our [user documentation](https://ibm.github.io/zopeneditor-about/Docs/cloud_custom_images.md).

## How to build the images

The images build without many prerequites. However, as they are all based on [Red Hat Universal Developer image](https://catalog.redhat.com/software/containers/devspaces/udi-rhel8/622bce914a14c05796114be4) you need a Red Hat account that is entitled to pull these images.

Once you logged in with `docker login registry.redhat.io` you can build the image simply by running

```bash
docker build -f devspaces-sidecar/wazi.Dockerfile -t idzee-devspaces-sidecar:5.0.0 ./devspaces-sidecar
```

Then you can push the image to your private image registry and use them from your devfiles as described in the documentation. We recommend that you run security scans on these images yourself to assess if you can live with any open security isues that are found. Red Hat and IBM are committed to keeping images up-to-date, but as new issues are discovered every day you want to scan and evaluate the risks as well.

## Images available

In addition to the source files for the image we also provide the images in the IBM Cloud image registry. Whenever we make an update to the image we will republish it and update the table below. If you are using our image then check back frequently to get the latest updates and security fixes.

| <sub>Registry and Image</sub> | <sub>Version</sub> |
| :--- | --- |
||<sub><b>5.x.</b></sub>|
|<sub>icr.io/cpopen/ibm-wazi-for-devspaces-sidecar@TBD</sub>||
||<sub><b>4.0.0</b><br>(deprecated)</sub>|
|<sub>icr.io/cpopen/ibm-wazi-for-devspaces-catalog@sha256:67657651a8e608cb8d37c1c157fe46e141dd61d4166a14c467354f3a36f46b8a</sub>||
|<sub>icr.io/cpopen/ibm-wazi-for-devspaces-operator@sha256:259ef8610f517b343228f244b5c2977639c95c7022f8047694a8a28e35da5909</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-dashboard@sha256:c5d54c928658fdf40e6c75f14e948f2d02fba2dd2188faa4f071bc903db2cab2</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-devfile@sha256:5146d93a261058fa8c9c3ce1632b4f4a5d4ed77490864ce780c72004d0193e02</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-plugin@sha256:84041eb62f881d6d49fdd99963687370fdb7df50aacf19a7c4c97487f08f6e3c</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-sidecar@sha256:ba14c3e89595ca93d189869d2d118739c2d435f5d502ff46c39130b86968eb32</sub>||

## Feedback

We would love to hear feedback from you about IBM Wazi for Dev Spaces.
File an issue or provide feedback here: [IBM Wazi for Dev Spaces Issues](https://github.com/IBM/wazi-devspaces-images/issues)
