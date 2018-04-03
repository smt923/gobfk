package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

// BrainfuckProgram represents everything needed to interpret a brainfuck program
//
// Instructions = parsed set of tokens to execute
// Tape = our tape of memory
// DP = Data Pointer = pointer inside our tape
// PC = Program Counter = where are we in our instructions
// Finished = is the program ready to exit?
type BrainfuckProgram struct {
	Instructions []Token
	Tape         []byte
	DP           int
	PC           int
	Finished     bool
}

// Token represents a token in Brainfuck
type Token = byte

// Tokens
const (
	COMMENT Token = iota

	RIGHT
	LEFT
	INC
	DEC
	PRINT
	READ
	LOOPL
	LOOPR
)

func main() {
	if len(os.Args) > 2 {
		// too many arguments, print error and exit
		fmt.Printf("Error, invalid number of arguments, usage:\n%[1]s \t\t- to run as REPL\n%[1]s filename.bf \t- to interpret a file", os.Args[0])
		os.Exit(1)
	} else if len(os.Args) == 2 {
		// program.exe filename.bf - run the file instead
		file, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Printf("Error reading file:\n    %v\n", err)
			os.Exit(1)
		}
		bf := CreateBrainfuckProgram(string(file))
		for bf.Finished != true {
			bf.Evaluate()
		}
	} else {
		// program.exe - interactive mode
		reader := bufio.NewReader(os.Stdin)
		prompt := ":: "
		for {
			fmt.Printf("\n%s", prompt)
			input, _ := reader.ReadString('\n')
			bf := CreateBrainfuckProgram(input)
			for bf.Finished != true {
				bf.Evaluate()
			}
		}
	}
}

// CreateBrainfuckProgram takes an input string and returns a set up BrainfuckProgram
func CreateBrainfuckProgram(input string) BrainfuckProgram {
	bf := BrainfuckProgram{
		Tape:     make([]byte, 64000),
		DP:       0,
		PC:       0,
		Finished: false,
	}
	bf.tokenize(input)
	return bf
}

// Tokenize method will tokenize a string of brainfuck and set up the BF struct with the tokens
func (bf *BrainfuckProgram) tokenize(input string) {
	tokenized := make([]Token, 0, len(input))
	for _, char := range input {
		switch char {
		case '>':
			tokenized = append(tokenized, RIGHT)
		case '<':
			tokenized = append(tokenized, LEFT)
		case '+':
			tokenized = append(tokenized, INC)
		case '-':
			tokenized = append(tokenized, DEC)
		case '.':
			tokenized = append(tokenized, PRINT)
		case ',':
			tokenized = append(tokenized, READ)
		case '[':
			tokenized = append(tokenized, LOOPL)
		case ']':
			tokenized = append(tokenized, LOOPR)
		default:
			tokenized = append(tokenized, COMMENT)
		}
	}
	bf.Instructions = tokenized
}

// Evaluate method will take a single step through our program, executing the intended
// instruction then increasing the program counter
func (bf *BrainfuckProgram) Evaluate() {
	switch bf.Instructions[bf.PC] {
	case RIGHT:
		bf.DP++
	case LEFT:
		bf.DP--
	case INC:
		bf.Tape[bf.DP]++
	case DEC:
		bf.Tape[bf.DP]--
	case PRINT:
		fmt.Printf("%c", bf.Tape[bf.DP])
	case READ:
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		bf.Tape[bf.DP] = byte(char)
	case LOOPL:
		bf.openLoop()
	case LOOPR:
		bf.closeLoop()
	}
	bf.PC++

	if bf.PC >= len(bf.Instructions) {
		bf.Finished = true
	}
}

func (bf *BrainfuckProgram) openLoop() {
	balance := 1
	if bf.Tape[bf.DP] == 0 {
		for balance != 0 {
			bf.PC++
			if bf.Instructions[bf.PC] == LOOPL {
				balance++
			} else if bf.Instructions[bf.PC] == LOOPR {
				balance--
			}
		}
	}
}

func (bf *BrainfuckProgram) closeLoop() {
	balance := 0
	for {
		if bf.Instructions[bf.PC] == LOOPL {
			balance++
		} else if bf.Instructions[bf.PC] == LOOPR {
			balance--
		}
		bf.PC--
		if balance == 0 {
			break
		}
	}
}

func (bf *BrainfuckProgram) debugTokens() []string {
	result := make([]string, 0, len(bf.Instructions))
	for _, token := range bf.Instructions {
		switch token {
		case RIGHT:
			result = append(result, "RIGHT")
		case LEFT:
			result = append(result, "LEFT")
		case INC:
			result = append(result, "INC")
		case DEC:
			result = append(result, "DEC")
		case PRINT:
			result = append(result, "PRINT")
		case READ:
			result = append(result, "READ")
		case LOOPL:
			result = append(result, "LOOPL")
		case LOOPR:
			result = append(result, "LOOPR")
		}
	}
	return result
}
