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

## Check mapping

```
docker exec -it fader_elasticsearch_1 curl -v 127.0.0.1:9200/fader/_mapping?pretty=true
```