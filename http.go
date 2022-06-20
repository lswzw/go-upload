package main

import (
	"io"
	"net/http"
	"os"
)

var html = `
<html>
<body>
<form enctype="multipart/form-data" action="/" method="post">
  <input type="file" name="file" />
  <input type="submit" value="上传" />
</form>
<form enctype="multipart/form-data" action="/files/" method="get">
  <input type="submit" value="查看文件" />
</form>
</body>
</html>
`

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		switch r.Method {
		case "GET":
			w.Write([]byte(html))
		case "POST":
			r.ParseMultipartForm(32 << 20)
			file, header, _ := r.FormFile("file")
			defer file.Close()
			f, _ := os.OpenFile("./files/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			io.Copy(f, file)

			w.Write([]byte(html))
			io.WriteString(w, string(header.Filename+" 上传成功！"))
		}
	} else if r.URL.Path == "/files/" {
		had := http.StripPrefix("/files/", http.FileServer(http.Dir("files")))
		had.ServeHTTP(w, r)
	}
}

func main() {
	_, err := os.Stat("files")
	if os.IsNotExist(err) {
		os.MkdirAll("files", 0777)
	}
	http.HandleFunc("/", ServeHTTP)
	http.ListenAndServe(":80", nil)
}
