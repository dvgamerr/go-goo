# go-goo

A [Google Photos](http://photos.google.com/) backup tool.

`go-goo` tool uses [Google Photos API](https://developers.google.com/photos/library/guides/get-started#enable-the-api) to continously download all photos on a local device.

It can be used as a daemon to keep in sync with a google-photos account.

## Build

![Go](https://github.com/dvgamerr/go-goo/workflows/Go/badge.svg)

## Downloading and Installing:

* Go to[Latest](https://github.com/dvgamerr/go-goo/releases/latest) release.
* Download go-goo.zip

Unzip and run, there are no other dependencies.

The zip contains Windows, Linux and MAC OS binaries.

## Usage:

### Enable Google-Photos API:

(There is more than one way to do the following):
* Go to [API](https://developers.google.com/photos/library/guides/get-started#enable-the-api)
* Click on `Enable the Google Photos API` button
* Create a new project, name it whatever you like.
* Write `go-goo` in the `product name` field, (can be whatever you like)
* On the `when are you calling from`, choose `Desktop app`, click `Create`.
* Download the client configuration.

### Configure go-goo:

* Copy the downloaded `credentials.json` to the same folder with `go-goo`.
* Run `go-goo`, and follow the provided link.
* Sign in, and click `Allow`.
* You will be redirected to a local address and it `This browser window can be now closed...`.
* `go-goo` will authorize and start downloading content (authorization code will be automically received). 
```
$ ./go-goo
2018/09/12 10:18:07 This is go-goo ver 0.1
Go to the following link in your browser then type the authorization code:
https://accounts.google.com/o/oauth2/auth?access_type=...
4/WACqgFeX5OTB8X4LWd5i2TFH....
Saving credential file to: token.json
2018/09/12 10:20:07 Connecting ...
2018/09/12 10:20:07 Processed: 0, Downloaded: 0, Errors: 0, Total Size: 0 B, Waiting 5s
```


This is probably not what you want, hit `crt-c` to stop it.

### Usage:

```
Usage of ./go-goo:
  - album
        download only from this album (use google album id)
  -folder string
        backup folder (default current working directory)
  -force
        ignore errors, and force working
  -version
        at startup, print the go-goo version
  -logfile string
        log to this file
  -credentials-file string
        filepath to where the credentials file can be found (default 'credentials.json')
  -token-file string
        filepath to where the token should be stored (default 'token.json')
  -loop
        loops forever (use as daemon)
  -max int
        max items to download (default 2147483647)
  -pagesize int
        number of items to download on per API call (default 50)
  -throttle int
        Time, in seconds, to wait between API calls (default 5)
  -folder-format string
        Time format used for folder paths based on https://golang.org/pkg/time/#Time.Format (default "2016/Janurary")
  -use-file-name
        Use file name when uploaded to Google Photos (default off)
  -include-exif
        Retain EXIF metadata on downloaded images. Location information is not included because Google does not include it. (default off)
  -download-throttle
        Rate in KB/sec, to limit downloading of items (default off)
  -concurrent-downloads
        Number of concurrent item downloads (default 5)
  -loopback-port
        Port number bound on `127.0.0.1` to receive auth code during authentication (default 8080)
```

On Linux, running the following is a good practice:

```
$ ./go-goo -folder archive -logfile gitmoo.log -use-file-name -include-exif -loop -throttle 45 &
```

This will start the process in background, making an API call every 45 seconds, looping forever on all items and saving them to `{pwd}/archive`. All images will be downloaded with a filename and metadata as close to original as Google offers through the api.

Logfile will be saved as `gitmoo.log`.

#### Naming

Files are created as follows:

`[folder][year][month][day]_[hash].json` and `.jpg`. The `json` file holds the metadata from `google-photos`. 

## Building:

To build you may need to specify that module download mode is using a vendor folder.  Failure to do this will mean that modified vendor files will not be used.

`go build -mod vendor`

## Testing:

`go test -mod vendor ./...`

## Docker (Linux only)

You can run go-goo in Docker. At the moment you have to build the image yourself. After cloning the repo run:

```
$ docker build -t dvgamerr/go-goo:latest .
```

Now run gitmoo-goo in Docker:

```
$ docker run -v $(pwd):/app --user=$(id -u):$(id -g) dvgamerr/go-goo:latest
```

Replace `$(pwd)` with the location of your storage directory on your computer.
Within the storage directory go-goo expects the `credentials.json` and will place all the downloaded files.

The part `--user=$(id -u):$(id -g)` ensures that the downloaded files are owned by the user launching the container.

Configuring additional settings is possible by adding command arguments like so:
```
$ docker run -v $(pwd):/app --user=$(id -u):$(id -g) dvgamerr/go-goo:latest -loop -throttle 45
```
