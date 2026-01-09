package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"sort"
	"strings"
)

type FuncMetrics struct {
	Name string
	SLOC int
	A    float64
	B    float64
	C    float64
	ABC  float64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: complexity <file.go>")
		os.Exit(1)
	}

	targetFile := os.Args[1]

	// Read source for SLOC calculation
	src, err := os.ReadFile(targetFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, targetFile, src, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(src), "\n")

	var funcs []FuncMetrics
	var totalABC float64

	// Collect function metrics
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		var a, b, c float64

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			switch n.(type) {
			case *ast.AssignStmt:
				a++
			case *ast.CallExpr:
				b++
			case *ast.IfStmt, *ast.ForStmt, *ast.CaseClause, *ast.BinaryExpr:
				c++
			}
			return true
		})

		abc := math.Sqrt(a*a + b*b + c*c)
		sloc := countSLOC(fset, lines, fn)

		funcs = append(funcs, FuncMetrics{
			Name: fn.Name.Name,
			SLOC: sloc,
			A:    a,
			B:    b,
			C:    c,
			ABC:  abc,
		})

		totalABC += abc
	}

	if len(funcs) == 0 {
		fmt.Println("No functions found.")
		return
	}

	// Sort functions by descending ABC
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].ABC > funcs[j].ABC
	})

	funcCount := len(funcs)
	avgABC := totalABC / float64(funcCount)

	maxABC := funcs[0].ABC
	maxFunc := funcs[0].Name

	// Output report
	fmt.Println("===========================================")
	fmt.Println("ABC Complexity Report")
	fmt.Println("===========================================")
	fmt.Printf("File: %s\n", targetFile)
	fmt.Printf("Functions: %d\n", funcCount)
	fmt.Printf("Total ABC: %.2f\n", totalABC)
	fmt.Printf("Average ABC: %.2f\n", avgABC)
	fmt.Printf("Max ABC: %.2f (%s)\n", maxABC, maxFunc)
	fmt.Println()

	fmt.Println("Assessment:")
	fmt.Printf("→ %s\n", designAssessment(funcCount, avgABC, maxABC))
	fmt.Println()

	fmt.Println("Function Breakdown:")
	fmt.Println("-------------------------------------------")
	fmt.Printf("%-25s %5s %6s %6s %6s %7s %7s\n",
		"Function", "SLOC", "A", "B", "C", "ABC", "%Total")

	for _, f := range funcs {
		percent := (f.ABC / totalABC) * 100
		fmt.Printf("%-25s %5d %6.0f %6.0f %6.0f %7.2f %6.1f%%\n",
			f.Name, f.SLOC, f.A, f.B, f.C, f.ABC, percent)
	}
}

func countSLOC(fset *token.FileSet, lines []string, fn *ast.FuncDecl) int {
	start := fset.Position(fn.Body.Pos()).Line
	end := fset.Position(fn.Body.End()).Line

	count := 0
	for i := start - 1; i < end && i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		count++
	}
	return count
}

func designAssessment(funcCount int, avgABC, maxABC float64) string {
	switch {
	case funcCount == 1 && maxABC > 15:
		return "God function detected (logic not decomposed)"
	case maxABC > avgABC*2:
		return "Complexity outlier detected (uneven distribution)"
	case funcCount > 15 && avgABC < 3:
		return "Possible over-fragmentation (too many trivial functions)"
	case avgABC >= 6 && avgABC <= 10 && maxABC < avgABC*1.5:
		return "Healthy, cohesive design"
	default:
		return "Review recommended"
	}
}
