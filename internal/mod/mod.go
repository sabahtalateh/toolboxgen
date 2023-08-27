package mod

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrNotModRoot  = errors.New("directory is not a go module root")
	ErrPkgNotFound = errors.New("package not found")
)

type Module struct {
	Path             string              // module package name
	Dir              string              // module root directory
	AltRoots         []string            // alternative module roots
	VendoredPkgs     map[string]struct{} // list of vendored packages from vendor/modules.txt
	Parent           *Module             // module may be found while lookup, in this case it will have a parent
	ResolvedPkgsInfo map[string]*PkgInfo // package pkgPath inside module -> package pkgName from go source code files
}

type PkgInfo struct {
	Path     string
	Name     string
	Dir      string
	Vendored bool
	StdLib   bool
}

// func (m *Module) GetRequireContainingPkg(pkg string) *string {
//     var require string
//     found := false
//
//     for req, pkgs := range m.Requires {
//         for _, p := range pkgs {
//             if p == pkg {
//                 require = req
//                 found = true
//                 break
//             }
//         }
//     }
//
//     if found {
//         return &require
//     } else {
//         return nil
//     }
// }

// func (m *Module) tryRoots(pInf *PkgInfo, pkgPath string) (*PkgInfo, error) {
//     roots := append([]string{m.Dir}, m.AltRoots...)
//     require := m.GetRequireContainingPkg(pkgPath)
//     found := false
//     for _, root := range roots {
//         if require != nil {
//             pInf.Dir = filepath.Join(root, "vendor", pkgPath)
//             pInf.Vendored = true
//         } else {
//             pInf.Dir = filepath.Join(root, strings.TrimPrefix(pkgPath, m.Path))
//             pInf.Vendored = false
//         }
//         inf, err := os.Stat(pInf.Dir)
//         if os.IsNotExist(err) || !inf.IsDir() {
//             continue
//         }
//         if err != nil {
//             return nil, err
//         }
//         found = true
//         break
//     }
//
//     if !found {
//         return nil, ErrPkgNotFound
//     }
//
//     return pInf, nil
// }
//
// func (m *Module) PkgInfo(pkgPath string) (*PkgInfo, error) {
//     pInf, ok := m.ResolvedPkgsInfo[pkgPath]
//     if ok {
//         return pInf, nil
//     }
//
//     if err := stdLibPkgs.load(); err != nil {
//         return nil, err
//     }
//     pInf = &PkgInfo{Path: pkgPath}
//
//     if p, ok := stdLibPkgs.pkgs[pkgPath]; ok {
//         dirParts := strings.Split(p.GoFiles[0], string(os.PathSeparator))
//         if dirParts[0] == "" {
//             dirParts[0] = string(os.PathSeparator)
//         }
//         dir := filepath.Join(dirParts[:len(dirParts)-1]...)
//         nameParts := strings.Split(pkgPath, "/")
//         pInf.name = nameParts[len(nameParts)-1]
//         pInf.Dir = dir
//         pInf.Vendored = false
//         pInf.StdLib = true
//         m.ResolvedPkgsInfo[pkgPath] = pInf
//         return pInf, nil
//     }
//
//     pInf, err := m.tryRoots(pInf, pkgPath)
//     if err != nil {
//         return nil, err
//     }
//
//     files := token.NewFileSet()
//     pkgs, err := parser.ParseDir(files, pInf.Dir, nil, parser.AllErrors)
//     if err != nil {
//         return nil, fmt.Errorf("failed to parse `%s` dir: %s", pInf.Dir, err)
//     }
//
//     // remove *_test packages
//     if len(pkgs) > 1 {
//         for p := range pkgs {
//             if strings.HasSuffix(p, "_test") {
//                 delete(pkgs, p)
//             }
//             if p == "main" {
//                 delete(pkgs, p)
//             }
//         }
//     }
//
//     // check if multiple packages in one dir. it is ok from ast perspective but golang not allows it
//     if len(pkgs) > 1 {
//         pkgNames := make([]string, 0)
//         for _, p := range pkgs {
//             pkgNames = append(pkgNames, p.name)
//         }
//
//         return nil, fmt.Errorf(
//             "multiple packages [%s] was found at `%s` dir should be only one",
//             strings.Join(pkgNames, ", "), pInf.Dir,
//         )
//     }
//
//     if len(pkgs) == 0 {
//         _, dir := filepath.Split(pInf.Dir)
//         pInf.name = dir
//     } else {
//         for _, p := range pkgs {
//             pInf.name = p.name
//             break
//         }
//     }
//
//     m.ResolvedPkgsInfo[pkgPath] = pInf
//
//     return pInf, err
// }

