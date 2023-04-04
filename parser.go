package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type result struct {
	header     header
	funcFrames []funcFrame
}

type header struct {
	packageName string
	types       []string
	imports     []string
}

type funcFrame struct {
	name       string
	statements []statement
}

type statement struct {
}

type parser struct {
	file    *ast.File
	builder *builder

	result result
}

func newParser(file *ast.File) *parser {
	return &parser{
		file:    file,
		builder: newBuilder(),
	}
}

func (p *parser) Run() error {
	p.result.header = p.parseHeader()

	//
	//	for _, v := range p.file.Decls {
	//		if fd, ok := v.(*ast.FuncDecl); ok {
	//			p.parseNode(fd)
	//		} else {
	//			fmt.Println("GOT not func")
	//		}
	//	}

	return nil
}

func (p *parser) Result() result {
	return p.result
}

func (p *parser) parseHeader() header {
	h := header{
		packageName: p.file.Name.Name,
	}

	for _, v := range p.file.Decls {
		if gd, ok := v.(*ast.GenDecl); ok {
			fmt.Println(gd.Tok)
			switch gd.Tok {
			case token.IMPORT:
				for _, s := range gd.Specs {
					h.imports = append(h.imports, p.parseNode(s))
				}
			case token.TYPE:
				for _, s := range gd.Specs {
					h.types = append(h.types, p.parseNode(s))
				}
			case token.CONST:
				// TODO: IMPLEMENT ME
			case token.VAR:
				// TODO: IMPLEMENT ME
			default:
				fmt.Println("UNKNOWN token ", gd.Tok)
			}
		} else {
			fmt.Printf("%#v\n", v)
		}

	}

	return h
}

func (p *parser) parseNode(node ast.Node) string {
	var output strings.Builder

	switch n := node.(type) {
	case *ast.FuncDecl:
		p.parseFunc(n)
	case *ast.Package:
		fmt.Printf("%#v\n", n.Name)

		output.WriteString(n.Name)
	case *ast.ImportSpec:
		output.WriteString(n.Path.Value)
	case *ast.TypeSpec:
		fmt.Println(n.Type)
		output.WriteString(n.Name.Name + p.parseNode(n.Type))
	case *ast.StructType:
		output.WriteString(n.Fields.List[0].Names[0].Name)
	default:
		fmt.Printf("UNKNOWN node type %#v\n", n)
	}

	return output.String()
}

func (p *parser) parseFunc(funcDecl *ast.FuncDecl) {
	name := funcDecl.Name

	fmt.Println("GOT func with name:", name)

	p.builder.WriteString("func " + name.Name + "\n")

	for _, v := range funcDecl.Body.List {
		ast.Inspect(v, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.Ident:
				fmt.Printf("IDENT: %#v\n", n)

				if n.Obj != nil {
					fmt.Println(n.Obj)
				}
			case *ast.AssignStmt:
				p.parseAssignStmt(n)

				return false
			}

			return true
		})
	}
}

func (p *parser) parseAssignStmt(stmt *ast.AssignStmt) {
	fmt.Println("AssignStmt: ", stmt)

	for _, v := range stmt.Lhs {
		switch n := v.(type) {
		case *ast.Ident:
			fmt.Printf("LEFT Ident %#v\n", n)

			p.builder.WriteString(n.Name)
		}
	}

	p.builder.WriteString(stmt.Tok.String())

	for _, v := range stmt.Rhs {
		switch n := v.(type) {
		case *ast.Ident:
			fmt.Printf("RIGHT Ident %#v\n", n)

			p.builder.WriteString(n.Name)
		case *ast.CallExpr:
			p.parseCallExpr(n)
		}
	}

	p.builder.WriteString("\n")
}

func (p *parser) parseCallExpr(callExpr *ast.CallExpr) {
	fmt.Printf("RIGHT CallExpr %#v\n", callExpr)

	args := callExpr.Args
	_ = args

	//	p.parseExpr([]ast.Expr{callExpr.Fun})
	p.parseExpr(callExpr.Args)

}

func (p *parser) parseExpr(exprs []ast.Expr) {
	fmt.Println("GOT EXPRS: ", exprs)

	for _, e := range exprs {
		switch n := e.(type) {
		case *ast.BasicLit:
			fmt.Print("EXPR BASICLIT ", n)

		case *ast.SelectorExpr:
			fmt.Println("EXPR SELECTOREXPR ", n.X)

			p.parseExpr([]ast.Expr{n.X})

			p.builder.WriteString(n.Sel.Name)
		case *ast.Ident:
			fmt.Println("EXPR IDENT ", n.Name)

			p.builder.WriteString(n.Name + ".")
		default:
			fmt.Printf("EXPR DEFAULT %#v", n)
		}
	}
}
