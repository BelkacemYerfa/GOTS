package main

type Person struct {
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Hobbies   []string  `json:"hobbies"`
	Relations Relations `json:"relations"`
}

type Relations struct {
	Parents  []string `json:"parents"`
	Siblings []string `json:"siblings"`
	Children []string `json:"children"`
}
