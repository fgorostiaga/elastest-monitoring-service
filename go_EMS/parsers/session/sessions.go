package session

import(
//	"errors"
	"fmt"
//	"log"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type Filter struct {
	Pred common.Predicate
	Tag dt.Channel
}

type Filters struct {
	Defs []Filter
}

type Monitor Filters

func newFiltersNode(defs interface{}) (Filters) {
	parsed_defs := common.ToSlice(defs)
	ds := make([]Filter, len(parsed_defs))
	for i,v := range parsed_defs {
		ds[i] = v.(Filter)
	}
	return Filters{ds}
}

// Monitoring Machines

type MoMVisitor interface {
	VisitFilter(Filter)
	VisitSession(Session)
	VisitStream(Stream)
	VisitTrigger(Trigger)
	VisitPredicateDecl(PredicateDecl)
}

type MoM interface {
    Accept(MoMVisitor)
}


//
// Action
//
type EmitAction struct {
	StreamName striverdt.StreamName
	TagName    common.Tag
}
// func (a EmitAction) Sprint() string {
// 	return fmt.Sprintf("emit %s on %s\n",a.StreamName,a.TagName.Tag)
// }

type Trigger struct {
	Pred    common.Predicate
	Action  EmitAction
}

func (this Trigger) Accept(visitor MoMVisitor) {
    visitor.VisitTrigger(this)
}

func newEmitAction(n, t interface{}) (EmitAction) {
	name := n.(common.Identifier).Val
	tag  := t.(common.Tag)

	return EmitAction{striverdt.StreamName(name),tag}
}

func newTrigger(p, a interface{}) (Trigger) {
	pred:= p.(common.Predicate)
	act := a.(EmitAction)

	return Trigger{pred,act}
}


type StreamType int
const (
	NumT  StreamType   = iota
	BoolT
	StringT
	LastType = StringT
)

// func (t StreamType) Sprint() string {
// 
// 	type_names := []string{"int","bool","string"}
// 
// 	// str string 
// 	// switch t {
// 	// case Int:
// 	// 	str = "int"
// 	// case Bool:
// 	// 	str = "bool"
// 	// case String:
// 	// 	str = "string"
// 	// }
// 	// return str
// 
// 	if t>=LastType { return "" }
// 	return fmt.Sprintf("%s",type_names[t])
// }


//
// We need a dictionary of streams (so all streams used are defined)
//
type Stream struct { // a Stream is a Name:=Expr
	Type StreamType
	Name striverdt.StreamName
	Expr StreamExpr
}

func (this Stream) Accept(visitor MoMVisitor) {
    visitor.VisitStream(this)
}

type Session struct {
	Name  striverdt.StreamName
	Begin common.Predicate
	End   common.Predicate
}

func (this Session) Accept(visitor MoMVisitor) {
    visitor.VisitSession(this)
}

//
// Expressions
//

type StreamExprVisitor interface {
	VisitAggregatorExpr(AggregatorExpr)
	VisitIfThenExpr(IfThenExpr)
	VisitIfThenElseExpr(IfThenElseExpr)
	// VisittringExpr(StringExpr)
	VisitNumStreamExpr(NumStreamExpr)
	VisitPredExpr(PredExpr)
	VisitStringPathExpr(StringPathExpr)
	VisitPrevExpr(PrevExpr)
}

type StreamExpr interface {
	// add functions here
    Accept (StreamExprVisitor)
}

type NumStreamExpr struct {
    Expr common.NumExpr
}
func (this NumStreamExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitNumStreamExpr(this)
}


type AggregatorExpr struct {
	Operation string
	Stream    striverdt.StreamName //StreamName
	Session   striverdt.StreamName //StreamName
}

func (this AggregatorExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitAggregatorExpr(this)
}

type IfThenExpr struct {
	If   common.Predicate
	Then StreamExpr
}

func (this IfThenExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitIfThenExpr(this)
}

type IfThenElseExpr struct {
	If   common.Predicate
	Then StreamExpr
	Else StreamExpr
}

func (this IfThenElseExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitIfThenElseExpr(this)
}

// Is this ever used?
//type StringExpr struct {
//	Path string// so far only e.get(path) claiming to return a string
//}

type PredExpr struct {
	Pred common.Predicate
}

func (this PredExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitPredExpr(this)
}

type PrevExpr struct {
	Stream string
}

func (this PrevExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitPrevExpr(this)
}

