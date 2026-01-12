# LightRay
Plug-and-Play LAN-Based Cluster Tool for Zero-config CPU/GPU Resource Pooling

## Purpose
In colleges most students nowadays have gaming laptops which come shipped with a decent gpu. This tool aims to let students build a GPU cluster on LAN so that can avoid buying gpu services online (like colab etc.). LightRay helps in pooling the GPU computing resources easily. It is built using go and very minimalistic. LightRay uses Ray.io 's docker images under the hood. The users can simply connect their laptops over ethernet and use LightRay's UI to setup and send jobs to the cluster.

## Installation
The requirements for using LightRay are:-
1. Ubuntu OS (or any linux OS. The reason Windows is not recommended because while we are developing we were facing a lot of issues because of window's firewall.)
2. Docker cli installed. (The required images will be pulled automatically.)
3.  Go installed.

## Cluster Setup
1. navigate to the folder where the new.go file is present. run that file using the command ```go run new.go```. This launches LightRay.
2. Open the LightRay's UI in your browser.
3. First the head has to be setup, then the subsequent computers can join as workers.


