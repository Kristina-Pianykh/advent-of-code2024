package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := readLinesFromStream(os.Stdin)
	registers, prog := parse(lines)
	registers, output := process(registers, prog)
	fmt.Printf("%v\n", registers)
	// fmt.Printf("%v\n", output)
	fmt.Printf("part 1 | output: %s\n", stringify(output))

	a := solve2(registers, prog)
	fmt.Printf("solved for A: %d\n", a)
}

func stringify(lst []int) string {
	var sb strings.Builder
	for i, v := range lst {
		if i != len(lst)-1 {
			sb.WriteString(fmt.Sprintf("%d,", v))
		} else {
			sb.WriteString(fmt.Sprintf("%d", v))
		}
	}
	return sb.String()
}

func calcCombo(registers map[string]int, operand int) int {
	var combo int
	regCombo := map[int]int{
		4: registers["A"],
		5: registers["B"],
		6: registers["C"],
	}
	if operand >= 0 && operand <= 3 {
		combo = operand
	} else if operand >= 4 && operand <= 6 {
		combo = regCombo[operand]
	} // FIXME: no control for 7 combo
	return combo
}

func solve2(origRegisters map[string]int, prog []int) int {
	var output []int
	// lower := 1
	// for i := 0; i < 14; i++ {
	// 	lower = lower * 8
	// }
	// upper := 1
	// for i := 0; i < 16; i++ {
	// 	upper = upper * 8
	// }
	a := 1
	for i := 0; i < 15; i++ {
		a = a * 8
	}

	// fmt.Printf("a: %d\n", lower)
	fmt.Printf("prog  : %s\n", stringify(prog))

	i := a * 3
	for ; i <= (a * 4); i++ {
		if i%100000 == 0 {
			fmt.Printf("A: %d\n", i)
		}
		registers := origRegisters
		registers["A"] = i
		_, output = process(registers, prog)
		// fmt.Printf("output: %s\n", stringify(output))
		// 281474976710664 too high
		// 105553129200000 somewhere around here
		// 140737488355328 obere schranke
		if stringify(output) == stringify(prog) {
			fmt.Printf("i: %d\n", i)
			fmt.Printf("prog: %s\n", stringify(prog))
			fmt.Printf("output: %s\n", stringify(output))
			break
		}
	}
	return i
}

func process(registers map[string]int, prog []int) (map[string]int, []int) {
	output := []int{}

	i := 0
	for i < len(prog)-1 {
		// fmt.Printf("inst pointer: %d\n", i)
		// fmt.Printf("registers: %v\n", registers)
		opcode := prog[i]
		operand := prog[i+1]

		combo := calcCombo(registers, operand)
		// fmt.Printf("Opcode: %d, operand: %d, combo: %d\n", opcode, operand, combo)

		switch opcode {
		case 0:
			registers["A"] = int(float64(registers["A"]) / math.Pow(2.0, float64(combo)))
			i = i + 2
		case 1:
			registers["B"] = registers["B"] ^ operand
			i = i + 2
		case 2:
			registers["B"] = combo % 8
			i = i + 2
		case 3:
			if registers["A"] == 0 {
				i++
			} else { // ?? increase by 1 or 2 if not jumping?
				i = operand
			}
			continue
		case 4:
			registers["B"] = registers["B"] ^ registers["C"]
			i = i + 2
		case 5:
			output = append(output, combo%8)
			if prog[len(output)-1] != output[len(output)-1] {
				return registers, output
			}
			// fmt.Printf("combo: %d, combo mod 8: %d\n", combo, combo%8)
			// fmt.Printf("output: %v\n", output)
			i = i + 2
		case 6:
			registers["B"] = int(float64(registers["A"]) / math.Pow(2.0, float64(combo)))
			i = i + 2
		case 7:
			registers["C"] = int(float64(registers["A"]) / math.Pow(2.0, float64(combo)))
			i = i + 2
		default:
			panic(fmt.Sprintf("unknown opcode: %d; expected opcode in range: 0-7\n", opcode))
		}
	}
	return registers, output
}

func parse(lines []string) (map[string]int, []int) {
	registers := map[string]int{}
	rgA := regexp.MustCompile("^Register A: ([[:digit:]]+)$")
	rgB := regexp.MustCompile("^Register B: ([[:digit:]]+)$")
	rgC := regexp.MustCompile("^Register C: ([[:digit:]]+)$")
	progReg := regexp.MustCompile("^Program: (([[:digit:]]+,)+[[:digit:]])$")
	prog := []int{}

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "Register A"):
			if v, err := strconv.Atoi(string(rgA.FindSubmatch([]byte(line))[1])); err == nil {
				registers["A"] = v
			}
		case strings.HasPrefix(line, "Register B"):
			if v, err := strconv.Atoi(string(rgB.FindSubmatch([]byte(line))[1])); err == nil {
				registers["B"] = v
			}
		case strings.HasPrefix(line, "Register C"):
			if v, err := strconv.Atoi(string(rgC.FindSubmatch([]byte(line))[1])); err == nil {
				registers["C"] = v
			}
		case strings.HasPrefix(line, "Program"):
			// for _, v := range progReg.FindSubmatch([]byte(line))[1] {
			// 	fmt.Printf("%v\n", v)
			// }
			for _, v := range strings.Split(string(progReg.FindSubmatch([]byte(line))[1]), ",") {
				if val, err := strconv.Atoi(v); err == nil {
					prog = append(prog, val)
				}
			}
		default:
			continue
		}
	}
	fmt.Printf("registers: %v\n", registers)
	fmt.Printf("prog: %v\n", prog)
	return registers, prog
}

func readLinesFromStream(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}
