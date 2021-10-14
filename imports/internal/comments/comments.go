package comments

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
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
			node.Doc = addFuncDoc(&updateDocNode{Doc: node.Doc, Name: node.Name})
		}
	}
}

func addFuncDoc(node *updateDocNode) *ast.CommentGroup {
	if !needDoc(node.Doc, node.Name) {
		return node.Doc
	}
	return genDoc(node.Name)
}

func genDoc(name *ast.Ident) *ast.CommentGroup {
	return &ast.CommentGroup{
		List: []*ast.Comment{{Text: "// " + name.Name + " ..."}},
	}
}

func needDoc(doc *ast.CommentGroup, name *ast.Ident) bool {
	if hasDoc(doc) {
		return false
	}

	if name == nil || len(name.Name) <= 0 {
		return false
	}

	if c := name.Name[0]; c < 'A' || c > 'Z' {
		return false
	}

	return true
}

func hasDoc(doc *ast.CommentGroup) bool {
	return doc != nil && len(doc.List) > 0
}
