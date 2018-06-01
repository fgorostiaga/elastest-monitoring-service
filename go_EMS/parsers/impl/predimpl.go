package impl

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type EvalVisitor struct {
    Result bool
    Event dt.Event
    Preds map[string]common.Predicate
    ArgsMap map[striverdt.StreamName]interface{}
}

type EvalNumVisitor struct {
    Result float32
    Event dt.Event
    ArgsMap map[striverdt.StreamName]interface{}
}

func (visitor *EvalVisitor) VisitAndPredicate(p common.AndPredicate) {
    p.Left.Accept(visitor)
    rLeft := visitor.Result
    p.Right.Accept(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft && rRight
}
func (visitor *EvalVisitor) VisitTruePredicate(p common.TruePredicate) {
    visitor.Result = true
}
func (visitor *EvalVisitor) VisitFalsePredicate(p common.FalsePredicate) {
    visitor.Result = false
}
func (visitor *EvalVisitor) VisitNotPredicate(p common.NotPredicate) {
	p.Inner.Accept(visitor)
    visitor.Result = !visitor.Result
}
func (visitor *EvalVisitor) VisitOrPredicate(p common.OrPredicate) {
    p.Left.Accept(visitor)
    rLeft := visitor.Result
    p.Right.Accept(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft || rRight
}
func (visitor *EvalVisitor) VisitPathPredicate(p common.PathPredicate) {
    _,err := jsonrw.ExtractFromMap(visitor.Event.Payload, dt.JSONPath(p.Path))
	visitor.Result = err == nil
}
func (visitor *EvalVisitor) VisitStrPredicate(p common.StrPredicate) {
    strif,err := jsonrw.ExtractFromMap(visitor.Event.Payload, dt.JSONPath(p.Path))
    if err != nil {
        visitor.Result = false
        //fmt.Println("No string found in event ", visitor.Event)
        return
    }
    visitor.Result = strif.(string) == p.Expected
    //fmt.Println("Comparing",strif, "with", p.Expected, "and the result is", visitor.Result)
}
func (visitor *EvalVisitor) VisitTagPredicate(p common.TagPredicate) {
    visitor.Result = sets.SetIn(p.Tag, visitor.Event.Channels)
}
func (visitor *EvalVisitor) VisitNamedPredicate(p common.NamedPredicate) {
    if thepred, ok := visitor.Preds[p.Name]; ok {
        thepred.Accept(visitor)
    } else {
        visitor.Result = visitor.ArgsMap[striverdt.StreamName(p.Name)].(bool)
    }
}
func (visitor *EvalVisitor) VisitNumComparisonPredicate(p common.NumComparisonPredicate) {
    p.NumComparison.Accept(visitor)
}

// It also visits numcomparisons!

func (visitor *EvalVisitor) VisitNumLess(exp common.NumLess) {
    numvisitor := EvalNumVisitor{0, visitor.Event, visitor.ArgsMap}
    exp.Left.Accept(&numvisitor)
    a := numvisitor.Result
    exp.Right.Accept(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a<b
}
func (visitor *EvalVisitor) VisitNumLessEq(exp common.NumLessEq) {
    numvisitor := EvalNumVisitor{0, visitor.Event, visitor.ArgsMap}
    exp.Left.Accept(&numvisitor)
    a := numvisitor.Result
    exp.Right.Accept(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a<=b
}
func (visitor *EvalVisitor) VisitNumEq(exp common.NumEq) {
    numvisitor := EvalNumVisitor{0, visitor.Event, visitor.ArgsMap}
    exp.Left.Accept(&numvisitor)
    a := numvisitor.Result
    exp.Right.Accept(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a==b
}
func (visitor *EvalVisitor) VisitNumGreater(exp common.NumGreater) {
    numvisitor := EvalNumVisitor{0, visitor.Event, visitor.ArgsMap}
    exp.Left.Accept(&numvisitor)
    a := numvisitor.Result
    exp.Right.Accept(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a>b
}
func (visitor *EvalVisitor) VisitNumGreaterEq(exp common.NumGreaterEq) {
    numvisitor := EvalNumVisitor{0, visitor.Event, visitor.ArgsMap}
    exp.Left.Accept(&numvisitor)
    a := numvisitor.Result
    exp.Right.Accept(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a>=b
}
func (visitor *EvalVisitor) VisitNumNotEq(exp common.NumNotEq) {
    numvisitor := EvalNumVisitor{0, visitor.Event, visitor.ArgsMap}
    exp.Left.Accept(&numvisitor)
    a := numvisitor.Result
    exp.Right.Accept(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a!=b
}

// And also visits NumExps!

func (visitor *EvalNumVisitor) VisitIntLiteralExpr(exp common.IntLiteralExpr) {
    visitor.Result = float32(exp.Num)
}
func (visitor *EvalNumVisitor) VisitFloatLiteralExpr(exp common.FloatLiteralExpr) {
    visitor.Result = exp.Num
}
func (visitor *EvalNumVisitor) VisitStreamNameExpr(exp common.StreamNameExpr) {
    visitor.Result = visitor.ArgsMap[exp.StreamName].(float32)
}
func (visitor *EvalNumVisitor) VisitNumMulExpr(exp common.NumMulExpr) {
    exp.Left.Accept(visitor)
    rLeft := visitor.Result
    exp.Right.Accept(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft * rRight
}
func (visitor *EvalNumVisitor) VisitNumDivExpr(exp common.NumDivExpr) {
    exp.Left.Accept(visitor)
    rLeft := visitor.Result
    exp.Right.Accept(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft / rRight
}
func (visitor *EvalNumVisitor) VisitNumPlusExpr(exp common.NumPlusExpr) {
    exp.Left.Accept(visitor)
    rLeft := visitor.Result
    exp.Right.Accept(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft + rRight
}
func (visitor *EvalNumVisitor) VisitNumMinusExpr(exp common.NumMinusExpr) {
    exp.Left.Accept(visitor)
    rLeft := visitor.Result
    exp.Right.Accept(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft - rRight
}
func (visitor *EvalNumVisitor) VisitNumPathExpr(exp common.NumPathExpr) {
    panic("not implemented")
}