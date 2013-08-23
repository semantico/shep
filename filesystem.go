package main

import (
	"errors"
	"github.com/libgit2/git2go"
	//"ioutil"
	"log"
	"net/http"
	"os"
	//"regexp"
)

type GitFileSystem struct {
	Repo *git.Repository
}

func (g *GitFileSystem) Open(name string) (http.File, error) {
	// var id string
	// reg, err := regexp.Compile("^/(.*?)\\..*$")
	// if matches := reg.FindStringSubmatch(name); len(matches) > 1 {
	// 	id = matches[1]
	// } else {
	// 	id = name
	// }
	// log.Println(id)
	odb, err := g.Repo.Odb()
	if err != nil {
		log.Fatal(err)
	}

	obj, err := g.Repo.RevparseSingle("HEAD")
	defer obj.Free()
	if err != nil {
		log.Fatal(err)
	}

	commit, err := g.Repo.LookupCommit(obj.Id())
	defer commit.Free()
	if err != nil {
		log.Fatal(err)
	}

	tree, err := commit.Tree()
	if err != nil {
		log.Fatal(err)
	}

	treeEntry := tree.EntryByName(name[1:])

	OdbObj, err := odb.Read(treeEntry.Id)
	defer OdbObj.Free()
	if err != nil {
		log.Fatal(err)
	}

	return &GitFile{data: OdbObj.Data()}, errors.New("Currently unimplemented")
}

func NewGitFileSystem(baseGitPath string) (*GitFileSystem, error) {
	repo, err := git.OpenRepository(baseGitPath)
	if err != nil {
		log.Fatal(err)
		return &GitFileSystem{}, err
	}
	return &GitFileSystem{
		Repo: repo,
	}, nil
}

type GitFile struct {
	data []byte
}

func (g *GitFile) Close() error {
	return errors.New("Currently unimplemented")
}

func (g *GitFile) Stat() (os.FileInfo, error) {
	return os.Lstat("unimplemented")
}

func (g *GitFile) Readdir(count int) (infos []os.FileInfo, err error) {
	var info os.FileInfo
	info, err = os.Lstat("unimplemented")
	infos = []os.FileInfo{info}
	return
}

func (g *GitFile) Read(bytes []byte) (int, error) {
	return copy(bytes, g.data), errors.New("Currently unimplemented")
}

func (g *GitFile) Seek(offset int64, whence int) (int64, error) {
	return int64(0), errors.New("Currently unimplemented")
}
