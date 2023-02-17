package SimpleXML

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/beevik/etree"
)

func NewXmlHandle(fileName string, rootName string) *Xml {
	// 初始化根节点
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(fileName); err != nil {
		Log.ErrorLog("初始化xml文件失败,err = %v, fileName = %v, rootName = %v", err, fileName, rootName)
		return nil
	}
	root := doc.SelectElement(rootName)
	if root == nil {
		Log.ErrorLog("初始化xml文件失败, 根节点为空fileName = %v, rootName = %v", fileName, rootName)
		return nil
	}
	return &Xml{
		doc:      doc,
		root:     root,
		fileName: fileName,
	}
}

type Xml struct {
	doc      *etree.Document
	root     *etree.Element
	fileName string
}

func (x *Xml) GetValue(nodeName string) string {
	node := x.root.FindElement(fmt.Sprintf("./%s", nodeName))
	if node == nil {
		return ""
	}
	return node.Text()
}

func (x *Xml) SetValue(nodeName string, value string) bool {
	node := x.root.FindElement(fmt.Sprintf("./%s", nodeName))
	if node == nil {
		return false
	}
	node.SetText(value)
	return true
}

func (x *Xml) Child(nodeName string) *Xml {
	node := x.root.FindElement(fmt.Sprintf("./%s", nodeName))
	if node == nil {
		return nil
	}
	return &Xml{
		root: node,
	}
}

func (x *Xml) OutPut(fileName string) {
	if fileName == "" {
		fileName = x.fileName
	}
	x.doc.WriteToFile(fileName)
}
