package debug

import (
	"fmt"
	"regexp"
	"strconv"
)

// Extract constraint index from gnark error string
func ExtractConstraintIndex(err error) (int, bool) {
	re := regexp.MustCompile(`constraint #(\d+)`)
	matches := re.FindStringSubmatch(err.Error())

	if len(matches) < 2 {
		return 0, false
	}

	idx, _ := strconv.Atoi(matches[1])
	return idx, true
}

// Map constraint → step
func ConstraintToStep(constraintIdx int, constraintsPerStep int) int {
	return constraintIdx / constraintsPerStep
}

// Pretty print failure
func PrintFailure(err error, constraintsPerStep int) {
	idx, ok := ExtractConstraintIndex(err)
	if !ok {
		fmt.Println("Could not parse constraint index")
		return
	}

	step := ConstraintToStep(idx, constraintsPerStep)

	fmt.Println("❌ Constraint Failure Analysis")
	fmt.Printf("Constraint #: %d\n", idx)
	fmt.Printf("Likely Step : %d\n", step)
}
