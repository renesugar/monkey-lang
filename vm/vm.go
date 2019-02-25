package vm

// package vm implements our virtual machine which is passed a sequence of
// bytecode instructions to execute which has been parsed and compiled by
// the lexer/parser and compiler in previous steps

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"unicode"

	"github.com/prologic/monkey-lang/builtins"
	"github.com/prologic/monkey-lang/code"
	"github.com/prologic/monkey-lang/compiler"
	"github.com/prologic/monkey-lang/lexer"
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/parser"
	"github.com/prologic/monkey-lang/utils"
)

const (
	StackSize  = 2048
	MaxFrames  = 1024
	MaxGlobals = 65536
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {

	case *object.Boolean:
		return obj.Value

	case *object.Null:
		return false

	default:
		return true
	}
}

// ExecModule compiles the named module and returns a *object.Module object
func ExecModule(name string, state *VMState) (object.Object, error) {
	filename := utils.FindModule(name)
	if filename == "" {
		return nil, fmt.Errorf("ImportError: no module named '%s'", name)
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("IOError: error reading module '%s': %s", name, err)
	}

	l := lexer.New(string(b))
	p := parser.New(l)

	module := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, fmt.Errorf("ParseError: %s", p.Errors())
	}

	c := compiler.NewWithState(state.Symbols, state.Constants)
	err = c.Compile(module)
	if err != nil {
		return nil, fmt.Errorf("CompileError: %s", err)
	}

	code := c.Bytecode()
	state.Constants = code.Constants

	machine := NewWithState(code, state)
	err = machine.Run()
	if err != nil {
		return nil, fmt.Errorf("RuntimeError: error loading module '%s'", err)
	}

	return state.ExportedHash(), nil
}

type VMState struct {
	Constants []object.Object
	Globals   []object.Object
	Symbols   *compiler.SymbolTable
}

func NewVMState() *VMState {
	symbolTable := compiler.NewSymbolTable()
	for i, builtin := range builtins.BuiltinsIndex {
		symbolTable.DefineBuiltin(i, builtin.Name)
	}

	return &VMState{
		Constants: []object.Object{},
		Globals:   make([]object.Object, MaxGlobals),
		Symbols:   symbolTable,
	}
}

// ExportedHash returns a new Hash with the names and values of every publically
// exported binding in the vm state. That is every binding that starts with a
// capital letter. This is used by the module import system to wrap up the
// compiled and evaulated module into an object.
func (s *VMState) ExportedHash() *object.Hash {
	pairs := make(map[object.HashKey]object.HashPair)
	for name, symbol := range s.Symbols.Store {
		if unicode.IsUpper(rune(name[0])) {
			if symbol.Scope == compiler.GlobalScope {
				obj := s.Globals[symbol.Index]
				s := &object.String{Value: name}
				pairs[s.HashKey()] = object.HashPair{Key: s, Value: obj}
			}
		}
	}
	return &object.Hash{Pairs: pairs}
}

