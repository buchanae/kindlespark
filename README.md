# kindlespark

kindlespark can generate [SparkNotes](http://www.sparknotes.com/) ebooks for Kindle.

More technically, it scrapes the SparkNotes website, does some basic formatting,
and uses Amazon's [KindleGen](http://www.amazon.com/gp/feature.html?docId=1000765211) to generate a .mobi file.

This was inspired by [this Reddit post](http://www.reddit.com/r/kindle/comments/g7l2h/sparknotes_for_kindle/).

# Download

First, you need to download [KindleGen](http://www.amazon.com/gp/feature.html?docId=1000765211).

Next, download a release of kindlespark:

- [Windows](https://github.com/abuchanan/kindlespark/releases/download/2.0.0/kindlespark-win-amd64)
- [Mac](https://github.com/abuchanan/kindlespark/releases/download/2.0.0/kindlespark-mac)
- [Linux](https://github.com/abuchanan/kindlespark/releases/download/2.0.0/kindlespark-linux-amd64)

(Note, I have only tested this on Mac 10.10 so far!)

# Usage

These programs require you to use the command line. These instructions assume you are familiar with that environment. If you're not, sorry, I'll try to make this simpler at some point.

Find the SparkNotes URL you want to convert, for example: http://www.sparknotes.com/lit/gatsby/

Run `./kindlespark gatsby http://www.sparknotes.com/lit/gatsby/`

Run `./kindlegen ./gatsby/gatsby.opf`

Upload `./gatsby/gatsby.mobi` to your Kindle.

# Build

```
GOPATH=`pwd`/gopackages/ go build kindlespark.go
```
