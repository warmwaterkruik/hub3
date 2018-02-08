// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hub3/fragments/fragments.proto

/*
Package fragments is a generated protocol buffer package.

It is generated from these files:
	hub3/fragments/fragments.proto

It has these top-level messages:
	FragmentSearchResponse
	FragmentUpdateResponse
	FragmentError
	FragmentRequest
	FragmentGraph
	Fragment
*/
package fragments

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ObjectType int32

const (
	ObjectType_RESOURCE ObjectType = 0
	ObjectType_LITERAL  ObjectType = 1
)

var ObjectType_name = map[int32]string{
	0: "RESOURCE",
	1: "LITERAL",
}
var ObjectType_value = map[string]int32{
	"RESOURCE": 0,
	"LITERAL":  1,
}

func (x ObjectType) String() string {
	return proto.EnumName(ObjectType_name, int32(x))
}
func (ObjectType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ObjectXSDType int32

const (
	ObjectXSDType_STRING             ObjectXSDType = 0
	ObjectXSDType_BOOLEAN            ObjectXSDType = 1
	ObjectXSDType_DECIMAL            ObjectXSDType = 2
	ObjectXSDType_FLOAT              ObjectXSDType = 3
	ObjectXSDType_DOUBLE             ObjectXSDType = 4
	ObjectXSDType_DATETIME           ObjectXSDType = 5
	ObjectXSDType_TIME               ObjectXSDType = 6
	ObjectXSDType_DATE               ObjectXSDType = 7
	ObjectXSDType_GYEARMONTH         ObjectXSDType = 8
	ObjectXSDType_GYEAR              ObjectXSDType = 9
	ObjectXSDType_GMONTHDAY          ObjectXSDType = 10
	ObjectXSDType_GDAY               ObjectXSDType = 11
	ObjectXSDType_GMONTH             ObjectXSDType = 12
	ObjectXSDType_HEXBINARY          ObjectXSDType = 13
	ObjectXSDType_BASE64BINARY       ObjectXSDType = 14
	ObjectXSDType_ANYURI             ObjectXSDType = 15
	ObjectXSDType_NORMALIZEDSTRING   ObjectXSDType = 16
	ObjectXSDType_TOKEN              ObjectXSDType = 17
	ObjectXSDType_LANGUAGE           ObjectXSDType = 18
	ObjectXSDType_NMTOKEN            ObjectXSDType = 19
	ObjectXSDType_NAME               ObjectXSDType = 20
	ObjectXSDType_NCNAME             ObjectXSDType = 21
	ObjectXSDType_INTEGER            ObjectXSDType = 22
	ObjectXSDType_NONPOSITIVEINTEGER ObjectXSDType = 23
	ObjectXSDType_NEGATIVEINTEGER    ObjectXSDType = 24
	ObjectXSDType_LONG               ObjectXSDType = 25
	ObjectXSDType_INT                ObjectXSDType = 26
	ObjectXSDType_SHORT              ObjectXSDType = 27
	ObjectXSDType_BYTE               ObjectXSDType = 28
	ObjectXSDType_NONNEGATIVEINTEGER ObjectXSDType = 29
	ObjectXSDType_UNSIGNEDLONG       ObjectXSDType = 30
	ObjectXSDType_UNSIGNEDINT        ObjectXSDType = 31
	ObjectXSDType_UNSIGNEDSHORT      ObjectXSDType = 32
	ObjectXSDType_UNSIGNEDBYTE       ObjectXSDType = 33
	ObjectXSDType_POSITIVEINTEGER    ObjectXSDType = 34
)

var ObjectXSDType_name = map[int32]string{
	0:  "STRING",
	1:  "BOOLEAN",
	2:  "DECIMAL",
	3:  "FLOAT",
	4:  "DOUBLE",
	5:  "DATETIME",
	6:  "TIME",
	7:  "DATE",
	8:  "GYEARMONTH",
	9:  "GYEAR",
	10: "GMONTHDAY",
	11: "GDAY",
	12: "GMONTH",
	13: "HEXBINARY",
	14: "BASE64BINARY",
	15: "ANYURI",
	16: "NORMALIZEDSTRING",
	17: "TOKEN",
	18: "LANGUAGE",
	19: "NMTOKEN",
	20: "NAME",
	21: "NCNAME",
	22: "INTEGER",
	23: "NONPOSITIVEINTEGER",
	24: "NEGATIVEINTEGER",
	25: "LONG",
	26: "INT",
	27: "SHORT",
	28: "BYTE",
	29: "NONNEGATIVEINTEGER",
	30: "UNSIGNEDLONG",
	31: "UNSIGNEDINT",
	32: "UNSIGNEDSHORT",
	33: "UNSIGNEDBYTE",
	34: "POSITIVEINTEGER",
}
var ObjectXSDType_value = map[string]int32{
	"STRING":             0,
	"BOOLEAN":            1,
	"DECIMAL":            2,
	"FLOAT":              3,
	"DOUBLE":             4,
	"DATETIME":           5,
	"TIME":               6,
	"DATE":               7,
	"GYEARMONTH":         8,
	"GYEAR":              9,
	"GMONTHDAY":          10,
	"GDAY":               11,
	"GMONTH":             12,
	"HEXBINARY":          13,
	"BASE64BINARY":       14,
	"ANYURI":             15,
	"NORMALIZEDSTRING":   16,
	"TOKEN":              17,
	"LANGUAGE":           18,
	"NMTOKEN":            19,
	"NAME":               20,
	"NCNAME":             21,
	"INTEGER":            22,
	"NONPOSITIVEINTEGER": 23,
	"NEGATIVEINTEGER":    24,
	"LONG":               25,
	"INT":                26,
	"SHORT":              27,
	"BYTE":               28,
	"NONNEGATIVEINTEGER": 29,
	"UNSIGNEDLONG":       30,
	"UNSIGNEDINT":        31,
	"UNSIGNEDSHORT":      32,
	"UNSIGNEDBYTE":       33,
	"POSITIVEINTEGER":    34,
}

func (x ObjectXSDType) String() string {
	return proto.EnumName(ObjectXSDType_name, int32(x))
}
func (ObjectXSDType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type RecordType int32

const (
	RecordType_NARTHEX    RecordType = 0
	RecordType_SCHEMA     RecordType = 1
	RecordType_VOCABULARY RecordType = 2
	RecordType_SOURCE     RecordType = 3
	RecordType_CACHE      RecordType = 4
)

var RecordType_name = map[int32]string{
	0: "NARTHEX",
	1: "SCHEMA",
	2: "VOCABULARY",
	3: "SOURCE",
	4: "CACHE",
}
var RecordType_value = map[string]int32{
	"NARTHEX":    0,
	"SCHEMA":     1,
	"VOCABULARY": 2,
	"SOURCE":     3,
	"CACHE":      4,
}

func (x RecordType) String() string {
	return proto.EnumName(RecordType_name, int32(x))
}
func (RecordType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type FragmentSearchResponse struct {
	Request   *FragmentRequest `protobuf:"bytes,1,opt,name=request" json:"request,omitempty"`
	NrTriples int32            `protobuf:"varint,2,opt,name=nrTriples" json:"nrTriples,omitempty"`
	Page      int32            `protobuf:"varint,3,opt,name=page" json:"page,omitempty"`
	Fragments []*Fragment      `protobuf:"bytes,4,rep,name=fragments" json:"fragments,omitempty"`
}

func (m *FragmentSearchResponse) Reset()                    { *m = FragmentSearchResponse{} }
func (m *FragmentSearchResponse) String() string            { return proto.CompactTextString(m) }
func (*FragmentSearchResponse) ProtoMessage()               {}
func (*FragmentSearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *FragmentSearchResponse) GetRequest() *FragmentRequest {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *FragmentSearchResponse) GetNrTriples() int32 {
	if m != nil {
		return m.NrTriples
	}
	return 0
}

func (m *FragmentSearchResponse) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *FragmentSearchResponse) GetFragments() []*Fragment {
	if m != nil {
		return m.Fragments
	}
	return nil
}

type FragmentUpdateResponse struct {
	GraphsStored int32            `protobuf:"varint,1,opt,name=graphsStored" json:"graphsStored,omitempty"`
	Spec         string           `protobuf:"bytes,2,opt,name=spec" json:"spec,omitempty"`
	HasErrors    bool             `protobuf:"varint,3,opt,name=hasErrors" json:"hasErrors,omitempty"`
	Errors       []*FragmentError `protobuf:"bytes,4,rep,name=errors" json:"errors,omitempty"`
}

func (m *FragmentUpdateResponse) Reset()                    { *m = FragmentUpdateResponse{} }
func (m *FragmentUpdateResponse) String() string            { return proto.CompactTextString(m) }
func (*FragmentUpdateResponse) ProtoMessage()               {}
func (*FragmentUpdateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *FragmentUpdateResponse) GetGraphsStored() int32 {
	if m != nil {
		return m.GraphsStored
	}
	return 0
}

func (m *FragmentUpdateResponse) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

func (m *FragmentUpdateResponse) GetHasErrors() bool {
	if m != nil {
		return m.HasErrors
	}
	return false
}

func (m *FragmentUpdateResponse) GetErrors() []*FragmentError {
	if m != nil {
		return m.Errors
	}
	return nil
}

type FragmentError struct {
}

func (m *FragmentError) Reset()                    { *m = FragmentError{} }
func (m *FragmentError) String() string            { return proto.CompactTextString(m) }
func (*FragmentError) ProtoMessage()               {}
func (*FragmentError) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type FragmentRequest struct {
	Subject   string `protobuf:"bytes,1,opt,name=subject" json:"subject,omitempty"`
	Predicate string `protobuf:"bytes,2,opt,name=predicate" json:"predicate,omitempty"`
	Object    string `protobuf:"bytes,3,opt,name=object" json:"object,omitempty"`
	Page      int32  `protobuf:"varint,4,opt,name=page" json:"page,omitempty"`
	Language  string `protobuf:"bytes,5,opt,name=language" json:"language,omitempty"`
	OrgID     string `protobuf:"bytes,6,opt,name=orgID" json:"orgID,omitempty"`
	Graph     string `protobuf:"bytes,7,opt,name=graph" json:"graph,omitempty"`
	Spec      string `protobuf:"bytes,8,opt,name=spec" json:"spec,omitempty"`
}

func (m *FragmentRequest) Reset()                    { *m = FragmentRequest{} }
func (m *FragmentRequest) String() string            { return proto.CompactTextString(m) }
func (*FragmentRequest) ProtoMessage()               {}
func (*FragmentRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *FragmentRequest) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *FragmentRequest) GetPredicate() string {
	if m != nil {
		return m.Predicate
	}
	return ""
}

func (m *FragmentRequest) GetObject() string {
	if m != nil {
		return m.Object
	}
	return ""
}

func (m *FragmentRequest) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *FragmentRequest) GetLanguage() string {
	if m != nil {
		return m.Language
	}
	return ""
}