type VM struct {
	Debug bool

	state *VMState

	frames      []*Frame
	framesIndex int

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	state := NewVMState()
	state.Constants = bytecode.Constants

	return &VM{
		state: state,

		frames:      frames,
		framesIndex: 1,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func NewWithState(bytecode *compiler.Bytecode, state *VMState) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		state: state,

		frames:      frames,
		framesIndex: 1,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {

	// {"a": 1} + {"b": 2}
	case op == code.Add && left.Type() == object.HASH && right.Type() == object.HASH:
		leftVal := left.(*object.Hash).Pairs
		rightVal := right.(*object.Hash).Pairs
		pairs := make(map[object.HashKey]object.HashPair)
		for k, v := range leftVal {
			pairs[k] = v
		}
		for k, v := range rightVal {
			pairs[k] = v
		}
		return vm.push(&object.Hash{Pairs: pairs})

	// [1] + [2]
	case op == code.Add && left.Type() == object.ARRAY && right.Type() == object.ARRAY:
		leftVal := left.(*object.Array).Elements
		rightVal := right.(*object.Array).Elements
		elements := make([]object.Object, len(leftVal)+len(rightVal))
		elements = append(leftVal, rightVal...)
		return vm.push(&object.Array{Elements: elements})

	// [1] * 3
	case op == code.Mul && left.Type() == object.ARRAY && right.Type() == object.INTEGER:
		leftVal := left.(*object.Array).Elements
		rightVal := int(right.(*object.Integer).Value)
		elements := leftVal
		for i := rightVal; i > 1; i-- {
			elements = append(elements, leftVal...)
		}
		return vm.push(&object.Array{Elements: elements})
	// 3 * [1]
	case op == code.Mul && left.Type() == object.INTEGER && right.Type() == object.ARRAY:
		leftVal := int(left.(*object.Integer).Value)
		rightVal := right.(*object.Array).Elements
		elements := rightVal
		for i := leftVal; i > 1; i-- {
			elements = append(elements, rightVal...)
		}
		return vm.push(&object.Array{Elements: elements})

	// " " * 4
	case op == code.Mul && left.Type() == object.STRING && right.Type() == object.INTEGER:
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.Integer).Value
		return vm.push(&object.String{Value: strings.Repeat(leftVal, int(rightVal))})
	// 4 * " "
	case op == code.Mul && left.Type() == object.INTEGER && right.Type() == object.STRING:
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.String).Value
		return vm.push(&object.String{Value: strings.Repeat(rightVal, int(leftVal))})

	case leftType == object.BOOLEAN && rightType == object.BOOLEAN:
		return vm.executeBinaryBooleanOperation(op, left, right)
	case leftType == object.INTEGER && rightType == object.INTEGER:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING && rightType == object.STRING:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s",
			leftType, rightType)
	}
}

func (vm *VM) executeBinaryStringOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	if op != code.Add {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
}

func (vm *VM) executeBinaryBooleanOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	var result bool

	switch op {
	case code.Or:
		result = leftValue || rightValue
	case code.And:
		result = leftValue && rightValue
	default:
		return fmt.Errorf("unknown boolean operator: %d", op)
	}

	return vm.push(&object.Boolean{Value: result})
}

func (vm *VM) executeBinaryIntegerOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.Add:
		result = leftValue + rightValue
	case code.Sub:
		result = leftValue - rightValue
	case code.Mul:
		result = leftValue * rightValue
	case code.Div:
		result = leftValue / rightValue
	case code.Mod:
		result = leftValue % rightValue
	case code.BitwiseOR:
		result = leftValue | rightValue
	case code.BitwiseXOR:
		result = leftValue ^ rightValue
	case code.BitwiseAND:
		result = leftValue & rightValue
	case code.LeftShift:
		result = leftValue << uint64(rightValue)
	case code.RightShift:
		result = leftValue >> uint64(rightValue)
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	switch op {
	case code.Equal:
		return vm.push(nativeBoolToBooleanObject(left.(object.Comparable).Compare(right) == 0))
	case code.NotEqual:
		return vm.push(nativeBoolToBooleanObject(left.(object.Comparable).Compare(right) != 0))
	case code.GreaterThanEqual:
		return vm.push(nativeBoolToBooleanObject(left.(object.Comparable).Compare(right) > -1))
	case code.GreaterThan:
		return vm.push(nativeBoolToBooleanObject(left.(object.Comparable).Compare(right) == 1))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func (vm *VM) executeBitwiseNotOperator() error {
	operand := vm.pop()
	if i, ok := operand.(*object.Integer); ok {
		return vm.push(&object.Integer{Value: ^i.Value})
	}
	return fmt.Errorf("expected int got=%T", operand)
}

func (vm *VM) executeNotOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()
	if i, ok := operand.(*object.Integer); ok {
		return vm.push(&object.Integer{Value: -i.Value})
	}
	return fmt.Errorf("expected int got=%T", operand)
}

func (vm *VM) executeSetItem(left, index, value object.Object) error {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return vm.executeArraySetItem(left, index, value)
	case left.Type() == object.HASH:
		return vm.executeHashSetItem(left, index, value)
	default:
		return fmt.Errorf(
			"set item operation not supported: left=%s index=%s",
			left.Type(), index.Type(),
		)
	}
}

