package DHTNet

import (
	"fmt"
	"strings"
	"testing"
)

func TestIntByte(t *testing.T) {
	fmt.Print(strings.Split("a.b.c.",".")[0:3])
}
func TestStringToDomain(t *testing.T){
	a := "QWERTYasdfgZXCVBN"
	domain := PeerIdToDomain(a)
	fmt.Println(domain)
	p := DomainToPeerId(domain)
	fmt.Println(p)
}
