package tests

import (
	"bytes"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"example.com/m/4_12"
	"github.com/spf13/afero"
)

func Equal(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}
var (
	FS                         afero.Fs
	FSUtil                     *afero.Afero
  )
var AppFs = afero.NewMemMapFs()

type MockFile struct{}
type MockHttp struct{}

func (MockFile) Create(name string)(file *os.File,err error) {
	return FSUtil.Create("temp")
}
func (MockFile) Write(file *os.File,b []byte)(n int, err error){
	data := b
	icondata:= ReadFile("icon.json")
	if !Equal(data, icondata){
		log.Fatalf("Want:%s,got:%s.",icondata,data)
	}
	return 1,nil
}
func (MockFile) Stat(file *os.File)(statistic fs.FileInfo,err error){return file.Stat()}
func (MockHttp) Get(url string) (resp *http.Response, err error) { 
	data:= ReadFile("icon.json") 
	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(data)),
	}
	return resp,err
}

func TestFetch (t *testing.T){
	data := ReadFile("icon.json")
	type Table struct {
		name   string
		fields xkcd.URLsFetcher
		args   int

		wantText  string
		wantError string
	}
	tests := []Table{
		// {name: "ok", args: args{urls: []string{"https://ifconfig.co/ip"}}, wantRespBody: "XXX.XXX.XXX.XXX"},
		{"fetch0-ok", xkcd.URLsFetcher{MockFile{}, MockHttp{},nil,nil}, 1,  string(data), "" },
		{"fetch0-error", xkcd.URLsFetcher{MockFile{}, MockHttp{}, nil, nil}, 9, "", "panic: open 9.json: The system cannot find the file specified." },
	}

	for _, tc :=  range tests {
		t.Run(tc.name , func(t *testing.T) {
			fetcher := xkcd.NewURLsFetcher(
				tc.fields.File,
				tc.fields.Http,
				outWriter,
				errWriter,
			)
			condition := tc.args

			xkcd.Fetch( xkcd.ReadFile(condition))
			xkcd.PrintData(os.Stdout,condition)
		})
	}
}

// ReadFile returns the contents of `filename`.
func ReadFile (filename string) ([]byte) {
	a := afero.Afero{
		Fs: AppFs,
	}
	fileBytes, err := a.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return fileBytes
}