package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Soontao/pdi-util/ast/token"
)

// NewProgram type
func NewProgram(eles ...interface{}) (interface{}, error) {
	rt := GrammerNode{}

	switch len(eles) {
	case 2:
		rt["Type"] = "Program"
		if eles[0] != nil {
			rt["Statements"] = eles[0]
		}
		if eles[1] != nil {
			rt["BODefination"] = eles[1]
		}
	}

	return &rt, nil
}

// NewCommonList type
func NewCommonList(l, e interface{}) (interface{}, error) {
	rt := []*GrammerNode{e.(*GrammerNode)}
	if l != nil {
		rt = append(*l.(*[]*GrammerNode), e.(*GrammerNode))
	}
	return &rt, nil
}

// NewKeyValueList type
func NewKeyValueList(i, l interface{}) (interface{}, error) {
	rt := []*GrammerNode{i.(*GrammerNode)}
	if l != nil {
		rt = append(*l.(*[]*GrammerNode), i.(*GrammerNode))
	}
	return &rt, nil
}

// NewComplexValue type
func NewComplexValue(v interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "ComplexValue"}
	rt["Value"] = v
	return &rt, nil
}

// NewKeyValuePair type
func NewKeyValuePair(k, v interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "KeyValuePair"}
	rt["Key"] = k
	rt["Value"] = v
	return &rt, nil
}

// NewRaiseExpr type
func NewRaiseExpr(ids interface{}) (interface{}, error) {
	return &GrammerNode{"Type": "RaiseExpr", "Messages": ids}, nil
}

// NewDataType type
func NewDataType(id, ns interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "DataType"}
	if id != nil {
		rt["Identifier"] = id
	}
	if ns != nil {
		rt["Namespace"] = ns
	}
	return &rt, nil
}

// NewAnnotation type
func NewAnnotation(id, paramNamespace, paramText interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "Annotation"}
	if id != nil {
		rt["AnnotationName"] = id
	}
	if paramNamespace != nil {
		rt["ParamIdentifier"] = paramNamespace
	}
	if paramText != nil {
		rt["ParamText"] = paramText
	}
	return &rt, nil
}

// NewAnnotationList type
func NewAnnotationList(e, l interface{}) (interface{}, error) {
	rt := []*GrammerNode{}
	if l != nil {
		rt = append((*(l.(*[]*GrammerNode))), e.(*GrammerNode))
	} else {
		rt = []*GrammerNode{e.(*GrammerNode)}
	}
	return &rt, nil
}

// NewActionItem type
func NewActionItem(action, raises interface{}) interface{} {
	return &GrammerNode{"Type": "ActionItem", "Action": action, "Raises": raises}
}

// NewBODefination type
func NewBODefination(annnotations, name, raises, elements interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "BusinessObjectDefination"}
	if annnotations != nil {
		rt["Annotations"] = annnotations
	}
	rt["BOName"] = name
	if raises != nil {
		rt["Raises"] = raises
	}
	if elements != nil {
		rt["Elements"] = elements
	}

	return &rt, nil
}

// NewObjectFieldList type
func NewObjectFieldList(l, e interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "BusinessObjectFieldList"}
	if l != nil {
		sList := (*l.(*GrammerNode))["Fields"]
		if sList == nil {
			sList = []*GrammerNode{}
		}
		rt["Fields"] = append(sList.([]*GrammerNode), e.(*GrammerNode))
	} else {
		rt["Fields"] = []*GrammerNode{e.(*GrammerNode)}
	}
	return &rt, nil
}

// NewCondition type
func NewCondition(v1, v2 interface{}) interface{} {
	return &GrammerNode{
		"Type":  "Condition",
		"Left":  v1,
		"Right": v2,
	}
}

// NewAssociationItem type
func NewAssociationItem(id, multiplicity, target, valuation interface{}) interface{} {
	return &GrammerNode{
		"Type":         "AssociationItem",
		"Name":         id,
		"Multiplicity": multiplicity,
		"Target":       target,
		"Valuation":    valuation,
	}
}

// NewForeachStmt type
func NewForeachStmt(loopObject, itemID, stmts interface{}) interface{} {
	return &GrammerNode{
		"Type":       "ForeachStatement",
		"LoopObject": loopObject,
		"ItemID":     itemID,
		"Statements": stmts,
	}
}

// NewIfStmt type
func NewIfStmt(exp, stmts interface{}) interface{} {
	return &GrammerNode{
		"Type":        "IfStatement",
		"IFCondition": exp,
		"Statements":  stmts,
	}
}

// NewFuncCallExpr type
func NewFuncCallExpr(function, args interface{}) interface{} {
	return &GrammerNode{
		"Type":      "FuncCallExpr",
		"Function":  function,
		"Arguments": args,
	}
}

// NewVariableDeclarationStmt type
func NewVariableDeclarationStmt(name interface{}) interface{} {
	return &GrammerNode{
		"Type":    "VariableDeclarationStmt",
		"VarName": name,
	}
}

