package analyzer

import (
    "golang.org/x/tools/go/analysis"
)

func RunAsPlugin() {
    var analyzers []*analysis.Analyzer
    for _, rule := range GetAllRules() {
        analyzers = append(analyzers, &analysis.Analyzer{
            Name: rule.Name(),
            Doc:  rule.Description(),
            Run: func(pass *analysis.Pass) (interface{}, error) {
                for _, file := range pass.Files {
                    for _, diagnostic := range rule.Check(pass.Fset, file, file) {
                        pass.Report(diagnostic)
                    }
                }
                return nil, nil
            },
        })
    }
}
