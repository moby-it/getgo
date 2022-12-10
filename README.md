# GetGo Continous Deployment Tool

GetGo is a tool that aims to help you deploy your Container Registry continious integrations to your virtual machine.

## Why GetGo? 

GetGo aims to provide a seemles integration with a Container Registry to your VM, whenever it resides. 

GetGo runs as a service on a given VM [See install intruction bellow](#install-and-configure). Making use of Docker Webhooks from your Container Registry you can easily configure the continious deployment of your service to the VM of your choice.

## Install and Configure
## What made me do this? Need behind the idea
As a developer using popular cloud providers (AWS and Azure mostly) I realised that I became more entangled with these services than I wished for. I realised that I was paying more than I was usually needing. This is why I started transfering my infrastructure to seperate providers, comparably smaller than the big fishes that I mentioned earlier, and noticed pretty easily that there are many alternatives out there.

While there are plenty of upsides using "smaller" cloud providers, a common downside of those alternatives is that on some occasions the tooling is not as easy-to-use as on the bigger fishes in the pond, which makes sense.

This is why I decided that maybe I can contribute on that matter. While there are (for me at least) really easy ways to build and deploy your images to a container registry, there is not a tool that will help me specifically newly pushed image to my Web Server. 
