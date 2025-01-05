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
	_, _, _, output := process(registers["A"], registers["B"], registers["C"], prog)
	fmt.Printf("part 1 | output: %s\n", stringify(output))

	a := solve2(prog)
	fmt.Printf("part 2 | solved for A: %d\n", a)
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

func calcCombo(A, B, C int, operand int) int {
	var combo int
	regCombo := map[int]int{
		4: A,
		5: B,
		6: C,
	}
	if operand >= 0 && operand <= 3 {
		combo = operand
	} else if operand >= 4 && operand <= 6 {
		combo = regCombo[operand]
	} // FIXME: no control for 7 combo
	return combo
}

func binPretty(bin string) string {
	var sb strings.Builder

	i := 0
	rest := len(bin) % 3
	padding := 3 - rest
	if rest > 0 {
		sb.WriteString(strings.Repeat("0", padding))

		for ; i < rest; i++ {
			sb.WriteByte(bin[i])
		}
	}

	for ; i <= len(bin)-1; i++ {
		if (i+padding)%3 == 0 {
			sb.WriteString(" ")
		}
		sb.WriteByte(bin[i])
	}
	return sb.String()
}

func intToBinString(v int) string {
	return strconv.FormatInt(int64(v), 2)
}

func crackSuffix(prog []int, nOctals int) int {
	var sb strings.Builder
	sb.WriteString("001")
	sb.WriteString(strings.Repeat("000", nOctals-1))
	suffix, err := strconv.ParseInt(sb.String(), 2, 64)

	if err != nil {
		panic(err)
	}

	i := int(suffix)
	for ; i < int(math.Pow(2.0, 27.0))-1; i++ {
		_, _, _, output := process(i, 0, 0, prog)
		if len(output) == nOctals {
			if stringify(prog[:nOctals]) == stringify(output) {
				break
			}
		}
	}
	return i
}

func solve2(prog []int) int {
	var output []int

	// the 9 least significant octals that are pre-computed in ok time
	// 001_000_101_111_001_100_111_111_101 = 18323965 in decimal
	suffixSize := 9
	suffix := crackSuffix(prog, suffixSize)
	fmt.Printf("Precomputed suffix of size %d:\n", suffixSize)
	fmt.Printf("  Decimal: %d\n  Binary: %s\n", suffix, binPretty(intToBinString(suffix)))
	prefixBin := intToBinString(suffix)

	// brute-force the 7 most significant octals starting from
	// 001_000_000_000_000_000_000 = 262144 in decimal
	prefix := 262144
	var a int
	for ; prefix < int(math.Pow(2.0, 33.0))-1; prefix++ {
		var sb strings.Builder
		sb.WriteString(intToBinString(prefix))
		sb.WriteString(prefixBin)
		new_v, err := strconv.ParseInt(sb.String(), 2, 64)
		if err != nil {
			panic(err)
		}
		a = int(new_v)
		_, _, _, output = process(a, 0, 0, prog)
		if len(prog) == len(output) {
			if stringify(prog) == stringify(output) {
				break
			}
		}
	}

	return a
}

func process(A, B, C int, prog []int) (int, int, int, []int) {
	output := []int{}

	i := 0
	for i < len(prog)-1 {
		opcode := prog[i]
		operand := prog[i+1]
		combo := calcCombo(A, B, C, operand)

		switch opcode {
		case 0:
			A = int(float64(A) / math.Pow(2.0, float64(combo)))
			i = i + 2
		case 1:
			B = B ^ operand
			i = i + 2
		case 2:
			B = combo % 8
			i = i + 2
		case 3:
			if A == 0 {
				i++
			} else {
				i = operand
			}
			continue
		case 4:
			B = B ^ C
			i = i + 2
		case 5:
			output = append(output, combo%8)
			// if prog[len(output)-1] != output[len(output)-1] {
			// 	return registers, output
			// }
			i = i + 2
		case 6:
			B = int(float64(A) / math.Pow(2.0, float64(combo)))
			i = i + 2
		case 7:
			C = int(float64(A) / math.Pow(2.0, float64(combo)))
			i = i + 2
		default:
			panic(fmt.Sprintf("unknown opcode: %d; expected opcode in range: 0-7\n", opcode))
		}
	}
	return A, B, C, output
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
			for _, v := range strings.Split(string(progReg.FindSubmatch([]byte(line))[1]), ",") {
				if val, err := strconv.Atoi(v); err == nil {
					prog = append(prog, val)
				}
			}
		default:
			continue
		}
	}
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
