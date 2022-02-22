### INSTAGRAM MEDIA UI FOR DOWNLOAD THEM

Extract Instagram Media images, reels,tv, albums etc. urls for download them only one UI.

- No login required
- Fetch nearly all media except single video posts
- Embed templates ready to use - build code and ready to use one single executable file

![Simple UI](https://raw.githubusercontent.com/uretgec/simple-instagram-api-ui/master/screenshot.png)

### Installation

Docker-Compose

```
export SERVICE_BUILD=$(date '+%Y%m%d%H%M') && export SERVICE_COMMIT_ID=$(git describe --always) &&  docker-compose -f docker-compose-builder.yml build --compress --progress plain
```

```
export SERVICE_BUILD=$(date '+%Y%m%d%H%M') && export SERVICE_COMMIT_ID=$(git describe --always) &&  docker-compose -f docker-compose.yml build --compress --progress plain
```

```
docker-compose up
```

```
docker-compose down
```

Local Use:

```
cd reweb && go run .
```


Local Build: Run after go to build folder

```
./build.sh reweb
```

Finished. Open the favourite browser and Go to "localhost:3001"

>Inspired by [Scrip7/simple-instagram-api](https://github.com/Scrip7/simple-instagram-api)