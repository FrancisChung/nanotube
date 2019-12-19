// Package rules provides primitives for working with routing rules.
package rules

import (
	"fmt"
	"nanotube/pkg/conf"
	"nanotube/pkg/target"
	"regexp"

	"github.com/pkg/errors"
)

// Rules represent all the routing rules/routing table.
type Rules []Rule

// Rule is a routing rule.
type Rule struct {
	Regexs     []string
	Targets    []*target.Cluster
	Continue   bool
	CompiledRE []*regexp.Regexp
}

// Build reads rules from config, compiles them.
func Build(crs conf.Rules, clusters target.Clusters) (Rules, error) {
	var rs Rules
	for _, cr := range crs.Rule {
		r := Rule{
			Regexs: cr.Regexs,
		}
		for _, clName := range cr.Clusters {
			cl, ok := clusters[clName]
			if !ok {
				return rs,
					fmt.Errorf("got non-existent cluster name %s in the rules config",
						clName)
			}
			r.Targets = append(r.Targets, cl)
		}
		r.Continue = cr.Continue

		rs = append(rs, r)
	}

	err := rs.Compile()
	if err != nil {
		return rs, errors.Wrap(err, "rules compilation failed :")
	}

	return rs, nil
}

// Compile precompiles regexps for perf and performs validation.
func (rs Rules) Compile() error {
	for i := range rs {
		rs[i].CompiledRE = make([]*regexp.Regexp, 0)
		for _, re := range rs[i].Regexs {
			cre, err := regexp.Compile(re)
			if err != nil {
				return errors.Wrapf(err, "compiling regex %s failed", cre)
			}
			rs[i].CompiledRE = append(rs[i].CompiledRE, cre)
		}
	}

	return nil
}