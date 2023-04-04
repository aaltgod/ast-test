package main

import (
	"fmt"
	goParser "go/parser"
	"go/token"
	"log"
)

type manager struct {
	parser  *parser
	builder *builder
}

func newManager(parser *parser, builder *builder) *manager {
	return &manager{
		parser:  parser,
		builder: builder,
	}
}

func (m *manager) Run() error {
	if err := m.parser.Run(); err != nil {
		return err
	}

	parserRusult := m.parser.Result()
	fmt.Println(parserRusult)

	if err := m.builder.RenderHeader(
		builderHeader{
			PackageName: parserRusult.header.packageName,
			Imports:     parserRusult.header.imports,
			Types:       parserRusult.header.types,
		},
	); err != nil {
		return err
	}

	fmt.Println(m.builder.Render())

	return nil
}

func main() {
	fileSet := token.NewFileSet()
	fileName := "example.go"

	f, err := goParser.ParseFile(fileSet, fileName, nil, goParser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	man := newManager(
		newParser(f),
		newBuilder(),
	)

	man.Run()

	//    err = ast.Print(fileSet, f)
	//    if err != nil {
	//        log.Fatal(err)
	//    }
}
