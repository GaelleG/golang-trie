package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Node struct {
	Value    string
	Children []*Node
}

func (node *Node) String() string {
	return fmt.Sprintf("%s %+v", node.Value, node.Children)
}

func main() {
	words := getWords()

	var root = Node{"", []*Node{}}

	for _, word := range words {
		currentNode := &root
		parsed := word
		for len(parsed) > 0 {
			value := parsed[0:1]
			err, child := getChild(currentNode, value)
			if err != nil {
				newChild := Node{value, []*Node{}}
				child = &newChild
				currentNode.Children = append(currentNode.Children, child)
			}
			currentNode = child
			parsed = parsed[1:len(parsed)]
		}
		endNode := Node{".", []*Node{}}
		currentNode.Children = append(currentNode.Children, &endNode)
	}

	// mergeNodes(&root)

	fmt.Printf(root.String() + "\n")

	leafs := getLeafs(&root)
	log.Print("sheets: ", len(leafs))
	log.Print("words: ", len(words))

	for true {
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')

		prefix := ""
		if len(answer) > 1 {
			prefix = answer[0:len(answer)-1]
		}
		log.Print("[" + prefix + "]")
		err, prefixNode := findPrefixNode(&root, prefix)
		if err != nil {
			log.Print(err)
		} else {
			nodeWords := []string{}
			log.Print(getNodeWords(prefixNode, prefix[0:len(prefix)-1], nodeWords))
		}
	}
}

func getWords() []string {
	data, err := ioutil.ReadFile("words")
	if err != nil {
		log.Fatal("Error reading words file: ", err)
	}
	return strings.Split(string(data), "\n")
}

func getChild(node *Node, value string) (error, *Node) {
	for _, child := range node.Children {
		if child.Value == value {
			return nil, child
		}
	}
	return errors.New("No child exists for value: " + value), node
}

func mergeNodes(node *Node) {
	if node == nil {
		return
	}
	if len(node.Children) == 0 {
		return
	}
	if len(node.Children) == 1 {
		node.Value = node.Value + node.Children[0].Value
		node.Children = node.Children[0].Children
		mergeNodes(node)
	} else {
		for _, child := range node.Children {
			mergeNodes(child)
		}
	}
}

func getLeafs(node *Node) []Node {
	leafs := []Node{}
	if len(node.Children) == 0 {
		return append(leafs, *node)
	} else {
		for _, child := range node.Children {
			leafs = append(leafs, getLeafs(child)...)
		}
		return leafs
	}
}

func findPrefixNode(node *Node, prefix string) (error, *Node) {
	word := ""
	currentNode := node
	parsed := prefix
	for len(parsed) > 0 {
		value := parsed[0:1]
		err, child := getChild(currentNode, value)
		if err != nil {
			return errors.New("Prefix [" + prefix + "] not found. Stopped at [" + word + "]."), node
		} else {
			word = word + value
			currentNode = child
			parsed = parsed[1:len(parsed)]
		}
	}
	return nil, currentNode
}

func getNodeWords(node *Node, prefix string, words []string) []string {
	if len(node.Children) == 0 {
		return append(words, prefix+node.Value)
	} else {
		words := []string{}
		for _, child := range node.Children {
			words = append(words, getNodeWords(child, prefix+node.Value, words)...)
		}
		return words
	}
}