// // ResolveImportAlias returns package pkgName (from package directory), package alias (from go file), error
// func (m *Module) ResolveImportAlias(impSpec *ast.ImportSpec) (string, error) {
//     pkgPath := UnquoteImport(impSpec.Path.Value)
//     pInf, err := m.PkgInfo(pkgPath)
//     if err != nil {
//         return "", err
//     }
//
//     var impAlias string
//     if impSpec.name != nil {
//         impAlias = impSpec.name.name
//     } else {
//         parts := strings.Split(pkgPath, "/")
//         if parts[len(parts)-1] != pInf.name {
//             impAlias = pInf.name
//         } else {
//             impAlias = parts[len(parts)-1]
//         }
//     }
//
//     return impAlias, nil
// }

func LookupDir(dir string, withVendored bool) (*Module, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("dir `%s` not exists", dir)
	}

	mod, err := FromGoMod(filepath.Join(dir, "go.mod"), withVendored)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func FromGoMod(goModPath string, withVendored bool) (*Module, error) {
	goModBytes, err := os.ReadFile(goModPath)
	if os.IsNotExist(err) {
		return nil, ErrNotModRoot
	}
	if err != nil {
		return nil, err
	}

	goMod, err := modfile.Parse("go.mod", goModBytes, nil)
	if err != nil {
		return nil, err
	}

	mod := Module{
		Path:             goMod.Module.Mod.Path,
		Dir:              filepath.Dir(goModPath),
		ResolvedPkgsInfo: make(map[string]*PkgInfo),
	}

	if withVendored && len(goMod.Require) > 0 {
		if err = validateModVendor(&mod); err != nil {
			return nil, err
		}

		if mod.VendoredPkgs, err = vendoredPkgsFromModulesTxt(mod.Dir, "vendor", "modules.txt"); err != nil {
			return nil, err
		}
	}

	return &mod, nil
}

func validateModVendor(mod *Module) error {
	vendor := filepath.Join(mod.Dir, "vendor")
	inf, err := os.Stat(vendor)
	if os.IsNotExist(err) {
		return fmt.Errorf("`vendor` dir not exists at `%s` run `go mod vendor` at dir", mod.Dir)
	}
	if !inf.IsDir() {
		return fmt.Errorf("`%s` is not a directory", vendor)
	}
	modulesFile := filepath.Join(vendor, "modules.txt")
	inf, err = os.Stat(modulesFile)
	if os.IsNotExist(err) {
		return fmt.Errorf(
			"`vendor` dir `%s` not contains at modules.txt. run `go mod vendor` at `%s`",
			modulesFile, mod.Dir,
		)
	}
	if inf.IsDir() {
		return fmt.Errorf("`%s` is directory. should be a file", modulesFile)
	}

	return nil
}

func vendoredPkgsFromModulesTxt(pathParts ...string) (map[string]struct{}, error) {
	vendorData, err := parseModulesTxt(pathParts...)
	if err != nil {
		return nil, err
	}

	vendors := make(map[string]struct{})
	for pkg, d := range vendorData {
		vendors[pkg] = struct{}{}
		for _, pkg2 := range d {
			vendors[pkg2] = struct{}{}
		}
	}

	return vendors, nil
}

func parseModulesTxt(pathParts ...string) (map[string][]string, error) {
	modulesTxtPath := filepath.Join(pathParts...)
	res := make(map[string][]string)

	file, err := os.Open(modulesTxtPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	mod := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			parts := strings.Split(line, " ")
			mod = parts[1]
			res[mod] = make([]string, 0)
		} else if strings.HasPrefix(line, "## ") {
			continue
		} else {
			res[mod] = append(res[mod], line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
