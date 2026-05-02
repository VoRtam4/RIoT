package internal

import (
	"fmt"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func buildTagFilterFlux(node *sharedModel.FilterNode) (string, bool) {
	if node == nil {
		return "", false
	}
	expr := buildFilterExpr(node)
	if expr == "" {
		return "", false
	}
	return fmt.Sprintf(`|> filter(fn: (r) => %s)`, expr), true
}

func buildFilterExpr(node *sharedModel.FilterNode) string {
	if node == nil {
		return ""
	}
	switch node.Type {
	case sharedModel.FilterNodeTypeRule:
		return buildRuleExpr(node.Rule)
	case sharedModel.FilterNodeTypeLogical:
		if node.Operator == sharedModel.LogicalNot {
			if len(node.Nodes) != 1 {
				return ""
			}
			inner := buildFilterExpr(&node.Nodes[0])
			if inner == "" {
				return ""
			}
			return fmt.Sprintf("not (%s)", inner)
		}
		parts := make([]string, 0, len(node.Nodes))
		for _, n := range node.Nodes {
			expr := buildFilterExpr(&n)
			if expr != "" {
				parts = append(parts, expr)
			}
		}
		if len(parts) == 0 {
			return ""
		}
		op := "and"
		if node.Operator == sharedModel.LogicalOr {
			op = "or"
		}
		return "(" + strings.Join(parts, " "+op+" ") + ")"
	}
	return ""
}

func buildRuleExpr(rule *sharedModel.FilterRule) string {
	if rule == nil {
		return ""
	}

	tag := fmt.Sprintf(`r["%s"]`, rule.Tag)

	switch rule.Operator {
	case sharedModel.OpEQ:
		return fmt.Sprintf(`%s == %q`, tag, rule.Value)

	case sharedModel.OpNEQ:
		return fmt.Sprintf(`%s != %q`, tag, rule.Value)

	case sharedModel.OpContains:
		return fmt.Sprintf(`strings.containsStr(v: %s, substr: %q)`, tag, rule.Value)

	case sharedModel.OpPrefix:
		return fmt.Sprintf(`strings.hasPrefix(v: %s, prefix: %q)`, tag, rule.Value)

	case sharedModel.OpSuffix:
		return fmt.Sprintf(`strings.hasSuffix(v: %s, suffix: %q)`, tag, rule.Value)

	case sharedModel.OpRegex:
		return fmt.Sprintf(`%s =~ /%s/`, tag, escapeFluxRegex(rule.Value))

	case sharedModel.OpIn:
		values := strings.Split(rule.Value, ",")
		parts := make([]string, 0, len(values))
		for _, v := range values {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			parts = append(parts, fmt.Sprintf(`%s == %q`, tag, v))
		}
		if len(parts) == 0 {
			return ""
		}
		return "(" + strings.Join(parts, " or ") + ")"
	}

	return ""
}

func escapeFluxRegex(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `/`, `\/`)
	return s
}
