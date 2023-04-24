package rules

import (
	"go/ast"

	"github.com/golangci/golangci-lint/pkg/nodepass"
)

type InvalidUsageOfModifiedVariable struct{}

func (r *InvalidUsageOfModifiedVariable) ID() string {
	return "INVALID_USAGE_OF_MODIFIED_VARIABLE"
}

func (r *InvalidUsageOfModifiedVariable) Check(_ *nodepass.File, node *ast.Node) (nodepass.Report, bool) {
	if assign, ok := (*node).(*ast.AssignStmt); ok {
		if len(assign.Lhs) == 1 && len(assign.Rhs) == 2 {
			if ident, ok := assign.Lhs[0].(*ast.Ident); ok {
				if ident.Name != "_" {
					if callExpr, ok := assign.Rhs[1].(*ast.CallExpr); ok {
						if ident2, ok := assign.Rhs[0].(*ast.Ident); ok {
							if ident2.Name == "err" {
								for _, arg := range callExpr.Args {
									if ident3, ok := arg.(*ast.Ident); ok {
										if ident3.Name == ident.Name {
											return nodepass.Report{
												Pos:     assign.Pos(),
												Message: "Variable " + ident.Name + " is likely modified and later used on error. In some cases this could result in panics due to a nil dereference",
												Severity: nodepass.SeverityHigh,
												Category: r.ID(),
											}, true
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil, false
}
