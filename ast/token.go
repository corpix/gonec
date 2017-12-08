package ast

import "github.com/corpix/yoptec/pos"

type Token struct {
	pos.PosImpl // StmtImpl provide Pos() function.
	Tok         int
	Lit         string
}
