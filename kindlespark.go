package main

import (
  "bytes"
  "fmt"
  "log"
  "os"
  "path"
  "strings"
  "text/template"
  "github.com/PuerkitoBio/goquery"
)


type BookMeta struct {
  Title string
  Author string
}

type Section struct {
  Title string
  Content string
  Filename string
  Index int
}

type TopLevelTemplateData struct {
  TocSections string
  ManifestSections string
  SpineSections string
  NcxSections string
  BookMeta BookMeta
}

type SectionOutput struct {
  Section Section
  Content string
}

type Output struct {
  Toc string
  Opf string
  Ncx string
  Sections []SectionOutput
}


func main() {

  if len(os.Args) != 2 {
    fmt.Println("Usage: kindlespark <book name>")
    return
  }

  name := os.Args[1]
  baseUrl := "http://www.sparknotes.com/lit/" + name + "/"
  doc, err := goquery.NewDocument(baseUrl)

  if err != nil {
    log.Fatal(err)
  }

  bookMeta := parseBookMeta(doc)
  sectionHrefs := parseSectionHrefs(doc)
  var sections []Section

  for idx, href := range sectionHrefs {
    sectionDoc, _ := goquery.NewDocument(baseUrl + href)
    section := parseSection(sectionDoc)
    section.Filename = strings.Replace(href, ".rhtml", ".html", 1)
    section.Index = idx
    sections = append(sections, section)
  }

  output := buildOutput(bookMeta, sections)
  writeOutput(name, output)

  fmt.Printf("Ebook content written to the \"%s\" directory. Now run: ./kindlegen %s/%s.opf\n", name, name, name)
}

func parseBookMeta(doc *goquery.Document) BookMeta {
  title := doc.Find("h1.title").Text()
  author := doc.Find("h2.author").Text()

  return BookMeta{title, author}
}

func parseSection(doc *goquery.Document) Section {
  title := doc.Find("div.left-menu li.active").Text()
  content, _ := doc.Find("div.studyGuideText").Html()
  return Section{ Title: title, Content: content }
}

func parseSectionHrefs(doc *goquery.Document) []string {
  return doc.Find("div.entry a").Map(func(i int, s *goquery.Selection) string {
    href, _ := s.Attr("href")
    return href;
  })
}


func buildOutput(bookMeta BookMeta, sections []Section) (output Output) {

  var tocSectionsBuffer bytes.Buffer
  var manifestSectionsBuffer bytes.Buffer
  var spineSectionsBuffer bytes.Buffer
  var ncxSectionsBuffer bytes.Buffer

  for _, section := range sections {
    var sectionContentBuffer bytes.Buffer

    err := CONTENT_TPL.Execute(&sectionContentBuffer, section)
    if err != nil { panic(err) }

    output.Sections = append(output.Sections, SectionOutput{ section, sectionContentBuffer.String() })

    err = TOC_SECTION_TPL.Execute(&tocSectionsBuffer, section)
    if err != nil { panic(err) }

    err = MANIFEST_SECTION_TPL.Execute(&manifestSectionsBuffer, section)
    if err != nil { panic(err) }

    err = SPINE_SECTION_TPL.Execute(&spineSectionsBuffer, section)
    if err != nil { panic(err) }

    err = NCX_SECTION_TPL.Execute(&ncxSectionsBuffer, section)
    if err != nil { panic(err) }
  }

  var tocBuffer bytes.Buffer
  var opfBuffer bytes.Buffer
  var ncxBuffer bytes.Buffer

  templateData := TopLevelTemplateData{
    tocSectionsBuffer.String(),
    manifestSectionsBuffer.String(),
    spineSectionsBuffer.String(),
    ncxSectionsBuffer.String(),
    bookMeta,
  }

  err := TOC_TPL.Execute(&tocBuffer, templateData)
  if err != nil { panic(err) }

  err = OPF_TPL.Execute(&opfBuffer, templateData)
  if err != nil { panic(err) }

  err = NCX_TPL.Execute(&ncxBuffer, templateData)
  if err != nil { panic(err) }

  output.Toc = tocBuffer.String()
  output.Opf = opfBuffer.String()
  output.Ncx = ncxBuffer.String()

  return output
}

