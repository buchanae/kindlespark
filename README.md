# kindlespark

kindlespark can generate [SparkNotes](http://www.sparknotes.com/lit/gatsby/) in Kindle format.

More technically, it scrapes the SparkNotes website, does some basic formatting,
and uses Amazon's KindleGen to generate a .mobi file.

This was inspired by [this Reddit post](http://www.reddit.com/r/kindle/comments/g7l2h/sparknotes_for_kindle/).

# Prerequisites

First, you need to download [KindleGen](http://www.amazon.com/gp/feature.html?docId=1000765211).
Unzip this into the kindlespark directory (or otherwise link the kindlegen executable there).

You'll need python 2.7 and you'll need to install [scrapy](http://scrapy.org/).

# Usage

```python kindlespark.py http://www.sparknotes.com/lit/gatsby/```

# Known Issues

Disclaimer: I whipped this up in a couple hours. There are parts that are very obviously bad.

Currently, this doesn't scrape the chapter/section analysis. Those are multiple pages, and rhtml,
and they didn't "just work". I'll fix that at some point though.
