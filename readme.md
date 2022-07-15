# postalcode-search

1. Download `x-ken-all.csv` from http://zipcloud.ibsnet.co.jp/
2. Encode the csv file from Shift-JIS to UTF8
3. Put it to the root directory of the project
4. Run `./run.sh` and test
5. Run `./deploy.sh` to deploy

## Run locally

```
$ ./run.sh
```

Then, open http://localhost:8080/?q=長野

## Deploy

```
$ ./deploy.sh
```

Then, open https://asia-northeast1-postalcode-firebase.cloudfunctions.net/postalcode-search/?q=ながの
