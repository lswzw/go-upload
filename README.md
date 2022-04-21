# go-upload
go gin 实现的文件上传

# build

go build go-upload.go

# start

./go-upload -h

./go-upload

./go-upload -p 8080 -d


# http.go
这是使用http实现的文件上传（html文件直接使用var定义。w.Write([]byte(html)) 使用）
