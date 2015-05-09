#!/usr/bin/env python

import argparse
import sys
import os
import re
import urlparse
import shutil

from twisted.internet import reactor
from scrapy import Spider, Item, Field, Request
from scrapy import log, signals
from scrapy.crawler import Crawler
from scrapy.settings import Settings
from scrapy.utils.project import get_project_settings

import templates


class Meta(Item):
    title = Field()
    author = Field()
    section_order = Field()


class Section(Item):
    title = Field()
    content = Field()
    url = Field()

    @property
    def file_name(self):
        return urlparse.urlparse(self['url']).path.split('/')[-1]



class SparknotesSpider(Spider):

    def __init__(self, name, url):
        self.name = name
        self.start_urls = [url]
        super(SparknotesSpider, self).__init__()

    def parse_section(self, response):
        try:
            yield Section(
                title=response.css("div.left-menu li.active::text")[0].extract(),
                content=response.css("div.studyGuideText")[0].extract(),
                url=response.url
            )
        except AttributeError:
            pass

    def parse(self, response):
        meta = Meta(
            title=response.css("h1.title::text")[0].extract(),
            author=response.css("h2.author::text")[0].extract()
        )

        section_order = []

        for section_url in response.css("div.entry a").xpath("@href").extract():
            # TODO temporarily ignoring these since they seem to cause lots of problems
            if 'rhtml' in section_url:
                continue

            u = urlparse.urljoin(response.url, section_url)
            section_order.append(u)
            print u
            yield Request(u, callback=self.parse_section)

        meta['section_order'] = section_order
        yield meta 


class SparkMobiPipeline(object):

    def process_item(self, item, spider):
        if isinstance(item, Meta):
            self.meta = item

        elif isinstance(item, Section):
            self.sections[item['url']] = item

        return item

    def open_spider(self, spider):
        self.name = spider.name
        self.meta = None
        self.sections = {}

    def get_sections_in_order(self):
        ordered = []

        for key in self.meta['section_order']:
            try:
                section = self.sections[key]
                ordered.append(section)
            except KeyError:
                log.msg('Missing section: {}'.format(key))

        return ordered

    def write_file(self, file_name, content):
        path = os.path.join(self.name, file_name)
        with open(path, 'w') as fh:
            fh.write(content.encode('utf8'))

    def close_spider(self, spider):
        toc_sections = ''
        manifest_sections = ''
        spine_sections = ''
        ncx_sections = ''

        sections = self.get_sections_in_order()

        for idx, section in enumerate(sections, start=2):
            file_name = section.file_name

            content = templates.CONTENT_TPL.format(
                title=self.meta['title'],
                content=section['content']
            )
            self.write_file(file_name, content)


            toc_sections += templates.TOC_SECTION_TPL.format(file_name=file_name,
                                                   title=section['title'])

            spine_sections += templates.SPINE_SECTION_TPL.format(idx=idx)

            manifest_sections += templates.MANIFEST_SECTION_TPL.format(idx=idx,
                                                             file_name=file_name)

            ncx_sections += templates.NCX_SECTION_TPL.format(file_name=file_name,
                                                   idx=idx, title=section['title'])

        toc = templates.TOC_TPL.format(
            toc_sections=toc_sections)

        self.write_file('toc.html', toc)

        opf = templates.OPF_TPL.format(
            title=self.meta['title'],
            manifest_sections=manifest_sections,
            spine_sections=spine_sections)

        self.write_file(self.name + '.opf', opf)

        ncx = templates.NCX_TPL.format(
            title=self.meta['title'],
            author=['author'],
            sections=ncx_sections)

        self.write_file('toc.ncx', ncx)


def main():
    #shutil.copy('styles.css', index.dirname)
    #shutil.copy('cover.jpg', index.dirname)


    parser = argparse.ArgumentParser()
    parser.add_argument('url')
    args = parser.parse_args()

    name = urlparse.urlparse(args.url).path.rstrip('/').split('/')[-1]
    print name

    if not os.path.exists(name):
        os.makedirs(name)

    spider = SparknotesSpider(name, args.url)

    settings = get_project_settings()
    settings.set('ITEM_PIPELINES', {'__main__.SparkMobiPipeline': 500})

    log.start(loglevel=log.DEBUG)

    crawler = Crawler(settings)
    crawler.signals.connect(reactor.stop, signal=signals.spider_closed)
    crawler.configure()
    crawler.crawl(spider)
    crawler.start()

    reactor.run()

    os.system('./kindlegen {0}/{0}.opf'.format(name))

main()

# TODO get section pages
