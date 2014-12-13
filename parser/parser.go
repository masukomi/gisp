package parser

import (
	"fmt"
	"github.com/masukomi/gisp/lexer"
	"go/token"
)

type Pos int

func (p Pos) Position() Pos {
	return p
}

var nilNode = NewIdentNode("nil")

func ParseFromString(name, program string) []Node {
	return Parse(lexer.Lex(name, program))
}

func Parse(l *lexer.Lexer) []Node {
	return parser(l, make([]Node, 0), ' ')
}

func parser(l *lexer.Lexer, tree []Node, lookingFor rune) []Node {
	for item := l.NextItem(); item.Type != lexer.ItemEOF; {
		switch t := item.Type; t {
		case lexer.ItemIdent:
			tree = append(tree, NewIdentNode(item.Value))
		case lexer.ItemString:
			tree = append(tree, newStringNode(item.Value))
		case lexer.ItemInt:
			tree = append(tree, newIntNode(item.Value))
		case lexer.ItemFloat:
			tree = append(tree, newFloatNode(item.Value))
		case lexer.ItemComplex:
			tree = append(tree, newComplexNode(item.Value))
		case lexer.ItemLeftParen:
			tree = append(tree, newCallNode(parser(l, make([]Node, 0), ')')))
		case lexer.ItemLeftVect:
			tree = append(tree, newVectNode(parser(l, make([]Node, 0), ']')))
		case lexer.ItemRightParen:
			if lookingFor != ')' {
				panic(fmt.Sprintf("unexpected \")\" [%d]", item.Pos))
			}
			return tree
		case lexer.ItemRightVect:
			if lookingFor != ']' {
				panic(fmt.Sprintf("unexpected \"]\" [%d]", item.Pos))
			}
			return tree
		case lexer.ItemError:
			println(item.Value)
		default:
			panic("Bad Item type")
		}
		item = l.NextItem()
	}

	return tree
}

func NewIdentNode(name string) *IdentNode {
	return &IdentNode{NodeType: NodeIdent, Ident: name}
}

func newStringNode(val string) *StringNode {
	return &StringNode{NodeType: NodeString, Value: val}
}

func newIntNode(val string) *NumberNode {
	return &NumberNode{NodeType: NodeNumber, Value: val, NumberType: token.INT}
}

func newFloatNode(val string) *NumberNode {
	return &NumberNode{NodeType: NodeNumber, Value: val, NumberType: token.FLOAT}
}

func newComplexNode(val string) *NumberNode {
	return &NumberNode{NodeType: NodeNumber, Value: val, NumberType: token.IMAG}
}

// We return Node here, because it could be that it's nil
func newCallNode(args []Node) Node {
	if len(args) > 0 {
		return &CallNode{NodeType: NodeCall, Callee: args[0], Args: args[1:]}
	} else {
		return nilNode
	}
}

func newVectNode(content []Node) *VectorNode {
	return &VectorNode{NodeType: NodeVector, Nodes: content}
}
