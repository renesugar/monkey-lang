package code

// package code implements the bytecode instruction set for our virtual
// machine that will execute monkey-lan source code.

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Definition struct {
	Name          string
	OperandWidths []int
}

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Opcode byte

const (
	LoadConstant Opcode = iota
	LoadBuiltin
	LoadGlobal
	BindGlobal
	LoadLocal
	BindLocal
	LoadFree
	LoadTrue
	LoadFalse
	LoadNull
	MakeArray
	MakeHash
	MakeClosure
	Pop
	Add
	Sub
	Mul
	Div
	Equal
	NotEqual
	GreaterThan
	Minus
	Bang
	Index
	JumpIfFalse
	Jump
	Call
	Return
	ReturnValue
)

var definitions = map[Opcode]*Definition{
	LoadConstant: {"LoadConstant", []int{2}},
	LoadBuiltin:  {"LoadBuiltin", []int{1}},
	LoadGlobal:   {"LoadGlobal", []int{2}},
	BindGlobal:   {"BindGlobal", []int{2}},
	LoadLocal:    {"LoadLocal", []int{1}},
	BindLocal:    {"BindLocal", []int{1}},
	LoadFree:     {"LoadFree", []int{1}},
	LoadTrue:     {"LoadTrue", []int{}},
	LoadFalse:    {"LoadFalse", []int{}},
	LoadNull:     {"LoadNull", []int{}},
	MakeArray:    {"MakeArray", []int{2}},
	MakeHash:     {"MakeHash", []int{2}},
	MakeClosure:  {"MakeClosure", []int{2, 1}},
	Pop:          {"Pop", []int{}},
	Add:          {"Add", []int{}},
	Sub:          {"Sub", []int{}},
	Mul:          {"Mul", []int{}},
	Div:          {"Div", []int{}},
	Equal:        {"Equal", []int{}},
	NotEqual:     {"NotEqual", []int{}},
	GreaterThan:  {"GreaterThan", []int{}},
	Minus:        {"Minus", []int{}},
	Bang:         {"Bang", []int{}},
	Index:        {"Index", []int{}},
	JumpIfFalse:  {"JumpIfFalse", []int{2}},
	Jump:         {"Jump", []int{2}},
	Call:         {"Call", []int{1}},
	Return:       {"Return", []int{}},
	ReturnValue:  {"ReturnValue", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint8(ins Instructions) uint8 { return uint8(ins[0]) }

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}