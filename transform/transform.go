package transform

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func AddImportToFile(filePath, importPath string) error {
	// Step 1: Read the file content
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Step 2: Parse the file into an AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, src, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return err
	}

	if hasImport(node, importPath) {
		return nil
	}

	// Step 3: Add the import statement
	addImport(node, importPath)

	// Step 4: Write the modified AST back to the file
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, node)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func hasImport(f *ast.File, importPath string) bool {
	for _, imp := range f.Imports {
		if imp.Path.Value == `"`+importPath+`"` {
			return true
		}
	}
	return false
}

func addImport(f *ast.File, importPath string) {
	newImport := &ast.ImportSpec{
		Name: ast.NewIdent("_"),
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"` + importPath + `"`,
		},
	}

	f.Imports = append(f.Imports, newImport)

	// Add the import to the declaration
	if f.Decls == nil {
		f.Decls = []ast.Decl{}
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}

		genDecl.Specs = append(genDecl.Specs, newImport)
		return
	}

	// If no import declaration exists, create one
	f.Decls = append([]ast.Decl{&ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{newImport},
	}}, f.Decls...)
}
