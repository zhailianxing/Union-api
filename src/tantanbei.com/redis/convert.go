package redis

import (
	"errors"
)

//one byte convert to value(0-9)
//-1 is error
func SingleInt64(b byte) (int, error) {
	if b < '0' || b > '9' {
		return -1, ERR_NOT_INT
	}
	return int(b - '0'), nil
}

//Bytes convert to int
func MyBstoI64(bs []byte) (int, error) {
	num := 0

	isPositive, err := sign(bs[0])
	if err != nil {
		return 0, err
	}

	if isPositive == true {
		for _, b := range bs {
			num = num * 10
			n, err := SingleInt64(b)
			if err != nil {
				return 0, err
			}
			num = num + n
		}
		return num, err
	} else {
		for _, b := range bs[1:] {
			num = num * 10
			n, err := SingleInt64(b)
			if err != nil {
				return 0, err
			}
			num = num + n
		}
		return -num, err
	}
}

//Determine the number is positive or negative
//return false is negative, true is positive
func sign(b byte) (bool, error) {
	if b == '-' {
		return false, nil
	} else if b >= '0' && b <= '9' {
		return true, nil
	} else {
		return false, errors.New("Not int")
	}
}

//Bytes convert to int32
func MyBstoI32(bs []byte) (int32, error) {

	var num int32

	isPositive, err := sign(bs[0])
	if err != nil {
		return 0, err
	}

	if isPositive == true {
		for _, b := range bs {
			num = num * 10
			n, err := SingleInt32(b)
			if err != nil {
				return 0, err
			}
			num = num + n
		}
		return num, err
	} else {
		for _, b := range bs[1:] {
			num = num * 10
			n, err := SingleInt32(b)
			if err != nil {
				return 0, err
			}
			num = num + n
		}
		return -num, err
	}
}

//one byte convert to value(0-9)
//-1 is error
func SingleInt32(b byte) (int32, error) {
	switch b {
	case '0':
		return 0, nil
	case '1':
		return 1, nil
	case '2':
		return 2, nil
	case '3':
		return 3, nil
	case '4':
		return 4, nil
	case '5':
		return 5, nil
	case '6':
		return 6, nil
	case '7':
		return 7, nil
	case '8':
		return 8, nil
	case '9':
		return 9, nil
	default:
		return -1, ERR_NOT_INT
	}
}

//Bytes convert to int16
func MyBstoI16(bs []byte) (int16, error) {

	var num int16

	isPositive, err := sign(bs[0])
	if err != nil {
		return 0, err
	}

	if isPositive {
		for _, b := range bs {
			num = num * 10
			n, err := SingleInt16(b)
			if err != nil {
				return 0, err
			}
			num = num + n
		}
		return num, err
	} else {
		for _, b := range bs[1:] {
			num = num * 10
			n, err := SingleInt16(b)
			if err != nil {
				return 0, err
			}
			num = num + n
		}
		return -num, err
	}
}

//one byte convert to value(0-9)
//-1 is error
func SingleInt16(b byte) (int16, error) {
	switch b {
	case '0':
		return 0, nil
	case '1':
		return 1, nil
	case '2':
		return 2, nil
	case '3':
		return 3, nil
	case '4':
		return 4, nil
	case '5':
		return 5, nil
	case '6':
		return 6, nil
	case '7':
		return 7, nil
	case '8':
		return 8, nil
	case '9':
		return 9, nil
	default:
		return -1, ERR_NOT_INT
	}
}

//Bytes convert to int16
func MyBstoI8(bs []byte) (int8, error) {
	var num int8

	isPositive, err := sign(bs[0])
	if err != nil {
		return 0, err
	}

	if isPositive {
		for _, b := range bs {

			num = num * 10
			n, err := SingleInt8(b)

			if err != nil {
				return 0, err
			}

			num = num + n
		}
		return num, nil
	} else {
		for _, b := range bs[1:] {

			num = num * 10
			n, err := SingleInt8(b)

			if err != nil {
				return 0, err
			}

			num = num + n
		}

		return -num, nil
	}
}

//one byte convert to value(0-9)
//-1 is error
func SingleInt8(b byte) (int8, error) {
	switch b {
	case '0':
		return 0, nil
	case '1':
		return 1, nil
	case '2':
		return 2, nil
	case '3':
		return 3, nil
	case '4':
		return 4, nil
	case '5':
		return 5, nil
	case '6':
		return 6, nil
	case '7':
		return 7, nil
	case '8':
		return 8, nil
	case '9':
		return 9, nil
	default:
		return -1, ERR_NOT_INT
	}
}

////int convert to []byte
//func MyItobs(num int) []byte {
//	exp := 1
//	num_exp := 0

//	//Calculate the size of int length--num_exp
//	//Calculate the order of magnitude--exp
//	for {
//		exp = exp * 10
//		num_exp++
//		if num < exp {
//			exp = exp / 10
//			break
//		}
//	}
//	bs := make([]byte, num_exp)

//	for n := 0; n < num_exp; n++ {

//		//Calculate the byte of highest level value
//		singleNum := num / exp

//		//Fill the bs
//		bs[n] = SingleItob(singleNum)

//		num = num - singleNum*exp
//		exp = exp / 10
//	}
//	return bs
//}

////value(0-9) convert to byte
//func SingleItob(num int) byte {
//	if num < 0 || num > 10 {
//		panic("SingleItob() input error!")
//	}
//	switch num {
//	case 1:
//		return '1'
//	case 2:
//		return '2'
//	case 3:
//		return '3'
//	case 4:
//		return '4'
//	case 5:
//		return '5'
//	case 6:
//		return '6'
//	case 7:
//		return '7'
//	case 8:
//		return '8'
//	case 9:
//		return '9'
//	case 0:
//		return '0'
//	default:
//		panic("error")
//	}
//}