func (vm *VM) executeGetItem(left, index object.Object) error {
	switch {
	case left.Type() == object.STRING && index.Type() == object.INTEGER:
		return vm.executeStringGetItem(left, index)
	case left.Type() == object.STRING && index.Type() == object.STRING:
		return vm.executeStringIndex(left, index)
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return vm.executeArrayGetItem(left, index)
	case left.Type() == object.HASH:
		return vm.executeHashGetItem(left, index)
	case left.Type() == object.MODULE:
		return vm.executeHashGetItem(left.(*object.Module).Attrs, index)
	default:
		return fmt.Errorf(
			"index operator not supported: left=%s index=%s",
			left.Type(), index.Type(),
		)
	}
}

func (vm *VM) executeStringGetItem(str, index object.Object) error {
	stringObject := str.(*object.String)
	i := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if i < 0 || i > max {
		return vm.push(&object.String{Value: ""})
	}

	return vm.push(&object.String{Value: string(stringObject.Value[i])})
}

func (vm *VM) executeStringIndex(str, index object.Object) error {
	stringObject := str.(*object.String)
	substr := index.(*object.String).Value

	return vm.push(
		&object.Integer{
			Value: int64(strings.Index(stringObject.Value, substr)),
		},
	)
}

func (vm *VM) executeArrayGetItem(array, index object.Object) error {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > max {
		return vm.push(Null)
	}

	return vm.push(arrayObject.Elements[i])
}

func (vm *VM) executeArraySetItem(array, index, value object.Object) error {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > max {
		return fmt.Errorf("index out of bounds: %d", i)
	}

	arrayObject.Elements[i] = value
	return vm.push(Null)
}

func (vm *VM) executeHashGetItem(hash, index object.Object) error {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) executeHashSetItem(hash, index, value object.Object) error {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	hashed := key.HashKey()
	hashObject.Pairs[hashed] = object.HashPair{Key: index, Value: value}

	return vm.push(Null)
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}

		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) executeCall(numArgs int) error {
	callee := vm.stack[vm.sp-1-numArgs]
	switch callee := callee.(type) {
	case *object.Closure:
		return vm.callClosure(callee, numArgs)
	case *object.Builtin:
		return vm.callBuiltin(callee, numArgs)
	default:
		return fmt.Errorf(
			"calling non-closure and non-builtin: %T %v",
			callee, callee,
		)
	}
}

func (vm *VM) callClosure(cl *object.Closure, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			cl.Fn.NumParameters, numArgs)
	}

	// Optimize tail calls and avoid creating a new frame
	if cl.Fn == vm.currentFrame().cl.Fn {
		nextOp := vm.currentFrame().NextOp()
		if nextOp == code.Return {
			for p := 0; p < numArgs; p++ {
				vm.stack[vm.currentFrame().basePointer+p] = vm.stack[vm.sp-numArgs+p]
			}
			vm.sp -= numArgs + 1
			vm.currentFrame().ip = -1 // reset IP to beginning of the frame
			return nil
		}
	}

	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func (vm *VM) callBuiltin(builtin *object.Builtin, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := builtin.Fn(args...)
	vm.sp = vm.sp - numArgs - 1

	if result != nil {
		vm.push(result)
	} else {
		vm.push(Null)
	}

	return nil
}

func (vm *VM) pushClosure(constIndex, numFree int) error {
	constant := vm.state.Constants[constIndex]
	function, ok := constant.(*object.CompiledFunction)
	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}

	free := make([]object.Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}
	vm.sp = vm.sp - numFree

	closure := &object.Closure{Fn: function, Free: free}
	return vm.push(closure)
}

