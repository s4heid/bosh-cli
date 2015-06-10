// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/cloudfoundry/bosh-init/cloud (interfaces: Cloud,Factory)

package mocks

import (
	gomock "code.google.com/p/gomock/gomock"
	cloud "github.com/cloudfoundry/bosh-init/cloud"
	property "github.com/cloudfoundry/bosh-utils/property"
	installation "github.com/cloudfoundry/bosh-init/installation"
)

// Mock of Cloud interface
type MockCloud struct {
	ctrl     *gomock.Controller
	recorder *_MockCloudRecorder
}

// Recorder for MockCloud (not exported)
type _MockCloudRecorder struct {
	mock *MockCloud
}

func NewMockCloud(ctrl *gomock.Controller) *MockCloud {
	mock := &MockCloud{ctrl: ctrl}
	mock.recorder = &_MockCloudRecorder{mock}
	return mock
}

func (_m *MockCloud) EXPECT() *_MockCloudRecorder {
	return _m.recorder
}

func (_m *MockCloud) AttachDisk(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "AttachDisk", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudRecorder) AttachDisk(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AttachDisk", arg0, arg1)
}

func (_m *MockCloud) CreateDisk(_param0 int, _param1 property.Map, _param2 string) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateDisk", _param0, _param1, _param2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCloudRecorder) CreateDisk(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateDisk", arg0, arg1, arg2)
}

func (_m *MockCloud) CreateStemcell(_param0 string, _param1 property.Map) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateStemcell", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCloudRecorder) CreateStemcell(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateStemcell", arg0, arg1)
}

func (_m *MockCloud) CreateVM(_param0 string, _param1 string, _param2 property.Map, _param3 map[string]property.Map, _param4 property.Map) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateVM", _param0, _param1, _param2, _param3, _param4)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCloudRecorder) CreateVM(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateVM", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockCloud) DeleteDisk(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteDisk", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudRecorder) DeleteDisk(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteDisk", arg0)
}

func (_m *MockCloud) DeleteStemcell(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteStemcell", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudRecorder) DeleteStemcell(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteStemcell", arg0)
}

func (_m *MockCloud) DeleteVM(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteVM", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudRecorder) DeleteVM(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteVM", arg0)
}

func (_m *MockCloud) DetachDisk(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DetachDisk", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudRecorder) DetachDisk(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DetachDisk", arg0, arg1)
}

func (_m *MockCloud) HasVM(_param0 string) (bool, error) {
	ret := _m.ctrl.Call(_m, "HasVM", _param0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCloudRecorder) HasVM(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "HasVM", arg0)
}

func (_m *MockCloud) SetVMMetadata(_param0 string, _param1 cloud.VMMetadata) error {
	ret := _m.ctrl.Call(_m, "SetVMMetadata", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCloudRecorder) SetVMMetadata(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetVMMetadata", arg0, arg1)
}

func (_m *MockCloud) String() string {
	ret := _m.ctrl.Call(_m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockCloudRecorder) String() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "String")
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

func (_m *MockFactory) NewCloud(_param0 installation.Installation, _param1 string) (cloud.Cloud, error) {
	ret := _m.ctrl.Call(_m, "NewCloud", _param0, _param1)
	ret0, _ := ret[0].(cloud.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFactoryRecorder) NewCloud(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewCloud", arg0, arg1)
}
