[![Build](https://app.travis-ci.com/IBM/wazi-devspaces-images.svg?branch=main)](https://app.travis-ci.com/IBM/wazi-devspaces-images)
[![Release](https://img.shields.io/github/release/IBM/wazi-devspaces-images.svg)](../../releases/latest)
[![License](https://img.shields.io/github/license/IBM/wazi-devspaces-images)](./LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-blue?color=1f618d)](https://ibm.biz/wazi-crw-doc)

## What's inside?

This repository contains the registry, operator, bundle, and index images for IBM Wazi for Dev Spaces.

* [IBM Wazi for Dev Spaces - Operator](/codeready-workspaces-operator)
  * This repository is based off of the upstream [Red Hat Dev Spaces Images](https://github.com/redhat-developer/devspaces-images) repository, where the generated output is located in the respective branch for each version.
* [IBM Wazi for Dev Spaces - Devfile](/codeready-workspaces-devfileregistry) / [IBM Wazi for Dev Spaces - Plug-in](/codeready-workspaces-pluginregistry)
  * You can use the custom devfile and plug-in registry images to tailor your development environment. After you rebuild the registry images with modifications, and push the images into a container registry, you can deploy these images when creating an instance of IBM Wazi for Dev Spaces or apply the images to an existing instance.
  * _Stacks_ or workspace definitions are pre-configured CodeReady Workspaces with a dedicated set of tools that cover different developer personas. For example, you can pre-configure a workbench for testers with only the tools needed for their purposes. You can add new stacks by customizing the IBM Wazi for Dev Spaces devfile and plug-in registry.


## IBM Wazi for Dev Spaces

IBM Wazi for Dev Spaces provides a modern experience for mainframe software developers working with z/OS applications in the cloud. Powered by the open source projects Zowe and Red Hat OpenShift Dev Spaces, IBM Wazi for Dev Spaces offers an easy, streamlined onboarding process to provide mainframe developers with the tools they need. Using container technology, IBM Wazi for Dev Spaces brings the necessary tools to the task at hand. IBM Wazi for Dev Spaces is a component shipped with IBM Developer for z/OS Enterprise Edition (IDzEE).

For more benefits of IBM Wazi for Dev Spaces, see the [IBM Wazi product page](https://www.ibm.com/products/wazi-developer) or [IDzEE product page](https://www.ibm.com/products/developer-for-zos).

## Details

IBM Wazi for Dev Spaces provides a custom stack with the all-in-one mainframe development package that enables mainframe developers to:

- Use a modern mainframe editor with rich language support for COBOL, JCL, Assembler (HLASM), REXX, and PL/I, which provides language-specific features such as syntax highlighting, outline view, declaration hovering, code completion, snippets, a preview of copybooks, copybook navigation, and basic refactoring using IBM Z Open Editor, a component of IBM Wazi for VS Code
- Integrate with any flavor of Git source code management (SCM)
- Perform a user build with IBM Dependency Based Build for any flavor of Git
- Work with z/OS resources such as MVS, UNIX files, and JES jobs
- Connect to the Z host with z/OSMF or IBM Remote System Explorer (RSE) API, using Zowe Explorer plus IBM Z Open Editor for a graphical user interface and Zowe CLI plus the RSE API plug-in for Zowe CLI for command line access
- Debug COBOL and PL/I applications using IBM Z Open Debug
- Use a mainframe development package with a custom plug-in and devfile registry support from the [IBM Wazi Code stack](https://github.com/IBM/wazi-devspaces-images)
- Leverage the full language support for Ansible and Red Hat Ansible Certified Content for IBM Z to author and execute playbooks to configure z/OS for development. Automate the deployment and configuration of tools and dependencies such as build, test and debug, to quickly build, deploy and run your applications.

## Documentation

For details of the features for IBM Wazi for Dev Spaces, see its [official documentation](https://ibm.biz/wazi-crw-doc).

### Image Information

* `catalog` Image - This is an index image used for the operator hub. 
* `operator` Image - This is the operator image inherited by Red Hat. 
* `dev-file` Image - This is a customized dev-file registry image. 
* `plugin` Image - This is a customized plugin registry image. 
* `codeready` Image - This is a Wazi sidecar image used to provide supporting resources. 

| <sub>Registry and Image</sub> | <sub>Version</sub> |
| :--- | --- |
||<sub><b>2.1.0</b></sub>|
|<sub>icr.io/cpopen/ibm-wazi-for-devspaces-catalog@sha256:93ad2466cb4315310be5766128cef744980ea0c7158ff0dc7356eb51e56e66a4</sub>||
|<sub>icr.io/cpopen/ibm-wazi-for-devspaces-operator@sha256:7715e865bb9ba22783e7743922307dba260b3032106cf361fc87918e0ca2998c</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-dev-file@sha256:894f23c01de10aca78c8acf33a80b76688441ce73060f1caaddebd6542fa2969</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-plugin@sha256:5def707aadbb9d0e048701b34b766bf421c85d0e90a51bc2e7d3b9443b645779</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-codeready@sha256:3f10050adbd00eb2fa8cb9ba58503b698079190cd56c5b2d735e409fa1ce8285</sub>||
||<sub><b>2.0.0</b></sub>|
|<sub>icr.io/cpopen/ibm-wazi-developer-for-workspaces-catalog@sha256:02f1cb4ea404d20cf80ad08d7ce6ab93c2fc7906fe53679ee7e37ad2719f2864</sub> ||
|<sub>icr.io/cpopen/ibm-wazi-developer-for-workspaces-operator@sha256:8d7bef2f925b34b3f2abdb1befadbfcf734a40edbc2342163a7c265d65217098</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-dev-file@sha256:47c8bc6edb671247be92fb5239e6233afee370168a324aeacce7b0d99949adfa</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-plugin@sha256:0f814ea0708d2af2c0b47b7ebbcaf5b3406d64056cab2fb2e1290191a8cdecbc</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-codeready@sha256:f78dfe340e0407f277cd6b8b0cb8803ddd3320a2753b74dc653d333004a12129</sub> ||
||<sub><b>1.4.1</b></sub>|
|<sub>icr.io/cpopen/ibm-wazi-developer-for-workspaces-catalog@sha256:dd6ed0f026ba1ef9abea6bcdf8b8f812ce6f79356ef085149e750b1ff0dd71fe</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-dev-file@sha256:70929334aa9bc85916dfcdd2a82572aad6ffd730f575a79888a8824f51199be9</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-plugin@sha256:96ab5a1363fccc4938c4e23ca61971db47a8602cebb5b1a0d9b49644e4479139</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-codeready@sha256:6a032f2ebdf19f6b491832a8d897a97351636fef33cd29fc57f1bf36b8e6cba3</sub> ||
||<sub><b>1.4.0</b></sub>|
|<sub>icr.io/cpopen/ibm-wazi-developer-for-workspaces-catalog@sha256:d7f26aa2f05fdc9812d5347f1f36eb1d5d3d77467ff34982e4d08fe89375b036</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-dev-file@sha256:8cf9c39cb8a45cd14ff4f7df09bc9ceb39768befe87761f7140e9c9d259ad397</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-plugin@sha256:6baada84cccbe4b1abcf560e317998efb12e4db238557d85a1ec1c63e65a396f</sub> ||
|<sub>icr.io/wazi-code/ibm-wazi-developer-for-workspaces-codeready@sha256:28beb27a2c9d6441cc8e16d19192e004d14701f31e0b9ffdf21cd9b5f0e8f0aa</sub> ||
||<sub><b>1.2.5</b></sub>|
|<sub>icr.io/cpopen/wazi-code-operator-catalog@sha256:42646eedcbf2fa09465023d7adda0220efa7daf3f7b10d8b17575617f13022aa</sub> ||
|<sub>icr.io/wazi-code/wazi-code-dev-file@sha256:fe366163f7e678ccac7c67b4ee86cd1d5be34b5287064c22447f5c58c1187333</sub> ||
|<sub>icr.io/wazi-code/wazi-code-plugin@sha256:6560e097e63a55e9757526253098a4ed1523c71218f4231eb5f3cd342f9afa8a</sub> ||
|<sub>icr.io/wazi-code/wazi-code-codeready@sha256:9894c02388f3b624730fe67ec7b376ffe498305b005d68423b0ef74db618c4a4</sub> ||
||<sub><b>1.2.0</b></sub>|
|<sub>icr.io/cpopen/wazi-code-operator-catalog@sha256:edfa9da664ac00712e8f4d10bcf06fd9d5c748f62521a4dd34bfc5fb07e9bf27</sub> ||
|<sub>icr.io/wazi-code/wazi-code-dev-file@sha256:9c7da5733618c09cad1b2dd613f3b8fa6bd0aacf8cc44bde28875b1f8bb541d3</sub> ||
|<sub>icr.io/wazi-code/wazi-code-plugin@sha256:e87382f99e74385c92ef550f5fa16bfc1f0778143a4e6f5aa75b5ed6f1f3c917</sub> ||
|<sub>icr.io/wazi-code/wazi-code-codeready@sha256:17c874e53b68b898e52e9693ef1113c35a86242e733519d119d363d59265096a</sub> ||
||<sub><b>1.1.0</b></sub>|
|<sub>icr.io/cpopen/wazi-code-operator-catalog@sha256:a436db8522674ac9dbd6b46188a2ceb03370a0a7987d5658a8ff383d7bce0bfe</sub> ||
|<sub>icr.io/wazi-code/wazi-code-dev-file@sha256:0aa798917576211c2ea59829f1db4d96462d9a42fb7cc2245a24938d15fbedc4</sub> ||
|<sub>icr.io/wazi-code/wazi-code-plugin@sha256:3997b50fd8a1940bb1328f4ea530156eecb6e34b16e8a26c9f4587187959a23f</sub> ||
|<sub>icr.io/wazi-code/wazi-code-codeready@sha256:fde1b1dc4117058f6e7d6f05876fb398c1f8b0b3291eea8c67134f6ea2688d26</sub> ||

## Feedback
  
We would love to hear feedback from you about IBM Wazi for Dev Spaces.  
File an issue or provide feedback here: [IBM Wazi for Dev Spaces Issues](https://github.com/IBM/wazi-devspaces-images/issues)

---
  
## Reference: Red Hat content

Links marked with this icon :door: are _internal to Red Hat_. This includes Jenkins servers, job configs in gitlab, and container sources in dist-git. 

Because these services are internal, in the interest of making all things open, we've copied as much as possible into this repo. Details below.

## Midstream code
This repo is used to house identical copies of the code used to build the **CodeReady Workspaces images** in Brew/OSBS, but made public to enable pull requests and easier contribution.

* Downstream code can be found in repos _internal to Red Hat_ at http://pkgs.devel.redhat.com/cgit/?q=codeready-workspaces :door:
    - select the `crw-2-rhel-8` branch for the latest `2.x` synced from upstream main branches, or 
    - select a branch like `crw-2.8-rhel-8` for a specific release, synced to a stable branch like `7.28.x`.

## Generated code

In some cases, where we need to house code generated by downstream processes in a public location, this repo will contain folders that end in `-generated` to differentiate from code that's synced from upstream and copied to downstream.

* [codeready-workspaces-operator-metadata-generated/manifests](https://github.com/redhat-developer/codeready-workspaces-images/tree/crw-2-rhel-8/codeready-workspaces-operator-metadata-generated/manifests/) contains CSVs with [pinned digests](http://pkgs.devel.redhat.com/cgit/containers/codeready-workspaces-operator-metadata/tree/container.yaml?h=crw-2-rhel-8#n24) :door:, generated from downstream Brew processes.

## Jenkins jobs

This repo also contains an identical copy of the [Jenkinsfiles and groovy](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/tree/master/jobs/CRW_CI) :door: sources used to configure the [jenkins-csb](https://gitlab.cee.redhat.com/ccit/jenkins-csb) :door: Configuration-as-Code (casc) Jenkins instance used to build the artifacts needed for Brew/OSBS builds. Since the server and config sources are _internal to Red Hat_, [this copy](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/) is provided to make it easier to see how CodeReady Workspaces is built. Hooray for open source!

* To run a local Jenkins, see [README](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/blob/master/README.md#first-time-user-setup) :door:
* [Job](https://main-jenkins-csb-crwqe.apps.ocp4.prod.psi.redhat.com/job/CRW_CI/job/Releng/job/sync-jenkins-gitlab-to-github_2.x/) :door: that performs the sync from [gitlab](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/blob/master/jobs/CRW_CI/Releng/sync-jenkins-gitlab-to-github.groovy) :door: to [github](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/sync-jenkins-gitlab-to-github.groovy) at intervals
* Other jobs are used to:
    * [build a series of artifact and container builds](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/build-all-images.groovy),
    * [sync midstream to downstream](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/crw-sync-to-downstream.groovy),
    * [build artifacts](https://github.com/redhat-developer/codeready-workspaces-images/tree/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/),
    * [orchestrate Brew builds](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/get-sources-rhpkg-container-build.groovy),
    * [copy containers to quay](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/push-latest-container-to-quay.groovy),
    * [check & update digests in registries/metadata images](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/update-digests-in-registries-and-metadata.groovy)
* Or, to:
    * [send email notifications](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/send-email-qe-build-list.groovy) of ER and RC builds
    * [tag sources & collect manifests](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/get-3rd-party-deps-manifests.groovy), [collect sources](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/get-3rd-party-sources.groovy) to create a release
    * set up subsequent releases ([branching](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/create-branches.groovy), [bumping versions](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/update-version-and-registry-tags.groovy))