func (m *FragmentRequest) GetOrgID() string {
	if m != nil {
		return m.OrgID
	}
	return ""
}

func (m *FragmentRequest) GetGraph() string {
	if m != nil {
		return m.Graph
	}
	return ""
}

func (m *FragmentRequest) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

type FragmentGraph struct {
	OrgID         string      `protobuf:"bytes,1,opt,name=orgID" json:"orgID,omitempty"`
	Spec          string      `protobuf:"bytes,2,opt,name=spec" json:"spec,omitempty"`
	Revision      int32       `protobuf:"varint,3,opt,name=revision" json:"revision,omitempty"`
	HubID         string      `protobuf:"bytes,4,opt,name=hubID" json:"hubID,omitempty"`
	EntryURI      string      `protobuf:"bytes,5,opt,name=entryURI" json:"entryURI,omitempty"`
	NamedGraphURI string      `protobuf:"bytes,6,opt,name=namedGraphURI" json:"namedGraphURI,omitempty"`
	RDF           []byte      `protobuf:"bytes,7,opt,name=RDF,proto3" json:"RDF,omitempty"`
	RdfMimeType   string      `protobuf:"bytes,8,opt,name=rdfMimeType" json:"rdfMimeType,omitempty"`
	RecordType    RecordType  `protobuf:"varint,9,opt,name=recordType,enum=fragments.RecordType" json:"recordType,omitempty"`
	Fragments     []*Fragment `protobuf:"bytes,10,rep,name=fragments" json:"fragments,omitempty"`
	Tags          []string    `protobuf:"bytes,11,rep,name=tags" json:"tags,omitempty"`
	DocType       string      `protobuf:"bytes,12,opt,name=docType" json:"docType,omitempty"`
}

