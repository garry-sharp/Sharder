package crypt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func ReadLang() string {
	lang := ""
	fmt.Print("Please select your language: ")
	fmt.Scan(&lang)
	return lang
}

type tree struct {
	leaf     bool
	children []tree
	value    string
}

func printRecursive(t *tree, depth int, prefix string) string {
	out := fmt.Sprintf("%s - %s (%v)\n", strings.Repeat(prefix, depth), t.value, t.leaf)
	if t.leaf {
		return out
	} else {
		for _, child := range t.children {
			out += printRecursive(&child, depth+1, prefix)
		}
	}
	return out
}

func (t *tree) String() string {
	return printRecursive(t, 0, "  ")
}

func (t1 *tree) Equal(t2 *tree) bool {
	if t1.leaf && t2.leaf {
		return t1.value == t2.value
	} else if !t1.leaf && !t2.leaf {
		if len(t1.children) != len(t2.children) {
			return false
		} else {
			childrenEquality := true
			for i := 0; i < len(t1.children); i++ {
				childrenEquality = childrenEquality && t1.children[i].Equal(&t2.children[i])
			}
			return childrenEquality && t1.value == t2.value
		}
	} else {
		return false
	}
}

// Used for autocomplete
func buildWordListTree(words []string) *tree {
	t := new(tree)
	for _, word := range words {
		tt := t
		for i, letter := range word {
			found := -1
			for nodeIndex, node := range tt.children {
				if node.value == string(letter) {
					found = nodeIndex
				}
			}
			if found != -1 {
				tt = &tt.children[found]
			} else {
				tt.children = append(tt.children, tree{value: string(letter)})
				tt = &tt.children[len(tt.children)-1]
			}

			if i == len(word)-1 {
				tt.children = []tree{{leaf: true, value: word}}
			}
		}
	}
	return t
}

func getOptions(word string, t *tree) []string {
	nodes := []tree{}
	ops := []string{}

	for _, letter := range word {
		found := -1
		for i, node := range t.children {
			if node.value == string(letter) {
				found = i
			}
		}
		if found == -1 {
			return ops
		} else {
			t = &t.children[found]
		}
	}

	nodes = getAllChildren(t)
	for _, node := range nodes {
		ops = append(ops, node.value)
	}

	return ops
}

func getAllChildren(t *tree) []tree {
	nodes := []tree{}
	if t.leaf {
		return []tree{*t}
	}
	for _, node := range t.children {
		nodes = append(nodes, getAllChildren(&node)...)
	}
	return nodes
}

func ReadMnemonic(wordList []string) (string, error) {
	l := 0
	fmt.Print("How long is your mnemonic: ")
	fmt.Scan(&l)

	if l == 12 || l == 15 || l == 18 || l == 21 || l == 24 {
		//fmt.Printf("Please enter your %d word mnemonic (esc to exit): ", l)
		wordTree := buildWordListTree(wordList)

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		reader := bufio.NewReader(os.Stdin)
		words := []string{}
		word := ""
		clearLength := 5
		fmt.Println("Please enter below your mnemonic, space separated\r\nSome words may autocomplete as you go (esc to exit)")
		for len(words) < l {

			fmt.Printf("\r: %s", strings.Repeat(" ", clearLength))
			fmt.Printf("\r: %s", strings.Join(append(words, word), " "))
			char, _ := reader.ReadByte()
			letter := string(char)
			if letter[0] == 127 {
				//handle backspace
				if len(word) > 0 {
					clearLength--
					word = word[:len(word)-1]
				} else {
					if len(words) > 0 {
						clearLength--
						word = words[len(words)-1]
						words = words[:len(words)-1]
					}
				}
			} else if letter[0] == 32 {
				//space
				clearLength++
				ops := getOptions(word, wordTree)
				for _, op := range ops {
					if word == op {
						words = append(words, word)
						word = ""
						if len(words) == l {
							break
						}
					}
				}
			} else if letter[0] == 27 {
				return "", fmt.Errorf("Exited")
			} else if letter[0] >= 97 && letter[0] <= 122 {
				clearLength++
				word += letter
				ops := getOptions(word, wordTree)
				if len(ops) == 1 {
					clearLength += len(ops[0]) - len(word)
					words = append(words, ops[0])
					word = ""
					if len(words) == l {
						break
					}
				}
			} else {
				//do nothing
			}
		}

		return strings.Join(words, " "), nil
	} else {
		return "", fmt.Errorf("Invalid mnemonic length")
	}
}
