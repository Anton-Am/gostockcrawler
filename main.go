package main

import (
	"golang.org/x/net/html"
	"github.com/joho/godotenv"
	"net/http"
	"strings"
	"os"
	"bytes"
	"io"
	"time"
	"strconv"
)

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := getAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}
	return false
}

func getElementById(n *html.Node, id string) *html.Node {
	if checkId(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := getElementById(c, id)
		if result != nil {
			return result
		}
	}

	return nil
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func saveData(body string) {
	pwd, err := os.Getwd()
	check(err)

	var filename = pwd + "/stocks-data/" + time.Now().Format("02-01-2006") + ".html"

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	check(err)
	_, err = f.WriteString("<div style='display: table; margin: 0 auto;'>" + body + "<hr/><hr/><hr/></div>")
	err = f.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Load config data
	err := godotenv.Load()
	check(err)

	// Declare http client
	client := &http.Client{}

	// Declare post data
	PostData := strings.NewReader("")

	// Declare HTTP Method and Url
	req, err := http.NewRequest("GET", os.Getenv("QUOTE_SOURCE"), PostData)
	check(err)

	// Set cookie
	req.Header.Set("Cookie", os.Getenv("COOKIE"))
	// Make response
	resp, err := client.Do(req)
	check(err)
	// Parse data
	doc, err := html.Parse(resp.Body)
	check(err)

	// Find an element
	elementData := getElementById(doc, os.Getenv("BLOCK_ID"))

	// First table data
	childFlag, err := strconv.ParseBool(os.Getenv("FIRST_CHILD_DATA"))
	check(err)
	if (childFlag) {
		firstTableData := renderNode(elementData.FirstChild.NextSibling)
		saveData(firstTableData)
	} else {
		saveData(renderNode(elementData))
	}

}