func writeOutput(name string, output Output) {
  os.Mkdir(name, 0777)

  tocFile, err := os.Create(path.Join(name, "toc.html"))
  if err != nil { panic(err) }
  defer tocFile.Close()

  tocFile.WriteString(output.Toc)
  tocFile.Sync()

  opfFile, err := os.Create(path.Join(name, name + ".opf"))
  if err != nil { panic(err) }
  defer opfFile.Close()

  opfFile.WriteString(output.Opf)
  opfFile.Sync()

  ncxFile, err := os.Create(path.Join(name, "toc.ncx"))
  if err != nil { panic(err) }
  defer ncxFile.Close()

  ncxFile.WriteString(output.Ncx)
  ncxFile.Sync()

  for _, sectionOutput := range output.Sections {
    sectionFile, err := os.Create(path.Join(name, sectionOutput.Section.Filename))
    if err != nil { panic(err) }
    defer sectionFile.Close()

    sectionFile.WriteString(sectionOutput.Content)
    sectionFile.Sync()
  }
}


var CONTENT_TPL, _ = template.New("CONTENT_TPL").Parse(`
<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<title>{{.Title}}</title>
<link rel="stylesheet" href="styles.css" type="text/css">
</head>
<body>
    {{.Content}}
</body>
</html>
`)

var TOC_SECTION_TPL, _ = template.New("TOC_SECTION_TPL").Parse(`
<li><a href="{{.Filename}}">{{.Title}}</a></li>
`)

var MANIFEST_SECTION_TPL, _ = template.New("MANIFEST_SECTION_TPL").Parse(`
<item id="item-{{.Index}}" media-type="application/xhtml+xml" href="{{.Filename}}"></item>
`)

var SPINE_SECTION_TPL, _ = template.New("SPINE_SECTION_TPL").Parse(`
<itemref idref="item-{{.Index}}"/>
`)

var NCX_SECTION_TPL, _ = template.New("NCX_SECTION_TPL").Parse(`
<navPoint class="{{.Filename}}" id="{{.Filename}}" playOrder="{{.Index}}"\>
  <navLabel><text>{{.Title}}</text></navLabel>
  <content src="{{.Filename}}"/>
</navPoint>
`)

var TOC_TPL, _ = template.New("TOC_TPL").Parse(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head><title>Table of Contents</title></head>
<body>
<div><h1><b>TABLE OF CONTENTS</b></h1></div><br />

<div><ul>
    {{.TocSections}}
</ul></div><br />
</body>
</html>
`)

var OPF_TPL, _ = template.New("OPF_TPL").Parse(`
<?xml version="1.0" encoding="utf-8"?>
<package unique-identifier="uid">
  <metadata>
    <dc-metadata xmlns:dc="http://purl.org/metadata/dublin_core" xmlns:oebpackage="http://openebook.org/namespaces/oeb-package/1.0/">
    <dc:Title>Sparknotes: {{.BookMeta.Title}}</dc:Title>
    <dc:Creator>Sparknotes</dc:Creator>
    <dc:Language>en-us</dc:Language>
    <dc:Identifier id="uid">9095C522E6</dc:Identifier>
    </dc-metadata>

    <x-metadata>
        <output encoding="utf-8"></output>
    </x-metadata>
  </metadata>

  <manifest>
    <item id="toc" media-type="application/xhtml+xml" href="toc.html"></item>
    {{.ManifestSections}}

    <item id="toc-ncx" media-type="application/x-dtbncx+xml" href="toc.ncx"/>
  </manifest>

  <spine toc="toc-ncx">
    <itemref idref="toc"/>
    {{.SpineSections}}

  </spine>
  <tours></tours>
  <guide></guide>
</package>
`)

var NCX_TPL, _ = template.New("NCX_TPL").Parse(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">

<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" xml:lang="en-US">
<head>
    <meta name="dtb:uid" content="BookId"/>
    <meta name="dtb:depth" content="2"/>
    <meta name="dtb:totalPageCount" content="0"/>
    <meta name="dtb:maxPageNumber" content="0"/>
</head>

<docTitle><text>{{.BookMeta.Title}}</text></docTitle>
<docAuthor><text>{{.BookMeta.Author}}</text></docAuthor>

<navMap>
  <navPoint class="toc" id="toc" playOrder="1">
    <navLabel>
      <text>Table of Contents</text>
    </navLabel>
    <content src="toc.html"/>
  </navPoint>

  {{.NcxSections}}
</navMap>
</ncx>
`)
