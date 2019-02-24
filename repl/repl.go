package repl

// Package repl implements the Read-Eval-Print-Loop or interactive console
// by lexing, parsing and evaluating the input in the interpreter

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/prologic/monkey-lang/compiler"
	"github.com/prologic/monkey-lang/eval"
	"github.com/prologic/monkey-lang/lexer"
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/parser"
	"github.com/prologic/monkey-lang/vm"
)

// PROMPT is the REPL prompt displayed for each input
const PROMPT = ">> "

// MonkeyFace is the REPL's face of shock and horror when you encounter a
// parser error :D
const MonkeyFace = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

type Options struct {
	Debug       bool
	Engine      string
	Interactive bool
}

type REPL struct {
	user string
	args []string
	opts *Options
}

func New(user string, args []string, opts *Options) *REPL {
	object.StandardInput = os.Stdin
	object.StandardOutput = os.Stdout
	object.ExitFunction = os.Exit

	return &REPL{user, args, opts}
}

// Eval parses and evalulates the program given by f and returns the resulting
// environment, any errors are printed to stderr
func (r *REPL) Eval(f io.Reader) (env *object.Environment) {
	env = object.NewEnvironment()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source file: %s", err)
		return
	}

	l := lexer.New(string(b))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		return
	}

	eval.Eval(program, env)
	return
}

// Exec parses, compiles and executes the program given by f and returns
// the resulting virtual machine, any errors are printed to stderr
func (r *REPL) Exec(f io.Reader) (state *vm.VMState) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source file: %s", err)
		return
	}

	state = vm.NewVMState()

	l := lexer.New(string(b))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		return
	}

	c := compiler.NewWithState(state.Symbols, state.Constants)
	c.Debug = r.opts.Debug
	err = c.Compile(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Woops! Compilation failed:\n %s\n", err)
		return
	}

	code := c.Bytecode()
	state.Constants = code.Constants

	machine := vm.NewWithGlobalsStore(code, state.Globals)
	machine.Debug = r.opts.Debug
	err = machine.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Woops! Executing bytecode failed:\n %s\n", err)
		return
	}

	return
}

// StartEvalLoop starts the REPL in a continious eval loop
func (r *REPL) StartEvalLoop(in io.Reader, out io.Writer, env *object.Environment) {
	scanner := bufio.NewScanner(in)

	if env == nil {
		env = object.NewEnvironment()
	}

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		obj := eval.Eval(program, env)
		if obj != nil {
			if _, ok := obj.(*object.Null); !ok {
				io.WriteString(out, obj.Inspect())
				io.WriteString(out, "\n")
			}
		}
	}
}

// StartExecLoop starts the REPL in a continious exec loop
func (r *REPL) StartExecLoop(in io.Reader, out io.Writer, state *vm.VMState) {
	scanner := bufio.NewScanner(in)

	if state == nil {
		state = vm.NewVMState()
	}

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		c := compiler.NewWithState(state.Symbols, state.Constants)
		c.Debug = r.opts.Debug
		err := c.Compile(program)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Woops! Compilation failed:\n %s\n", err)
			return
		}

		code := c.Bytecode()
		state.Constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, state.Globals)
		machine.Debug = r.opts.Debug
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Woops! Executing bytecode failed:\n %s\n", err)
			return
		}

		obj := machine.LastPopped()
		if _, ok := obj.(*object.Null); !ok {
			io.WriteString(out, obj.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func (r *REPL) Run() {
	object.Arguments = make([]string, len(r.args))
	copy(object.Arguments, r.args)

	if len(r.args) == 0 {
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", r.user)
		fmt.Printf("Feel free to type in commands\n")
		if r.opts.Engine == "eval" {
			r.StartEvalLoop(os.Stdin, os.Stdout, nil)
		} else {
			r.StartExecLoop(os.Stdin, os.Stdout, nil)
		}
		return
	}

	if len(r.args) > 0 {
		f, err := os.Open(r.args[0])
		if err != nil {
			log.Fatalf("could not open source file %s: %s", r.args[0], err)
		}

		// Remove program argument (zero)
		r.args = r.args[1:]
		object.Arguments = object.Arguments[1:]

		if r.opts.Engine == "eval" {
			env := r.Eval(f)
			if r.opts.Interactive {
				r.StartEvalLoop(os.Stdin, os.Stdout, env)
			}
		} else {
			state := r.Exec(f)
			if r.opts.Interactive {
				r.StartExecLoop(os.Stdin, os.Stdout, state)
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MonkeyFace)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
