docker compose up -d
docker compose exec -T roach-node-1 ./cockroach init --insecure --host=roach-node-1:26357 2>/dev/null
