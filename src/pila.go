package main

import "fmt"

type Stack struct{
  stack []TokenAnalysis
  len int
}





func (s *Stack)getStack() []TokenAnalysis{
  return s.stack
}

func NewStack() *Stack{
  return &Stack{
    stack: []TokenAnalysis{},
    len: 0,
  }
}

func (s *Stack) push(val TokenAnalysis){
  s.stack = append(s.stack,val)
  s.len++
}


func (s Stack) peek() *TokenAnalysis{
  temp := s.stack[s.len-1]
  return &temp
}

func (s *Stack) pop() TokenAnalysis{
  if s.len!=0{
    s.len--
    s.stack = s.stack[: s.len-1]
    return s.stack[s.len-1]
  }else{
    panic("")
  }
}

func (s *Stack) pushAtPos (pos int, arr []TokenAnalysis){
  tempIni := s.stack[:pos:pos+1]
  tempFin := s.stack[pos+1:]


  s.stack = append(tempIni,arr...)

  s.stack = append(s.stack,tempFin...)

  s.len=len(s.stack)
}


func (s *Stack) ToString(){
  for _,i := range s.stack{
      i.ToString()
  }
  fmt.Println("\n");

}

func (s *Stack)deletePos(pos int){
  tempIni := s.stack[:pos:pos+1]
  tempFin := s.stack[pos+1:]
  s.stack = append(tempIni,tempFin...)
}

func (s *Stack) pos(i int) *TokenAnalysis{
  temp := s.stack[i]
  return &temp
}
