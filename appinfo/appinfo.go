package appinfo

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
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
	//there are lotta types, but only string, uint32, and subtree markers are used
	//btw float values are strings and time values are uint32
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

var (
	collection    []AppInfo
	Reading       = true
	maxProgress   int64
	progress      int64
	progressMutex sync.Mutex
)

func GetProgress() float32 {
	progressMutex.Lock()
	defer progressMutex.Unlock()
	return float32(progress) / float32(maxProgress)
}

func setProgress(value int64) {
	progressMutex.Lock()
	defer progressMutex.Unlock()
	progress = value
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

// GetValue gets value from tree, path separator is ':'
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

// GetValue gets value from root of tree
func (i *AppInfo) GetValue(path string) string {
	return i.Root.GetValue(path)
}

func (i *AppInfo) GetAlign() string {
	if Reading {
		return "CenterCenter"
	}
	return i.GetValue("appinfo:common:library_assets:logo_position:pinned_position")
}

func (i *AppInfo) GetWidth() string {
	if Reading {
		return fmt.Sprintf("%f", GetProgress())
	}
	return i.GetValue("appinfo:common:library_assets:logo_position:width_pct")
}

func (i *AppInfo) GetHeight() string {
	if Reading {
		return fmt.Sprintf("%f", GetProgress())
	}
	return i.GetValue("appinfo:common:library_assets:logo_position:height_pct")
}

//GetName returns name of game
func (i *AppInfo) GetName() string {
	if Reading {
		return "_VDF_READING"
	}
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

//readCString reads C string of variable width from file
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

//readNode reads tree recursively
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

func GetAppInfo(appId uint32) AppInfo {
	if Reading {
		return AppInfo{
			Header: appHeader{},
			Root:   nil,
		}
	}

	num := sort.Search(len(collection), func(i int) bool {
		return collection[i].Header.AppId >= appId
	})
	if num < len(collection) {
		return collection[num]
	} else {
		return AppInfo{}
	}
}

func ParseAsync(path string) {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	info, err := f.Stat()
	if err != nil {
		panic("cannot read size of file:" + err.Error())
	}
	maxProgress = info.Size()

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
			if err == io.ErrUnexpectedEOF { // actually, we should read uint32 appID first, check it
				break // for 0x00000000 value, breaking out if true, but it includes extra Seek call
			} // every appid or changing AppInfo header struct, so we just bounce when file is ended
			panic(err)
		}

		curApp.Root = readNode(f)

		nowLength, err := f.Seek(0, 1)
		if err != nil {
			panic(err)
		}

		setProgress(nowLength)

		_ = readNode(f) //extra typeEndArray node at end of tree data, returns nil
		collection = append(collection, curApp)
	}
	Reading = false
}
