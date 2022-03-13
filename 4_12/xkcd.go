package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"github.com/spf13/afero"
)

type ParseStruct struct {
	Num              int
	Year, Month, Day string
	Title            string `json:"title"`
	Transcript       string
	Alt              string
	Img              string // url
}

var mass []ParseStruct

type HttpIface interface {
	Get(url string) (resp *http.Response, err error)
}
type FileIface interface {
	Create(name string)(file afero.File,err error)
	Stat(file *afero.File)(statistic fs.FileInfo,err error)
	Write(file *afero.File,b []byte)( err error)
}

var (
	FS                         afero.Fs
	FSUtil                     *afero.Afero
  )

type RealFile struct {}
type RealHttp struct{}

func (RealFile) Create(name string)(file afero.File,err error) {return FSUtil.Create(name)}
func (RealFile) Write(file *afero.File,b []byte)( err error){return FSUtil.WriteFile((*file).Name(),b,0600)}
func (RealFile) Stat(file *afero.File)(statistic fs.FileInfo,err error){return FSUtil.Stat((*file).Name())}
func (RealHttp) Get(url string) (resp *http.Response, err error) { return http.Get(url) }

type URLsFetcher struct {
	File FileIface
	Http HttpIface
	outWrite io.Writer
	errWrite io.Writer
}

var ( // default <nil> values
	DefaultHttp     HttpIface
	DefaultOutWrite io.Writer
	DefaultErrWrite io.Writer
	DefaultFile FileIface
)

func NewURLsFetcher(File FileIface, Http HttpIface, outWrite io.Writer, errWrite io.Writer) *URLsFetcher {
	fetcher := URLsFetcher{File, Http,  outWrite, errWrite}
	// notest
	
	if Http == DefaultHttp {
		fetcher.Http = RealHttp{}
	}
	if outWrite == DefaultOutWrite {
		fetcher.outWrite = os.Stdout
	}
	if errWrite == DefaultErrWrite {
		fetcher.errWrite = os.Stderr
	}
	if File == DefaultFile {
		fetcher.File = RealFile{}
	}

	return &fetcher
}



func(f *URLsFetcher) getInfo(num int, wg *sync.WaitGroup, file *afero.File) {
	resp, err := f.Http.Get(fmt.Sprintf("https://xkcd.com/%v/info.0.json", num))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	var p ParseStruct
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	mass = append(mass, p)
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	statistic,err := f.File.Stat(file)
	if err != nil {
		panic(err)
	}
	if statistic.Size() > 1 {
		newfile,err := f.File.Create(fmt.Sprintf("%d.json",num))
		if err != nil {
			log.Fatal(err)
		}
		defer newfile.Close()
		file = &newfile
	}
	f.File.Write(file,data)
	wg.Done()
}

func getNum() int {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("search query failed: %s", resp.Status)
		os.Exit(1)
	}
	var parseres ParseStruct
	err = json.NewDecoder(resp.Body).Decode(&parseres)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return parseres.Num
}
func (f *URLsFetcher)CreateFiles(file *afero.File,st fs.FileInfo){
		if st.Size() < 1{
			var wg sync.WaitGroup
		num := getNum()
		wg.Add(num)
		for i := 1; i < 5; i++ {
			go f.getInfo(i, &wg, file)
		}
		wg.Wait()
	}
}

func (f *URLsFetcher) ReadFile(condition string)([]byte){
	file, err := FSUtil.OpenFile(fmt.Sprintf("%v.json",condition), os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	nfile := f.CheckAndFill(&file)
	data, err := FSUtil.ReadFile((*nfile).Name())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return data
}

func(f *URLsFetcher) CheckAndFill (file *afero.File)(*afero.File){
	statistic, err := f.File.Stat(file)
	if err != nil {
		panic(err)
	}
	f.CreateFiles(file,statistic)
	return file
}

func main() {
	f := NewURLsFetcher(DefaultFile,DefaultHttp,os.Stdout,os.Stderr)
	condition := os.Args[1]
	Fetch(f.ReadFile(condition))
	PrintData(os.Stdout,condition)
}

func PrintData(thread io.Writer,condition string){
	for _, elem := range mass {
		if elem.Title == condition || strconv.Itoa(elem.Num) == condition || elem.Year == condition {
			fmt.Fprintf(thread,"%d. Title: %s \n  Data: %s-%s-%s \n Text: %s ", elem.Num, elem.Title, elem.Day, elem.Month, elem.Year, elem.Transcript)
		}
	}
}

func Fetch (data []byte) {
	var pas ParseStruct
	if len(mass) == 0 {
		err := json.Unmarshal(data,&pas)
		mass = append(mass, pas)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}