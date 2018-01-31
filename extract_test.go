package extract

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestHTML(t *testing.T) {
	table := []struct {
		query  string
		expect string
		html   string
	}{
		{"", "", ""},
		{"h1", "<h1>Hello World</h1>", "<h1>Hello World</h1>"},
		{"h1", "<h1>Hello World</h1>", "<p>Foo <h1>Hello World</h1> Bar</p>"},
		{"a.bar", "<a class=\"bar\">2</a>", "<a class=\"foo\">1</a><a class=\"bar\">2</a><a class=\"baz\">3</a>"},
		{"#b", "<a id=\"b\">2</a>", "<a id=\"a\">1</a><a id=\"b\">2</a><a id=\"c\">3</a>"},
	}

	for i, entry := range table {
		e, err := HTML(strings.NewReader(entry.html), entry.query)
		if err != nil {
			t.Errorf("Expected no error for entry %d", i)
			return
		}
		if e != entry.expect {
			t.Errorf("Expected '%s' but got '%s'", entry.expect, e)
			return
		}
	}
}

func TestNodeToString(t *testing.T) {
	table := []struct {
		html string
	}{
		{""},
		{"<h1>Hello World</h1>"},
		{"<article><h2>Foo</h2><p>Bar <a href=\"#\">Baz</a></p></article>"},
	}

	for i, entry := range table {
		n, err := html.Parse(strings.NewReader(entry.html))
		if err != nil {
			t.Errorf("Expected no error for entry %d", i)
			return
		}

		s := nodeToString(n)
		if s != entry.html {
			t.Errorf("Expected '%s' but got '%s'", entry.html, s)
			return
		}
	}
}

func TestParseQuery(t *testing.T) {
	q := parseQuery("h1")
	if q.name != "h1" {
		t.Errorf("Expected name '%s' but got '%s'", "h1", q.name)
		return
	}

	q = parseQuery("h2.foo.bar")
	if q.name != "h2" {
		t.Errorf("Expected name '%s' but got '%s'", "h2", q.name)
		return
	}
	if len(q.classes) != 2 || !stringSliceContains(q.classes, "foo") || !stringSliceContains(q.classes, "bar") {
		t.Errorf("Expected h2 to have two classes 'foo' and 'bar'")
		return
	}

	q = parseQuery("#foo.bar")
	if q.id != "foo" {
		t.Errorf("Expected id '%s' but got '%s'", "foo", q.id)
		return
	}
	if len(q.classes) != 1 || !stringSliceContains(q.classes, "bar") {
		t.Errorf("Expected id foo to have one class 'bar'")
		return
	}
}
