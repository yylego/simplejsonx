package sure

import (
	"testing"

	"github.com/yylego/runpath"
	"github.com/yylego/sure"
	"github.com/yylego/sure/sure_pkg_gen"
)

func TestGen(t *testing.T) {
	pkgPath := "github.com/yylego/simplejsonx"
	t.Log(pkgPath)
	srcRoot := runpath.PARENT.Up(1)
	t.Log(srcRoot)
	outRoot := runpath.PARENT.Path()
	t.Log(outRoot)
	sure_pkg_gen.GenerateSurePackageFiles(t, sure_pkg_gen.NewSurePackageConfig(srcRoot, sure.SOFT, pkgPath).WithOutputRoot(outRoot).WithNewPkgName("simplejsons"))
	sure_pkg_gen.GenerateSurePackageFiles(t, sure_pkg_gen.NewSurePackageConfig(srcRoot, sure.MUST, pkgPath).WithOutputRoot(outRoot).WithNewPkgName("simplejsonm"))
	sure_pkg_gen.GenerateSurePackageFiles(t, sure_pkg_gen.NewSurePackageConfig(srcRoot, sure.OMIT, pkgPath).WithOutputRoot(outRoot).WithNewPkgName("simplejsono"))
}
