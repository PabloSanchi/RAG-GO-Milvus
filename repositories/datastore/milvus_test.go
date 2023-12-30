package datastore_test

import (
    "fmt"
    "context"
    "testing"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/assert"
    "github.com/pablosanchi/datastore/repositories/datastore"
    "github.com/pablosanchi/datastore/helpers"
    "github.com/milvus-io/milvus-sdk-go/v2/client"
    "github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MockMilvusClient struct {
    mock.Mock
    client.Client
}

func (m *MockMilvusClient) ListCollections(ctx context.Context) ([]*entity.Collection, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*entity.Collection), args.Error(1)
}

type MockEncoder struct {
    mock.Mock
}

func (m *MockEncoder) Encode(text string) ([]float32, error) {
    args := m.Called(text)
    return args.Get(0).([]float32), args.Error(1)
}

func TestListSuccess(t *testing.T) {
    realClient := helpers.NewMilvusClient()
    mockClient := &MockMilvusClient{Client: realClient}

    expectedCollections := []*entity.Collection{{Name: "coll1"}, {Name: "coll2"}}
    mockClient.On("ListCollections", mock.Anything).Return(expectedCollections, nil)

    mockEncoder := new(MockEncoder)
    mockEncoder.On("Encode", "test text").Return([]float32{0.1, 0.2, 0.3}, nil)

    repo := datastore.DatastoreMilvusRepository{
        Client: mockClient,
        Encoder: mockEncoder,
    }

    collectionNames, err := repo.List()

    expectedCollectionNames := []string{"coll1", "coll2"}

    assert.Nil(t, err)
    assert.Equal(t, expectedCollectionNames, collectionNames)

    mockClient.AssertExpectations(t)
}

func TestListError(t *testing.T) {
    realClient := helpers.NewMilvusClient()
    mockClient := &MockMilvusClient{Client: realClient}

    mockClient.On("ListCollections", mock.Anything).Return([]*entity.Collection{}, fmt.Errorf("mock error"))

    mockEncoder := new(MockEncoder)
    mockEncoder.On("Encode", "test text").Return([]float32{0.1, 0.2, 0.3}, nil)

    repo := datastore.DatastoreMilvusRepository{
        Client: mockClient,
        Encoder: mockEncoder,
    }

    collectionNames, err := repo.List()

    assert.NotNil(t, err)
    assert.Nil(t, collectionNames)

    mockClient.AssertExpectations(t)
}