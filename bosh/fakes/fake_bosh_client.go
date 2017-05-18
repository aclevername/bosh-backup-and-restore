// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/bosh-backup-and-restore/bosh"
	"github.com/pivotal-cf/bosh-backup-and-restore/orchestrator"
)

type FakeBoshClient struct {
	FindInstancesStub        func(deploymentName string) ([]orchestrator.Instance, error)
	findInstancesMutex       sync.RWMutex
	findInstancesArgsForCall []struct {
		deploymentName string
	}
	findInstancesReturns struct {
		result1 []orchestrator.Instance
		result2 error
	}
	findInstancesReturnsOnCall map[int]struct {
		result1 []orchestrator.Instance
		result2 error
	}
	GetManifestStub        func(deploymentName string) (string, error)
	getManifestMutex       sync.RWMutex
	getManifestArgsForCall []struct {
		deploymentName string
	}
	getManifestReturns struct {
		result1 string
		result2 error
	}
	getManifestReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBoshClient) FindInstances(deploymentName string) ([]orchestrator.Instance, error) {
	fake.findInstancesMutex.Lock()
	ret, specificReturn := fake.findInstancesReturnsOnCall[len(fake.findInstancesArgsForCall)]
	fake.findInstancesArgsForCall = append(fake.findInstancesArgsForCall, struct {
		deploymentName string
	}{deploymentName})
	fake.recordInvocation("FindInstances", []interface{}{deploymentName})
	fake.findInstancesMutex.Unlock()
	if fake.FindInstancesStub != nil {
		return fake.FindInstancesStub(deploymentName)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.findInstancesReturns.result1, fake.findInstancesReturns.result2
}

func (fake *FakeBoshClient) FindInstancesCallCount() int {
	fake.findInstancesMutex.RLock()
	defer fake.findInstancesMutex.RUnlock()
	return len(fake.findInstancesArgsForCall)
}

func (fake *FakeBoshClient) FindInstancesArgsForCall(i int) string {
	fake.findInstancesMutex.RLock()
	defer fake.findInstancesMutex.RUnlock()
	return fake.findInstancesArgsForCall[i].deploymentName
}

func (fake *FakeBoshClient) FindInstancesReturns(result1 []orchestrator.Instance, result2 error) {
	fake.FindInstancesStub = nil
	fake.findInstancesReturns = struct {
		result1 []orchestrator.Instance
		result2 error
	}{result1, result2}
}

func (fake *FakeBoshClient) FindInstancesReturnsOnCall(i int, result1 []orchestrator.Instance, result2 error) {
	fake.FindInstancesStub = nil
	if fake.findInstancesReturnsOnCall == nil {
		fake.findInstancesReturnsOnCall = make(map[int]struct {
			result1 []orchestrator.Instance
			result2 error
		})
	}
	fake.findInstancesReturnsOnCall[i] = struct {
		result1 []orchestrator.Instance
		result2 error
	}{result1, result2}
}

func (fake *FakeBoshClient) GetManifest(deploymentName string) (string, error) {
	fake.getManifestMutex.Lock()
	ret, specificReturn := fake.getManifestReturnsOnCall[len(fake.getManifestArgsForCall)]
	fake.getManifestArgsForCall = append(fake.getManifestArgsForCall, struct {
		deploymentName string
	}{deploymentName})
	fake.recordInvocation("GetManifest", []interface{}{deploymentName})
	fake.getManifestMutex.Unlock()
	if fake.GetManifestStub != nil {
		return fake.GetManifestStub(deploymentName)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getManifestReturns.result1, fake.getManifestReturns.result2
}

func (fake *FakeBoshClient) GetManifestCallCount() int {
	fake.getManifestMutex.RLock()
	defer fake.getManifestMutex.RUnlock()
	return len(fake.getManifestArgsForCall)
}

func (fake *FakeBoshClient) GetManifestArgsForCall(i int) string {
	fake.getManifestMutex.RLock()
	defer fake.getManifestMutex.RUnlock()
	return fake.getManifestArgsForCall[i].deploymentName
}

func (fake *FakeBoshClient) GetManifestReturns(result1 string, result2 error) {
	fake.GetManifestStub = nil
	fake.getManifestReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeBoshClient) GetManifestReturnsOnCall(i int, result1 string, result2 error) {
	fake.GetManifestStub = nil
	if fake.getManifestReturnsOnCall == nil {
		fake.getManifestReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getManifestReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeBoshClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findInstancesMutex.RLock()
	defer fake.findInstancesMutex.RUnlock()
	fake.getManifestMutex.RLock()
	defer fake.getManifestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBoshClient) recordInvocation(key string, args []interface{}) {
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

var _ bosh.BoshClient = new(FakeBoshClient)
