package jbin

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/ir"
)

const artifactVersion = "jbin-v0.1"

type artifact struct {
	Version string
	Program ir.ProgramIR
}

func init() {
	registerGobTypes()
}

func registerGobTypes() {
	gob.Register(ast.IdentifierExpr{})
	gob.Register(ast.NumberExpr{})
	gob.Register(ast.StringExpr{})
	gob.Register(ast.BoolExpr{})
	gob.Register(ast.CollectionLiteralExpr{})
	gob.Register(ast.KeyValueExpr{})
	gob.Register(ast.MemberExpr{})
	gob.Register(ast.IndexExpr{})
	gob.Register(ast.CallExpr{})
	gob.Register(ast.BinaryExpr{})

	gob.Register(ir.VarDeclInstr{})
	gob.Register(ir.AssignInstr{})
	gob.Register(ir.ExprInstr{})
	gob.Register(ir.ReturnInstr{})
	gob.Register(ir.IfInstr{})
	gob.Register(ir.IfBranchIR{})
	gob.Register(ir.WhileInstr{})
	gob.Register(ir.ForInstr{})
	gob.Register(ir.WithinInstr{})
}

// WriteFile encodes a lowered IR program to a .jbin artifact.
func WriteFile(path string, program *ir.ProgramIR) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	payload := artifact{Version: artifactVersion, Program: *program}
	if err := enc.Encode(payload); err != nil {
		return err
	}
	return nil
}

// ReadFile decodes a .jbin artifact into an IR program.
func ReadFile(path string) (*ir.ProgramIR, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	var payload artifact
	if err := dec.Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Version != artifactVersion {
		return nil, fmt.Errorf("unsupported jbin version: %s", payload.Version)
	}
	return &payload.Program, nil
}