func (m *FragmentGraph) Reset()                    { *m = FragmentGraph{} }
func (m *FragmentGraph) String() string            { return proto.CompactTextString(m) }
func (*FragmentGraph) ProtoMessage()               {}
func (*FragmentGraph) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *FragmentGraph) GetOrgID() string {
	if m != nil {
		return m.OrgID
	}
	return ""
}

func (m *FragmentGraph) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

func (m *FragmentGraph) GetRevision() int32 {
	if m != nil {
		return m.Revision
	}
	return 0
}

func (m *FragmentGraph) GetHubID() string {
	if m != nil {
		return m.HubID
	}
	return ""
}

func (m *FragmentGraph) GetEntryURI() string {
	if m != nil {
		return m.EntryURI
	}
	return ""
}

func (m *FragmentGraph) GetNamedGraphURI() string {
	if m != nil {
		return m.NamedGraphURI
	}
	return ""
}

func (m *FragmentGraph) GetRDF() []byte {
	if m != nil {
		return m.RDF
	}
	return nil
}

func (m *FragmentGraph) GetRdfMimeType() string {
	if m != nil {
		return m.RdfMimeType
	}
	return ""
}

func (m *FragmentGraph) GetRecordType() RecordType {
	if m != nil {
		return m.RecordType
	}
	return RecordType_NARTHEX
}

func (m *FragmentGraph) GetFragments() []*Fragment {
	if m != nil {
		return m.Fragments
	}
	return nil
}

func (m *FragmentGraph) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *FragmentGraph) GetDocType() string {
	if m != nil {
		return m.DocType
	}
	return ""
}