// GeneralAddProperty to grammer node
func GeneralAddProperty(n interface{}, name string, value interface{}) interface{} {
	(*n.(*GrammerNode))[name] = value
	return n
}

// NewMessageItem type
func NewMessageItem(id, template, types interface{}) interface{} {
	return &GrammerNode{
		"Type":       "MessageItem",
		"Name":       id,
		"Template":   template,
		"ValueTypes": types,
	}
}

// NewAnnotatedBOItem type
func NewAnnotatedBOItem(annotations, element interface{}) (interface{}, error) {
	rt := *element.(*GrammerNode)
	if annotations != nil {
		rt["Annotations"] = annotations
	}

	return &rt, nil
}

// NewElementItem type
func NewElementItem(tokens ...interface{}) (interface{}, error) {

	rt := GrammerNode{"Type": "ElementItem"}
	rt["FieldIdentifier"] = tokens[0]
	rt["FieldType"] = tokens[1]

	if len(tokens) == 3 && tokens[2] != nil {
		rt["DefaultValue"] = tokens[2]
	}

	return &rt, nil
}

// NewExpression type
func NewExpression(op interface{}, left, right interface{}) interface{} {
	return &GrammerNode{
		"Type":     "Expression",
		"Operator": op,
		"Left":     left,
		"Right":    right,
	}
}

// NewBoolValue type
func NewBoolValue(t interface{}) (interface{}, error) {
	v, e := strconv.ParseBool(string((t.(*token.Token).Lit)))
	return &GrammerNode{"Type": "BoolValue", "Value": v}, e
}

// NewObjectReference type
func NewObjectReference(v interface{}) (interface{}, error) {
	return &GrammerNode{"Type": "ObjectReference", "Value": v}, nil
}

// NewFunctionCallExpr type
func NewFunctionCallExpr(v interface{}) (interface{}, error) {
	return &GrammerNode{"Type": "FunctionCallExpression", "Value": v}, nil
}

// NewStringValue value
func NewStringValue(t interface{}) (interface{}, error) {
	return &GrammerNode{"Type": "StringValue", "Value": strings.Trim(string((t.(*token.Token).Lit)), `"`)}, nil
}

// ConvertToString type
func ConvertToString(t interface{}) interface{} {
	return string(t.(*token.Token).Lit)
}

// NewNumberValue value
func NewNumberValue(t interface{}) (interface{}, error) {
	v, e := strconv.ParseFloat(string((t.(*token.Token).Lit)), 64)
	return &GrammerNode{"Type": "NumberValue", "Value": v}, e
}

// NewMultiplicity type
func NewMultiplicity(i1, i2 interface{}) (interface{}, error) {
	return &GrammerNode{"Type": "Multiplicity", "i1": i1, "i2": i2}, nil
}

// NewBusinessObjectNode type
func NewBusinessObjectNode(tokens ...interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "BusinessObjectNode"}

	if len(tokens) == 1 {
		return nil, fmt.Errorf("Node must have at least one element")
	}

	rt["NodeName"] = tokens[0]
	rt["NodeElements"] = tokens[1]

	if len(tokens) == 3 && tokens[2] != nil {
		rt["Multiplicity"] = tokens[2]
	}
	if len(tokens) == 4 && tokens[3] != nil {
		rt["Raises"] = tokens[3]
	}

	return &rt, nil
}

// NewStatementList type
func NewStatementList(s, l interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "StatementList"}
	if l != nil {
		sList := (*l.(*GrammerNode))["Statements"]
		if sList == nil {
			sList = []*GrammerNode{}
		}
		rt["Statements"] = append(sList.([]*GrammerNode), s.(*GrammerNode))
	} else {
		rt["Statements"] = []*GrammerNode{s.(*GrammerNode)}
	}
	return &rt, nil
}

// NewStatement type
func NewStatement(n interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "Statement"}
	rt["Content"] = n
	return &rt, nil
}

// NewImportDeclaration type
func NewImportDeclaration(n, i interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "ImportDeclaration"}
	if n != nil {
		rt["Namespace"] = n
	}
	if i != nil {
		rt["Alias"] = i
	}
	return &rt, nil
}

// NewSelector type
func NewSelector(id, sub interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "Selector"}
	if sub != nil {
		rt["SubSelector"] = sub
	}
	if id != nil {
		rt["Identifier"] = id
	}
	return &rt, nil
}

// NewNamespace type
func NewNamespace(id, sub interface{}) (interface{}, error) {
	rt := GrammerNode{"Type": "Namespace"}
	if sub != nil {
		rt["SubNamespace"] = sub
	}
	if id != nil {
		rt["Identifier"] = id
	}
	return &rt, nil
}

// NewIdentifier type
func NewIdentifier(i interface{}) (interface{}, error) {
	return &GrammerNode{
		"Type": "Identifier",
		"ID":   string(i.(*token.Token).Lit),
	}, nil
}
