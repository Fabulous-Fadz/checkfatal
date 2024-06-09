package checkfatal

import (
	"go/ast"
	"go/token"
)

var funcDecls map[string]*ast.FuncDecl

func callsOsExit(call *ast.CallExpr, visited map[*ast.FuncDecl]bool) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if ok && selector.X.(*ast.Ident).Name == "os" && selector.Sel.Name == "Exit" {
		return true
	}

	// Check all function calls excpet those already visited.
	for _, arg := range call.Args {
		if funcExpr, ok := arg.(*ast.CallExpr); ok {
			funcDecl := funcExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name
			if !visited[lookupFuncDecl(funcDecl)] {
				visited[lookupFuncDecl(funcDecl)] = true {
					if mightCallOsExit(funcExpr, visited) {
						return true
					}
				}
			}
		}
	}

	return false
}

func lookupFuncDecl(name string) *ast.FuncDecl {
	fd, _ := funcDecls[name]
	return fd
}

func checkForFollowingLine(node *ast.CallExpr, fset *token.FileSet) (bool, int) {
	callLine := fset.Position(node.Pos()).callLine
	if node.Parent.(*ast.BlockStmt).List != nil {
		nextLine := node.Parent.(*ast.BlockStmt).List[node.Pos()+1]
		if nextLine != nil {
			return true, fset.Position(nextLine.Pos()).Line - callLine
		}
	}

	return false, 0
}

func checkForOsExit(node ast.Node, fset *token.FileSet) bool {
	switch n := node.(type) {
	case *ast.CallExpr:
		visited := make(map[*ast.FuncDecl]bool)
		if callsOsExit(n, visited) {
			hasLineAfter, lineDiff := checkForFollowingLine(n, fset)
			if hasLineAfter && lineDiff > 1 {
				// Found a call that calls os.Exit that has another line following it.
				return true
			}
		}
	case *ast.BlockStmt, *ast.ForStmt, *ast.IfStmt:
		for _, child := range n.List {
			if checkForOsExit(child, fset) {
				return true
			}
		}
	}

	return false
}