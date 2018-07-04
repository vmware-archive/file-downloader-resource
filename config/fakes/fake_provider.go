// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/calebwashburn/file-downloader/config"
	"github.com/calebwashburn/file-downloader/types"
)

type FakeProvider struct {
	LatestVersionStub        func() (*types.Version, error)
	latestVersionMutex       sync.RWMutex
	latestVersionArgsForCall []struct{}
	latestVersionReturns     struct {
		result1 *types.Version
		result2 error
	}
	GetVersionInfoStub        func(revision, productName string) (*types.VersionInfo, error)
	getVersionInfoMutex       sync.RWMutex
	getVersionInfoArgsForCall []struct {
		revision    string
		productName string
	}
	getVersionInfoReturns struct {
		result1 *types.VersionInfo
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProvider) LatestVersion() (*types.Version, error) {
	fake.latestVersionMutex.Lock()
	fake.latestVersionArgsForCall = append(fake.latestVersionArgsForCall, struct{}{})
	fake.recordInvocation("LatestVersion", []interface{}{})
	fake.latestVersionMutex.Unlock()
	if fake.LatestVersionStub != nil {
		return fake.LatestVersionStub()
	} else {
		return fake.latestVersionReturns.result1, fake.latestVersionReturns.result2
	}
}

func (fake *FakeProvider) LatestVersionCallCount() int {
	fake.latestVersionMutex.RLock()
	defer fake.latestVersionMutex.RUnlock()
	return len(fake.latestVersionArgsForCall)
}

func (fake *FakeProvider) LatestVersionReturns(result1 *types.Version, result2 error) {
	fake.LatestVersionStub = nil
	fake.latestVersionReturns = struct {
		result1 *types.Version
		result2 error
	}{result1, result2}
}

func (fake *FakeProvider) GetVersionInfo(revision string, productName string) (*types.VersionInfo, error) {
	fake.getVersionInfoMutex.Lock()
	fake.getVersionInfoArgsForCall = append(fake.getVersionInfoArgsForCall, struct {
		revision    string
		productName string
	}{revision, productName})
	fake.recordInvocation("GetVersionInfo", []interface{}{revision, productName})
	fake.getVersionInfoMutex.Unlock()
	if fake.GetVersionInfoStub != nil {
		return fake.GetVersionInfoStub(revision, productName)
	} else {
		return fake.getVersionInfoReturns.result1, fake.getVersionInfoReturns.result2
	}
}

func (fake *FakeProvider) GetVersionInfoCallCount() int {
	fake.getVersionInfoMutex.RLock()
	defer fake.getVersionInfoMutex.RUnlock()
	return len(fake.getVersionInfoArgsForCall)
}

func (fake *FakeProvider) GetVersionInfoArgsForCall(i int) (string, string) {
	fake.getVersionInfoMutex.RLock()
	defer fake.getVersionInfoMutex.RUnlock()
	return fake.getVersionInfoArgsForCall[i].revision, fake.getVersionInfoArgsForCall[i].productName
}

func (fake *FakeProvider) GetVersionInfoReturns(result1 *types.VersionInfo, result2 error) {
	fake.GetVersionInfoStub = nil
	fake.getVersionInfoReturns = struct {
		result1 *types.VersionInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.latestVersionMutex.RLock()
	defer fake.latestVersionMutex.RUnlock()
	fake.getVersionInfoMutex.RLock()
	defer fake.getVersionInfoMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeProvider) recordInvocation(key string, args []interface{}) {
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

var _ config.Provider = new(FakeProvider)
