package gui

import (
	"errors"
	"testing"

	"fyne.io/fyne/test"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextSelector(t *testing.T) {
	spec.Run(t, "Test Context Selector", testContextSelector)
}

func testContextSelector(t *testing.T, when spec.G, it spec.S) {
	var (
		testWindow        = test.NewApp().NewWindow("testingwindow")
		contextGetterStub = contextGetterStub{}
		callbacksSpy      = contextSelectorCallbacks{}
	)

	when("can retrieve the contexts", func() {
		it("calls callback with context name, when context is pressed", func() {
			contextGetterStub.getAllReturnsList = []string{"first-context", "second-context"}
			selector := NewContextSelector()
			selector.Show(testWindow, contextGetterStub, callbacksSpy.onContextSelected, callbacksSpy.onError)

			require.Len(t, selector.contextButtons, 2)
			test.Tap(selector.contextButtons[0])
			assert.Equal(t, callbacksSpy.onContextSelectedWasCalledWith, "first-context")
			assert.NoError(t, callbacksSpy.onErrorWasCalledWith)
		})
	})

	when("cannot retrieve the contexts", func() {
		it("calls error callback with error", func() {
			contextGetterStub.getAllReturnsError = errors.New("some error")
			selector := NewContextSelector()
			selector.Show(testWindow, contextGetterStub, callbacksSpy.onContextSelected, callbacksSpy.onError)

			require.Len(t, selector.contextButtons, 0)
			assert.EqualError(t, callbacksSpy.onErrorWasCalledWith, "on context select: some error")
		})
	})
}

type contextGetterStub struct {
	getAllReturnsList  []string
	getAllReturnsError error
}

func (c contextGetterStub) GetAll() ([]string, error) {
	return c.getAllReturnsList, c.getAllReturnsError
}

type contextSelectorCallbacks struct {
	onContextSelectedWasCalledWith string
	onErrorWasCalledWith           error
}

func (s *contextSelectorCallbacks) onContextSelected(contextUsed string) {
	s.onContextSelectedWasCalledWith = contextUsed
}

func (s *contextSelectorCallbacks) onError(err error) {
	s.onErrorWasCalledWith = err
}
