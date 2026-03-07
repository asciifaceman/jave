package runtime

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/asciifaceman/jave/internal/ast"
	"github.com/asciifaceman/jave/internal/ir"
)

// ExecuteOptions configures runtime IO and argument behavior.
type ExecuteOptions struct {
	Stdout io.Writer
	Stderr io.Writer
	Args   []string
}

// ProgramExitError carries an explicit program-requested exit code.
type ProgramExitError struct {
	Code int
}

func (e *ProgramExitError) Error() string {
	return fmt.Sprintf("program requested exit with code %d", e.Code)
}

// ExitCodeForError returns deterministic process exit codes for runtime failures.
func ExitCodeForError(err error) int {
	var pe *ProgramExitError
	if errors.As(err, &pe) {
		return pe.Code
	}
	return 1
}

// Execute runs a lowered IR program.
func Execute(program *ir.ProgramIR, out io.Writer) error {
	return ExecuteWithOptions(program, ExecuteOptions{Stdout: out})
}

// ExecuteWithOptions runs a lowered IR program with explicit IO/argv configuration.
func ExecuteWithOptions(program *ir.ProgramIR, opts ExecuteOptions) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}
	r := &runner{out: stdout, err: stderr, args: append([]string(nil), opts.Args...), program: program, frames: []map[string]any{{}}, modules: []string{""}}
	return r.execute(program)
}

type runner struct {
	out     io.Writer
	err     io.Writer
	args    []string
	program *ir.ProgramIR
	frames  []map[string]any
	modules []string
}

func (r *runner) execute(program *ir.ProgramIR) error {
	for _, foreward := range program.Forewards {
		_, _, err := r.executeInstructions(foreward.Instructions)
		if err != nil {
			return err
		}
	}
	_, _, err := r.executeInstructions(program.Foremost.Instructions)
	return err
}

func (r *runner) executeInstructions(instructions []ir.Instruction) (bool, any, error) {
	for _, instr := range instructions {
		switch in := instr.(type) {
		case ir.VarDeclInstr:
			val, err := r.eval(in.Value)
			if err != nil {
				return false, nil, err
			}
			r.currentFrame()[in.Name] = val
		case ir.AssignInstr:
			val, err := r.eval(in.Value)
			if err != nil {
				return false, nil, err
			}
			r.assignVar(in.Name, val)
		case ir.ExprInstr:
			if _, err := r.eval(in.Expr); err != nil {
				return false, nil, err
			}
		case ir.IfInstr:
			executed := false
			for _, b := range in.Branches {
				condAny, err := r.eval(b.Condition)
				if err != nil {
					return false, nil, err
				}
				cond, ok := condAny.(bool)
				if !ok {
					return false, nil, fmt.Errorf("if condition must evaluate to truther")
				}
				if !cond {
					continue
				}
				executed = true
				returned, value, err := r.executeInstructions(b.Body)
				if err != nil {
					return false, nil, err
				}
				if returned {
					return true, value, nil
				}
				break
			}
			if !executed && len(in.ElseBody) > 0 {
				returned, value, err := r.executeInstructions(in.ElseBody)
				if err != nil {
					return false, nil, err
				}
				if returned {
					return true, value, nil
				}
			}
		case ir.WhileInstr:
			for {
				condAny, err := r.eval(in.Condition)
				if err != nil {
					return false, nil, err
				}
				cond, ok := condAny.(bool)
				if !ok {
					return false, nil, fmt.Errorf("given while condition must evaluate to truther")
				}
				if !cond {
					break
				}
				returned, value, err := r.executeInstructions(in.Body)
				if err != nil {
					return false, nil, err
				}
				if returned {
					return true, value, nil
				}
			}
		case ir.ForInstr:
			if _, _, err := r.executeInstructions(in.Init); err != nil {
				return false, nil, err
			}
			for {
				condAny, err := r.eval(in.Condition)
				if err != nil {
					return false, nil, err
				}
				cond, ok := condAny.(bool)
				if !ok {
					return false, nil, fmt.Errorf("given for condition must evaluate to truther")
				}
				if !cond {
					break
				}
				returned, value, err := r.executeInstructions(in.Body)
				if err != nil {
					return false, nil, err
				}
				if returned {
					return true, value, nil
				}
				if _, _, err := r.executeInstructions(in.Step); err != nil {
					return false, nil, err
				}
			}
		case ir.WithinInstr:
			iterAny, err := r.eval(in.Iterable)
			if err != nil {
				return false, nil, err
			}
			var items []any
			switch v := iterAny.(type) {
			case []any:
				items = v
			default:
				return false, nil, fmt.Errorf("within iterable is not an ordered collection")
			}
			for _, item := range items {
				r.currentFrame()[in.VarName] = item
				returned, value, err := r.executeInstructions(in.Body)
				if err != nil {
					return false, nil, err
				}
				if returned {
					return true, value, nil
				}
			}
		case ir.ReturnInstr:
			var value any
			if in.Value != nil {
				evaled, err := r.eval(in.Value)
				if err != nil {
					return false, nil, err
				}
				value = evaled
			}
			return true, value, nil
		}
	}
	return false, nil, nil
}

