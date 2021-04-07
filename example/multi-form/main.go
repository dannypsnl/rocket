package main

import (
	"io"
	"os"

	"github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/response"
)

func failPage() response.Html {
	return response.Html(`
<html>
<body>
	<h1>Failed</h1>
</body>
</html>
`)
}

func home() response.Html {
	return response.Html(`
<html>
<body>
	<form action="/upload" method="POST" class="form-group" enctype="multipart/form-data">
        <input type="file" name="file" id="file">
        <input type="submit" value="Upload File" name="submit" class="btn btn-success">
	</form>
</body>
</html>
`)
}

type File struct {
	ReadCloser io.ReadCloser `multiform:"file" file:"yes"`
}

func upload(f *File) response.Html {
	file, err := os.Create("temp.jpg")
	if err != nil {
		return failPage()
	}
	_, err = file.ReadFrom(f.ReadCloser)
	if err != nil {
		return failPage()
	}
	err = f.ReadCloser.Close()
	if err != nil {
		return failPage()
	}
	return home()
}

func main() {
	rocket.Ignite(":8080").
		Mount(
			rocket.Get("/", home),
			rocket.Post("upload", upload),
		).
		Launch()
}
