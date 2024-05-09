docker stop $(docker ps -a -q)
docker container prune -f
docker volume prune -af
