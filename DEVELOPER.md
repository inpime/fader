# Developer guid

Only the first time 
```
cd state/app
ln -s Developer.Dockerfile Dockerfile
cd ../../
```

Apply updates in dockerimage
```
make build-linux-dev
docker-compose build
docker-compose up -d
```