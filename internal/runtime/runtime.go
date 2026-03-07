package runtime

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/ir"
)

// Execute runs a lowered IR program.
func Execute(program *ir.ProgramIR, out io.Writer) error {
	r := &runner{out: out, vars: map[string]any{}}
	return r.execute(program)
}

type runner struct {
	out  io.Writer
	vars map[string]any
}

func (r *runner) execute(program *ir.ProgramIR) error {
	_, err := r.executeInstructions(program.Foremost.Instructions)
	return err
}

func (r *runner) executeInstructions(instructions []ir.Instruction) (bool, error) {
	for _, instr := range instructions {
		switch in := instr.(type) {
		case ir.VarDeclInstr:
			val, err := r.eval(in.Value)
			if err != nil {
				return false, err
			}
			r.vars[in.Name] = val
		case ir.AssignInstr:
			val, err := r.eval(in.Value)
			if err != nil {
				return false, err
			}
			r.vars[in.Name] = val
		case ir.ExprInstr:
			if _, err := r.eval(in.Expr); err != nil {
				return false, err
			}
		case ir.IfInstr:
			executed := false
			for _, b := range in.Branches {
				condAny, err := r.eval(b.Condition)
				if err != nil {
					return false, err
				}
				cond, ok := condAny.(bool)
				if !ok {
					return false, fmt.Errorf("if condition must evaluate to truther")
				}
				if !cond {
					continue
				}
				executed = true
				returned, err := r.executeInstructions(b.Body)
				if err != nil {
					return false, err
				}
				if returned {
					return true, nil
				}
				break
			}
			if !executed && len(in.ElseBody) > 0 {
				returned, err := r.executeInstructions(in.ElseBody)
				if err != nil {
					return false, err
				}
				if returned {
					return true, nil
				}
			}
		case ir.WhileInstr:
			for {
				condAny, err := r.eval(in.Condition)
				if err != nil {
					return false, err
				}
				cond, ok := condAny.(bool)
				if !ok {
					return false, fmt.Errorf("given while condition must evaluate to truther")
				}
				if !cond {
					break
				}
				returned, err := r.executeInstructions(in.Body)
				if err != nil {
					return false, err
				}
				if returned {
					return true, nil
				}
			}
		case ir.ForInstr:
			if _, err := r.executeInstructions(in.Init); err != nil {
				return false, err
			}
			for {
				condAny, err := r.eval(in.Condition)
				if err != nil {
					return false, err
				}
				cond, ok := condAny.(bool)
				if !ok {
					return false, fmt.Errorf("given for condition must evaluate to truther")
				}
				if !cond {
					break
				}
				returned, err := r.executeInstructions(in.Body)
				if err != nil {
					return false, err
				}
				if returned {
					return true, nil
				}
				if _, err := r.executeInstructions(in.Step); err != nil {
					return false, err
				}
			}
		case ir.WithinInstr:
			iterAny, err := r.eval(in.Iterable)
			if err != nil {
				return false, err
			}
			var items []any
			switch v := iterAny.(type) {
			case []any:
				items = v
			default:
				return false, fmt.Errorf("within iterable is not an ordered collection")
			}
			for _, item := range items {
				r.vars[in.VarName] = item
				returned, err := r.executeInstructions(in.Body)
				if err != nil {
					return false, err
				}
				if returned {
					return true, nil
				}
			}
		case ir.ReturnInstr:
			if in.Value != nil {
				if _, err := r.eval(in.Value); err != nil {
					return false, err
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func (r *runner) eval(expr ast.Expr) (any, error) {
	switch e := expr.(type) {
	case ast.IdentifierExpr:
		v, ok := r.vars[e.Name]
		if !ok {
			return nil, fmt.Errorf("undefined identifier at runtime: %s", e.Name)
		}
		return v, nil
	case ast.StringExpr:
		return e.Value, nil
	case ast.NumberExpr:
		if strings.Contains(e.Value, ".") {
			f, err := strconv.ParseFloat(e.Value, 64)
			if err != nil {
				return nil, err
			}
			return f, nil
		}
		i, err := strconv.ParseInt(e.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	case ast.BoolExpr:
		return e.Value, nil
	case ast.BinaryExpr:
		return r.evalBinary(e)
	case ast.MemberExpr:
		return e, nil
	case ast.IndexExpr:
		target, err := r.eval(e.Target)
		if err != nil {
			return nil, err
		}
		idxAny, err := r.eval(e.Index)
		if err != nil {
			return nil, err
		}
		idxInt, ok := toInt64(idxAny)
		if !ok {
			return nil, fmt.Errorf("index is not an integer")
		}
		switch t := target.(type) {
		case []any:
			if idxInt < 0 || int(idxInt) >= len(t) {
				return nil, fmt.Errorf("index out of range")
			}
			return t[idxInt], nil
		default:
			return nil, fmt.Errorf("value is not indexable")
		}
	case ast.CallExpr:
		return r.evalCall(e)
	case ast.CollectionLiteralExpr:
		return r.evalCollection(e)
	default:
		return nil, fmt.Errorf("unsupported runtime expression")
	}
}

func (r *runner) evalCollection(e ast.CollectionLiteralExpr) (any, error) {
	switch e.Form {
	case "table", "enumeration":
		out := make([]any, 0, len(e.Items))
		for _, item := range e.Items {
			v, err := r.eval(item)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case "lexis":
		m := map[string]any{}
		for _, p := range e.Pairs {
			k, err := r.eval(p.Key)
			if err != nil {
				return nil, err
			}
			ks, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("lexis key must be string")
			}
			v, err := r.eval(p.Value)
			if err != nil {
				return nil, err
			}
			m[ks] = v
		}
		return m, nil
	default:
		return nil, fmt.Errorf("unsupported collection form: %s", e.Form)
	}
}

func (r *runner) evalCall(e ast.CallExpr) (any, error) {
	if ident, ok := e.Callee.(ast.IdentifierExpr); ok {
		switch ident.Name {
		case "pront":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("pront expects one argument")
			}
			v, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			_, _ = fmt.Fprintln(r.out, toDisplay(v))
			return nil, nil
		case "girth":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("girth expects one argument")
			}
			v, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			switch x := v.(type) {
			case string:
				return int64(len(x)), nil
			case []any:
				return int64(len(x)), nil
			case map[string]any:
				return int64(len(x)), nil
			default:
				return nil, fmt.Errorf("girth unsupported for value")
			}
		}
	}

	if member, ok := e.Callee.(ast.MemberExpr); ok {
		if target, ok := member.Target.(ast.IdentifierExpr); ok && target.Name == "Strangs" && member.Name == "Combobulate" {
			if len(e.Args) == 0 {
				return "", nil
			}
			tmplAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			tmpl, ok := tmplAny.(string)
			if !ok {
				return nil, fmt.Errorf("combobulate template must be string")
			}
			for i := 1; i < len(e.Args); i++ {
				arg, err := r.eval(e.Args[i])
				if err != nil {
					return nil, err
				}
				tmpl = replaceFirstDirective(tmpl, toDisplay(arg))
			}
			return tmpl, nil
		}
	}

	return nil, fmt.Errorf("unsupported call expression")
}

func (r *runner) evalBinary(e ast.BinaryExpr) (any, error) {
	leftAny, err := r.eval(e.Left)
	if err != nil {
		return nil, err
	}
	rightAny, err := r.eval(e.Right)
	if err != nil {
		return nil, err
	}

	lf, lok := toFloat64(leftAny)
	rf, rok := toFloat64(rightAny)
	switch e.Op {
	case "+":
		if lok && rok {
			if isWhole(lf) && isWhole(rf) {
				return int64(lf + rf), nil
			}
			return lf + rf, nil
		}
		return toDisplay(leftAny) + toDisplay(rightAny), nil
	case "-":
		if lok && rok {
			if isWhole(lf) && isWhole(rf) {
				return int64(lf - rf), nil
			}
			return lf - rf, nil
		}
	case "*":
		if lok && rok {
			if isWhole(lf) && isWhole(rf) {
				return int64(lf * rf), nil
			}
			return lf * rf, nil
		}
	case "/":
		if lok && rok {
			return lf / rf, nil
		}
	case "bigly":
		if lok && rok {
			return lf > rf, nil
		}
	case "lessly":
		if lok && rok {
			return lf < rf, nil
		}
	case "biglysame":
		if lok && rok {
			return lf >= rf, nil
		}
	case "lesslysame":
		if lok && rok {
			return lf <= rf, nil
		}
	case "samewise":
		return toDisplay(leftAny) == toDisplay(rightAny), nil
	case "notsamewise":
		return toDisplay(leftAny) != toDisplay(rightAny), nil
	case "plusalso":
		lb, lok := leftAny.(bool)
		rb, rok := rightAny.(bool)
		if lok && rok {
			return lb && rb, nil
		}
	case "orelse":
		lb, lok := leftAny.(bool)
		rb, rok := rightAny.(bool)
		if lok && rok {
			return lb || rb, nil
		}
	}

	return nil, fmt.Errorf("unsupported binary operation: %s", e.Op)
}

func toDisplay(v any) string {
	switch x := v.(type) {
	case nil:
		return "naw"
	case string:
		return x
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		if x {
			return "yee"
		}
		return "nee"
	default:
		return fmt.Sprintf("%v", x)
	}
}

func toFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case int64:
		return float64(x), true
	case float64:
		return x, true
	default:
		return 0, false
	}
}

func toInt64(v any) (int, bool) {
	switch x := v.(type) {
	case int64:
		return int(x), true
	default:
		return 0, false
	}
}

func isWhole(v float64) bool {
	return v == float64(int64(v))
}

func replaceFirstDirective(s, value string) string {
	start := strings.Index(s, "%")
	if start == -1 {
		return s
	}
	end := start + 1
	for end < len(s) {
		ch := s[end]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			end++
			continue
		}
		break
	}
	if end == start+1 {
		return s
	}
	return s[:start] + value + s[end:]
}
