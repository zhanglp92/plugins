package comments

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

// Process ...
func Process(data []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parse(fset, "", data)
	if err != nil {
		return nil, err
	}

	addDoc(file)

	file.Comments = nil
	return out(fset, file)
}

func out(fset *token.FileSet, node interface{}) ([]byte, error) {
	printerMode := printer.UseSpaces
	printerMode |= printer.TabIndent

	printConfig := &printer.Config{Mode: printerMode, Tabwidth: 8}

	var buf bytes.Buffer
	err := printConfig.Fprint(&buf, fset, node)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func parse(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
	parserMode := parser.Mode(0)
	parserMode |= parser.ParseComments
	parserMode |= parser.AllErrors

	return parser.ParseFile(fset, filename, src, parserMode)
}

type updateDocNode struct {
	Doc  *ast.CommentGroup
	Name *ast.Ident
}

func addDoc(f *ast.File) {
	for _, decl := range f.Decls {
		switch node := decl.(type) {
		case *ast.FuncDecl:
			node.Doc = updateDoc(&updateDocNode{Doc: node.Doc, Name: node.Name})
		}
	}
}

func updateDoc(node *updateDocNode) *ast.CommentGroup {
	if node.Name == nil || len(node.Name.Name) <= 0 || (node.Name.Name[0] >= 'a' &&  node.Name.Name[0] <= 'z') {
		return node.Doc
	}

	if node.Doc == nil || len(node.Doc.List) <= 0 {
		return genDoc(node.Name)
	}

	var invalidDocIdx = -1
	for i, c := range node.Doc.List {
		invalidPos, isSpace := getInvalidPos(c)

		if isSpace {
			continue
		}

		if invalidPos < 0 {
			node.Doc.List = append(genDoc(node.Name).List, node.Doc.List...)
			return node.Doc
		}

		invalidDocIdx = i
		if needInsert(invalidPos, c) {
			insertNameToDoc(invalidPos, node.Name.Name, c)
		} else {
			replaceNameToDoc(invalidPos, node.Name.Name, c)
		}
	}

	if invalidDocIdx < 0 {
		return genDoc(node.Name)
	}

	node.Doc.List = node.Doc.List[invalidDocIdx:]
	return node.Doc
}

func needInsert(invalidPos int, c *ast.Comment) bool {
	a := c.Text[invalidPos]
	return a >= 'a' && a <= 'z'
}

func getInvalidPos(c *ast.Comment) (int, bool) {
	var (
		status  = 0
		lineCnt = 0
	)

	for i, h := range c.Text {
		if status == 0 {
			if h == ' ' {
				continue
			} else {
				status++
			}
		}

		if status == 1 {
			if h == '/' {
				lineCnt++
				continue
			} else if lineCnt < 2 {
				break
			} else {
				status++
			}
		}

		if status == 2 {
			if h == ' ' {
				continue
			} else {
				return i, false
			}
		}
	}

	if (status == 1 && lineCnt >= 2) || status == 2 {
		return -1, true
	}
	return -1, false
}

func replaceNameToDoc(pos int, name string, c *ast.Comment) {
	spaceIdx := strings.Index(strings.TrimSpace(c.Text[pos:]), " ")
	if spaceIdx < 0 {
		c.Text = strings.TrimSpace(c.Text[:pos]) + " " + name + " ..."
	} else {
		c.Text = strings.TrimSpace(c.Text[:pos]) + " " + name + c.Text[pos+spaceIdx:]
	}
}

func insertNameToDoc(pos int, name string, c *ast.Comment) {
	c.Text = strings.TrimSpace(c.Text[:pos]) + " " + name + " " + c.Text[pos:]
}

func genDoc(name *ast.Ident) *ast.CommentGroup {
	return &ast.CommentGroup{
		List: []*ast.Comment{{Text: "// " + name.Name + " ..."}},
	}
}
