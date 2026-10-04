package sure

import (
	"testing"

	"github.com/yylego/runpath"
	"github.com/yylego/sure"
	"github.com/yylego/sure/sure_pkg_gen"
)

// No -tags=sure_generate here: sure_pkg_gen emits standalone packages that this
// package does not depend on, so a broken artifact leaves this package alone
//
// 这里不带 -tags=sure_generate：sure_pkg_gen 生成的是独立新包，本包不依赖它，
// 产物坏掉也不会连坐本包
//
//go:generate go test -run ^TestGen$
func TestGen(t *testing.T) {
	pkgPath := "github.com/yylego/simplejsonx"
	t.Log(pkgPath)
	srcRoot := runpath.PARENT.Up(1)
	t.Log(srcRoot)
	outRoot := runpath.PARENT.Path()
	t.Log(outRoot)
	sure_pkg_gen.GenerateSurePackageFiles(sure_pkg_gen.NewSurePackageConfig(srcRoot, sure.SOFT, pkgPath).WithOutputRoot(outRoot).WithNewPkgName("simplejsons"))
	sure_pkg_gen.GenerateSurePackageFiles(sure_pkg_gen.NewSurePackageConfig(srcRoot, sure.MUST, pkgPath).WithOutputRoot(outRoot).WithNewPkgName("simplejsonm"))
	sure_pkg_gen.GenerateSurePackageFiles(sure_pkg_gen.NewSurePackageConfig(srcRoot, sure.OMIT, pkgPath).WithOutputRoot(outRoot).WithNewPkgName("simplejsono"))
}