func (vm *VM) loadModule(name object.Object) error {
	s, ok := name.(*object.String)
	if !ok {
		return fmt.Errorf(
			"TypeError: import() expected argument #1 to be `str` got `%s`",
			name.Type(),
		)
	}

	attrs, err := ExecModule(s.Value, vm.state)
	if err != nil {
		return err
	}

	module := &object.Module{Name: s.Value, Attrs: attrs}
	return vm.push(module)
}

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	var (
		ip  int
		ins code.Instructions
		op  code.Opcode
	)

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])

		if vm.Debug {
			log.Printf(
				"%-25s %-20s\n",
				fmt.Sprintf(
					"%04d %s", ip,
					strings.Split(ins[ip:].String(), "\n")[0][4:],
				),
				fmt.Sprintf(
					"[ip=%02d fp=%02d, sp=%02d]",
					ip, vm.framesIndex-1, vm.sp,
				),
			)
		}

		switch op {

		case code.LoadBuiltin:
			builtinIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			builtin := builtins.BuiltinsIndex[builtinIndex]

			err := vm.push(builtin)
			if err != nil {
				return err
			}

		case code.LoadConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.state.Constants[constIndex])
			if err != nil {
				return err
			}

		case code.AssignGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.state.Globals[globalIndex] = vm.pop()

			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.AssignLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()

			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.BindGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			ref := vm.pop()
			if immutable, ok := ref.(object.Immutable); ok {
				vm.state.Globals[globalIndex] = immutable.Clone()
			} else {
				vm.state.Globals[globalIndex] = ref
			}

			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.LoadGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.state.Globals[globalIndex])
			if err != nil {
				return err
			}

		case code.BindLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()

			ref := vm.pop()
			if immutable, ok := ref.(object.Immutable); ok {
				vm.stack[frame.basePointer+int(localIndex)] = immutable.Clone()
			} else {
				vm.stack[frame.basePointer+int(localIndex)] = ref
			}

			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.LoadLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()

			err := vm.push(vm.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}

		case code.LoadFree:
			freeIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			currentClosure := vm.currentFrame().cl
			err := vm.push(currentClosure.Free[freeIndex])
			if err != nil {
				return err
			}

		case code.LoadModule:
			name := vm.pop()

			err := vm.loadModule(name)
			if err != nil {
				return err
			}

		case code.SetSelf:
			freeIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			currentClosure := vm.currentFrame().cl
			currentClosure.Free[freeIndex] = currentClosure

		case code.LoadTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.LoadFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.LoadNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.MakeHash:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElements

			err = vm.push(hash)
			if err != nil {
				return err
			}

		case code.MakeArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}

		case code.MakeClosure:
			constIndex := code.ReadUint16(ins[ip+1:])
			numFree := code.ReadUint8(ins[ip+3:])
			vm.currentFrame().ip += 3

			err := vm.pushClosure(int(constIndex), int(numFree))
			if err != nil {
				return err
			}

		case code.Add, code.Sub, code.Mul, code.Div, code.Mod,
			code.Or, code.And,
			code.BitwiseOR, code.BitwiseXOR, code.BitwiseAND,
			code.LeftShift, code.RightShift:

			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}

		case code.Equal, code.NotEqual, code.GreaterThan, code.GreaterThanEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.Not:
			err := vm.executeNotOperator()
			if err != nil {
				return err
			}

		case code.BitwiseNOT:
			err := vm.executeBitwiseNotOperator()
			if err != nil {
				return err
			}

		case code.SetItem:
			value := vm.pop()
			index := vm.pop()
			left := vm.pop()

			err := vm.executeSetItem(left, index, value)
			if err != nil {
				return err
			}

		case code.GetItem:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeGetItem(left, index)
			if err != nil {
				return err
			}

		case code.Minus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}

		case code.Call:
			numArgs := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			err := vm.executeCall(int(numArgs))
			if err != nil {
				return err
			}

		case code.Return:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				return err
			}

		case code.JumpIfFalse:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}

		case code.Jump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1

		case code.Pop:
			vm.pop()

		}

		if vm.Debug {
			log.Printf(
				"%-25s [ip=%02d fp=%02d, sp=%02d]",
				"", ip, vm.framesIndex-1, vm.sp,
			)
		}

	}

	return nil
}
