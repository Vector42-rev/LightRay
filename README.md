# LightRay
Plug-and-Play LAN-Based Cluster Tool for Zero-config CPU/GPU Resource Pooling

## Purpose
In colleges most students nowadays have gaming laptops which come shipped with a decent gpu. This tool aims to let students build a GPU cluster on LAN so that can avoid buying gpu services online (like colab etc.). LightRay helps in pooling the GPU computing resources easily. It is built using go and very minimalistic. LightRay uses Ray.io 's docker images under the hood. The users can simply connect their laptops over ethernet and use LightRay's UI to setup and send jobs to the cluster.

## Installation
The requirements for using LightRay are:-
1. Ubuntu OS (or any linux OS. The reason Windows is not recommended because while we are developing we were facing a lot of issues because of window's firewall.)
2. Docker cli installed. (The required images will be pulled automatically.)
3.  Go installed.
4.  Clone the repository to get the cluster setup file (new.go) and the sample training job (python file).

## Cluster Setup
1. navigate to the folder where the new.go file is present. run that file using the command ```go run new.go```. This launches LightRay.
2. Open the LightRay's UI in your browser.
3. First the head has to be setup, then the subsequent computers can join as workers.

## Execute Jobs
1. keep your python file in the path as expected in the source code. when you click on "Copy Training File" , the file will be copied inside the container and executed on the workers.


## Screenshots
1. Two Laptops connected over Ethernet, with their CPU-GPU resources pooled.
<p align="center">
  <img width="1794" height="876" alt="i1" src="https://github.com/user-attachments/assets/f8b33dcb-0f84-4416-a72a-7720f7d603d9" />
</p>

<p align="center"><i>UI of LightRay</i></p>

2. Loss plots
<p align="center">
<img width="1843" height="869" alt="i2" src="https://github.com/user-attachments/assets/b89906e9-6c69-4fbd-936d-7cc585d121b3" />
<p align="center"><i>Tensorboard after completion of Model training Job</i></p>

