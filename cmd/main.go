package main

func main() {
	fset := token.NewFileSet()

	funcDecls := map[string]*ast.FuncDecl {
		// Add function declarations here...
	}

	// Analyze the AST
	for _, file := range fset.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if checkfatal.checkForOsExit(funcDecl.Body, fset) {
					// We have a winner. Prize is an error message.
				}
			}
		}
	}
}