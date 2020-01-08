sudo yum update -y --quiet
sudo yum install docker -y --quiet

echo "$2" > ./gitlabtoken

sudo systemctl start docker

sudo cat ./gitlabtoken | sudo docker login -u="$1" --password-stdin gitlab-registry.tubit.tu-berlin.de

sudo docker pull gitlab-registry.tubit.tu-berlin.de/mcc-fred/fred/fred:latest

sudo docker run -it \
      --name=fred \
      -d \
      --restart=unless-stopped \
      -p 80:9001 \
      -p 5555:5555 \
      -v /tmp/config.toml:/config.toml \
      gitlab-registry.tubit.tu-berlin.de/mcc-fred/fred/fred:latest --config config.toml