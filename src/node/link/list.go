package link

import (
	"container/list"
)

type List struct {
	links *list.List
}

type Callback func(*Link)

func NewList() *List {
	l := new(List)

	l.links = list.New()

	return l
}

func (list *List) Len() int {
	return list.links.Len()
}

func (list *List) Add(link *Link) {
	list.links.PushBack(link)
}

func (list *List) Remove(link *Link) {
	e := list.findElement(link)

	list.links.Remove(e)
}

func (list *List) Each(callback Callback) {
	for e := list.links.Front(); e != nil; e = e.Next() {
		l := e.Value.(*Link)
		callback(l)
	}
}

func (list *List) findElement(link *Link) *list.Element {
	for e := list.links.Front(); e != nil; e = e.Next() {
		l := e.Value.(*Link)

		if link == l {
			return e
		}
	}

	return nil
}
