// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lightstep/otel-launcher-go/lightstep/sdk/metric (interfaces: PushExporter,Producer)

// Package metric is a generated GoMock package.
package metric

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	data "github.com/lightstep/otel-launcher-go/lightstep/sdk/metric/data"
)

// MockPushExporter is a mock of PushExporter interface.
type MockPushExporter struct {
	ctrl     *gomock.Controller
	recorder *MockPushExporterMockRecorder
}

// MockPushExporterMockRecorder is the mock recorder for MockPushExporter.
type MockPushExporterMockRecorder struct {
	mock *MockPushExporter
}

// NewMockPushExporter creates a new mock instance.
func NewMockPushExporter(ctrl *gomock.Controller) *MockPushExporter {
	mock := &MockPushExporter{ctrl: ctrl}
	mock.recorder = &MockPushExporterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPushExporter) EXPECT() *MockPushExporterMockRecorder {
	return m.recorder
}

// ExportMetrics mocks base method.
func (m *MockPushExporter) ExportMetrics(arg0 context.Context, arg1 data.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportMetrics", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExportMetrics indicates an expected call of ExportMetrics.
func (mr *MockPushExporterMockRecorder) ExportMetrics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportMetrics", reflect.TypeOf((*MockPushExporter)(nil).ExportMetrics), arg0, arg1)
}

// ForceFlushMetrics mocks base method.
func (m *MockPushExporter) ForceFlushMetrics(arg0 context.Context, arg1 data.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForceFlushMetrics", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForceFlushMetrics indicates an expected call of ForceFlushMetrics.
func (mr *MockPushExporterMockRecorder) ForceFlushMetrics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForceFlushMetrics", reflect.TypeOf((*MockPushExporter)(nil).ForceFlushMetrics), arg0, arg1)
}

// ShutdownMetrics mocks base method.
func (m *MockPushExporter) ShutdownMetrics(arg0 context.Context, arg1 data.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShutdownMetrics", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ShutdownMetrics indicates an expected call of ShutdownMetrics.
func (mr *MockPushExporterMockRecorder) ShutdownMetrics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShutdownMetrics", reflect.TypeOf((*MockPushExporter)(nil).ShutdownMetrics), arg0, arg1)
}

// String mocks base method.
func (m *MockPushExporter) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockPushExporterMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockPushExporter)(nil).String))
}

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// Produce mocks base method.
func (m *MockProducer) Produce(arg0 *data.Metrics) data.Metrics {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Produce", arg0)
	ret0, _ := ret[0].(data.Metrics)
	return ret0
}

// Produce indicates an expected call of Produce.
func (mr *MockProducerMockRecorder) Produce(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockProducer)(nil).Produce), arg0)
}
