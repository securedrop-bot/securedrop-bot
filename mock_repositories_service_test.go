// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package securedropbot

import (
	"context"
	"github.com/google/go-github/github"
	"sync"
)

var (
	lockrepositoriesServiceMockGetCombinedStatus sync.RWMutex
)

// repositoriesServiceMock is a mock implementation of repositoriesService.
//
//     func TestSomethingThatUsesrepositoriesService(t *testing.T) {
//
//         // make and configure a mocked repositoriesService
//         mockedrepositoriesService := &repositoriesServiceMock{
//             GetCombinedStatusFunc: func(in1 context.Context, in2 string, in3 string, in4 string, in5 *github.ListOptions) (*github.CombinedStatus, *github.Response, error) {
// 	               panic("TODO: mock out the GetCombinedStatus method")
//             },
//         }
//
//         // TODO: use mockedrepositoriesService in code that requires repositoriesService
//         //       and then make assertions.
//
//     }
type repositoriesServiceMock struct {
	// GetCombinedStatusFunc mocks the GetCombinedStatus method.
	GetCombinedStatusFunc func(in1 context.Context, in2 string, in3 string, in4 string, in5 *github.ListOptions) (*github.CombinedStatus, *github.Response, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetCombinedStatus holds details about calls to the GetCombinedStatus method.
		GetCombinedStatus []struct {
			// In1 is the in1 argument value.
			In1 context.Context
			// In2 is the in2 argument value.
			In2 string
			// In3 is the in3 argument value.
			In3 string
			// In4 is the in4 argument value.
			In4 string
			// In5 is the in5 argument value.
			In5 *github.ListOptions
		}
	}
}

// GetCombinedStatus calls GetCombinedStatusFunc.
func (mock *repositoriesServiceMock) GetCombinedStatus(in1 context.Context, in2 string, in3 string, in4 string, in5 *github.ListOptions) (*github.CombinedStatus, *github.Response, error) {
	if mock.GetCombinedStatusFunc == nil {
		panic("moq: repositoriesServiceMock.GetCombinedStatusFunc is nil but repositoriesService.GetCombinedStatus was just called")
	}
	callInfo := struct {
		In1 context.Context
		In2 string
		In3 string
		In4 string
		In5 *github.ListOptions
	}{
		In1: in1,
		In2: in2,
		In3: in3,
		In4: in4,
		In5: in5,
	}
	lockrepositoriesServiceMockGetCombinedStatus.Lock()
	mock.calls.GetCombinedStatus = append(mock.calls.GetCombinedStatus, callInfo)
	lockrepositoriesServiceMockGetCombinedStatus.Unlock()
	return mock.GetCombinedStatusFunc(in1, in2, in3, in4, in5)
}

// GetCombinedStatusCalls gets all the calls that were made to GetCombinedStatus.
// Check the length with:
//     len(mockedrepositoriesService.GetCombinedStatusCalls())
func (mock *repositoriesServiceMock) GetCombinedStatusCalls() []struct {
	In1 context.Context
	In2 string
	In3 string
	In4 string
	In5 *github.ListOptions
} {
	var calls []struct {
		In1 context.Context
		In2 string
		In3 string
		In4 string
		In5 *github.ListOptions
	}
	lockrepositoriesServiceMockGetCombinedStatus.RLock()
	calls = mock.calls.GetCombinedStatus
	lockrepositoriesServiceMockGetCombinedStatus.RUnlock()
	return calls
}