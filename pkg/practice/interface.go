package practice

import "fmt"

type Practice interface {
	Print()
}

type A struct {
	Name string
}

type B struct {
	Name2 string
}

func (a *A) Print() {
	fmt.Printf("This is A, Name is a %s\n", a.Name)
}

func (b *B) Print() {
	fmt.Printf("This is B, Name is a %s\n", b.Name2)
}

func Test() {
	Interface()
}

func Interface() {
	fmt.Println("testting interface...")
	test := []Practice{&A{"arman"}, &B{"banana"}}
	for _, v := range test {
		v.Print()
	}
}
