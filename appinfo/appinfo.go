package appinfo

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"steamview-go/steam"
	"strings"
)

type magic struct {
	Magic   uint32
	Version uint32
}

type appHeader struct {
	AppId        uint32
	Size         uint32
	InfoState    uint32
	LastUpdated  uint32
	PicsToken    uint64
	SHA1         [20]byte
	ChangeNumber uint32
}

const (
	typeArray byte = iota
	typeString
	typeUint32
	typeEndArray byte = 8
)

type Node struct {
	Name  string
	Value string
	Array []*Node
}

type AppInfo struct {
	Header appHeader
	Root   *Node
}

func (n *Node) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%q: ", n.Name))
	if n.Array == nil {
		sb.WriteString(fmt.Sprintf("%q", n.Value))
	} else {
		sb.WriteString("{\n")
		for _, node := range n.Array {
			sb.WriteString(fmt.Sprintf("%s\n", node.String()))
		}
		sb.WriteString("}\n")
	}
	return sb.String()
}

func (n *Node) GetValue(path string) string {
	if path == n.Name {
		return n.Value
	}
	spPath := strings.Split(path, ":")
	if len(spPath) == 1 || n.Name != spPath[0] || n.Array == nil {
		return ""
	}

	for i := 0; i < len(n.Array); i++ {
		if n.Array[i].Name == spPath[1] {
			return n.Array[i].GetValue(strings.Join(spPath[1:], ":"))
		}
	}

	return ""
}

func (i *AppInfo) GetValue(path string) string {
	return i.Root.GetValue(path)
}

func (i *AppInfo) GetAlign() string {
	return i.GetValue("appinfo:common:library_assets:logo_position:pinned_position")
}

func (i *AppInfo) GetWidth() string {
	return i.GetValue("appinfo:common:library_assets:logo_position:width_pct")
}

func (i *AppInfo) GetHeight() string {
	return i.GetValue("appinfo:common:library_assets:logo_position:height_pct")
}

func (i *AppInfo) GetName() string {
	return i.GetValue("appinfo:common:name")
}

func readUint32(f io.Reader) string {
	var ch uint32
	err := binary.Read(f, binary.LittleEndian, &ch)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", ch)
}

func readCString(f io.Reader) string {
	var sb strings.Builder
	ch := make([]byte, 1)
	for {
		_, err := f.Read(ch)
		if err != nil {
			panic(err)
		}
		if ch[0] == 0 {
			return sb.String()
		}
		sb.WriteByte(ch[0])
	}
}

func readNode(r io.Reader) *Node {
	curNode := new(Node)
	var nodeType byte
	err := binary.Read(r, binary.LittleEndian, &nodeType)
	if err != nil {
		panic(err)
	}

	switch nodeType {
	case typeEndArray:
		return nil
	case typeArray:
		curNode.Name = readCString(r)
		for subNode := readNode(r); subNode != nil; subNode = readNode(r) {
			curNode.Array = append(curNode.Array, subNode)
		}
	case typeString:
		curNode.Name = readCString(r)
		curNode.Value = readCString(r)
	case typeUint32:
		curNode.Name = readCString(r)
		curNode.Value = readUint32(r)
	default:
		panic("unfinished parse!")
	}

	return curNode
}

var collection []AppInfo

func GetAppInfo(appId uint32) AppInfo {
	num := sort.Search(len(collection), func(i int) bool { return collection[i].Header.AppId >= appId })
	if num < len(collection) {
		return collection[num]
	} else {
		return AppInfo{}
	}
}

func Parse() {
	f, err := os.Open(path.Join(steam.CacheRoot, "appinfo.vdf"))

	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	err = binary.Read(f, binary.LittleEndian, new(magic))
	if err != nil {
		panic(err)
	}
	for {
		var curApp AppInfo
		err = binary.Read(f, binary.LittleEndian, &curApp.Header)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				break
			}
			panic(err)
		}

		curApp.Root = readNode(f)

		_ = readNode(f) //extra nil node at end
		collection = append(collection, curApp)
	}
}
