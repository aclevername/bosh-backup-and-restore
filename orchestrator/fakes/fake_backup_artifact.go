// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"io"
	"sync"

	"github.com/pivotal-cf/bosh-backup-and-restore/orchestrator"
)

type FakeBackupArtifact struct {
	InstanceNameStub        func() string
	instanceNameMutex       sync.RWMutex
	instanceNameArgsForCall []struct{}
	instanceNameReturns     struct {
		result1 string
	}
	instanceNameReturnsOnCall map[int]struct {
		result1 string
	}
	InstanceIndexStub        func() string
	instanceIndexMutex       sync.RWMutex
	instanceIndexArgsForCall []struct{}
	instanceIndexReturns     struct {
		result1 string
	}
	instanceIndexReturnsOnCall map[int]struct {
		result1 string
	}
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct{}
	nameReturns     struct {
		result1 string
	}
	nameReturnsOnCall map[int]struct {
		result1 string
	}
	HasCustomNameStub        func() bool
	hasCustomNameMutex       sync.RWMutex
	hasCustomNameArgsForCall []struct{}
	hasCustomNameReturns     struct {
		result1 bool
	}
	hasCustomNameReturnsOnCall map[int]struct {
		result1 bool
	}
	SizeStub        func() (string, error)
	sizeMutex       sync.RWMutex
	sizeArgsForCall []struct{}
	sizeReturns     struct {
		result1 string
		result2 error
	}
	sizeReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	ChecksumStub        func() (orchestrator.BackupChecksum, error)
	checksumMutex       sync.RWMutex
	checksumArgsForCall []struct{}
	checksumReturns     struct {
		result1 orchestrator.BackupChecksum
		result2 error
	}
	checksumReturnsOnCall map[int]struct {
		result1 orchestrator.BackupChecksum
		result2 error
	}
	StreamFromRemoteStub        func(io.Writer) error
	streamFromRemoteMutex       sync.RWMutex
	streamFromRemoteArgsForCall []struct {
		arg1 io.Writer
	}
	streamFromRemoteReturns struct {
		result1 error
	}
	streamFromRemoteReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func() error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct{}
	deleteReturns     struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	StreamToRemoteStub        func(io.Reader) error
	streamToRemoteMutex       sync.RWMutex
	streamToRemoteArgsForCall []struct {
		arg1 io.Reader
	}
	streamToRemoteReturns struct {
		result1 error
	}
	streamToRemoteReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBackupArtifact) InstanceName() string {
	fake.instanceNameMutex.Lock()
	ret, specificReturn := fake.instanceNameReturnsOnCall[len(fake.instanceNameArgsForCall)]
	fake.instanceNameArgsForCall = append(fake.instanceNameArgsForCall, struct{}{})
	fake.recordInvocation("InstanceName", []interface{}{})
	fake.instanceNameMutex.Unlock()
	if fake.InstanceNameStub != nil {
		return fake.InstanceNameStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.instanceNameReturns.result1
}

func (fake *FakeBackupArtifact) InstanceNameCallCount() int {
	fake.instanceNameMutex.RLock()
	defer fake.instanceNameMutex.RUnlock()
	return len(fake.instanceNameArgsForCall)
}

func (fake *FakeBackupArtifact) InstanceNameReturns(result1 string) {
	fake.InstanceNameStub = nil
	fake.instanceNameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeBackupArtifact) InstanceNameReturnsOnCall(i int, result1 string) {
	fake.InstanceNameStub = nil
	if fake.instanceNameReturnsOnCall == nil {
		fake.instanceNameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.instanceNameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeBackupArtifact) InstanceIndex() string {
	fake.instanceIndexMutex.Lock()
	ret, specificReturn := fake.instanceIndexReturnsOnCall[len(fake.instanceIndexArgsForCall)]
	fake.instanceIndexArgsForCall = append(fake.instanceIndexArgsForCall, struct{}{})
	fake.recordInvocation("InstanceIndex", []interface{}{})
	fake.instanceIndexMutex.Unlock()
	if fake.InstanceIndexStub != nil {
		return fake.InstanceIndexStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.instanceIndexReturns.result1
}

func (fake *FakeBackupArtifact) InstanceIndexCallCount() int {
	fake.instanceIndexMutex.RLock()
	defer fake.instanceIndexMutex.RUnlock()
	return len(fake.instanceIndexArgsForCall)
}

func (fake *FakeBackupArtifact) InstanceIndexReturns(result1 string) {
	fake.InstanceIndexStub = nil
	fake.instanceIndexReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeBackupArtifact) InstanceIndexReturnsOnCall(i int, result1 string) {
	fake.InstanceIndexStub = nil
	if fake.instanceIndexReturnsOnCall == nil {
		fake.instanceIndexReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.instanceIndexReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeBackupArtifact) Name() string {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct{}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.nameReturns.result1
}

func (fake *FakeBackupArtifact) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *FakeBackupArtifact) NameReturns(result1 string) {
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeBackupArtifact) NameReturnsOnCall(i int, result1 string) {
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeBackupArtifact) HasCustomName() bool {
	fake.hasCustomNameMutex.Lock()
	ret, specificReturn := fake.hasCustomNameReturnsOnCall[len(fake.hasCustomNameArgsForCall)]
	fake.hasCustomNameArgsForCall = append(fake.hasCustomNameArgsForCall, struct{}{})
	fake.recordInvocation("HasCustomName", []interface{}{})
	fake.hasCustomNameMutex.Unlock()
	if fake.HasCustomNameStub != nil {
		return fake.HasCustomNameStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.hasCustomNameReturns.result1
}

func (fake *FakeBackupArtifact) HasCustomNameCallCount() int {
	fake.hasCustomNameMutex.RLock()
	defer fake.hasCustomNameMutex.RUnlock()
	return len(fake.hasCustomNameArgsForCall)
}

func (fake *FakeBackupArtifact) HasCustomNameReturns(result1 bool) {
	fake.HasCustomNameStub = nil
	fake.hasCustomNameReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeBackupArtifact) HasCustomNameReturnsOnCall(i int, result1 bool) {
	fake.HasCustomNameStub = nil
	if fake.hasCustomNameReturnsOnCall == nil {
		fake.hasCustomNameReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.hasCustomNameReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeBackupArtifact) Size() (string, error) {
	fake.sizeMutex.Lock()
	ret, specificReturn := fake.sizeReturnsOnCall[len(fake.sizeArgsForCall)]
	fake.sizeArgsForCall = append(fake.sizeArgsForCall, struct{}{})
	fake.recordInvocation("Size", []interface{}{})
	fake.sizeMutex.Unlock()
	if fake.SizeStub != nil {
		return fake.SizeStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.sizeReturns.result1, fake.sizeReturns.result2
}

func (fake *FakeBackupArtifact) SizeCallCount() int {
	fake.sizeMutex.RLock()
	defer fake.sizeMutex.RUnlock()
	return len(fake.sizeArgsForCall)
}

func (fake *FakeBackupArtifact) SizeReturns(result1 string, result2 error) {
	fake.SizeStub = nil
	fake.sizeReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeBackupArtifact) SizeReturnsOnCall(i int, result1 string, result2 error) {
	fake.SizeStub = nil
	if fake.sizeReturnsOnCall == nil {
		fake.sizeReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.sizeReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeBackupArtifact) Checksum() (orchestrator.BackupChecksum, error) {
	fake.checksumMutex.Lock()
	ret, specificReturn := fake.checksumReturnsOnCall[len(fake.checksumArgsForCall)]
	fake.checksumArgsForCall = append(fake.checksumArgsForCall, struct{}{})
	fake.recordInvocation("Checksum", []interface{}{})
	fake.checksumMutex.Unlock()
	if fake.ChecksumStub != nil {
		return fake.ChecksumStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.checksumReturns.result1, fake.checksumReturns.result2
}

func (fake *FakeBackupArtifact) ChecksumCallCount() int {
	fake.checksumMutex.RLock()
	defer fake.checksumMutex.RUnlock()
	return len(fake.checksumArgsForCall)
}

func (fake *FakeBackupArtifact) ChecksumReturns(result1 orchestrator.BackupChecksum, result2 error) {
	fake.ChecksumStub = nil
	fake.checksumReturns = struct {
		result1 orchestrator.BackupChecksum
		result2 error
	}{result1, result2}
}

func (fake *FakeBackupArtifact) ChecksumReturnsOnCall(i int, result1 orchestrator.BackupChecksum, result2 error) {
	fake.ChecksumStub = nil
	if fake.checksumReturnsOnCall == nil {
		fake.checksumReturnsOnCall = make(map[int]struct {
			result1 orchestrator.BackupChecksum
			result2 error
		})
	}
	fake.checksumReturnsOnCall[i] = struct {
		result1 orchestrator.BackupChecksum
		result2 error
	}{result1, result2}
}

func (fake *FakeBackupArtifact) StreamFromRemote(arg1 io.Writer) error {
	fake.streamFromRemoteMutex.Lock()
	ret, specificReturn := fake.streamFromRemoteReturnsOnCall[len(fake.streamFromRemoteArgsForCall)]
	fake.streamFromRemoteArgsForCall = append(fake.streamFromRemoteArgsForCall, struct {
		arg1 io.Writer
	}{arg1})
	fake.recordInvocation("StreamFromRemote", []interface{}{arg1})
	fake.streamFromRemoteMutex.Unlock()
	if fake.StreamFromRemoteStub != nil {
		return fake.StreamFromRemoteStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.streamFromRemoteReturns.result1
}

func (fake *FakeBackupArtifact) StreamFromRemoteCallCount() int {
	fake.streamFromRemoteMutex.RLock()
	defer fake.streamFromRemoteMutex.RUnlock()
	return len(fake.streamFromRemoteArgsForCall)
}

func (fake *FakeBackupArtifact) StreamFromRemoteArgsForCall(i int) io.Writer {
	fake.streamFromRemoteMutex.RLock()
	defer fake.streamFromRemoteMutex.RUnlock()
	return fake.streamFromRemoteArgsForCall[i].arg1
}

func (fake *FakeBackupArtifact) StreamFromRemoteReturns(result1 error) {
	fake.StreamFromRemoteStub = nil
	fake.streamFromRemoteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBackupArtifact) StreamFromRemoteReturnsOnCall(i int, result1 error) {
	fake.StreamFromRemoteStub = nil
	if fake.streamFromRemoteReturnsOnCall == nil {
		fake.streamFromRemoteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.streamFromRemoteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBackupArtifact) Delete() error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct{}{})
	fake.recordInvocation("Delete", []interface{}{})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteReturns.result1
}

func (fake *FakeBackupArtifact) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeBackupArtifact) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBackupArtifact) DeleteReturnsOnCall(i int, result1 error) {
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBackupArtifact) StreamToRemote(arg1 io.Reader) error {
	fake.streamToRemoteMutex.Lock()
	ret, specificReturn := fake.streamToRemoteReturnsOnCall[len(fake.streamToRemoteArgsForCall)]
	fake.streamToRemoteArgsForCall = append(fake.streamToRemoteArgsForCall, struct {
		arg1 io.Reader
	}{arg1})
	fake.recordInvocation("StreamToRemote", []interface{}{arg1})
	fake.streamToRemoteMutex.Unlock()
	if fake.StreamToRemoteStub != nil {
		return fake.StreamToRemoteStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.streamToRemoteReturns.result1
}

func (fake *FakeBackupArtifact) StreamToRemoteCallCount() int {
	fake.streamToRemoteMutex.RLock()
	defer fake.streamToRemoteMutex.RUnlock()
	return len(fake.streamToRemoteArgsForCall)
}

func (fake *FakeBackupArtifact) StreamToRemoteArgsForCall(i int) io.Reader {
	fake.streamToRemoteMutex.RLock()
	defer fake.streamToRemoteMutex.RUnlock()
	return fake.streamToRemoteArgsForCall[i].arg1
}

func (fake *FakeBackupArtifact) StreamToRemoteReturns(result1 error) {
	fake.StreamToRemoteStub = nil
	fake.streamToRemoteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBackupArtifact) StreamToRemoteReturnsOnCall(i int, result1 error) {
	fake.StreamToRemoteStub = nil
	if fake.streamToRemoteReturnsOnCall == nil {
		fake.streamToRemoteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.streamToRemoteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBackupArtifact) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.instanceNameMutex.RLock()
	defer fake.instanceNameMutex.RUnlock()
	fake.instanceIndexMutex.RLock()
	defer fake.instanceIndexMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.hasCustomNameMutex.RLock()
	defer fake.hasCustomNameMutex.RUnlock()
	fake.sizeMutex.RLock()
	defer fake.sizeMutex.RUnlock()
	fake.checksumMutex.RLock()
	defer fake.checksumMutex.RUnlock()
	fake.streamFromRemoteMutex.RLock()
	defer fake.streamFromRemoteMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.streamToRemoteMutex.RLock()
	defer fake.streamToRemoteMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBackupArtifact) recordInvocation(key string, args []interface{}) {
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

var _ orchestrator.BackupArtifact = new(FakeBackupArtifact)