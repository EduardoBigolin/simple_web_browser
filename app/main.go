package main

import (
	"fmt"
	"regexp"
	"strings"
)

// HTML Structs
type Node struct {
	children []*Node
	nodeType NodeType
}

type NodeType struct {
	Text  string
	Elemt ElementData
}

type ElementData struct {
	tag_name   string
	attributes AttrMap
}

type AttrMap = map[string][]string

// CSS Structs
type StyleSheets struct {
	selector string
	styles   []Style
}

type Style struct {
	property string
	value    string
}

// Will create a new node with the text
func text(data string) *Node {
	return &Node{nodeType: NodeType{Text: data}}
}

// Will create a new node with the tag name and the attributes
func elem(name string, attrs AttrMap, children ...*Node) *Node {
	return &Node{children: children, nodeType: NodeType{Elemt: ElementData{tag_name: name, attributes: attrs}}}
}

// This function will print the html without the styles for testing
func PrintHtml(node *Node, ident string) {
	if node.nodeType.Text != "" {
		fmt.Println(ident + node.nodeType.Text)
		return
	}
	for _, child := range node.children {
		PrintHtml(child, ident+" ")
	}
}

// Will parse the elementAttr string and return a AttrMap
func parseToClass(elementAttr string) AttrMap {
	attr := make(AttrMap)
	re := regexp.MustCompile(`(\S+)\s*=\s*["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(elementAttr, -1)

	for _, match := range matches {
		attrName := match[1]
		attrValue := match[2]

		switch attrName {
		case "class":
			classes := strings.Split(attrValue, " ")
			attr["class"] = classes
		case "id":
			attr["id"] = strings.Split(attrValue, " ")
		}
	}

	return attr
}

// Will parse the html string and return a Node
func parseHTMLString(htmlString string) *Node {
	var root *Node
	var stack []*Node

	for _, token := range strings.Split(htmlString, "<") {
		parts := strings.SplitN(token, ">", 2)
		tagName := parts[0]

		if tagName != "" {
			if tagName[0] != '/' {
				parseClass := AttrMap{}
				if strings.Contains(tagName, "\"") {
					parseClass = parseToClass(tagName)
				}

				elementAttr := strings.Split(tagName, " ")

				node := elem(elementAttr[0], parseClass)
				if len(stack) > 0 {
					parent := stack[len(stack)-1]

					parent.children = append(parent.children, node)
				} else {
					root = node
				}
				stack = append(stack, node)
			} else {
				if len(stack) == 0 {
					panic("closing tag without opening tag")
				}

				stack = stack[:len(stack)-1]
			}

			if len(parts) == 2 && parts[1] != "" {
				textNode := text(parts[1])
				current := stack[len(stack)-1]
				current.children = append(current.children, textNode)
			}
		}
	}

	return root
}

// Will parse the css string and return a slice of StyleSheets
func parseCSSString(cssString string) []StyleSheets {
	var styles []StyleSheets

	re := regexp.MustCompile(`\s*(\w+)\s*{([^{}]*)}`)
	matches := re.FindAllStringSubmatch(cssString, -1)

	for _, match := range matches {
		selector := match[1]
		propertiesStr := match[2]

		style := StyleSheets{selector: selector, styles: []Style{}}

		reProp := strings.Split(propertiesStr, ";")

		for _, prop := range reProp {
			propSplit := strings.Split(prop, ":")

			if len(propSplit) == 2 {
				style.styles = append(style.styles, Style{property: propSplit[0], value: propSplit[1]})
			}
		}

		styles = append(styles, style)
	}
	return styles
}

// Verify if the node has the style
func match(node *Node, selector string) bool {
	// TODO: Verify if the node has the class and id
	// TODO: Create a order to verify the selectors

	if node.nodeType.Text != "" {
		return false
	}

	if node.nodeType.Elemt.tag_name == selector {
		return true
	}

	for _, child := range node.children {
		if match(child, selector) {
			return true
		}
	}

	return false
}

// This function will verify all selectors for a node and run the match function to verify if the node has the style
func VerifyAllStyles(style []StyleSheets, node *Node) []Style {
	for _, child := range style {
		if match(node, child.selector) {
			return child.styles
		}
	}
	return nil
}

// This function will print the html with the styles for testing
func printHtmlWithStyle(node *Node, ident string, styles []StyleSheets) {
	if node.nodeType.Text != "" {
		fmt.Println(ident + " " + node.nodeType.Text)
		return
	}
	s := VerifyAllStyles(styles, node)

	if s != nil {
		fmt.Print(ident + "<" + node.nodeType.Elemt.tag_name + " style=\"")
		for _, child := range s {
			fmt.Print(child.property + ":" + child.value + "; ")
		}
		fmt.Print("\" ")
		fmt.Print("class=\"")
		for _, child := range node.nodeType.Elemt.attributes["class"] {
			fmt.Print(child + " ")
		}
		fmt.Print("\"")

		if len(node.nodeType.Elemt.attributes["id"]) > 0 {
			fmt.Print("id=\"")
			for _, child := range node.nodeType.Elemt.attributes["id"] {
				fmt.Print(child + " ")
			}
			fmt.Print("\"")
		}

		fmt.Println(">")

	} else {
		fmt.Println(ident + "<" + node.nodeType.Elemt.tag_name + ">")
	}

	for _, child := range node.children {
		printHtmlWithStyle(child, ident+" ", styles)
	}

	fmt.Println(ident + "</" + node.nodeType.Elemt.tag_name + ">")
}

func main() {
	html := `<html><head><title>My test page</title></head><body><h1 id="oi" class="title page h1">Hello world!</h1><div class="oi">OI</div></body></html>`
	css := `h1 { color: #ffffff; } h2 { color: #000000; } div { background-color: #000000; width: 100px; height: 100px; } .oi { color: #000000;}`

	// TODO: add go routines

	htmlParse := parseHTMLString(html)
	cssParse := parseCSSString(css)

	printHtmlWithStyle(htmlParse, "", cssParse)
}
