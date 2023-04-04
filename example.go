package main

import (
    "log"
    _ "fmt"
)

type Some struct {
    A int
}

func (s *Some) Ok(ids []int) bool {
    if len(ids) == 0 {
        return false
    }

    log.Println(1)

    tryResult := s.Try(s.A)
    _ = tryResult

    return true
}

func (s *Some) Try(input int) error {
    return nil
}