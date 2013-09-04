package main

import (
	"errors"
	"github.com/libgit2/git2go"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type GitFileSystem struct {
	Repo *git.Repository
}

func (g *GitFileSystem) Open(name string) (file http.File, err error) {
	commitish := "HEAD"

	commitishRegexp, err := regexp.Compile("^(.*?)@(.*)$")
	if err != nil {
		panic(err)
	}

	nameSlice := strings.Split(name, "/")[1:]
	nameSliceLen := len(nameSlice)
	var targetName string
	if nameSliceLen == 1 {
		targetName = name[1:]
	} else {
		targetName = nameSlice[nameSliceLen-1]
	}
	if match := commitishRegexp.FindStringSubmatch(targetName); match != nil {
		targetName = match[1]
		commitish = match[2]
	}
	prefix := strings.Join(nameSlice[:nameSliceLen-1], "/")
	if prefix != "" {
		prefix += "/"
	}

	odb, err := g.Repo.Odb()
	if err != nil {
		log.Fatal(err)
	}

	obj, err := g.Repo.RevparseSingle(commitish)
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

	err = tree.Walk(git.TreeWalkCallback(func(s string, entry *git.TreeEntry) int {
		if entry.Type == git.OBJ_BLOB && entry.Name == targetName && s == prefix {
			OdbObj, err := odb.Read(entry.Id)
			if err != nil {
				panic(err)
			}
			file = &GitFile{name: targetName, obj: OdbObj, when: commit.Committer().When}
		}
		return 0
	}))
	if file == nil {
		err = errors.New("File Not Found Exception")
	}
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
	name   string
	obj    *git.OdbObject
	when   time.Time
	offset int
	data   []byte
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
	bytesLength := len(bytes)
	if g.data == nil {
		g.data = g.obj.Data()
	}
	end := g.offset + bytesLength
	i, e = copy(bytes, g.data[g.offset:end]), nil
	g.offset = end
	return
}

func (g *GitFile) Seek(offset int64, whence int) (int64, error) {
	log.Println("offset" + string(offset))
	return int64(0), nil
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
