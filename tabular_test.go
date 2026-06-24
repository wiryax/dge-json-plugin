package jp

import (
	"reflect"
	"testing"

	dge "github.com/wiryax/direct-graph-engine"
)

func TestJsonToTabular(t *testing.T) {
	tests := []struct {
		title,
		payload string
		wantErr  bool
		expected func() dge.Tabular
	}{
		{
			title:   "single object",
			payload: `{"a":"a","b":"b","c":"c"}`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular()
				e.AddOrSetColumn("a", dge.ParseVariable([]byte("a")))
				e.AddOrSetColumn("b", dge.ParseVariable([]byte("b")))
				e.AddOrSetColumn("c", dge.ParseVariable([]byte("c")))
				return *e
			},
		}, {
			title:   "invalid JSON",
			payload: `{"a"}`,
			wantErr: true,
			expected: func() dge.Tabular {
				return dge.Tabular{}
			},
		}, {
			title:   "object array",
			payload: `{"a": ["1","2"]}`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular()
				e.AddOrSetColumn("a", dge.ParseVariable([]byte("1")), dge.ParseVariable([]byte("2")))
				return *e
			},
		}, {
			title:   "object multi dimensional array with object",
			payload: `{"a": [{"b": "1"},{"b": "2"}]}`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular()
				e.AddOrSetColumn("ab", dge.ParseVariable([]byte("1")), dge.ParseVariable([]byte("2")))
				return *e
			},
		}, {
			title:   "nested object",
			payload: `{"a": {"a":"1", "b": "2"}, "b" : "3"}`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular()
				e.AddOrSetColumn("aa", dge.ParseVariable([]byte("1")))
				e.AddOrSetColumn("ab", dge.ParseVariable([]byte("2")))
				e.AddOrSetColumn("b", dge.ParseVariable([]byte("3")))
				return *e
			},
		}, {
			title:   "inconsistency structure",
			payload: `[{"a":"1","b":"2"},"3", "4"]`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular()
				e.AddOrSetColumn("a", dge.ParseVariable([]byte("1")), dge.ParseVariable([]byte("1")))
				e.AddOrSetColumn("b", dge.ParseVariable([]byte("2")), dge.ParseVariable([]byte("2")))
				e.AddOrSetColumn("", dge.ParseVariable([]byte("3")), dge.ParseVariable([]byte("4")))
				return *e
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			e := test.expected()
			result, err := parseJsonToTabular([]byte(test.payload))
			if test.wantErr != (err != nil) {
				t.Fatalf("unexpected err result. want %v, got %v. err: %v", test.wantErr, (err != nil), err)
			}

			if !reflect.DeepEqual(e, result) {
				t.Errorf("unexpected result. want %v, got %v", e.String(), result.String())
			}
		})
	}
}

func BenchmarkLargeJSON(b *testing.B) {
	payload := []byte(`{"web-app": {
  "servlet": [   
    {
      "servlet-name": "cofaxCDS",
      "servlet-class": "org.cofax.cds.CDSServlet",
      "init-param": {
        "configGlossary:installationAt": "Philadelphia, PA",
        "configGlossary:adminEmail": "ksm@pobox.com",
        "configGlossary:poweredBy": "Cofax",
        "configGlossary:poweredByIcon": "/images/cofax.gif",
        "configGlossary:staticPath": "/content/static",
        "templateProcessorClass": "org.cofax.WysiwygTemplate",
        "templateLoaderClass": "org.cofax.FilesTemplateLoader",
        "templatePath": "templates",
        "templateOverridePath": "",
        "defaultListTemplate": "listTemplate.htm",
        "defaultFileTemplate": "articleTemplate.htm",
        "useJSP": false,
        "jspListTemplate": "listTemplate.jsp",
        "jspFileTemplate": "articleTemplate.jsp",
        "cachePackageTagsTrack": 200,
        "cachePackageTagsStore": 200,
        "cachePackageTagsRefresh": 60,
        "cacheTemplatesTrack": 100,
        "cacheTemplatesStore": 50,
        "cacheTemplatesRefresh": 15,
        "cachePagesTrack": 200,
        "cachePagesStore": 100,
        "cachePagesRefresh": 10,
        "cachePagesDirtyRead": 10,
        "searchEngineListTemplate": "forSearchEnginesList.htm",
        "searchEngineFileTemplate": "forSearchEngines.htm",
        "searchEngineRobotsDb": "WEB-INF/robots.db",
        "useDataStore": true,
        "dataStoreClass": "org.cofax.SqlDataStore",
        "redirectionClass": "org.cofax.SqlRedirection",
        "dataStoreName": "cofax",
        "dataStoreDriver": "com.microsoft.jdbc.sqlserver.SQLServerDriver",
        "dataStoreUrl": "jdbc:microsoft:sqlserver://LOCALHOST:1433;DatabaseName=goon",
        "dataStoreUser": "sa",
        "dataStorePassword": "dataStoreTestQuery",
        "dataStoreTestQuery": "SET NOCOUNT ON;select test='test';",
        "dataStoreLogFile": "/usr/local/tomcat/logs/datastore.log",
        "dataStoreInitConns": 10,
        "dataStoreMaxConns": 100,
        "dataStoreConnUsageLimit": 100,
        "dataStoreLogLevel": "debug",
        "maxUrlLength": 500}},
    {
      "servlet-name": "cofaxEmail",
      "servlet-class": "org.cofax.cds.EmailServlet",
      "init-param": {
      "mailHost": "mail1",
      "mailHostOverride": "mail2"}},
    {
      "servlet-name": "cofaxAdmin",
      "servlet-class": "org.cofax.cds.AdminServlet"},
 
    {
      "servlet-name": "fileServlet",
      "servlet-class": "org.cofax.cds.FileServlet"},
    {
      "servlet-name": "cofaxTools",
      "servlet-class": "org.cofax.cms.CofaxToolsServlet",
      "init-param": {
        "templatePath": "toolstemplates/",
        "log": 1,
        "logLocation": "/usr/local/tomcat/logs/CofaxTools.log",
        "logMaxSize": "",
        "dataLog": 1,
        "dataLogLocation": "/usr/local/tomcat/logs/dataLog.log",
        "dataLogMaxSize": "",
        "removePageCache": "/content/admin/remove?cache=pages&id=",
        "removeTemplateCache": "/content/admin/remove?cache=templates&id=",
        "fileTransferFolder": "/usr/local/tomcat/webapps/content/fileTransferFolder",
        "lookInContext": 1,
        "adminGroupID": 4,
        "betaServer": true}}],
  "servlet-mapping": {
    "cofaxCDS": "/",
    "cofaxEmail": "/cofaxutil/aemail/*",
    "cofaxAdmin": "/admin/*",
    "fileServlet": "/static/*",
    "cofaxTools": "/tools/*"},
 
  "taglib": {
    "taglib-uri": "cofax.tld",
    "taglib-location": "/WEB-INF/tlds/cofax.tld"}}}`)
	b.ResetTimer()

	for b.Loop() {
		_, err := parseJsonToTabular(payload)
		if err != nil {
			b.Errorf("%v", err)
		}
	}
}
