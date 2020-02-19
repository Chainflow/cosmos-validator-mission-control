#!/bin/sh

install_docker(){
    if command_exists docker && [ -e /var/run/docker.sock ]; then
        echo "Docker exists..."
    else
        curl https://get.docker.com | bash
    fi
}

create_docker_nw() {
  sudo docker network create -d bridge monitoring-nw
}

run_grafana() {
  sudo docker volume create grafana-storage
  sudo docker run -d --name=grafana --network=monitoring-nw -p 3000:3000 -v grafana-storage:/var/lib/grafana grafana/grafana
}

install_docker

create_docker_nw

run_grafana
