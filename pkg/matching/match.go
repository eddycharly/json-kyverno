package matching

import (
	"strings"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/templating"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

//nolint:errname
type Result struct {
	field.ErrorList
	Message string
}

func (r Result) Error() string {
	var lines []string
	if r.Message != "" {
		lines = append(lines, "-> "+r.Message)
	}
	for _, err := range r.ErrorList {
		lines = append(lines, " -> "+err.Error())
	}
	return strings.Join(lines, "\n")
}

//nolint:errname
type Results []Result

func (r Results) Error() string {
	var lines []string
	for _, err := range r {
		lines = append(lines, err.Error())
	}
	return strings.Join(lines, "\n")
}

func MatchAssert(path *field.Path, match v1alpha1.Assert, actual any, bindings binding.Bindings, compiler templating.Compiler) ([]Result, error) {
	if len(match.Any) == 0 && len(match.All) == 0 {
		return nil, field.Invalid(path, match, "an empty assert is not valid")
	} else {
		if len(match.Any) != 0 {
			var fails []Result
			path := path.Child("any")
			for i, assertion := range match.Any {
				path := path.Index(i).Child("check")
				parsed, err := assertion.Check.Assertion(compiler)
				if err != nil {
					return fails, err
				}
				checkFails, err := parsed.Assert(path, actual, bindings)
				if err != nil {
					return fails, err
				}
				if len(checkFails) == 0 {
					fails = nil
					break
				}
				fail := Result{
					ErrorList: checkFails,
				}
				if assertion.Message != nil {
					fail.Message = assertion.Message.Format(actual, bindings, compiler.Options().Jp...)
				}
				fails = append(fails, fail)
			}
			if fails != nil {
				return fails, nil
			}
		}
		if len(match.All) != 0 {
			var fails []Result
			path := path.Child("all")
			for i, assertion := range match.All {
				path := path.Index(i).Child("check")
				parsed, err := assertion.Check.Assertion(compiler)
				if err != nil {
					return fails, err
				}
				checkFails, err := parsed.Assert(path, actual, bindings)
				if err != nil {
					return fails, err
				}
				if len(checkFails) > 0 {
					fail := Result{
						ErrorList: checkFails,
					}
					if assertion.Message != nil {
						fail.Message = assertion.Message.Format(actual, bindings, compiler.Options().Jp...)
					}
					fails = append(fails, fail)
				}
			}
			return fails, nil
		}
		return nil, nil
	}
}

func Match(path *field.Path, match *v1alpha1.Match, actual any, bindings binding.Bindings, compiler templating.Compiler) (field.ErrorList, error) {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return nil, field.Invalid(path, match, "an empty match is not valid")
	} else {
		var errs field.ErrorList
		if len(match.Any) != 0 {
			_errs, err := MatchAny(path.Child("any"), match.Any, actual, bindings, compiler)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		if len(match.All) != 0 {
			_errs, err := MatchAll(path.Child("all"), match.All, actual, bindings, compiler)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		return errs, nil
	}
}

func MatchAny(path *field.Path, assertions []v1alpha1.AssertionTree, actual any, bindings binding.Bindings, compiler templating.Compiler) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range assertions {
		path := path.Index(i)
		assertion, err := assertion.Assertion(compiler)
		if err != nil {
			return errs, err
		}
		_errs, err := assertion.Assert(path, actual, bindings)
		if err != nil {
			return errs, err
		}
		if len(_errs) == 0 {
			return nil, nil
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}

func MatchAll(path *field.Path, assertions []v1alpha1.AssertionTree, actual any, bindings binding.Bindings, compiler templating.Compiler) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range assertions {
		path := path.Index(i)
		assertion, err := assertion.Assertion(compiler)
		if err != nil {
			return errs, err
		}
		_errs, err := assertion.Assert(path, actual, bindings)
		if err != nil {
			return errs, err
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}
