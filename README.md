# Hoerbert Playlist Generator

This tool can be used to convert a bunch of mp3 files and distribute them equally on a set of buttons (defined range) on a Hoerbert children's mp3 player (see https://www.hoerbert.com/).

It automates the conversion as described [here](https://www.hoerbert.com/service/linux-benutzer/) and groups those titles as needed by hoerbert - in subfolders with wav files. Each folder thereby contains a equal portion of the converted files. So if your input is one album with 9 mp3 files, you can convert those and configure Hoerbert to have the first row of buttons each carry 3 of the titles.

## Dependencies

```
$ sudo apt-get install libsox-fmt-mp3 sox
```

## build and run

```
#build
$ make

#get help
$ ./hoerbert -h

#convert mp3s in ./path/with/mp3s, distributed equally over all 9 Hoerbert buttons
$ ./hoerbert -in ./path/with/mp3s -out /path/to/sdcard

#convert mp3s in ./path/with/mp3s, distributed equally over the first 3 Hoerbert buttons
$ ./hoerbert -in ./path/with/mp3s -out /path/to/sdcard -start 0 -end 2

#convert mp3s in ./path/with/mp3s, distributed equally over the bottom 2 rows of Hoerbert buttons
$ ./hoerbert -in ./path/with/mp3s -out /path/to/sdcard -start 3 -end 8
```
