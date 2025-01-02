package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"

	"github.com/easonzhuyi/go-code-analyzer/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

type CommentRule struct {
	*analyzer.BaseRule
	patterns []*regexp.Regexp
}

func NewCommentRule() *CommentRule {
	return &CommentRule{
		BaseRule: analyzer.NewBaseRule(
			"gocomment",
			"检查注释中的不规范内容",
			"style",
		),
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?m)^[\s]*(?:func|type|var|const|package|import)\s+\w+`),
			regexp.MustCompile(`[\p{Han}]+.*[a-zA-Z]+|[a-zA-Z]+.*[\p{Han}]+`),
		},
	}
}

func (r *CommentRule) Check(fset *token.FileSet, file *ast.File, node ast.Node) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic

	// 如果传入的是整个文件，遍历所有注释组
	if file, ok := node.(*ast.File); ok {
		for _, cg := range file.Comments {
			diagnostics = append(diagnostics, r.checkCommentGroup(fset, cg)...)
		}
		return diagnostics
	}

	// 处理单个注释组
	if cg, ok := node.(*ast.CommentGroup); ok {
		return r.checkCommentGroup(fset, cg)
	}

	return diagnostics
}

// 抽取注释检查逻辑到单独的方法
func (r *CommentRule) checkCommentGroup(fset *token.FileSet, cg *ast.CommentGroup) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic
	for _, comment := range cg.List {
		for _, pattern := range r.patterns {
			if matches := pattern.FindStringSubmatch(comment.Text); len(matches) > 0 {
				pos := fset.Position(comment.Pos())
				diagnostics = append(diagnostics, analysis.Diagnostic{
					Pos: comment.Pos(),
					End: comment.End(),
					Message: fmt.Sprintf("%s:%d:%d: 发现不规范注释：%s",
						pos.Filename, pos.Line, pos.Column, matches[0]),
					Category: r.Category(),
				})
			}
		}
	}
	return diagnostics
}

func init() {
	analyzer.Register(NewCommentRule())
}
