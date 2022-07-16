# postalcode-search

Download latest `x-ken-all.csv` from http://zipcloud.ibsnet.co.jp/
Put it to the root directory of the project

## Run locally

```
$ go run cmd/main.go
```

- http://localhost:8080/?q=長野
- http://localhost:8080/?postalcode=1000000

## Test

```
$ go test ./...
```

## Deploy

```
$ ./deploy.sh
```

Then, open https://asia-northeast1-postalcode-firebase.cloudfunctions.net/postalcode-search/?q=ながの
