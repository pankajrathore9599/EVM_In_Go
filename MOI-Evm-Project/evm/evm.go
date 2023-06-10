package evm

import (
	"log"
	"math/big"
)

type EVM struct {
	stack   []uint64
	memory  []byte
	gas     uint64
	lastMem uint64
}

func NewEVM() *EVM {
	return &EVM{
		stack:   make([]uint64, 0),
		memory:  make([]byte, 0),
		gas:     0,
		lastMem: 0,
	}
}

func (e *EVM) Push1(value uint64) {
	e.stack = append(e.stack, value)
	e.gas += 3
}

func (e *EVM) Push2(value uint64) {
	e.stack = append(e.stack, value)
	e.gas += 3
}

func (e *EVM) Push3(value uint64) {
	e.stack = append(e.stack, value)
	e.gas += 3
}

func (e *EVM) Push32(value []byte) {
	// Assuming value is little-endian byte array
	// Convert it to big-endian for big.Int
	reverseBytes(value)
	bigInt := new(big.Int).SetBytes(value)
	e.stack = append(e.stack, bigInt.Uint64())
	e.gas += 3
}

func (e *EVM) Mstore() {
	if len(e.stack) < 2 {
		log.Fatal("Insufficient stack elements for MSTORE")
	}

	offset := e.stack[len(e.stack)-2]
	value := e.stack[len(e.stack)-1]

	if offset+32 > uint64(len(e.memory)) {
		// Expand memory size if necessary
		newSize := (offset+32+31)/32 * 32
		e.memory = append(e.memory, make([]byte, newSize-uint64(len(e.memory)))...)
		e.gas += (3*newSize + newSize*newSize/512) - e.lastMem
		e.lastMem = 3 * newSize
	}

	bigInt := new(big.Int).SetUint64(value)
	bytes := bigInt.Bytes()
	copy(e.memory[offset:], bytes)
	e.stack = e.stack[:len(e.stack)-2]
}

func (e *EVM) Mstore8() {
	if len(e.stack) < 2 {
		log.Fatal("Insufficient stack elements for MSTORE8")
	}

	offset := e.stack[len(e.stack)-2]
	value := e.stack[len(e.stack)-1] & 0xff

	if offset >= uint64(len(e.memory)) {
		// Expand memory size if necessary
		newSize := (offset+1+31)/32 * 32
		e.memory = append(e.memory, make([]byte, newSize-uint64(len(e.memory)))...)
		e.gas += (3*newSize + newSize*newSize/512) - e.lastMem
		e.lastMem = 3 * newSize
	}

	e.memory[offset] = byte(value)
	e.stack = e.stack[:len(e.stack)-2]
}

func (e *EVM) Add() {
	if len(e.stack) < 2 {
		log.Fatal("Insufficient stack elements for ADD")
	}

	a := e.stack[len(e.stack)-2]
	b := e.stack[len(e.stack)-1]
	sum := a + b

	e.stack[len(e.stack)-2] = sum
	e.stack = e.stack[:len(e.stack)-1]
	e.gas += 3
}

func (e *EVM) Mul() {
	if len(e.stack) < 2 {
		log.Fatal("Insufficient stack elements for MUL")
	}

	a := e.stack[len(e.stack)-2]
	b := e.stack[len(e.stack)-1]
	product := a * b

	e.stack[len(e.stack)-2] = product
	e.stack = e.stack[:len(e.stack)-1]
	e.gas += 5
}

func (e *EVM) Sdiv() {
	if len(e.stack) < 2 {
		log.Fatal("Insufficient stack elements for SDIV")
	}

	a := int64(e.stack[len(e.stack)-2])
	b := int64(e.stack[len(e.stack)-1])

	if b == 0 {
		log.Fatal("Division by zero")
	}

	quotient := uint64(a / b)

	e.stack[len(e.stack)-2] = quotient
	e.stack = e.stack[:len(e.stack)-1]
	e.gas += 5
}

func (e *EVM) Exp() {
	if len(e.stack) < 2 {
		log.Fatal("Insufficient stack elements for EXP")
	}

	base := e.stack[len(e.stack)-2]
	exponent := e.stack[len(e.stack)-1]

	bigBase := new(big.Int).SetUint64(base)
	bigExponent := new(big.Int).SetUint64(exponent)
	result := new(big.Int).Exp(bigBase, bigExponent, nil)

	resultBytes := result.Bytes()
	// Pad the resultBytes with leading zeros if necessary
	if len(resultBytes) < 32 {
		resultBytes = append(make([]byte, 32-len(resultBytes)), resultBytes...)
	}

	// Reverse the resultBytes to little-endian format
	reverseBytes(resultBytes)

	// Convert the little-endian resultBytes back to a uint64 value
	resultUint64 := new(big.Int).SetBytes(resultBytes).Uint64()

	e.stack[len(e.stack)-2] = resultUint64
	e.stack = e.stack[:len(e.stack)-1]
	e.gas += 50 * uint64(len(resultBytes))
}

func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func (e *EVM) GetStack() []uint64 {
	return e.stack
}

func (e *EVM) GetGas() uint64 {
	return e.gas
}
