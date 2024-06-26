package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var (
	inProgress atomic.Int64

	normalizeSpace = regexp.MustCompile(`\s+`)

	isDebug = flag.Bool("d", false, "print all warnings if set")
)

func entity(src string, node ast.Node) string {
	isNil := reflect.ValueOf(node).Kind() == reflect.Ptr && reflect.ValueOf(node).IsNil()
	if node == nil || isNil {
		return ""
	}
	return src[node.Pos()-1 : node.End()-1]
}

func normalize(s string) string {
	return normalizeSpace.ReplaceAllString(strings.ReplaceAll(s, "\n", ""), " ")
}

func printTokens(w io.Writer, filename string, lineNumber int, kind string, tokens ...string) {
	for i, tok := range tokens {
		tokens[i] = strings.TrimSpace(normalize(tok))
	}
	fmt.Fprintf(w, "%s:%d:%s: %s\n",
		filename, lineNumber, kind, normalize(strings.Join(tokens, " ")))
}

func joinNonEmpty(sep string, tokens ...string) string {
	j := 0
	for _, tok := range tokens {
		if tok == "" {
			continue
		}
		tokens[j] = tok
		j++
	}
	return strings.Join(tokens[:j], sep)
}

func parseSource(filename string, src string, w io.Writer) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		lnum := fset.Position(n.Pos()).Line
		switch nd := n.(type) {
		case *ast.TypeSpec:
			str := ""
			kind := ""
			name := nd.Name.Name
			var members []*ast.Field
			switch decl := nd.Type.(type) {
			case *ast.InterfaceType:
				kind = "interface"
				members = decl.Methods.List
			case *ast.StructType:
				kind = "struct"
				members = decl.Fields.List
			case *ast.FuncType:
				str := entity(src, decl)
				printTokens(w, filename, lnum, "func_type", str)
				return false
			default:
				if *isDebug {
					log.Printf("warn: type %T %s is unsupported", nd.Type, entity(src, nd.Type))
				}
				return true
			}
			for _, fld := range members {
				lnum = fset.Position(fld.Pos()).Line
				for _, fnam := range fld.Names {
					str = fmt.Sprintf("%s.%s %s",
						name, entity(src, fnam), entity(src, fld.Type))
					printTokens(w, filename, lnum, kind, str)
				}
			}
			return false
		case *ast.FuncDecl:
			str := entity(src, nd.Type)
			if nd.Recv != nil && len(nd.Recv.List) > 0 {
				recvType := nd.Recv.List[0].Type
				structName := strings.TrimPrefix(entity(src, recvType), "*")
				funcName := structName + "." + nd.Name.String()
				funcSig := fmt.Sprintf("%s%s %s",
					funcName, entity(src, nd.Type.Params), entity(src, nd.Type.Results))
				printTokens(w, filename, lnum, "method", funcSig)
				return true
			}
			printTokens(w, filename, lnum, "", str)
		case *ast.CallExpr:
			str := entity(src, nd.Fun)
			args := make([]string, len(nd.Args))
			for i, arg := range nd.Args {
				switch fl := arg.(type) {
				case *ast.FuncLit:
					args[i] = fmt.Sprintf("func (%s) %s {...}",
						entity(src, fl.Type.Params),
						entity(src, fl.Type.Results))
				default:
					args[i] = entity(src, arg)
				}
			}
			str += fmt.Sprintf("(%s)", strings.Join(args, ", "))
			printTokens(w, filename, lnum, "call", str)
			return false
		case *ast.IfStmt:
			initStr := entity(src, nd.Init)
			condStr := entity(src, nd.Cond)
			printTokens(w, filename, lnum, "stmt", "if",
				joinNonEmpty(";", initStr, condStr))
		case *ast.ForStmt:
			initStr := entity(src, nd.Init)
			condStr := entity(src, nd.Cond)
			postStr := entity(src, nd.Post)
			printTokens(w, filename, lnum, "stmt", "for",
				joinNonEmpty(";", initStr, condStr, postStr))
		}
		return true
	})
	return err
}

func parseFile(name string, w io.Writer) error {
	name, err := filepath.Abs(name)
	if err != nil {
		return err
	}
	src, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return parseSource(name, string(src), w)
}

func parseFiles(names <-chan string, w io.Writer) {
	for name := range names {
		if err := parseFile(name, w); err != nil {
			log.Printf("error: %s - %s", name, err)
		}
		inProgress.Add(-1)
	}
}

func genFilenames(args []string, r io.Reader) <-chan string {
	res := make(chan string)
	go func() {
		defer func() { close(res) }()
		if len(args) > 0 {
			for _, arg := range args {
				res <- arg
			}
			return
		}
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			res <- scanner.Text()
		}
		dieIf(scanner.Err())
	}()
	return res
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: gogr [<filename1>, ..., <filenameN>]")
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(),
			"If no filenames provided as arguments, gogr reads filenames from stdin")
		flag.PrintDefaults()
	}
	log.SetFlags(0)
}

func main() {
	flag.Parse()
	names := make(chan string)
	for i := 0; i < runtime.NumCPU()*2+1; i++ {
		go parseFiles(names, os.Stdout)
	}
	for name := range genFilenames(flag.Args(), os.Stdin) {
		inProgress.Add(1)
		names <- name
	}
	close(names)
	for inProgress.Load() > 0 {
		time.Sleep(time.Millisecond)
	}
}

func dieIf(err error) {
	if err != nil {
		log.Fatal("fatal:", err)
	}
}
