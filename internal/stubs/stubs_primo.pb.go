// Code generated by protoc-gen-go-primo. DO NOT EDIT.
// versions:
// 	protoc-gen-go-primo v
// 	protoc              v5.28.2
// source: github.com/dogmatiq/enginekit/internal/stubs/stubs.proto

package stubs

type ProtoStubABuilder struct {
	prototype ProtoStubA
}

// NewProtoStubABuilder returns a builder that constructs [ProtoStubA] messages.
func NewProtoStubABuilder() *ProtoStubABuilder {
	return &ProtoStubABuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *ProtoStubABuilder) From(x *ProtoStubA) *ProtoStubABuilder {
	b.prototype.Value = x.Value
	return b
}

// Build returns a new [ProtoStubA] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *ProtoStubABuilder) Build() *ProtoStubA {
	return &ProtoStubA{
		Value: b.prototype.Value,
	}
}

// WithValue configures the builder to set the Value field to v,
// then returns b.
func (b *ProtoStubABuilder) WithValue(v string) *ProtoStubABuilder {
	b.prototype.Value = v
	return b
}

type ProtoStubBBuilder struct {
	prototype ProtoStubB
}

// NewProtoStubBBuilder returns a builder that constructs [ProtoStubB] messages.
func NewProtoStubBBuilder() *ProtoStubBBuilder {
	return &ProtoStubBBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *ProtoStubBBuilder) From(x *ProtoStubB) *ProtoStubBBuilder {
	b.prototype.Value = x.Value
	return b
}

// Build returns a new [ProtoStubB] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *ProtoStubBBuilder) Build() *ProtoStubB {
	return &ProtoStubB{
		Value: b.prototype.Value,
	}
}

// WithValue configures the builder to set the Value field to v,
// then returns b.
func (b *ProtoStubBBuilder) WithValue(v string) *ProtoStubBBuilder {
	b.prototype.Value = v
	return b
}

type ProtoStubCBuilder struct {
	prototype ProtoStubC
}

// NewProtoStubCBuilder returns a builder that constructs [ProtoStubC] messages.
func NewProtoStubCBuilder() *ProtoStubCBuilder {
	return &ProtoStubCBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *ProtoStubCBuilder) From(x *ProtoStubC) *ProtoStubCBuilder {
	b.prototype.Value = x.Value
	return b
}

// Build returns a new [ProtoStubC] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *ProtoStubCBuilder) Build() *ProtoStubC {
	return &ProtoStubC{
		Value: b.prototype.Value,
	}
}

// WithValue configures the builder to set the Value field to v,
// then returns b.
func (b *ProtoStubCBuilder) WithValue(v string) *ProtoStubCBuilder {
	b.prototype.Value = v
	return b
}

// SetValue sets the x.Value field to v, then returns x.
func (x *ProtoStubA) SetValue(v string) {
	x.Value = v
}

// SetValue sets the x.Value field to v, then returns x.
func (x *ProtoStubB) SetValue(v string) {
	x.Value = v
}

// SetValue sets the x.Value field to v, then returns x.
func (x *ProtoStubC) SetValue(v string) {
	x.Value = v
}