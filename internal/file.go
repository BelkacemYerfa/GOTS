package internal

import "os"

type File struct {
	pathname *string
	filename string
	content  string
}

func NewFile(filename string, pathname *string, content string) *File {
	return &File{
		pathname: pathname,
		filename: filename,
		content:  content,
	}
}

func (f *File) GetPathname() *string {
	return f.pathname
}

func (f *File) GetFilename() string {
	return f.filename
}

func (f *File) GetContent() string {
	return f.content
}

func (f *File) CreateFile() *File {
	file, err := os.Create(f.filename)
	if err != nil {
		return nil
	}
	defer file.Close()
	return f
}

func (f *File) ReadFile() {
	content, err := os.ReadFile(f.filename)
	if err != nil {
		return
	}
	f.content = string(content)
}

func (f *File) WriteFile(content string) {
	err := os.WriteFile(f.filename, []byte(content), 0644)
	if err != nil {
		return
	}
	f.content = content
}

func (f *File) DeleteFile() (bool, error) {
	err := os.Remove(f.filename)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (f *File) AddContentToFile(moreContent string) {
	f.ReadFile()
	content := f.GetContent()
	content += "\n" + moreContent
	f.WriteFile(content)
}
