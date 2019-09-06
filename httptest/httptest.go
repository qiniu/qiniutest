package httptest

import (
	"github.com/qiniu/httptest"
	"github.com/qiniu/httptest/exec"

	_ "github.com/qiniu/qiniutest/httptest/exec/plugin"
)

// ---------------------------------------------------------------------------

type Context struct {
	*httptest.Context
	Ectx *exec.Context
}

func New(t httptest.TestingT) Context {

	ctx := httptest.New(t)
	ectx := exec.New()
	return Context{ctx, ectx}
}

func (p Context) Exec(code string) Context {

	p.Context.Exec(p.Ectx, code)
	return p
}

// ---------------------------------------------------------------------------
