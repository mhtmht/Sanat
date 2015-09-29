package windows

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"hasseg.org/sanat/model"
)

const xmlHeaderString string = `  <xsd:schema id="root" xmlns="" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:msdata="urn:schemas-microsoft-com:xml-msdata">
    <xsd:import namespace="http://www.w3.org/XML/1998/namespace" />
    <xsd:element name="root" msdata:IsDataSet="true">
      <xsd:complexType>
        <xsd:choice maxOccurs="unbounded">
          <xsd:element name="metadata">
            <xsd:complexType>
              <xsd:sequence>
                <xsd:element name="value" type="xsd:string" minOccurs="0" />
              </xsd:sequence>
              <xsd:attribute name="name" use="required" type="xsd:string" />
              <xsd:attribute name="type" type="xsd:string" />
              <xsd:attribute name="mimetype" type="xsd:string" />
              <xsd:attribute ref="xml:space" />
            </xsd:complexType>
          </xsd:element>
          <xsd:element name="assembly">
            <xsd:complexType>
              <xsd:attribute name="alias" type="xsd:string" />
              <xsd:attribute name="name" type="xsd:string" />
            </xsd:complexType>
          </xsd:element>
          <xsd:element name="data">
            <xsd:complexType>
              <xsd:sequence>
                <xsd:element name="value" type="xsd:string" minOccurs="0" msdata:Ordinal="1" />
                <xsd:element name="comment" type="xsd:string" minOccurs="0" msdata:Ordinal="2" />
              </xsd:sequence>
              <xsd:attribute name="name" type="xsd:string" use="required" msdata:Ordinal="1" />
              <xsd:attribute name="type" type="xsd:string" msdata:Ordinal="3" />
              <xsd:attribute name="mimetype" type="xsd:string" msdata:Ordinal="4" />
              <xsd:attribute ref="xml:space" />
            </xsd:complexType>
          </xsd:element>
          <xsd:element name="resheader">
            <xsd:complexType>
              <xsd:sequence>
                <xsd:element name="value" type="xsd:string" minOccurs="0" msdata:Ordinal="1" />
              </xsd:sequence>
              <xsd:attribute name="name" type="xsd:string" use="required" />
            </xsd:complexType>
          </xsd:element>
        </xsd:choice>
      </xsd:complexType>
    </xsd:element>
  </xsd:schema>
  <resheader name="resmimetype">
    <value>text/microsoft-resx</value>
  </resheader>
  <resheader name="version">
    <value>2.0</value>
  </resheader>
  <resheader name="reader">
    <value>System.Resources.ResXResourceReader, System.Windows.Forms, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089</value>
  </resheader>
  <resheader name="writer">
    <value>System.Resources.ResXResourceWriter, System.Windows.Forms, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089</value>
  </resheader>
`

func FormatSpecifierStringForFormatSpecifier(segment model.FormatSpecifierSegment, index int) string {
	outputOrderIndex := index
	if 0 < segment.SemanticOrderIndex {
		outputOrderIndex = segment.SemanticOrderIndex - 1
	}
	ret := "{" + strconv.Itoa(outputOrderIndex)
	if segment.DataType == model.DataTypeFloat && 0 <= segment.NumberOfDecimals {
		ret += ":F" + strconv.Itoa(segment.NumberOfDecimals)
	}
	ret += "}"
	return ret
}

func xmlEscaped(text string) string {
	var b bytes.Buffer
	xml.EscapeText(&b, []byte(text))
	return b.String()
}

func SanitizedForStringValue(text string) string {
	return xmlEscaped(text)
}

func stringFromSegments(segments []model.Segment) string {
	ret := ""
	for index, segment := range segments {
		switch segment.(type) {
		case model.TextSegment:
			ret += SanitizedForStringValue(segment.(model.TextSegment).Text)
		case model.FormatSpecifierSegment:
			ret += FormatSpecifierStringForFormatSpecifier(segment.(model.FormatSpecifierSegment), index)
		}
	}
	return ret
}

func GetStringsFileContents(set model.TranslationSet, language string) string {
	ret := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"<!--\n" +
		".NET Resource File\n" +
		"Generated by Sanat\n" +
		"Language: " + language + "\n" +
		"-->\n" +
		"<root>\n" +
		xmlHeaderString

	for _, section := range set.Sections {
		if 0 < len(section.Name) {
			sanitizedSectionName := strings.Replace(section.Name, "--", "- -", -1)
			ret += "\n  <!-- ********** " + sanitizedSectionName + " ********** -->\n\n"
		}
		for _, translation := range section.Translations {
			if !translation.IsForPlatform(model.PlatformWindows) {
				continue
			}
			for _, value := range translation.Values {
				if value.Language == language {
					ret += fmt.Sprintf("  <data name=\"%s\" xml:space=\"preserve\">\n",
						xmlEscaped(translation.Key))
					ret += fmt.Sprintf("    <value>%s</value>\n",
						stringFromSegments(value.Segments))
					ret += "  </data>\n"
				}
			}
		}
	}
	ret += "</root>\n"
	return ret
}

func WriteStringsFiles(set model.TranslationSet, outDirPath string) {
	for language, _ := range set.Languages {
		os.MkdirAll(outDirPath, 0777)

		f, err := os.Create(path.Join(outDirPath, "AppResources-"+language+".resx"))
		if err != nil {
			panic(err)
		}

		_, err = f.WriteString(GetStringsFileContents(set, language))
		if err != nil {
			panic(err)
		}
	}
}
