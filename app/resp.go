package main

//import (
//	"errors"
//	"fmt"
//	"strconv"
//)
//
//type Type byte
//
//const (
//	String  Type = '+'
//	Error   Type = '-'
//	Integer Type = ':'
//	Bulk    Type = '$'
//	Array   Type = '*'
//)
//
//type Resp struct {
//	Type  Type
//	Raw   []byte
//	Data  []byte
//	Count int
//}
//
//func Parse(command []byte) (*Resp, error) {
//	if len(command) == 0 {
//		return nil, errors.New("no data to read")
//	}
//
//	if n < 1 {
//		return nil, errors.New("not enough data")
//	}
//
//	if buff[n-1] != '\n' || buff[n-2] != '\r' {
//		return nil, errors.New("invalid resp termination, missing CRLF characters")
//	}
//
//	resp := &Resp{
//		Type: Type(buff[0]),
//		Raw:  buff,
//		Data: buff[1 : n-2],
//	}
//
//	switch resp.Type {
//	case String, Error:
//		break
//	case Integer:
//		if len(resp.Data) == 0 {
//			return nil, errors.New("invalid integer")
//		}
//
//		var i int
//		if resp.Data[0] == '-' || resp.Data[0] == '+' {
//			if len(resp.Data) == 1 {
//				return nil, errors.New("invalid integer sign")
//			}
//			i++
//		}
//
//		for ; i < len(resp.Data); i++ {
//			if resp.Data[i] < '0' || resp.Data[i] > '9' {
//				return nil, errors.New("invalid integer")
//			}
//		}
//		break
//	case Bulk:
//		resp.Count, err = strconv.Atoi(string(resp.Data))
//		if err != nil {
//			return nil, errors.New("invalid number of bytes")
//		}
//
//		if resp.Count < 0 {
//			resp.Data = nil
//			resp.Count = 0
//		}
//
//		break
//	case Array:
//		resp.Count, err = strconv.Atoi(string(resp.Data))
//		fmt.Println(buff, resp)
//		break
//	default:
//		return nil, errors.New("invalid kind")
//	}
//
//	return resp, nil
//}
