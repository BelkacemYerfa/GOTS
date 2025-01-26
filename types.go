package main

type Person struct {
	Name      string     `json:"name"`
	Age       int        `json:"age"`
	Hobbies   []string   `json:"hobbies"`
	Relations *Relations `json:"relations"`
	Happy     bool       `json:"happy"`
}

type Relations struct {
	Parents  []string `json:"parents"`
	Siblings []string `json:"siblings"`
	Children []string `json:"children"`
}

type Address struct {
	Street string `json:"street"`
}

type (
	HAHA struct {
		Name string `json:"name"`
		hehe HEHE   `json:"hehe"`
	}

	HEHE struct {
		Age int `json:"age"`
	}
)
