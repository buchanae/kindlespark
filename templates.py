from __future__ import unicode_literals

CONTENT_TPL = u'''
<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<title>{title}</title>
<link rel="stylesheet" href="styles.css" type="text/css">
</head>
<body>
    {content}
</body>
</html>
'''

TOC_TPL = u'''
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Table of Contents</title></head>
<body>
<div><h1><b>TABLE OF CONTENTS</b></h1></div><br />

<div><ul>
    {toc_sections}
</ul></div><br />
</body>
</html>
'''

TOC_SECTION_TPL = u'''
<li><a href="{file_name}">{title}</a></li>
'''


COVER = '''        <EmbeddedCover>cover.jpg</EmbeddedCover>
'''

OPF_TPL = u'''
<?xml version="1.0" encoding="utf-8"?>
<package unique-identifier="uid">
  <metadata>
    <dc-metadata xmlns:dc="http://purl.org/metadata/dublin_core" xmlns:oebpackage="http://openebook.org/namespaces/oeb-package/1.0/">
    <dc:Title>Sparknotes: {title}</dc:Title>
    <dc:Creator>Sparknotes</dc:Creator>
    <dc:Language>en-us</dc:Language>
    <dc:Identifier id="uid">9095C522E6</dc:Identifier>
    </dc-metadata>

    <x-metadata>
        <output encoding="utf-8"></output>
    </x-metadata>
  </metadata>

  <manifest>
    <item id="item-1" media-type="application/xhtml+xml" href="toc.html"></item>
    {manifest_sections}

    <item id="My_Table_of_Contents" media-type="application/x-dtbncx+xml" href="toc.ncx"/>
  </manifest>

  <spine toc="My_Table_of_Contents">
    <itemref idref="item-1"/>
    {spine_sections}

  </spine>
  <tours></tours>
  <guide></guide>
</package>
'''

MANIFEST_SECTION_TPL = u'''
<item id="item-{idx}" media-type="application/xhtml+xml" href="{file_name}"></item>
'''

SPINE_SECTION_TPL = u'''
<itemref idref="item-{idx}"/>
'''


NCX_TPL = u'''
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">

<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" xml:lang="en-US">
<head>
    <meta name="dtb:uid" content="BookId"/>
    <meta name="dtb:depth" content="2"/>
    <meta name="dtb:totalPageCount" content="0"/>
    <meta name="dtb:maxPageNumber" content="0"/>
</head>

<docTitle><text>{title}</text></docTitle>
<docAuthor><text>{author}</text></docAuthor>

<navMap>
  <navPoint class="toc" id="toc" playOrder="1">
    <navLabel>
      <text>Table of Contents</text>
    </navLabel>
    <content src="toc.html"/>
  </navPoint>

  {sections}
</navMap>

</ncx>
'''


NCX_SECTION_TPL = u'''
<navPoint class="{file_name}" id="{file_name}" playOrder="{idx}"\>
  <navLabel><text>{title}</text></navLabel>
  <content src="{file_name}"/>
</navPoint>
'''