type StringPathExpr struct {
	Path dt.JSONPath
}

func (this StringPathExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitStringPathExpr(this)
}

//
// Expression Node constructors
//
func newNumStreamExpr(n interface{}) NumStreamExpr {
	return NumStreamExpr{n.(common.NumExpr)}
}
func newAggregatorExpr(op, str, ses interface{}) AggregatorExpr {
	operation := op.(string)
	stream    := str.(common.Identifier).Val
	session   := ses.(common.Identifier).Val

	return AggregatorExpr{operation,striverdt.StreamName(stream),striverdt.StreamName(session)}
}

func newIfThenExpr(p,e interface{}) IfThenExpr {
	if_part   := p.(common.Predicate)
	then_part := e.(StreamExpr)
	return IfThenExpr{if_part,then_part}
}
func newIfThenElseExpr(p,a,b interface{}) IfThenElseExpr {
	if_part   := p.(common.Predicate)
	then_part := a.(StreamExpr)
	else_part := b.(StreamExpr)
	return IfThenElseExpr{if_part, then_part, else_part}
}
func newPredExpr(p interface{}) PredExpr {
	return PredExpr{p.(common.Predicate)}
}

func newStringPathExpr(p interface{}) StringPathExpr {
	path := p.(common.PathName).Val
	return StringPathExpr{dt.JSONPath(path)}
}
func newPrevExpr(p interface{}) (PrevExpr) {
	return PrevExpr{p.(string)}
}

//
// Declaration Node constructors
//

func newStreamDeclaration(t,n,e interface{}) Stream {
	the_type := t.(StreamType)
	name     := n.(common.Identifier).Val
	expr     := e.(StreamExpr)
	return Stream{the_type,striverdt.StreamName(name),expr}
}

func newSessionDeclaration(n,b,e interface{}) Session {
	name  := n.(common.Identifier).Val
	begin := b.(common.Predicate)
	end   := e.(common.Predicate)
	return Session{striverdt.StreamName(name),begin,end}
}

func newPredicateDeclaration(n,p interface{}) PredicateDecl {
	name := n.(common.Identifier).Val
	pred := p.(common.Predicate)
	return PredicateDecl{name,pred}
}

type PredicateDecl struct {
	Name string
	Pred common.Predicate
}

func (this PredicateDecl) Accept(visitor MoMVisitor) {
    visitor.VisitPredicateDecl(this)
}

type MonitorMachine struct {
	Stampers []Filter
	Sessions []Session
	Streams  []Stream
	Triggers []Trigger
	Preds    []PredicateDecl
}

//
// after ParseFile returns a []interafce{} all the elements in the slice
//    are Filer or Session or Stream or Trigger
//    this function creates a MonitorMachine from such a mixed slice
//

func ProcessDeclarations(ds []interface{}) MonitorMachine {
	machine := MonitorMachine{}
	for _,v := range ds {
		switch val := v.(type) {
		case Filter:
			machine.Stampers = append(machine.Stampers,val)
		case Session:
			machine.Sessions = append(machine.Sessions,val)
		case Stream:
			machine.Streams  = append(machine.Streams,val)
		case Trigger:
			machine.Triggers = append(machine.Triggers,val)
		case PredicateDecl:
			machine.Preds    = append(machine.Preds,val)
		}
	}

	// FIXME
	// Additionally, we should check that all
	// used streams are defined and that there is no circularity

	return machine
}

func Print(mon MonitorMachine) {
	fmt.Printf("There are %d stampers\n",len(mon.Stampers))
	fmt.Printf("There are %d sessions\n",len(mon.Sessions))
	fmt.Printf("There are %d streams\n", len(mon.Streams))
	fmt.Printf("There are %d triggers\n", len(mon.Triggers))
	fmt.Printf("There are %d predicates\n", len(mon.Preds))

	for _,v := range mon.Stampers {
		fmt.Printf("when %s do %s\n", v.Pred, v.Tag)
	}
	for _,v := range mon.Sessions {
		fmt.Printf("session %s := (begin=>%s,end=>%s)\n",v.Name,v.Begin,v.End)
	}
	for _,v := range mon.Streams {
		fmt.Printf("stream %s %s := %s\n",v.Type,v.Name,v.Expr)
	}
	for _,v := range mon.Triggers {
		fmt.Printf("trigger %s do %s\n",v.Pred,v.Action)
	}

}