func (r *runner) currentFrame() map[string]any {
	if len(r.frames) == 0 {
		r.frames = append(r.frames, map[string]any{})
	}
	return r.frames[len(r.frames)-1]
}

func (r *runner) assignVar(name string, value any) {
	for i := len(r.frames) - 1; i >= 0; i-- {
		if _, ok := r.frames[i][name]; ok {
			r.frames[i][name] = value
			return
		}
	}
	r.currentFrame()[name] = value
}

func (r *runner) resolveVar(name string) (any, bool) {
	for i := len(r.frames) - 1; i >= 0; i-- {
		if v, ok := r.frames[i][name]; ok {
			return v, true
		}
	}
	return nil, false
}

func (r *runner) currentModule() string {
	if len(r.modules) == 0 {
		return ""
	}
	return r.modules[len(r.modules)-1]
}

func (r *runner) eval(expr ast.Expr) (any, error) {
	switch e := expr.(type) {
	case ast.IdentifierExpr:
		v, ok := r.resolveVar(e.Name)
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
		switch t := target.(type) {
		case []any:
			idxInt, ok := toInt64(idxAny)
			if !ok {
				return nil, fmt.Errorf("index is not an integer")
			}
			if idxInt < 0 || int(idxInt) >= len(t) {
				return nil, fmt.Errorf("index out of range")
			}
			return t[idxInt], nil
		case map[string]any:
			idxStr, ok := idxAny.(string)
			if !ok {
				return nil, fmt.Errorf("lexis index key is not a strang")
			}
			v, exists := t[idxStr]
			if !exists {
				return nil, fmt.Errorf("lexis key not found: %s", idxStr)
			}
			return v, nil
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
		case "Pront":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("Pront expects one argument")
			}
			v, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			_, _ = fmt.Fprintln(r.out, toDisplay(v))
			return nil, nil
		case "ProntOops":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("ProntOops expects one argument")
			}
			v, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			_, _ = fmt.Fprintln(r.err, toDisplay(v))
			return nil, nil
		case "FeudGirth":
			if len(e.Args) != 0 {
				return nil, fmt.Errorf("FeudGirth expects no arguments")
			}
			return int64(len(r.args)), nil
		case "FeudAt":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("FeudAt expects one argument")
			}
			idxAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			idx, ok := toInt64(idxAny)
			if !ok {
				return nil, fmt.Errorf("FeudAt index must be exact")
			}
			if idx < 0 || idx >= len(r.args) {
				return nil, fmt.Errorf("FeudAt index out of range")
			}
			return r.args[idx], nil
		case "Exeunt":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("Exeunt expects one argument")
			}
			codeAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			code, ok := toInt64(codeAny)
			if !ok {
				return nil, fmt.Errorf("Exeunt code must be exact")
			}
			if code < 0 || code > 255 {
				return nil, fmt.Errorf("Exeunt code must be between 0 and 255")
			}
			return nil, &ProgramExitError{Code: code}
		case "TrailJunction":
			if len(e.Args) == 0 {
				return nil, fmt.Errorf("TrailJunction expects at least one argument")
			}
			parts := make([]string, 0, len(e.Args))
			for i := range e.Args {
				pAny, err := r.eval(e.Args[i])
				if err != nil {
					return nil, err
				}
				p, ok := pAny.(string)
				if !ok {
					return nil, fmt.Errorf("TrailJunction arguments must be strang")
				}
				parts = append(parts, p)
			}
			return filepath.Join(parts...), nil
		case "TrailNormify":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("TrailNormify expects one argument")
			}
			pAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			p, ok := pAny.(string)
			if !ok {
				return nil, fmt.Errorf("TrailNormify argument must be strang")
			}
			return filepath.Clean(filepath.FromSlash(p)), nil
		case "HomeStead":
			if len(e.Args) != 0 {
				return nil, fmt.Errorf("HomeStead expects no arguments")
			}
			cwd, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("HomeStead failed: %w", err)
			}
			return cwd, nil
		case "DossierPresent":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("DossierPresent expects one argument")
			}
			pAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			p, ok := pAny.(string)
			if !ok {
				return nil, fmt.Errorf("DossierPresent path must be strang")
			}
			if _, err := os.Stat(p); err != nil {
				if os.IsNotExist(err) {
					return false, nil
				}
				return nil, fmt.Errorf("DossierPresent failed: %w", err)
			}
			return true, nil
		case "DossierPeruseStrang":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("DossierPeruseStrang expects one argument")
			}
			pAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			p, ok := pAny.(string)
			if !ok {
				return nil, fmt.Errorf("DossierPeruseStrang path must be strang")
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return nil, fmt.Errorf("DossierPeruseStrang failed: %w", err)
			}
			return string(b), nil
		case "DossierJotStrang":
			if len(e.Args) != 2 {
				return nil, fmt.Errorf("DossierJotStrang expects two arguments")
			}
			pAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			p, ok := pAny.(string)
			if !ok {
				return nil, fmt.Errorf("DossierJotStrang path must be strang")
			}
			cAny, err := r.eval(e.Args[1])
			if err != nil {
				return nil, err
			}
			content, ok := cAny.(string)
			if !ok {
				return nil, fmt.Errorf("DossierJotStrang content must be strang")
			}
			if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
				return nil, fmt.Errorf("DossierJotStrang failed: %w", err)
			}
			return true, nil
		case "DossierAffixStrang":
			if len(e.Args) != 2 {
				return nil, fmt.Errorf("DossierAffixStrang expects two arguments")
			}
			pAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			p, ok := pAny.(string)
			if !ok {
				return nil, fmt.Errorf("DossierAffixStrang path must be strang")
			}
			cAny, err := r.eval(e.Args[1])
			if err != nil {
				return nil, err
			}
			content, ok := cAny.(string)
			if !ok {
				return nil, fmt.Errorf("DossierAffixStrang content must be strang")
			}
			f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return nil, fmt.Errorf("DossierAffixStrang failed: %w", err)
			}
			defer f.Close()
			if _, err := f.WriteString(content); err != nil {
				return nil, fmt.Errorf("DossierAffixStrang failed: %w", err)
			}
			return true, nil
		case "Girth":
			if len(e.Args) != 1 {
				return nil, fmt.Errorf("Girth expects one argument")
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
				return nil, fmt.Errorf("Girth unsupported for value")
			}
		case "Slotify":
			if len(e.Args) != 2 {
				return nil, fmt.Errorf("Slotify expects two arguments")
			}
			tmplAny, err := r.eval(e.Args[0])
			if err != nil {
				return nil, err
			}
			tmpl, ok := tmplAny.(string)
			if !ok {
				return nil, fmt.Errorf("Slotify template must be string")
			}
			v, err := r.eval(e.Args[1])
			if err != nil {
				return nil, err
			}
			return replaceFirstDirective(tmpl, toDisplay(v)), nil
		case "Prontulate":
			comb, err := r.formatCallTemplate(e.Args, "Prontulate")
			if err != nil {
				return nil, err
			}
			_, _ = fmt.Fprintln(r.out, toDisplay(comb))
			return nil, nil
		}

		if seq, exists := r.program.Sequences[ident.Name]; exists {
			return r.invokeSequence(seq, e.Args, ident.Name)
		}

		if overloads, exists := r.program.SequenceOverloads[ident.Name]; exists {
			if seq, arityExists := overloads[len(e.Args)]; arityExists {
				return r.invokeSequence(seq, e.Args, ident.Name)
			}
		}
		if seq, exists := r.program.SequenceVariadics[ident.Name]; exists {
			if len(e.Args) >= seq.FixedParams {
				return r.invokeSequence(seq, e.Args, ident.Name)
			}
		}

		if moduleName := r.currentModule(); moduleName != "" {
			if moduleSet, exists := r.program.ModuleSequenceOverloads[moduleName]; exists {
				if nameSet, nameExists := moduleSet[ident.Name]; nameExists {
					if seq, seqExists := nameSet[len(e.Args)]; seqExists {
						return r.invokeSequence(seq, e.Args, moduleName+"."+ident.Name)
					}
				}
			}
			if moduleSet, exists := r.program.ModuleSequenceVariadics[moduleName]; exists {
				if seq, exists := moduleSet[ident.Name]; exists && len(e.Args) >= seq.FixedParams {
					return r.invokeSequence(seq, e.Args, moduleName+"."+ident.Name)
				}
			}

			if moduleSet, exists := r.program.ModuleSequences[moduleName]; exists {
				if seq, seqExists := moduleSet[ident.Name]; seqExists {
					return r.invokeSequence(seq, e.Args, moduleName+"."+ident.Name)
				}
			}
		}
	}

	if member, ok := e.Callee.(ast.MemberExpr); ok {
		if target, ok := member.Target.(ast.IdentifierExpr); ok {
			if moduleSet, exists := r.program.ModuleSequenceOverloads[target.Name]; exists {
				if nameSet, nameExists := moduleSet[member.Name]; nameExists {
					if seq, arityExists := nameSet[len(e.Args)]; arityExists {
						return r.invokeSequence(seq, e.Args, target.Name+"."+member.Name)
					}
					if varSet, hasVar := r.program.ModuleSequenceVariadics[target.Name]; hasVar {
						if seq, varExists := varSet[member.Name]; varExists && len(e.Args) >= seq.FixedParams {
							return r.invokeSequence(seq, e.Args, target.Name+"."+member.Name)
						}
					}
					return nil, fmt.Errorf("sequence call arity mismatch for %s.%s", target.Name, member.Name)
				}
			}

			if varSet, hasVar := r.program.ModuleSequenceVariadics[target.Name]; hasVar {
				if seq, varExists := varSet[member.Name]; varExists {
					if len(e.Args) < seq.FixedParams {
						return nil, fmt.Errorf("sequence call arity mismatch for %s.%s", target.Name, member.Name)
					}
					return r.invokeSequence(seq, e.Args, target.Name+"."+member.Name)
				}
			}

			if moduleSet, exists := r.program.ModuleSequences[target.Name]; exists {
				seq, seqExists := moduleSet[member.Name]
				if !seqExists {
					return nil, fmt.Errorf("undefined module sequence: %s.%s", target.Name, member.Name)
				}
				return r.invokeSequence(seq, e.Args, target.Name+"."+member.Name)
			}
			return nil, fmt.Errorf("undefined module sequence: %s.%s", target.Name, member.Name)
		}

	}

	return nil, fmt.Errorf("unsupported call expression")
}

