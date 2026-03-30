package emiters

import (
	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
)

type ProgramEmitter interface {
	Emit(ir *generator.ProgramIR) (string, exception.IException)
}
