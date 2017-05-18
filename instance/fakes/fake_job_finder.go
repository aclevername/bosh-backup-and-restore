// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/bosh-backup-and-restore/instance"
)

type FakeJobFinder struct {
	FindJobsStub        func(hostIdentifier string, connection instance.SSHConnection) (instance.Jobs, error)
	findJobsMutex       sync.RWMutex
	findJobsArgsForCall []struct {
		hostIdentifier string
		connection     instance.SSHConnection
	}
	findJobsReturns struct {
		result1 instance.Jobs
		result2 error
	}
	findJobsReturnsOnCall map[int]struct {
		result1 instance.Jobs
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeJobFinder) FindJobs(hostIdentifier string, connection instance.SSHConnection) (instance.Jobs, error) {
	fake.findJobsMutex.Lock()
	ret, specificReturn := fake.findJobsReturnsOnCall[len(fake.findJobsArgsForCall)]
	fake.findJobsArgsForCall = append(fake.findJobsArgsForCall, struct {
		hostIdentifier string
		connection     instance.SSHConnection
	}{hostIdentifier, connection})
	fake.recordInvocation("FindJobs", []interface{}{hostIdentifier, connection})
	fake.findJobsMutex.Unlock()
	if fake.FindJobsStub != nil {
		return fake.FindJobsStub(hostIdentifier, connection)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.findJobsReturns.result1, fake.findJobsReturns.result2
}

func (fake *FakeJobFinder) FindJobsCallCount() int {
	fake.findJobsMutex.RLock()
	defer fake.findJobsMutex.RUnlock()
	return len(fake.findJobsArgsForCall)
}

func (fake *FakeJobFinder) FindJobsArgsForCall(i int) (string, instance.SSHConnection) {
	fake.findJobsMutex.RLock()
	defer fake.findJobsMutex.RUnlock()
	return fake.findJobsArgsForCall[i].hostIdentifier, fake.findJobsArgsForCall[i].connection
}

func (fake *FakeJobFinder) FindJobsReturns(result1 instance.Jobs, result2 error) {
	fake.FindJobsStub = nil
	fake.findJobsReturns = struct {
		result1 instance.Jobs
		result2 error
	}{result1, result2}
}

func (fake *FakeJobFinder) FindJobsReturnsOnCall(i int, result1 instance.Jobs, result2 error) {
	fake.FindJobsStub = nil
	if fake.findJobsReturnsOnCall == nil {
		fake.findJobsReturnsOnCall = make(map[int]struct {
			result1 instance.Jobs
			result2 error
		})
	}
	fake.findJobsReturnsOnCall[i] = struct {
		result1 instance.Jobs
		result2 error
	}{result1, result2}
}

func (fake *FakeJobFinder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findJobsMutex.RLock()
	defer fake.findJobsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeJobFinder) recordInvocation(key string, args []interface{}) {
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

var _ instance.JobFinder = new(FakeJobFinder)
