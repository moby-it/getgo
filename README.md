# GetGo Continous Deployment Tool
<p align="center">
<img src="https://github.com/moby-it/getgo/assets/27289923/d36c026a-92a4-4f88-9b9d-f2bcbcf3431f" width="400" />
</p>

GetGo is a tool that aims to help you deploy your Dockerhub Repositories to your virtual machine. It is based on [Dockerhub Webhooks](https://docs.docker.com/docker-hub/webhooks).

# Motivation
After moving resources away from mainstream cloud provider infrastructure and started using VPS on multiple providers online, I was in lack of a tool to automaticaly deploy my container from Dockerhub to my VPS and since I wanted to do something in Go either way, I decided to create this tool for myself. After a short while I realized that it's definitely easier to update your container via a remote ssh through Github Actions for solving my CD issue but nevertheless I created a stable version for this tool to sharpen my Go skills, before archiving it. It also has some value if someone want to completely decouple his Deployment Circle from his Source Control repository.

# How it Works

- GetGo runs as a **systemd** service inside your Debian-Based Virtual Machice. It exposes a simple HTTP POST endpoint at localhost:32041/deploy. The endpoint expects a json body similar to [what dockerhub uses for its webhooks](https://docs.docker.com/docker-hub/webhooks/#example-webhook-payload).

- If there is a running container with a name that matches the `repository.repo_name-push_tag.tag` then GetGo pulls the new image, destroys and recreates the container, exposing the same ports.

# Install and Configure

## Requirements

For GetGo to work, it is expected that you have at least ssh access to a remote machine with the following:

1. Docker installed. [Install docker](https://docs.docker.com/engine/install/ubuntu/)
2. Go installed (sudo snap install go --classic)
3. make installed (sudo apt-get install build-essential)
4. listening to http calls. ([Install and configure NGINX](https://ubuntu.com/tutorials/install-and-configure-nginx#1-overview))

## Install instructions

1. Connect yo your remote machine.
2. Clone this repository - [How to clone a repository](https://git-scm.com/book/en/v2/Git-Basics-Getting-a-Git-Repository)
3. Open a terminal and navigate inside the above repository folder.
4. run `sudo make install`.

Now the service should already be active. To check this, run `sudo systemctl status getgo`. If you want to enable this service to run on startup, run `sudo systemctl enable getgo`.

For the installation to be complete, you need to provide a DOCKER_USERNAME and DOCKER_PASSWORD as an environment variable for your systemd service. To do that run `sudo systemctl edit getgo`. You should add the following lines :

```
[Service]
Environment="DOCKER_USERNAME=YOUR_DOCKER_USERNAME"
Environment="DOCKER_PASSWORD=YOUR_DOCKER_ACCESS_TOKEN"`
```

[Suggested way of adding environment variables to systemd](https://serverfault.com/questions/413397/how-to-set-environment-variable-in-systemd-service)

After saving make sure to run `sudo systemctl daemon-reload`

## Configure

After installing you have a service running in your system that exposes a simple HTTP POST endpoint at `http://locahost:31042/deploy`.

Since this is a Continious Deployment tool, it only makes sense if you can hit this endpoint via the internet. At this point you should make sure that http://locahost:31042/deploy is accessible from the outside world. If you use NGINX you can read more on how to do this [here](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/).

The last piece of the puzzle is configuring the source of the "to-be-deployed" container. GetGo was written for consuming calls made from [DockerHub Webhooks](https://docs.docker.com/docker-hub/webhooks).
If you've wired up your web-server to correctly forward calls from `https://your-web-server-address.com/deploy` to `http://localhost:31042/deploy` then GetGo will start a deployment process if the following criteria are met:

1. The HTTP POST request has a body that follows the [DockerHub Webhook schema](https://docs.docker.com/docker-hub/webhooks/#example-webhook-payload). 
2. There is a container running on your machine, with a container name that matches the `repo_name-pushed_tag`.


## What NOT to expect from GetGo

- GetGo does not deploy your containers for the first time. **It expects of you to set the container up for the first time, with a container-name=repository.repo_name-push_data.tag and networking properties of your choice.**
- It is your responsibility to expose the port of this service to the internet on any of your VMs if you want to trigger it from an online container registry (like Dockerhub)

---

<h3 align="center">
<img src="https://moby-it.com/assets/MobyIT_Icon_50.png" width="30"/>
Build by Moby IT
<img src="https://moby-it.com/assets/MobyIT_Icon_50.png"  width="30"/>

</h3>

