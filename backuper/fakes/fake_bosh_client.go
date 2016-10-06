// This file was generated by counterfeiter
package fakes

import (
	"sync"

	foo "github.com/pivotal-cf/pcf-backup-and-restore/backuper"
)

type FakeBoshClient struct {
	CheckDeploymentExistsStub        func(name string) (bool, error)
	checkDeploymentExistsMutex       sync.RWMutex
	checkDeploymentExistsArgsForCall []struct {
		name string
	}
	checkDeploymentExistsReturns struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBoshClient) CheckDeploymentExists(name string) (bool, error) {
	fake.checkDeploymentExistsMutex.Lock()
	fake.checkDeploymentExistsArgsForCall = append(fake.checkDeploymentExistsArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("CheckDeploymentExists", []interface{}{name})
	fake.checkDeploymentExistsMutex.Unlock()
	if fake.CheckDeploymentExistsStub != nil {
		return fake.CheckDeploymentExistsStub(name)
	} else {
		return fake.checkDeploymentExistsReturns.result1, fake.checkDeploymentExistsReturns.result2
	}
}

func (fake *FakeBoshClient) CheckDeploymentExistsCallCount() int {
	fake.checkDeploymentExistsMutex.RLock()
	defer fake.checkDeploymentExistsMutex.RUnlock()
	return len(fake.checkDeploymentExistsArgsForCall)
}

func (fake *FakeBoshClient) CheckDeploymentExistsArgsForCall(i int) string {
	fake.checkDeploymentExistsMutex.RLock()
	defer fake.checkDeploymentExistsMutex.RUnlock()
	return fake.checkDeploymentExistsArgsForCall[i].name
}

func (fake *FakeBoshClient) CheckDeploymentExistsReturns(result1 bool, result2 error) {
	fake.CheckDeploymentExistsStub = nil
	fake.checkDeploymentExistsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeBoshClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkDeploymentExistsMutex.RLock()
	defer fake.checkDeploymentExistsMutex.RUnlock()
	return fake.invocations
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

var _ foo.BoshClient = new(FakeBoshClient)
