package errns

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrs(t *testing.T) {
	t.Parallel()

	errNs := NewErrorNamespace("TestErrs")
	errType := errNs.NewType("SomeType")
	errNew := errType.New("error creation with new")
	require.Equal(t, "TestErrs: SomeType: error creation with new", errNew.Error())
	require.True(t, errors.Is(errNew, errType))

	externalError := errors.New("some external error")
	errWrap := errType.Wrap(externalError, "Context")
	require.Equal(t, "TestErrs: SomeType: Context: some external error", errWrap.Error())
	require.True(t, errors.Is(errWrap, externalError))
	require.True(t, errors.Is(errWrap, errType))
}
