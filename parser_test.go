package main

import (
	"testing"
)

func TestSimpleMode_Gen(t *testing.T) {
	type fields struct {
		FuncList []baseFunc
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		m := &SimpleMode{
			FuncList: tt.fields.FuncList,
		}
		got, err := m.Gen()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SimpleMode.Gen() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. SimpleMode.Gen() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestHandlerMode_Gen(t *testing.T) {
	type fields struct {
		MainFunc    baseFunc
		SubFuncList []baseFunc
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "t01",
			fields: fields{
				MainFunc: baseFunc{FuncName: "main", FuncNote: "mainnote", FuncReqJson: "{\"cc\":1}", FuncRespJson: "{\"dd\":1}"},
				SubFuncList: []baseFunc{
					{FuncName: "hi", FuncNote: "hinote", FuncReqJson: "{\"a\":1}", FuncRespJson: "{\"b\":1}"},
					{FuncName: "hi", FuncNote: "hinote", FuncReqJson: "{\"a\":1}", FuncRespJson: "{\"b\":1}"},
					{FuncName: "hi", FuncNote: "hinote", FuncReqJson: "{\"a\":1}", FuncRespJson: "{\"b\":1}"},
				},
			},
		}}
	for _, tt := range tests {
		m := &HandlerMode{
			MainFunc:    tt.fields.MainFunc,
			SubFuncList: tt.fields.SubFuncList,
		}
		got, err := m.Gen()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. HandlerMode.Gen() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. HandlerMode.Gen() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestProxyMode_Gen(t *testing.T) {
	type fields struct {
		FuncList []ProxyFunc
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "t01",
			fields: fields{
				FuncList: []ProxyFunc{
					{baseFunc: baseFunc{FuncName: "hi", FuncNote: "hinote", FuncReqJson: "{\"a\":1}", FuncRespJson: "{\"b\":1}"}, FuncReqUrl: "/a/b/c", FuncReqMethod: "get"},
					{baseFunc: baseFunc{FuncName: "hi", FuncNote: "hinote", FuncReqJson: "{\"a\":1}", FuncRespJson: "{\"b\":1}"}, FuncReqUrl: "/a/b/c", FuncReqMethod: "post"},
				},
			},
		},
	}
	for _, tt := range tests {
		m := &ProxyMode{
			FuncList: tt.fields.FuncList,
		}
		got, err := m.Gen()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ProxyMode.Gen() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. ProxyMode.Gen() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
