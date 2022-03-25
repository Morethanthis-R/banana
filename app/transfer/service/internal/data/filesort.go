package data

import (
	"banana/app/transfer/service/internal/biz"
	"sort"
)

type FileMetaWrapper struct {
	File []*biz.File
	by func(p,q *biz.File) bool
}

type FileSortBy func(p,q *biz.File) bool

func (fw FileMetaWrapper) Len() int{
	return len(fw.File)
}

func (fw FileMetaWrapper) Swap(i,j int) {
	fw.File[i],fw.File[j] = fw.File[j],fw.File[i]
}

func (fw FileMetaWrapper) Less(i,j int) bool {
	return fw.by(fw.File[i],fw.File[j])
}

func SortFile(file []*biz.File,by FileSortBy){
	sort.Sort(FileMetaWrapper{file,by})
}


type DirMetaWrapper struct {
	Dir []*biz.UserDirectory
	by func(p,q *biz.UserDirectory) bool
}
type DirSortBy func(p,q *biz.UserDirectory) bool

func (dw DirMetaWrapper) Len() int{
	return len(dw.Dir)
}

func (dw DirMetaWrapper) Swap(i,j int){
	dw.Dir[i],dw.Dir[j] = dw.Dir[j],dw.Dir[i]
}

func (dw DirMetaWrapper) Less(i,j int) bool{
	return dw.by(dw.Dir[i],dw.Dir[j])
}

func SortDir(dir []*biz.UserDirectory,by DirSortBy){
	sort.Sort(DirMetaWrapper{dir,by})
}