package errors

import "fmt"

func ReportMultipleErrors(errs []error) error {
	wrappedErrs := errs[0]
	for _, err := range errs {
		if err != nil {
			wrappedErrs = fmt.Errorf("%w; %v", wrappedErrs, err)
		}
	}
	return wrappedErrs
}
