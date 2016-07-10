#!/bin/sh
set -e

echo "Create of workspace."
echo ""
mkdir -p state/app

echo "Download files: Dockerfile, entrypoint.sh, docker-compose.yml."
echo ""
curl https://raw.githubusercontent.com/inpime/fader/master/state/app/Production.Dockerfile > state/app/Dockerfile
curl https://raw.githubusercontent.com/inpime/fader/master/state/app/entrypoint.sh > state/app/entrypoint.sh
curl https://raw.githubusercontent.com/inpime/fader/master/docker-compose.yml > docker-compose.yml

echo "Done."
echo ""
echo "For run follow the command 'docker-compose up -d'."
echo ""