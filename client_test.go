package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	mocks "github.com/tuannguyensn2001/aurora-go/mocks"
)

func TestNewClient(t *testing.T) {
	mockFetcher := new(mocks.MockFetcher)
	mockFetcher.On("IsStatic").Return(true)
	mockFetcher.On("Fetch", mock.Anything).Return(make(map[string]auroratype.Parameter), nil)
	mockFetcher.On("FetchExperiments", mock.Anything).Return(nil, nil)

	client := NewClient(NewFetcherStorage(mockFetcher), ClientOptions{})

	assert.NotNil(t, client)
	assert.NotNil(t, client.engine)
	assert.NotNil(t, client.storage)
}

func TestClientGetParameter(t *testing.T) {
	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue: "matched_value",
				Constraints: []auroratype.Constraint{
					{Field: "env", Operator: "equal", Value: "production"},
				},
			},
		},
	}

	mockFetcher := new(mocks.MockFetcher)
	mockFetcher.On("IsStatic").Return(true)
	mockFetcher.On("Fetch", mock.Anything).Return(map[string]auroratype.Parameter{"test_param": param}, nil)
	mockFetcher.On("FetchExperiments", mock.Anything).Return(nil, nil)

	s := NewFetcherStorage(mockFetcher)
	err := s.Start(context.Background())
	assert.NoError(t, err)

	client := NewClient(s, ClientOptions{})

	t.Run("parameter matches constraint", func(t *testing.T) {
		attr := NewAttribute()
		attr.Set("env", "production")

		result := client.GetParameter(context.Background(), "test_param", attr)

		assert.True(t, result.matched)
		assert.Equal(t, "matched_value", result.value)
	})

	t.Run("parameter does not match constraint", func(t *testing.T) {
		attr := NewAttribute()
		attr.Set("env", "development")

		result := client.GetParameter(context.Background(), "test_param", attr)

		assert.False(t, result.matched)
		assert.Equal(t, "default", result.value)
	})
}

func TestClientGetParameterStorageError(t *testing.T) {
	mockFetcher := new(mocks.MockFetcher)
	mockFetcher.On("IsStatic").Return(true)
	mockFetcher.On("Fetch", mock.Anything).Return(nil, assert.AnError)

	s := NewFetcherStorage(mockFetcher)
	err := s.Start(context.Background())
	assert.Error(t, err)

	client := NewClient(s, ClientOptions{})

	attr := NewAttribute()
	result := client.GetParameter(context.Background(), "test_param", attr)

	assert.False(t, result.matched)
	assert.Nil(t, result.value)
}

func TestClientRegisterOperator(t *testing.T) {
	mockFetcher := new(mocks.MockFetcher)
	mockFetcher.On("IsStatic").Return(true)
	mockFetcher.On("Fetch", mock.Anything).Return(make(map[string]auroratype.Parameter), nil)
	mockFetcher.On("FetchExperiments", mock.Anything).Return(nil, nil)

	s := NewFetcherStorage(mockFetcher)
	err := s.Start(context.Background())
	assert.NoError(t, err)

	client := NewClient(s, ClientOptions{})

	customOperator := func(a, b any) bool {
		s1, ok1 := a.(string)
		s2, ok2 := b.(string)
		if !ok1 || !ok2 {
			return false
		}
		return len(s1) == len(s2)
	}

	client.RegisterOperator("sameLength", customOperator)

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue: "same_length",
				Constraints: []auroratype.Constraint{
					{Field: "key", Operator: "sameLength", Value: "test"},
				},
			},
		},
	}

	mockFetcher2 := new(mocks.MockFetcher)
	mockFetcher2.On("IsStatic").Return(true)
	mockFetcher2.On("Fetch", mock.Anything).Return(map[string]auroratype.Parameter{"test_param": param}, nil)
	mockFetcher2.On("FetchExperiments", mock.Anything).Return(nil, nil)

	s2 := NewFetcherStorage(mockFetcher2)
	err = s2.Start(context.Background())
	assert.NoError(t, err)
	client2 := NewClient(s2, ClientOptions{})
	client2.RegisterOperator("sameLength", customOperator)

	t.Run("custom operator matches", func(t *testing.T) {
		attr := NewAttribute()
		attr.Set("key", "abcd")

		result := client2.GetParameter(context.Background(), "test_param", attr)

		assert.True(t, result.matched)
		assert.Equal(t, "same_length", result.value)
	})

	t.Run("custom operator does not match", func(t *testing.T) {
		attr := NewAttribute()
		attr.Set("key", "abc")

		result := client2.GetParameter(context.Background(), "test_param", attr)

		assert.False(t, result.matched)
		assert.Equal(t, "default", result.value)
	})
}

func TestClientGetParameterWithDefaultValue(t *testing.T) {
	param := auroratype.Parameter{
		DefaultValue: "fallback",
		Rules:        []auroratype.Rule{},
	}

	mockFetcher := new(mocks.MockFetcher)
	mockFetcher.On("IsStatic").Return(true)
	mockFetcher.On("Fetch", mock.Anything).Return(map[string]auroratype.Parameter{"test_param": param}, nil)
	mockFetcher.On("FetchExperiments", mock.Anything).Return(nil, nil)

	s := NewFetcherStorage(mockFetcher)
	err := s.Start(context.Background())
	assert.NoError(t, err)

	client := NewClient(s, ClientOptions{})

	attr := NewAttribute()
	result := client.GetParameter(context.Background(), "test_param", attr)

	assert.False(t, result.matched)
	assert.Equal(t, "fallback", result.value)
}