type Fragment struct {
	// meta block
	OrgID    string `protobuf:"bytes,1,opt,name=orgID" json:"orgID,omitempty"`
	Spec     string `protobuf:"bytes,2,opt,name=spec" json:"spec,omitempty"`
	Revision int32  `protobuf:"varint,3,opt,name=revision" json:"revision,omitempty"`
	HubID    string `protobuf:"bytes,4,opt,name=hubID" json:"hubID,omitempty"`
	// RDF core
	Subject           string        `protobuf:"bytes,5,opt,name=subject" json:"subject,omitempty"`
	SubjectClass      []string      `protobuf:"bytes,6,rep,name=subjectClass" json:"subjectClass,omitempty"`
	Predicate         string        `protobuf:"bytes,7,opt,name=predicate" json:"predicate,omitempty"`
	SearchLabel       string        `protobuf:"bytes,8,opt,name=searchLabel" json:"searchLabel,omitempty"`
	Object            string        `protobuf:"bytes,9,opt,name=object" json:"object,omitempty"`
	ObjectType        ObjectType    `protobuf:"varint,10,opt,name=objectType,enum=fragments.ObjectType" json:"objectType,omitempty"`
	Language          string        `protobuf:"bytes,11,opt,name=language" json:"language,omitempty"`
	ObjectContentType string        `protobuf:"bytes,12,opt,name=objectContentType" json:"objectContentType,omitempty"`
	DataType          ObjectXSDType `protobuf:"varint,13,opt,name=dataType,enum=fragments.ObjectXSDType" json:"dataType,omitempty"`
	XSDRaw            string        `protobuf:"bytes,14,opt,name=XSDRaw" json:"XSDRaw,omitempty"`
	ObjectTypeRaw     string        `protobuf:"bytes,29,opt,name=objectTypeRaw" json:"objectTypeRaw,omitempty"`
	NamedGraphURI     string        `protobuf:"bytes,15,opt,name=namedGraphURI" json:"namedGraphURI,omitempty"`
	Triple            string        `protobuf:"bytes,16,opt,name=triple" json:"triple,omitempty"`
	// RDF graph position
	Level                int32    `protobuf:"varint,17,opt,name=level" json:"level,omitempty"`
	ReferrerSubject      string   `protobuf:"bytes,18,opt,name=referrerSubject" json:"referrerSubject,omitempty"`
	ReferrerPredicate    string   `protobuf:"bytes,19,opt,name=referrerPredicate" json:"referrerPredicate,omitempty"`
	ReferrerSearchLabel  string   `protobuf:"bytes,20,opt,name=referrerSearchLabel" json:"referrerSearchLabel,omitempty"`
	ReferrerSubjectClass []string `protobuf:"bytes,21,rep,name=referrerSubjectClass" json:"referrerSubjectClass,omitempty"`
	// Content Index
	LatLong                   string `protobuf:"bytes,22,opt,name=latLong" json:"latLong,omitempty"`
	Date                      string `protobuf:"bytes,23,opt,name=date" json:"date,omitempty"`
	DateRange                 string `protobuf:"bytes,24,opt,name=dateRange" json:"dateRange,omitempty"`
	Integer                   int32  `protobuf:"varint,25,opt,name=integer" json:"integer,omitempty"`
	IntegerRange              int32  `protobuf:"varint,26,opt,name=integerRange" json:"integerRange,omitempty"`
	ReferrerResourceSortOrder int32  `protobuf:"varint,27,opt,name=referrerResourceSortOrder" json:"referrerResourceSortOrder,omitempty"`
	// content tags
	// example values linkGraphExternal prefLabel linkDomainExternal thumbnail date
	Tags    []string `protobuf:"bytes,28,rep,name=tags" json:"tags,omitempty"`
	DocType string   `protobuf:"bytes,30,opt,name=docType" json:"docType,omitempty"`
}

func (m *Fragment) Reset()                    { *m = Fragment{} }
func (m *Fragment) String() string            { return proto.CompactTextString(m) }
func (*Fragment) ProtoMessage()               {}
func (*Fragment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Fragment) GetOrgID() string {
	if m != nil {
		return m.OrgID
	}
	return ""
}

func (m *Fragment) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

func (m *Fragment) GetRevision() int32 {
	if m != nil {
		return m.Revision
	}
	return 0
}

func (m *Fragment) GetHubID() string {
	if m != nil {
		return m.HubID
	}
	return ""
}

func (m *Fragment) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Fragment) GetSubjectClass() []string {
	if m != nil {
		return m.SubjectClass
	}
	return nil
}

func (m *Fragment) GetPredicate() string {
	if m != nil {
		return m.Predicate
	}
	return ""
}

func (m *Fragment) GetSearchLabel() string {
	if m != nil {
		return m.SearchLabel
	}
	return ""
}

func (m *Fragment) GetObject() string {
	if m != nil {
		return m.Object
	}
	return ""
}

func (m *Fragment) GetObjectType() ObjectType {
	if m != nil {
		return m.ObjectType
	}
	return ObjectType_RESOURCE
}

func (m *Fragment) GetLanguage() string {
	if m != nil {
		return m.Language
	}
	return ""
}

func (m *Fragment) GetObjectContentType() string {
	if m != nil {
		return m.ObjectContentType
	}
	return ""
}

func (m *Fragment) GetDataType() ObjectXSDType {
	if m != nil {
		return m.DataType
	}
	return ObjectXSDType_STRING
}

func (m *Fragment) GetXSDRaw() string {
	if m != nil {
		return m.XSDRaw
	}
	return ""
}

func (m *Fragment) GetObjectTypeRaw() string {
	if m != nil {
		return m.ObjectTypeRaw
	}
	return ""
}

func (m *Fragment) GetNamedGraphURI() string {
	if m != nil {
		return m.NamedGraphURI
	}
	return ""
}

func (m *Fragment) GetTriple() string {
	if m != nil {
		return m.Triple
	}
	return ""
}

func (m *Fragment) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *Fragment) GetReferrerSubject() string {
	if m != nil {
		return m.ReferrerSubject
	}
	return ""
}

func (m *Fragment) GetReferrerPredicate() string {
	if m != nil {
		return m.ReferrerPredicate
	}
	return ""
}

