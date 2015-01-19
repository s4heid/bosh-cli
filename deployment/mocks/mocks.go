// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/cloudfoundry/bosh-micro-cli/deployment (interfaces: Deployment,Factory,Deployer,Manager,ManagerFactory)

package mocks

import (
	gomock "code.google.com/p/gomock/gomock"
	cloud "github.com/cloudfoundry/bosh-micro-cli/cloud"
	deployment "github.com/cloudfoundry/bosh-micro-cli/deployment"
	agentclient "github.com/cloudfoundry/bosh-micro-cli/deployment/agentclient"
	disk "github.com/cloudfoundry/bosh-micro-cli/deployment/disk"
	instance "github.com/cloudfoundry/bosh-micro-cli/deployment/instance"
	manifest0 "github.com/cloudfoundry/bosh-micro-cli/deployment/manifest"
	stemcell "github.com/cloudfoundry/bosh-micro-cli/deployment/stemcell"
	vm "github.com/cloudfoundry/bosh-micro-cli/deployment/vm"
	eventlogger "github.com/cloudfoundry/bosh-micro-cli/eventlogger"
	manifest "github.com/cloudfoundry/bosh-micro-cli/installation/manifest"
)

// Mock of Deployment interface
type MockDeployment struct {
	ctrl     *gomock.Controller
	recorder *_MockDeploymentRecorder
}

// Recorder for MockDeployment (not exported)
type _MockDeploymentRecorder struct {
	mock *MockDeployment
}

func NewMockDeployment(ctrl *gomock.Controller) *MockDeployment {
	mock := &MockDeployment{ctrl: ctrl}
	mock.recorder = &_MockDeploymentRecorder{mock}
	return mock
}

func (_m *MockDeployment) EXPECT() *_MockDeploymentRecorder {
	return _m.recorder
}

func (_m *MockDeployment) Delete(_param0 eventlogger.Stage) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDeploymentRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

// Mock of Factory interface
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockFactoryRecorder
}

// Recorder for MockFactory (not exported)
type _MockFactoryRecorder struct {
	mock *MockFactory
}

func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &_MockFactoryRecorder{mock}
	return mock
}

func (_m *MockFactory) EXPECT() *_MockFactoryRecorder {
	return _m.recorder
}

func (_m *MockFactory) NewDeployment(_param0 []instance.Instance, _param1 []disk.Disk, _param2 []stemcell.CloudStemcell) deployment.Deployment {
	ret := _m.ctrl.Call(_m, "NewDeployment", _param0, _param1, _param2)
	ret0, _ := ret[0].(deployment.Deployment)
	return ret0
}

func (_mr *_MockFactoryRecorder) NewDeployment(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewDeployment", arg0, arg1, arg2)
}

// Mock of Deployer interface
type MockDeployer struct {
	ctrl     *gomock.Controller
	recorder *_MockDeployerRecorder
}

// Recorder for MockDeployer (not exported)
type _MockDeployerRecorder struct {
	mock *MockDeployer
}

func NewMockDeployer(ctrl *gomock.Controller) *MockDeployer {
	mock := &MockDeployer{ctrl: ctrl}
	mock.recorder = &_MockDeployerRecorder{mock}
	return mock
}

func (_m *MockDeployer) EXPECT() *_MockDeployerRecorder {
	return _m.recorder
}

func (_m *MockDeployer) Deploy(_param0 cloud.Cloud, _param1 manifest0.Manifest, _param2 stemcell.ExtractedStemcell, _param3 manifest.Registry, _param4 manifest.SSHTunnel, _param5 vm.Manager) (deployment.Deployment, error) {
	ret := _m.ctrl.Call(_m, "Deploy", _param0, _param1, _param2, _param3, _param4, _param5)
	ret0, _ := ret[0].(deployment.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDeployerRecorder) Deploy(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Deploy", arg0, arg1, arg2, arg3, arg4, arg5)
}

// Mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *_MockManagerRecorder
}

// Recorder for MockManager (not exported)
type _MockManagerRecorder struct {
	mock *MockManager
}

func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &_MockManagerRecorder{mock}
	return mock
}

func (_m *MockManager) EXPECT() *_MockManagerRecorder {
	return _m.recorder
}

func (_m *MockManager) Cleanup(_param0 eventlogger.Stage) error {
	ret := _m.ctrl.Call(_m, "Cleanup", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockManagerRecorder) Cleanup(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Cleanup", arg0)
}

func (_m *MockManager) FindCurrent() (deployment.Deployment, bool, error) {
	ret := _m.ctrl.Call(_m, "FindCurrent")
	ret0, _ := ret[0].(deployment.Deployment)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockManagerRecorder) FindCurrent() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FindCurrent")
}

// Mock of ManagerFactory interface
type MockManagerFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockManagerFactoryRecorder
}

// Recorder for MockManagerFactory (not exported)
type _MockManagerFactoryRecorder struct {
	mock *MockManagerFactory
}

func NewMockManagerFactory(ctrl *gomock.Controller) *MockManagerFactory {
	mock := &MockManagerFactory{ctrl: ctrl}
	mock.recorder = &_MockManagerFactoryRecorder{mock}
	return mock
}

func (_m *MockManagerFactory) EXPECT() *_MockManagerFactoryRecorder {
	return _m.recorder
}

func (_m *MockManagerFactory) NewManager(_param0 cloud.Cloud, _param1 agentclient.AgentClient, _param2 string) deployment.Manager {
	ret := _m.ctrl.Call(_m, "NewManager", _param0, _param1, _param2)
	ret0, _ := ret[0].(deployment.Manager)
	return ret0
}

func (_mr *_MockManagerFactoryRecorder) NewManager(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewManager", arg0, arg1, arg2)
}