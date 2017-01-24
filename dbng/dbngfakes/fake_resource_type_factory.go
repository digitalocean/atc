// This file was generated by counterfeiter
package dbngfakes

import (
	"sync"

	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
)

type FakeResourceTypeFactory struct {
	FindResourceTypeStub        func(pipelineID int, resourceType atc.ResourceType) (*dbng.UsedResourceType, bool, error)
	findResourceTypeMutex       sync.RWMutex
	findResourceTypeArgsForCall []struct {
		pipelineID   int
		resourceType atc.ResourceType
	}
	findResourceTypeReturns struct {
		result1 *dbng.UsedResourceType
		result2 bool
		result3 error
	}
	CreateResourceTypeStub        func(pipelineID int, resourceType atc.ResourceType, version atc.Version) (*dbng.UsedResourceType, error)
	createResourceTypeMutex       sync.RWMutex
	createResourceTypeArgsForCall []struct {
		pipelineID   int
		resourceType atc.ResourceType
		version      atc.Version
	}
	createResourceTypeReturns struct {
		result1 *dbng.UsedResourceType
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeResourceTypeFactory) FindResourceType(pipelineID int, resourceType atc.ResourceType) (*dbng.UsedResourceType, bool, error) {
	fake.findResourceTypeMutex.Lock()
	fake.findResourceTypeArgsForCall = append(fake.findResourceTypeArgsForCall, struct {
		pipelineID   int
		resourceType atc.ResourceType
	}{pipelineID, resourceType})
	fake.recordInvocation("FindResourceType", []interface{}{pipelineID, resourceType})
	fake.findResourceTypeMutex.Unlock()
	if fake.FindResourceTypeStub != nil {
		return fake.FindResourceTypeStub(pipelineID, resourceType)
	} else {
		return fake.findResourceTypeReturns.result1, fake.findResourceTypeReturns.result2, fake.findResourceTypeReturns.result3
	}
}

func (fake *FakeResourceTypeFactory) FindResourceTypeCallCount() int {
	fake.findResourceTypeMutex.RLock()
	defer fake.findResourceTypeMutex.RUnlock()
	return len(fake.findResourceTypeArgsForCall)
}

func (fake *FakeResourceTypeFactory) FindResourceTypeArgsForCall(i int) (int, atc.ResourceType) {
	fake.findResourceTypeMutex.RLock()
	defer fake.findResourceTypeMutex.RUnlock()
	return fake.findResourceTypeArgsForCall[i].pipelineID, fake.findResourceTypeArgsForCall[i].resourceType
}

func (fake *FakeResourceTypeFactory) FindResourceTypeReturns(result1 *dbng.UsedResourceType, result2 bool, result3 error) {
	fake.FindResourceTypeStub = nil
	fake.findResourceTypeReturns = struct {
		result1 *dbng.UsedResourceType
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeResourceTypeFactory) CreateResourceType(pipelineID int, resourceType atc.ResourceType, version atc.Version) (*dbng.UsedResourceType, error) {
	fake.createResourceTypeMutex.Lock()
	fake.createResourceTypeArgsForCall = append(fake.createResourceTypeArgsForCall, struct {
		pipelineID   int
		resourceType atc.ResourceType
		version      atc.Version
	}{pipelineID, resourceType, version})
	fake.recordInvocation("CreateResourceType", []interface{}{pipelineID, resourceType, version})
	fake.createResourceTypeMutex.Unlock()
	if fake.CreateResourceTypeStub != nil {
		return fake.CreateResourceTypeStub(pipelineID, resourceType, version)
	} else {
		return fake.createResourceTypeReturns.result1, fake.createResourceTypeReturns.result2
	}
}

func (fake *FakeResourceTypeFactory) CreateResourceTypeCallCount() int {
	fake.createResourceTypeMutex.RLock()
	defer fake.createResourceTypeMutex.RUnlock()
	return len(fake.createResourceTypeArgsForCall)
}

func (fake *FakeResourceTypeFactory) CreateResourceTypeArgsForCall(i int) (int, atc.ResourceType, atc.Version) {
	fake.createResourceTypeMutex.RLock()
	defer fake.createResourceTypeMutex.RUnlock()
	return fake.createResourceTypeArgsForCall[i].pipelineID, fake.createResourceTypeArgsForCall[i].resourceType, fake.createResourceTypeArgsForCall[i].version
}

func (fake *FakeResourceTypeFactory) CreateResourceTypeReturns(result1 *dbng.UsedResourceType, result2 error) {
	fake.CreateResourceTypeStub = nil
	fake.createResourceTypeReturns = struct {
		result1 *dbng.UsedResourceType
		result2 error
	}{result1, result2}
}

func (fake *FakeResourceTypeFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findResourceTypeMutex.RLock()
	defer fake.findResourceTypeMutex.RUnlock()
	fake.createResourceTypeMutex.RLock()
	defer fake.createResourceTypeMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeResourceTypeFactory) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ dbng.ResourceTypeFactory = new(FakeResourceTypeFactory)