func (m *Fragment) GetReferrerSearchLabel() string {
	if m != nil {
		return m.ReferrerSearchLabel
	}
	return ""
}

func (m *Fragment) GetReferrerSubjectClass() []string {
	if m != nil {
		return m.ReferrerSubjectClass
	}
	return nil
}

func (m *Fragment) GetLatLong() string {
	if m != nil {
		return m.LatLong
	}
	return ""
}

func (m *Fragment) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

func (m *Fragment) GetDateRange() string {
	if m != nil {
		return m.DateRange
	}
	return ""
}

func (m *Fragment) GetInteger() int32 {
	if m != nil {
		return m.Integer
	}
	return 0
}

func (m *Fragment) GetIntegerRange() int32 {
	if m != nil {
		return m.IntegerRange
	}
	return 0
}

func (m *Fragment) GetReferrerResourceSortOrder() int32 {
	if m != nil {
		return m.ReferrerResourceSortOrder
	}
	return 0
}

func (m *Fragment) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *Fragment) GetDocType() string {
	if m != nil {
		return m.DocType
	}
	return ""
}

func init() {
	proto.RegisterType((*FragmentSearchResponse)(nil), "fragments.FragmentSearchResponse")
	proto.RegisterType((*FragmentUpdateResponse)(nil), "fragments.FragmentUpdateResponse")
	proto.RegisterType((*FragmentError)(nil), "fragments.FragmentError")
	proto.RegisterType((*FragmentRequest)(nil), "fragments.FragmentRequest")
	proto.RegisterType((*FragmentGraph)(nil), "fragments.FragmentGraph")
	proto.RegisterType((*Fragment)(nil), "fragments.Fragment")
	proto.RegisterEnum("fragments.ObjectType", ObjectType_name, ObjectType_value)
	proto.RegisterEnum("fragments.ObjectXSDType", ObjectXSDType_name, ObjectXSDType_value)
	proto.RegisterEnum("fragments.RecordType", RecordType_name, RecordType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for FragmentService service

type FragmentServiceClient interface {
	Search(ctx context.Context, in *FragmentRequest, opts ...grpc.CallOption) (*FragmentSearchResponse, error)
	Update(ctx context.Context, in *FragmentGraph, opts ...grpc.CallOption) (*FragmentUpdateResponse, error)
}

type fragmentServiceClient struct {
	cc *grpc.ClientConn
}

func NewFragmentServiceClient(cc *grpc.ClientConn) FragmentServiceClient {
	return &fragmentServiceClient{cc}
}

func (c *fragmentServiceClient) Search(ctx context.Context, in *FragmentRequest, opts ...grpc.CallOption) (*FragmentSearchResponse, error) {
	out := new(FragmentSearchResponse)
	err := grpc.Invoke(ctx, "/fragments.FragmentService/Search", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fragmentServiceClient) Update(ctx context.Context, in *FragmentGraph, opts ...grpc.CallOption) (*FragmentUpdateResponse, error) {
	out := new(FragmentUpdateResponse)
	err := grpc.Invoke(ctx, "/fragments.FragmentService/Update", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for FragmentService service

type FragmentServiceServer interface {
	Search(context.Context, *FragmentRequest) (*FragmentSearchResponse, error)
	Update(context.Context, *FragmentGraph) (*FragmentUpdateResponse, error)
}

func RegisterFragmentServiceServer(s *grpc.Server, srv FragmentServiceServer) {
	s.RegisterService(&_FragmentService_serviceDesc, srv)
}

func _FragmentService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FragmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FragmentServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fragments.FragmentService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FragmentServiceServer).Search(ctx, req.(*FragmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FragmentService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FragmentGraph)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FragmentServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fragments.FragmentService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FragmentServiceServer).Update(ctx, req.(*FragmentGraph))
	}
	return interceptor(ctx, in, info, handler)
}

var _FragmentService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fragments.FragmentService",
	HandlerType: (*FragmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _FragmentService_Search_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _FragmentService_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hub3/fragments/fragments.proto",
}

func init() { proto.RegisterFile("hub3/fragments/fragments.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0xcf, 0x72, 0xdb, 0xb6,
	0x13, 0x8e, 0xac, 0xff, 0x2b, 0xc9, 0x82, 0x61, 0xc7, 0x61, 0x9c, 0x3f, 0x3f, 0x45, 0xf3, 0x9b,
	0xa9, 0x27, 0xd3, 0x49, 0x53, 0x27, 0xe9, 0xa9, 0x17, 0x5a, 0x42, 0x64, 0x4e, 0x29, 0x32, 0x03,
	0x52, 0x19, 0xbb, 0x37, 0x5a, 0x42, 0x64, 0x75, 0x14, 0x51, 0x05, 0xe9, 0x74, 0x72, 0xef, 0x83,
	0xb4, 0xcf, 0xd0, 0x27, 0xe9, 0xe3, 0xe4, 0xd6, 0x59, 0x80, 0xa4, 0x48, 0x59, 0xed, 0xf4, 0xd2,
	0xdb, 0xee, 0xb7, 0x8b, 0xc5, 0xb7, 0x8b, 0x0f, 0x24, 0xe0, 0xe9, 0xcd, 0xed, 0xf5, 0xab, 0x6f,
	0x3e, 0xc8, 0x60, 0xfe, 0x51, 0xac, 0xe2, 0x68, 0x63, 0xbd, 0x58, 0xcb, 0x30, 0x0e, 0x69, 0x33,
	0x03, 0xfa, 0x7f, 0x94, 0xe0, 0xf8, 0x6d, 0xe2, 0x79, 0x22, 0x90, 0xd3, 0x1b, 0x2e, 0xa2, 0x75,
	0xb8, 0x8a, 0x04, 0x7d, 0x0d, 0x75, 0x29, 0x7e, 0xbe, 0x15, 0x51, 0x6c, 0x94, 0x7a, 0xa5, 0xd3,
	0xd6, 0xd9, 0xc9, 0x8b, 0x4d, 0xa1, 0x74, 0x0d, 0xd7, 0x19, 0x3c, 0x4d, 0xa5, 0x8f, 0xa1, 0xb9,
	0x92, 0xbe, 0x5c, 0xac, 0x97, 0x22, 0x32, 0xf6, 0x7a, 0xa5, 0xd3, 0x2a, 0xdf, 0x00, 0x94, 0x42,
	0x65, 0x1d, 0xcc, 0x85, 0x51, 0x56, 0x01, 0x65, 0xd3, 0x6f, 0x61, 0xc3, 0xc7, 0xa8, 0xf4, 0xca,
	0xa7, 0xad, 0xb3, 0xc3, 0x5d, 0x3b, 0xe5, 0x58, 0xff, 0x96, 0x63, 0x3d, 0x59, 0xcf, 0x82, 0x58,
	0x64, 0xac, 0xfb, 0xd0, 0x9e, 0xcb, 0x60, 0x7d, 0x13, 0x79, 0x71, 0x28, 0xc5, 0x4c, 0x51, 0xaf,
	0xf2, 0x02, 0x86, 0x2c, 0xa2, 0xb5, 0x98, 0x2a, 0x7a, 0x4d, 0xae, 0x6c, 0xe4, 0x7d, 0x13, 0x44,
	0x4c, 0xca, 0x50, 0x46, 0x8a, 0x5e, 0x83, 0x6f, 0x00, 0xfa, 0x12, 0x6a, 0x42, 0x87, 0x34, 0x41,
	0x63, 0x07, 0x41, 0x95, 0xca, 0x93, 0xbc, 0x7e, 0x17, 0x3a, 0x85, 0x40, 0xff, 0xcf, 0x12, 0x74,
	0xb7, 0xa6, 0x46, 0x0d, 0xa8, 0x47, 0xb7, 0xd7, 0x3f, 0x89, 0xa9, 0x1e, 0x71, 0x93, 0xa7, 0x2e,
	0xd2, 0x59, 0x4b, 0x31, 0x5b, 0x4c, 0x83, 0x58, 0x24, 0x3c, 0x37, 0x00, 0x3d, 0x86, 0x5a, 0xa8,
	0x97, 0x95, 0x55, 0x28, 0xf1, 0xb2, 0xf1, 0x56, 0x72, 0xe3, 0x3d, 0x81, 0xc6, 0x32, 0x58, 0xcd,
	0x6f, 0x11, 0xaf, 0xaa, 0xec, 0xcc, 0xa7, 0x47, 0x50, 0x0d, 0xe5, 0xdc, 0x1a, 0x1a, 0x35, 0x15,
	0xd0, 0x0e, 0xa2, 0x6a, 0x5c, 0x46, 0x5d, 0xa3, 0xca, 0xc9, 0x86, 0xd6, 0xd8, 0x0c, 0xad, 0xff,
	0x65, 0x6f, 0xd3, 0xe5, 0x48, 0x65, 0x65, 0x15, 0x4b, 0xf9, 0x8a, 0xbb, 0x06, 0x7e, 0x02, 0x0d,
	0x29, 0x3e, 0x2d, 0xa2, 0x45, 0xb8, 0x4a, 0xe4, 0x90, 0xf9, 0x58, 0xe5, 0xe6, 0xf6, 0xda, 0x1a,
	0xaa, 0x46, 0x9a, 0x5c, 0x3b, 0xb8, 0x42, 0xac, 0x62, 0xf9, 0x79, 0xc2, 0xad, 0xb4, 0x93, 0xd4,
	0xa7, 0xff, 0x87, 0xce, 0x2a, 0xf8, 0x28, 0x66, 0x8a, 0x05, 0x26, 0xe8, 0x8e, 0x8a, 0x20, 0x25,
	0x50, 0xe6, 0xc3, 0xb7, 0xaa, 0xaf, 0x36, 0x47, 0x93, 0xf6, 0xa0, 0x25, 0x67, 0x1f, 0xc6, 0x8b,
	0x8f, 0xc2, 0xff, 0xbc, 0x16, 0x49, 0x73, 0x79, 0x88, 0xbe, 0x01, 0x90, 0x62, 0x1a, 0xca, 0x99,
	0x4a, 0x68, 0xf6, 0x4a, 0xa7, 0xfb, 0x67, 0xf7, 0x73, 0xc7, 0xcf, 0xb3, 0x20, 0xcf, 0x25, 0x16,
	0x55, 0x0d, 0xff, 0x46, 0xd5, 0x38, 0xa5, 0x38, 0x98, 0x47, 0x46, 0xab, 0x57, 0xc6, 0x29, 0xa1,
	0x8d, 0x0a, 0x99, 0x85, 0x53, 0xb5, 0x75, 0x5b, 0x2b, 0x24, 0x71, 0xfb, 0x5f, 0xea, 0xd0, 0x48,
	0xab, 0xfc, 0xa7, 0x63, 0xcf, 0x89, 0xb4, 0x5a, 0x14, 0x69, 0x1f, 0xda, 0x89, 0x39, 0x58, 0x06,
	0x51, 0x64, 0xd4, 0x14, 0xf1, 0x02, 0x56, 0x14, 0x72, 0x7d, 0x5b, 0xc8, 0x3d, 0x68, 0x45, 0xea,
	0xab, 0x63, 0x07, 0xd7, 0x62, 0x99, 0x8e, 0x3f, 0x07, 0xe5, 0xa4, 0xde, 0x2c, 0x48, 0xfd, 0x0d,
	0x80, 0xb6, 0xd4, 0x6c, 0xe0, 0xce, 0xb1, 0xb8, 0x59, 0x90, 0xe7, 0x12, 0x0b, 0xb7, 0xa1, 0xb5,
	0x75, 0x1b, 0xbe, 0x86, 0x03, 0x9d, 0x39, 0x08, 0x57, 0xb1, 0x58, 0xc5, 0xb9, 0xa9, 0xdf, 0x0d,
	0xd0, 0xd7, 0xd0, 0x98, 0x05, 0x71, 0xa0, 0x92, 0x3a, 0x6a, 0x7b, 0xe3, 0xce, 0xf6, 0x97, 0xde,
	0x50, 0x31, 0xc8, 0x32, 0xb1, 0x9d, 0x4b, 0x6f, 0xc8, 0x83, 0x5f, 0x8c, 0x7d, 0xdd, 0x8e, 0xf6,
	0x50, 0xbf, 0x1b, 0x96, 0x18, 0x7e, 0xa2, 0xf5, 0x5b, 0x00, 0xef, 0xaa, 0xbc, 0xbb, 0x4b, 0xe5,
	0xc7, 0x50, 0x8b, 0xd5, 0xf7, 0xd6, 0x20, 0x7a, 0x0f, 0xed, 0xe1, 0xf1, 0x2e, 0xc5, 0x27, 0xb1,
	0x34, 0x0e, 0xd4, 0xb9, 0x6b, 0x87, 0x9e, 0x42, 0x57, 0x8a, 0x0f, 0x42, 0x4a, 0x21, 0xbd, 0xe4,
	0x98, 0xa9, 0x5a, 0xb6, 0x0d, 0xe3, 0x7c, 0x52, 0xe8, 0x5d, 0x76, 0xa4, 0x87, 0x7a, 0x3e, 0x77,
	0x02, 0xf4, 0x25, 0x1c, 0x66, 0x05, 0x72, 0x47, 0x7c, 0xa4, 0xf2, 0x77, 0x85, 0xe8, 0x19, 0x1c,
	0x6d, 0x6d, 0xa9, 0x65, 0x75, 0x5f, 0xc9, 0x6a, 0x67, 0x0c, 0xc5, 0xb9, 0x0c, 0x62, 0x3b, 0x5c,
	0xcd, 0x8d, 0x63, 0x2d, 0xce, 0xc4, 0x45, 0xf1, 0xe3, 0x8f, 0xc1, 0x78, 0xa0, 0xc5, 0x8f, 0x36,
	0x8a, 0x51, 0xfd, 0x2c, 0x82, 0xd5, 0x5c, 0x18, 0x86, 0x16, 0x63, 0x06, 0x60, 0xad, 0xc5, 0x2a,
	0x16, 0x73, 0x21, 0x8d, 0x87, 0x6a, 0x42, 0xa9, 0x8b, 0x42, 0x4f, 0x4c, 0xbd, 0xf4, 0x44, 0xff,
	0x54, 0xf2, 0x18, 0xfd, 0x1e, 0x1e, 0xa6, 0x0c, 0xb9, 0x88, 0xc2, 0x5b, 0x39, 0x15, 0x5e, 0x28,
	0x63, 0x57, 0xce, 0x84, 0x34, 0x1e, 0xa9, 0x05, 0x7f, 0x9f, 0x90, 0xdd, 0xfd, 0xc7, 0xbb, 0xef,
	0xfe, 0xd3, 0xc2, 0xdd, 0x7f, 0xfe, 0x15, 0xc0, 0x46, 0xdf, 0xb4, 0x0d, 0x0d, 0xce, 0x3c, 0x77,
	0xc2, 0x07, 0x8c, 0xdc, 0xa3, 0x2d, 0xa8, 0xdb, 0x96, 0xcf, 0xb8, 0x69, 0x93, 0xd2, 0xf3, 0x5f,
	0x2b, 0xd0, 0x29, 0x48, 0x91, 0x02, 0xd4, 0x3c, 0x9f, 0x5b, 0xce, 0x48, 0xa7, 0x9e, 0xbb, 0xae,
	0xcd, 0x4c, 0x87, 0x94, 0xd0, 0x19, 0xb2, 0x81, 0x35, 0x36, 0x6d, 0xb2, 0x47, 0x9b, 0x50, 0x7d,
	0x6b, 0xbb, 0xa6, 0x4f, 0xca, 0xb8, 0x60, 0xe8, 0x4e, 0xce, 0x6d, 0x46, 0x2a, 0xb8, 0xd3, 0xd0,
	0xf4, 0x99, 0x6f, 0x8d, 0x19, 0xa9, 0xd2, 0x06, 0x54, 0x94, 0x55, 0x43, 0x0b, 0x71, 0x52, 0xa7,
	0xfb, 0x00, 0xa3, 0x2b, 0x66, 0xf2, 0xb1, 0xeb, 0xf8, 0x17, 0xa4, 0x81, 0x85, 0x94, 0x4f, 0x9a,
	0xb4, 0x03, 0xcd, 0x91, 0x82, 0x87, 0xe6, 0x15, 0x01, 0x5c, 0x33, 0x42, 0xab, 0x85, 0x3b, 0xe8,
	0x00, 0x69, 0x63, 0xd2, 0x05, 0xbb, 0x3c, 0xb7, 0x1c, 0x93, 0x5f, 0x91, 0x0e, 0x25, 0xd0, 0x3e,
	0x37, 0x3d, 0xf6, 0xdd, 0xeb, 0x04, 0xd9, 0xc7, 0x64, 0xd3, 0xb9, 0x9a, 0x70, 0x8b, 0x74, 0xe9,
	0x11, 0x10, 0xc7, 0xe5, 0x63, 0xd3, 0xb6, 0x7e, 0x64, 0xc3, 0xa4, 0x2b, 0x82, 0x5b, 0xfa, 0xee,
	0x0f, 0xcc, 0x21, 0x07, 0xc8, 0xd7, 0x36, 0x9d, 0xd1, 0xc4, 0x1c, 0x31, 0x42, 0xb1, 0x43, 0x67,
	0xac, 0x43, 0x87, 0xb8, 0xbd, 0x63, 0x8e, 0x19, 0x39, 0xc2, 0x8a, 0xce, 0x40, 0xd9, 0xf7, 0x31,
	0xc5, 0x72, 0x7c, 0x36, 0x62, 0x9c, 0x1c, 0xd3, 0x63, 0xa0, 0x8e, 0xeb, 0xbc, 0x73, 0x3d, 0xcb,
	0xb7, 0xde, 0xb3, 0x14, 0x7f, 0x40, 0x0f, 0xa1, 0xeb, 0xb0, 0x91, 0x99, 0x07, 0x0d, 0xac, 0x67,
	0xbb, 0xce, 0x88, 0x3c, 0xa4, 0x75, 0x28, 0x5b, 0x8e, 0x4f, 0x4e, 0x90, 0x88, 0x77, 0xe1, 0x72,
	0x9f, 0x3c, 0xc2, 0xe8, 0xf9, 0x95, 0xcf, 0xc8, 0xe3, 0xa4, 0xe8, 0xf6, 0xfa, 0x27, 0xd8, 0xe9,
	0xc4, 0xf1, 0xac, 0x91, 0xc3, 0x86, 0xaa, 0xce, 0x53, 0xda, 0x85, 0x56, 0x8a, 0x60, 0xbd, 0xff,
	0xd1, 0x03, 0xe8, 0xa4, 0x80, 0xae, 0xdb, 0xcb, 0xaf, 0x52, 0xf5, 0x9f, 0x21, 0xb9, 0x6d, 0xc6,
	0xfd, 0xe7, 0x36, 0xc0, 0xe6, 0x37, 0xa5, 0xe6, 0x60, 0x72, 0xff, 0x82, 0x5d, 0x92, 0x7b, 0x4a,
	0x0f, 0x83, 0x0b, 0x36, 0x36, 0x49, 0x09, 0x0f, 0xef, 0xbd, 0x3b, 0x30, 0xcf, 0x27, 0x36, 0xce,
	0x7a, 0x4f, 0xc5, 0xb4, 0xac, 0xca, 0xd8, 0xcc, 0xc0, 0x1c, 0x5c, 0x30, 0x52, 0x39, 0xfb, 0x3d,
	0xf7, 0x92, 0xf1, 0x84, 0xfc, 0xb4, 0x98, 0x0a, 0x3a, 0x82, 0x9a, 0xbe, 0xca, 0xf4, 0x1f, 0x5e,
	0x89, 0x27, 0xcf, 0x76, 0xc4, 0xb6, 0x5e, 0x9d, 0x0c, 0x6a, 0xfa, 0x45, 0x47, 0x77, 0xbd, 0xb1,
	0xd4, 0x37, 0x6e, 0x67, 0x99, 0xe2, 0x33, 0xf0, 0xba, 0xa6, 0x5e, 0xba, 0xaf, 0xfe, 0x0a, 0x00,
	0x00, 0xff, 0xff, 0xc5, 0xcf, 0xfc, 0x49, 0x0b, 0x0b, 0x00, 0x00,
}
