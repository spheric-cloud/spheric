// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"tools/utils"

	"github.com/spf13/pflag"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
)

type VersionPackage struct {
	Prefix  string
	Name    string
	Version string
}

func ParseVersionPackage(s string) (*VersionPackage, error) {
	version := path.Base(s)
	versionDir := path.Dir(s)
	name := path.Base(versionDir)
	prefix := path.Dir(versionDir)
	if name == "" {
		return nil, fmt.Errorf("could not determine name of versioned package %q", s)
	}
	if version == "" {
		return nil, fmt.Errorf("could not determine version of versioned package %q", s)
	}

	return &VersionPackage{
		Prefix:  prefix,
		Version: version,
		Name:    name,
	}, nil
}

func (e *VersionPackage) Path() string {
	return path.Join(e.Prefix, e.Name, e.Version)
}

type InternalPackage struct {
	Prefix string
	Name   string
}

func (p *InternalPackage) Path() string {
	return path.Join(p.Prefix, p.Name)
}

func RegexRewriteComment(f *ast.File, r *regexp.Regexp, repl string) {
	for _, cGroup := range f.Comments {
		for _, comment := range cGroup.List {
			comment.Text = r.ReplaceAllString(comment.Text, repl)
		}
	}
}

func RegexInsertCommentBefore(f *ast.File, r *regexp.Regexp, adds ...string) {
	for _, cGroup := range f.Comments {
		for i, comment := range cGroup.List {
			if !r.MatchString(comment.Text) {
				continue
			}

			newComments := make([]*ast.Comment, len(adds))
			for i, add := range adds {
				newComments[i] = &ast.Comment{Text: add}
			}

			cGroup.List = slices.Insert(cGroup.List, i, newComments...)
			return
		}
	}
}

func ParseVersionPackages(ss []string) ([]VersionPackage, error) {
	res := make([]VersionPackage, 0, len(ss))
	for _, s := range ss {
		pkg, err := ParseVersionPackage(s)
		if err != nil {
			return nil, err
		}
		res = append(res, *pkg)
	}
	return res, nil
}

func RemoveTags(f *ast.File) {
	astutil.Apply(f, func(cursor *astutil.Cursor) bool {
		switch x := cursor.Node().(type) {
		case *ast.Field:
			x.Tag = nil
		}
		return true
	}, nil)
}

func RewriteImport(fSet *token.FileSet, f *ast.File, modPath string, vPkg VersionPackage, iPkg InternalPackage) {
	oldImport := path.Join(modPath, vPkg.Path())
	newImport := path.Join(modPath, iPkg.Path())

	oldBase := vPkg.Name + vPkg.Version
	newBase := iPkg.Name

	if !astutil.DeleteNamedImport(fSet, f, oldBase, oldImport) {
		return
	}

	astutil.AddImport(fSet, f, newImport)

	astutil.Apply(f, func(cursor *astutil.Cursor) bool {
		switch x := cursor.Node().(type) {
		case *ast.SelectorExpr:
			if ident, ok := x.X.(*ast.Ident); ok && ident.Name == oldBase {
				ident.Name = newBase
			}
		}
		return true
	}, nil)
}

func RewriteRegister(fSet *token.FileSet, f *ast.File, name string) {
	groupName, ok := GetGroupName(f)
	if !ok {
		panic("could not get group name")
	}

	f.Doc.List[0].Text = fmt.Sprintf("// Package %[1]s contains API Schema definitions for the %[1]s internal API group", name)

	astutil.DeleteNamedImport(fSet, f, "metav1", "k8s.io/apimachinery/pkg/apis/meta/v1")

	astutil.Apply(f, func(cursor *astutil.Cursor) bool {
		switch x := cursor.Node().(type) {
		case *ast.GenDecl:
			if x.Tok != token.VAR {
				return true
			}

			for i, spec := range x.Specs {
				vSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					return true
				}

				for _, name := range vSpec.Names {
					if name.Name != "SchemeGroupVersion" {
						continue
					}

					vSpec.Values[i] = &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("schema"),
							Sel: ast.NewIdent("GroupVersion"),
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Group"),
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: strconv.Quote(groupName),
								},
							},
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Version"),
								Value: &ast.SelectorExpr{
									X:   ast.NewIdent("runtime"),
									Sel: ast.NewIdent("APIVersionInternal"),
								},
							},
						},
					}
					return true
				}
			}
		case *ast.FuncDecl:
			if x.Name.Name != "addKnownTypes" {
				return true
			}
			x.Body.List = slices.Delete(x.Body.List, len(x.Body.List)-2, len(x.Body.List)-1)
		}
		return true
	}, nil)
}

var (
	importCommentRegexp    = regexp.MustCompile(`import "[A-z0-9/.]+"`)
	packageCommentRegexp   = regexp.MustCompile(`Package \w+ is the \w+ version of the API.`)
	groupNameCommentRegexp = regexp.MustCompile(`// \+groupName`)
)

