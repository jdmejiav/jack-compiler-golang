package main

import "fmt"

type Node struct{
  next *Node
  data [][]string
}


type LifeCycleList struct {
  start *Node
  fin *Node
  len int
}

func NewNode (data [][]string) *Node{
  return &Node{
    data: data,
  }
}


func NewList (node *Node) *LifeCycleList{
  return &LifeCycleList{
    start: node,
    len: 1,
  }
}

func (l *LifeCycleList) push (node *Node){
  if (l.fin==nil){
    l.fin=node
    l.start.next=l.fin
    l.len++
  }else{
    temp := l.start
    for temp.next != nil{
      temp = temp.next
    }
    temp.next = node
    l.len++
  }
}

func (l *LifeCycleList) peek() *Node{
  if (l.fin==nil){
    l.len--
    return l.start
  }else{
    temp := l.start
    for temp.next != nil{
      temp = temp.next
    }
    l.fin = temp
    retorno := temp.next
    l.fin.next = nil
    l.len--
    return retorno
  }
}

func (l *LifeCycleList) pop () *Node{
  if (l.fin==nil){
    return l.start
  }else{
    return l.fin
  }
}


func (l *LifeCycleList) getLen() int{
  return l.len
}

func (n *Node) imprimirLista (){
  temp := n
  for temp.next!=nil{
    fmt.Println(temp.data)
    temp=temp.next
  }
}
