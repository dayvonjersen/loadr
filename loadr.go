package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type fileInfo struct {
	Name string
	Size int
}

type fileInfoSlice []*fileInfo

func (f fileInfoSlice) Len() int { return len(f) }
func (f fileInfoSlice) Less(i, j int) bool {
	return f[i].Size < f[j].Size
}
func (f fileInfoSlice) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f fileInfoSlice) Find(maxSize int) (index int) {
	return sort.Search(len(f), func(i int) bool {
		return f[i].Size >= maxSize
	})
}

func main() {
	fmt.Println("Generating playlist")
	rand.Seed(time.Now().Unix())
	ext := ".mp3"
	size := 3435973836 //int64(3.2 * 1024 * 1024 * 1024)
	path := "Z:\\Music"
	result := pickr(path, ext, size)
	f, err := os.Create("playlist.txt")
	checkErr(err)
	s := 0
	for _, r := range result {
		fmt.Fprintln(f, r.Name)
		s += r.Size
	}
	checkErr(f.Close())

	fmt.Println(len(result), "files", formatBytes(s))

	fmt.Println("Copying files")
	cmd := exec.Command("python", "upload.py", fmt.Sprintf("%d", s))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	checkErr(cmd.Run())
}

func pickr(path, ext string, size int) []*fileInfo {
	files := fileInfoSlice(readdir(path, ext))
	if len(files) == 0 {
		return nil
	}
	sort.Sort(files)
	min := files[0].Size
	max := files[len(files)-1].Size
	if min > size {
		return nil
	}
	files = files[:files.Find(max)]
	result := []*fileInfo{}
	for size > min && len(files) > 0 {
		i := rand.Intn(len(files))
		f := files[i]
		size -= f.Size
		result = append(result, f)
		files = append(files[0:i], files[i+1:]...)
		files = files[:files.Find(size)]
	}
	return result
}

func readdir(path, ext string) []*fileInfo {
	dir, err := os.Open(path)
	if err != nil {
		return nil
	}
	files, err := dir.Readdir(-1)
	checkErr(err)
	ret := []*fileInfo{}
	for _, f := range files {
		if f.IsDir() {
			ret = append(ret, readdir(join(path, f.Name()), ext)...)
		} else if filepath.Ext(f.Name()) == ext && f.Size() > 0 && !strings.HasPrefix(f.Name(), "INCOMPLETE~") {
			ret = append(ret, &fileInfo{Name: join(path, f.Name()), Size: int(f.Size())})
		}
	}
	dir.Close()
	return ret
}

var pathSeparator = string(os.PathSeparator)

func join(parts ...string) string {
	for i, part := range parts {
		if strings.HasSuffix(part, pathSeparator) {
			parts[i] = strings.TrimSuffix(part, pathSeparator)
		}
	}
	return strings.Join(parts, pathSeparator)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var byteUnits = [6]string{"", "K", "M", "G", "T", "P"}

func formatBytes(i int) string {
	k := 0
	b := float64(i)
	for b >= 1024 {
		k++
		b /= 1024
	}
	return fmt.Sprintf("%.3f %sB", b, byteUnits[k])
}