func RewriteDoc(f *ast.File, modPath string, iPkg InternalPackage) {
	RegexInsertCommentBefore(f, groupNameCommentRegexp,
		"// +k8s:defaulter-gen=TypeMeta",
		"// +k8s:protobuf-gen=package",
	)
	RegexRewriteComment(f, importCommentRegexp,
		fmt.Sprintf(`import "%s"`, path.Join(modPath, iPkg.Path())),
	)
	RegexRewriteComment(f, packageCommentRegexp,
		fmt.Sprintf("Package %s is the internal version of the API.", iPkg.Name),
	)
}

func GetGroupName(file *ast.File) (string, bool) {
	for _, cGroup := range file.Comments {
		for _, cItem := range cGroup.List {
			groupName, ok := strings.CutPrefix(cItem.Text, "// +groupName=")
			if ok {
				return groupName, true
			}
		}
	}
	return "", false
}

func RewriteFile(fSet *token.FileSet, modPath string, vPkgs []VersionPackage, iPkg InternalPackage, filename string, file *ast.File) error {
	file.Name = ast.NewIdent(iPkg.Name)

	// Rewrite all external imports to internal imports
	for _, vPkg := range vPkgs {
		vPkgIPkg := InternalPackage{Prefix: iPkg.Prefix, Name: vPkg.Name}

		RewriteImport(fSet, file, modPath, vPkg, vPkgIPkg)
		RemoveTags(file)
	}

	switch filepath.Base(filename) {
	case "register.go":
		RewriteRegister(fSet, file, iPkg.Name)
	case "doc.go":
		RewriteDoc(file, modPath, iPkg)
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fSet, file); err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, formatted, 0666); err != nil {
		return err
	}

	return nil
}

func Process(modPath string, vPkgs []VersionPackage, iPkg InternalPackage, version string) error {
	fSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fSet, iPkg.Path(), nil, parser.ParseComments)
	if err != nil {
		return err
	}

	pkg := pkgs[version]

	for filename, file := range pkg.Files {
		if err := RewriteFile(fSet, modPath, vPkgs, iPkg, filename, file); err != nil {
			return err
		}
	}

	return nil
}

var (
	ignoreFilenames = map[string]struct{}{
		"conversions.go": {},
	}
)

func GatherFilenames(dir string) ([]string, error) {
	ctx := build.Default
	ctx.BuildTags = append(ctx.BuildTags, "ignore_autogenerated")

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if _, ok := ignoreFilenames[entry.Name()]; ok {
			continue
		}

		ok, err := ctx.MatchFile(dir, entry.Name())
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		files = append(files, entry.Name())
	}
	return files, nil
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	set := make(map[string]struct{}, len(b))
	for _, x := range b {
		set[x] = struct{}{}
	}

	var diff []string
	for _, x := range a {
		if _, ok := set[x]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}

func SyncFiles(srcDir, dstDir string) error {
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return err
	}

	ctx := build.Default
	ctx.BuildTags = append(ctx.BuildTags, "ignore_autogenerated")

	srcFiles, err := GatherFilenames(srcDir)
	if err != nil {
		return err
	}

	dstFiles, err := GatherFilenames(dstDir)
	if err != nil {
		return err
	}

	deleteFiles := difference(dstFiles, srcFiles)
	for _, deleteFile := range deleteFiles {
		if err := os.Remove(filepath.Join(dstDir, deleteFile)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	for _, copyFile := range srcFiles {
		if err := utils.CopyFile(filepath.Join(srcDir, copyFile), filepath.Join(dstDir, copyFile)); err != nil {
			return err
		}
	}
	return nil
}

func Rewrite(vPkgs []VersionPackage, iPkgDir string) error {
	const modFileName = "go.mod"
	modFileData, err := os.ReadFile(modFileName)
	if err != nil {
		return err
	}

	mod, err := modfile.Parse(modFileName, modFileData, nil)
	if err != nil {
		return err
	}

	modPath := mod.Module.Mod.Path

	for _, vPkg := range vPkgs {
		iPkg := InternalPackage{Prefix: iPkgDir, Name: vPkg.Name}
		if err := SyncFiles(vPkg.Path(), iPkg.Path()); err != nil {
			return err
		}

		if err := Process(modPath, vPkgs, iPkg, vPkg.Version); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var (
		dir      string
		packages []string
	)

	pflag.StringVar(&dir, "dir", ".", "Output directory of the internal version packages.")
	pflag.StringSliceVar(&packages, "packages", packages, "Packages to generate internal versions for (e.g. api/networking/v1alpha1).")

	pflag.Parse()

	vPkgs, err := ParseVersionPackages(packages)
	if err != nil {
		slog.Error("Error parsing versioned packages", "error", err)
		os.Exit(1)
	}

	if err := Rewrite(vPkgs, dir); err != nil {
		slog.Error("Error rewriting", "error", err)
		os.Exit(1)
	}
}
