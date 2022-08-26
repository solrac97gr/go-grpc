package models

type Student struct {
	Id   string
	Name string
	Age  int32
}

type Test struct {
	Id   string
	Name string
}

type Question struct {
	Id       string
	Question string
	Answer   string
	TestId   string
}

type Enrollment struct {
	StudentId string
	TestId    string
}
