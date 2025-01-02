package analyzer

import (
    "go/ast"
    "go/token"
    "golang.org/x/tools/go/analysis"
)

// Rule 定义检查规则接口
type Rule interface {
    Name() string
    Description() string
    Category() string
    Check(fset *token.FileSet, file *ast.File, node ast.Node) []analysis.Diagnostic
}

// BaseRule 提供基础实现
type BaseRule struct {
    name        string
    description string
    category    string
}

func NewBaseRule(name, description, category string) *BaseRule {
    return &BaseRule{
        name:        name,
        description: description,
        category:    category,
    }
}

func (r *BaseRule) Name() string        { return r.name }
func (r *BaseRule) Description() string { return r.description }
func (r *BaseRule) Category() string    { return r.category }
