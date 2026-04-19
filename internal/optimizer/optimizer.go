package optimizer

import (
	"fmt"
	"strconv"

	"go_compiler/internal/quad"
)

func isLiteral(s string) bool {
	_, e1 := strconv.ParseInt(s, 10, 64)
	_, e2 := strconv.ParseFloat(s, 64)
	return e1 == nil || e2 == nil
}

func isTemp(s string) bool {
	return len(s) > 0 && s[0] == 'T'
}

func eval(op, a, b string) (string, bool) {
	fa, e1 := strconv.ParseFloat(a, 64)
	fb, e2 := strconv.ParseFloat(b, 64)
	if e1 != nil || e2 != nil {
		return "", false
	}

	switch op {
	case "+":
		return fmtFloat(fa + fb), true
	case "-":
		return fmtFloat(fa - fb), true
	case "*":
		return fmtFloat(fa * fb), true
	case "/":
		if fb == 0 {
			return "", false
		}
		return fmtFloat(fa / fb), true
	case "<":
		if fa < fb {
			return "1", true
		}
		return "0", true
	case ">":
		if fa > fb {
			return "1", true
		}
		return "0", true
	case "<=":
		if fa <= fb {
			return "1", true
		}
		return "0", true
	case ">=":
		if fa >= fb {
			return "1", true
		}
		return "0", true
	case "==":
		if fa == fb {
			return "1", true
		}
		return "0", true
	case "!=":
		if fa != fb {
			return "1", true
		}
		return "0", true
	}
	return "", false
}

func fmtFloat(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func propagate(quads []quad.Quad) []quad.Quad {
	env := map[string]string{}

	resolve := func(x string) string {
		for {
			v, ok := env[x]
			if !ok {
				return x
			}
			x = v
		}
	}

	out := []quad.Quad{}

	for _, q := range quads {
		a1 := resolve(q.Arg1)
		a2 := resolve(q.Arg2)

		if q.Op == ":=" && isLiteral(a1) {
			env[q.Result] = a1
		} else if q.Op == ":=" && isTemp(a1) {
			env[q.Result] = a1
		} else if q.Result != "" {
			delete(env, q.Result)
		}

		out = append(out, quad.Quad{Op: q.Op, Arg1: a1, Arg2: a2, Result: q.Result})
	}

	return out
}

func fold(quads []quad.Quad) []quad.Quad {
	out := []quad.Quad{}

	for _, q := range quads {

		if q.Op == "NEG" && isLiteral(q.Arg1) {
			f, _ := strconv.ParseFloat(q.Arg1, 64)
			out = append(out, quad.Quad{Op: ":=", Arg1: fmtFloat(-f), Arg2: "", Result: q.Result})
			continue
		}

		if isLiteral(q.Arg1) && isLiteral(q.Arg2) {
			if val, ok := eval(q.Op, q.Arg1, q.Arg2); ok {
				out = append(out, quad.Quad{Op: ":=", Arg1: val, Arg2: "", Result: q.Result})
				continue
			}
		}

		out = append(out, q)
	}

	return out
}

func simplify(quads []quad.Quad) []quad.Quad {
	out := []quad.Quad{}

	for _, q := range quads {

		switch q.Op {

		case "+":
			if q.Arg2 == "0" {
				out = append(out, quad.Quad{Op: ":=", Arg1: q.Arg1, Arg2: "", Result: q.Result})
				continue
			}
			if q.Arg1 == "0" {
				out = append(out, quad.Quad{Op: ":=", Arg1: q.Arg2, Arg2: "", Result: q.Result})
				continue
			}

		case "-":
			if q.Arg2 == "0" {
				out = append(out, quad.Quad{Op: ":=", Arg1: q.Arg1, Arg2: "", Result: q.Result})
				continue
			}

		case "*":
			if q.Arg1 == "0" || q.Arg2 == "0" {
				out = append(out, quad.Quad{Op: ":=", Arg1: "0", Arg2: "", Result: q.Result})
				continue
			}
			if q.Arg1 == "1" {
				out = append(out, quad.Quad{Op: ":=", Arg1: q.Arg2, Arg2: "", Result: q.Result})
				continue
			}
			if q.Arg2 == "1" {
				out = append(out, quad.Quad{Op: ":=", Arg1: q.Arg1, Arg2: "", Result: q.Result})
				continue
			}

		case "/":
			if q.Arg1 == "0" {
				out = append(out, quad.Quad{Op: ":=", Arg1: "0", Arg2: "", Result: q.Result})
				continue
			}
			if q.Arg2 == "1" {
				out = append(out, quad.Quad{Op: ":=", Arg1: q.Arg1, Arg2: "", Result: q.Result})
				continue
			}
		}

		out = append(out, q)
	}

	return out
}

func cse(quads []quad.Quad) []quad.Quad {
	exprMap := map[string]string{}
	out := []quad.Quad{}

	key := func(op, a, b string) string {
		return op + "|" + a + "|" + b
	}

	for _, q := range quads {

		if q.Op == "+" || q.Op == "*" || q.Op == "-" || q.Op == "/" {
			k := key(q.Op, q.Arg1, q.Arg2)

			if existing, ok := exprMap[k]; ok {
				out = append(out, quad.Quad{Op: ":=", Arg1: existing, Arg2: "", Result: q.Result})
				continue
			}

			exprMap[k] = q.Result
		}

		out = append(out, q)
	}

	return out
}

func deadCode(quads []quad.Quad) []quad.Quad {
	uses := map[string]int{}

	for _, q := range quads {
		if q.Arg1 != "" {
			uses[q.Arg1]++
		}
		if q.Arg2 != "" {
			uses[q.Arg2]++
		}
	}

	out := []quad.Quad{}

	for _, q := range quads {
		if q.Op == ":=" && isTemp(q.Result) && uses[q.Result] == 0 {
			continue
		}
		out = append(out, q)
	}

	return out
}

func Optimize(quads []quad.Quad) []quad.Quad {
	q := propagate(quads)
	q = fold(q)
	q = simplify(q)
	q = cse(q)
	q = deadCode(q)

	return q
}
