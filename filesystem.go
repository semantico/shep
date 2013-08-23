package main

import (
	"errors"
	"github.com/libgit2/git2go"
	//"ioutil"
	"log"
	"net/http"
	"os"
	"time"
	//"regexp"
)

type GitFileSystem struct {
	Repo *git.Repository
}

func (g *GitFileSystem) Open(name string) (file http.File, err error) {
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

	tree.Walk(git.TreeWalkCallback(func(s string, entry *git.TreeEntry) int {
		if entry.Name == name[1:] {
			OdbObj, err := odb.Read(entry.Id)
			if err != nil {
				panic(err)
			}
			file = &GitFile{name: name[1:], obj: OdbObj, when: commit.Committer().When}
		}
		return 0
	}))

	return
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
	name string
	obj  *git.OdbObject
	when time.Time
}

func (g *GitFile) Close() error {
	g.obj.Free()
	return nil
}

func (g *GitFile) Stat() (os.FileInfo, error) {
	return &GitFileInfo{name: g.name, size: int64(g.obj.Len()), modTime: g.when}, nil
}

func (g *GitFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{}, nil
}

func (g *GitFile) Read(bytes []byte) (i int, e error) {
	i, e = copy(bytes, g.obj.Data()), nil
	log.Println(i)
	return
}

func (g *GitFile) Seek(offset int64, whence int) (int64, error) {
	return int64(0), errors.New("Currently unimplemented")
}

type GitFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (g *GitFileInfo) Name() string {
	return g.name
}

func (g *GitFileInfo) Size() int64 {
	return g.size
}

func (g *GitFileInfo) Mode() os.FileMode {
	return os.FileMode(uint32(755))
}

func (g *GitFileInfo) ModTime() time.Time {
	return g.modTime
}

func (g *GitFileInfo) IsDir() bool {
	return false
}

func (g *GitFileInfo) Sys() interface{} {
	return nil
}
