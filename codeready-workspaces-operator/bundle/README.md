### Operator Package Manager Bundles
These are instructions for pre-generating a bundle.

#### IBM Wazi for Dev Spaces - Bundle
Execute the following to generate the necessary wazi bundle files to perform a simple multi-architecture `docker build`.
Change directory to the bundles folder and execute the following.

```terminal
cd wazi
opm alpha bundle generate --channels v2.1 --default v2.1 --directory ./manifests --package ibm-wazi-for-devspaces
cd ..
```

#### IBM Developer for z/OS Enterprise Edition - Bundle
Execute the following to generate the necessary wazi bundle files to perform a simple multi-architecture `docker build`.
Change directory to the bundles folder and execute the following.

```terminal
cd idzee
opm alpha bundle generate --channels v2.1 --default v2.1 --directory ./manifests --package ibm-developer-for-zos-enterprise-edition
cd ..
```