func (r *runner) formatCallTemplate(args []ast.Expr, caller string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("%s expects at least one argument", caller)
	}
	tmplAny, err := r.eval(args[0])
	if err != nil {
		return "", err
	}
	tmpl, ok := tmplAny.(string)
	if !ok {
		return "", fmt.Errorf("%s template must be string", caller)
	}
	for i := 1; i < len(args); i++ {
		v, err := r.eval(args[i])
		if err != nil {
			return "", err
		}
		tmpl = replaceFirstDirective(tmpl, toDisplay(v))
	}
	return tmpl, nil
}

func (r *runner) invokeSequence(seq ir.SequenceIR, argExprs []ast.Expr, displayName string) (any, error) {
	if seq.Variadic {
		if len(argExprs) < seq.FixedParams {
			return nil, fmt.Errorf("sequence call arity mismatch for %s", displayName)
		}
	} else if len(argExprs) != len(seq.Params) {
		return nil, fmt.Errorf("sequence call arity mismatch for %s", displayName)
	}
	args := make([]any, 0, len(argExprs))
	for _, argExpr := range argExprs {
		v, err := r.eval(argExpr)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}

	frame := map[string]any{}
	if seq.Variadic {
		for i := 0; i < seq.FixedParams; i++ {
			frame[seq.Params[i]] = args[i]
		}
		varName := seq.Params[len(seq.Params)-1]
		rest := make([]any, 0, len(args)-seq.FixedParams)
		for i := seq.FixedParams; i < len(args); i++ {
			rest = append(rest, args[i])
		}
		frame[varName] = rest
	} else {
		for i, name := range seq.Params {
			frame[name] = args[i]
		}
	}
	r.frames = append(r.frames, frame)
	r.modules = append(r.modules, seq.Module)
	returned, value, err := r.executeInstructions(seq.Instructions)
	r.frames = r.frames[:len(r.frames)-1]
	r.modules = r.modules[:len(r.modules)-1]
	if err != nil {
		return nil, err
	}
	if returned {
		return value, nil
	}
	return nil, nil
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
		// TODO(v0.2+): extend exact runtime representation to include additional widths/sign modes.
		// Candidate families: exactly8/exactly32 and unsigned forms such as exactlyposi8/exactlyposi32.
		// Floating family alignment should include width variants such as vagly32.
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
