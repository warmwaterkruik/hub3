// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hub3/fragments/viewconfig.proto

package fragments

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type DataSetType int32

const (
	DataSetType_SINGLE   DataSetType = 0
	DataSetType_MULTIPLE DataSetType = 1
	DataSetType_BYQUERY  DataSetType = 2
)

var DataSetType_name = map[int32]string{
	0: "SINGLE",
	1: "MULTIPLE",
	2: "BYQUERY",
}

var DataSetType_value = map[string]int32{
	"SINGLE":   0,
	"MULTIPLE": 1,
	"BYQUERY":  2,
}

func (x DataSetType) String() string {
	return proto.EnumName(DataSetType_name, int32(x))
}

func (DataSetType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{0}
}

type ResultType int32

const (
	ResultType_GRID    ResultType = 0
	ResultType_TABLE   ResultType = 1
	ResultType_MAP     ResultType = 2
	ResultType_ARCHIVE ResultType = 3
)

var ResultType_name = map[int32]string{
	0: "GRID",
	1: "TABLE",
	2: "MAP",
	3: "ARCHIVE",
}

var ResultType_value = map[string]int32{
	"GRID":    0,
	"TABLE":   1,
	"MAP":     2,
	"ARCHIVE": 3,
}

func (x ResultType) String() string {
	return proto.EnumName(ResultType_name, int32(x))
}

func (ResultType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{1}
}

type FieldType int32

const (
	FieldType_LITERAL        FieldType = 0
	FieldType_RESOURCE       FieldType = 1
	FieldType_DATE           FieldType = 2
	FieldType_POINT          FieldType = 3
	FieldType_DIGITAL_OBJECT FieldType = 4
	FieldType_MANIFEST       FieldType = 5
)

var FieldType_name = map[int32]string{
	0: "LITERAL",
	1: "RESOURCE",
	2: "DATE",
	3: "POINT",
	4: "DIGITAL_OBJECT",
	5: "MANIFEST",
}

var FieldType_value = map[string]int32{
	"LITERAL":        0,
	"RESOURCE":       1,
	"DATE":           2,
	"POINT":          3,
	"DIGITAL_OBJECT": 4,
	"MANIFEST":       5,
}

func (x FieldType) String() string {
	return proto.EnumName(FieldType_name, int32(x))
}

func (FieldType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{2}
}

type InlineType int32

const (
	InlineType_NONE                InlineType = 0
	InlineType_URI_ONLY            InlineType = 4
	InlineType_LABEL               InlineType = 1
	InlineType_MODAL               InlineType = 2
	InlineType_INLINE_DETAIL_BLOCK InlineType = 3
)

var InlineType_name = map[int32]string{
	0: "NONE",
	4: "URI_ONLY",
	1: "LABEL",
	2: "MODAL",
	3: "INLINE_DETAIL_BLOCK",
}

var InlineType_value = map[string]int32{
	"NONE":                0,
	"URI_ONLY":            4,
	"LABEL":               1,
	"MODAL":               2,
	"INLINE_DETAIL_BLOCK": 3,
}

func (x InlineType) String() string {
	return proto.EnumName(InlineType_name, int32(x))
}

func (InlineType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{3}
}

type DataSetConfig struct {
	ID                   string            `protobuf:"bytes,11,opt,name=ID,proto3" json:"ID,omitempty"`
	Title                string            `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Markdown             string            `protobuf:"bytes,2,opt,name=markdown,proto3" json:"markdown,omitempty"`
	DataSetType          DataSetType       `protobuf:"varint,3,opt,name=dataSetType,proto3,enum=fragments.DataSetType" json:"dataSetType,omitempty"`
	Facets               []*FacetField     `protobuf:"bytes,4,rep,name=facets,proto3" json:"facets,omitempty"`
	Spec                 []string          `protobuf:"bytes,5,rep,name=spec,proto3" json:"spec,omitempty"`
	ExcludeSpec          []string          `protobuf:"bytes,10,rep,name=excludeSpec,proto3" json:"excludeSpec,omitempty"`
	ViewConfig           *DetailViewConfig `protobuf:"bytes,6,opt,name=viewConfig,proto3" json:"viewConfig,omitempty"`
	ResultConfig         *ResultViewConfig `protobuf:"bytes,7,opt,name=resultConfig,proto3" json:"resultConfig,omitempty"`
	Filter               *DataSetFilter    `protobuf:"bytes,8,opt,name=filter,proto3" json:"filter,omitempty"`
	OrgID                string            `protobuf:"bytes,9,opt,name=orgID,proto3" json:"orgID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *DataSetConfig) Reset()         { *m = DataSetConfig{} }
func (m *DataSetConfig) String() string { return proto.CompactTextString(m) }
func (*DataSetConfig) ProtoMessage()    {}
func (*DataSetConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{0}
}

func (m *DataSetConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataSetConfig.Unmarshal(m, b)
}
func (m *DataSetConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataSetConfig.Marshal(b, m, deterministic)
}
func (m *DataSetConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataSetConfig.Merge(m, src)
}
func (m *DataSetConfig) XXX_Size() int {
	return xxx_messageInfo_DataSetConfig.Size(m)
}
func (m *DataSetConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DataSetConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DataSetConfig proto.InternalMessageInfo

func (m *DataSetConfig) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *DataSetConfig) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *DataSetConfig) GetMarkdown() string {
	if m != nil {
		return m.Markdown
	}
	return ""
}

func (m *DataSetConfig) GetDataSetType() DataSetType {
	if m != nil {
		return m.DataSetType
	}
	return DataSetType_SINGLE
}

func (m *DataSetConfig) GetFacets() []*FacetField {
	if m != nil {
		return m.Facets
	}
	return nil
}

func (m *DataSetConfig) GetSpec() []string {
	if m != nil {
		return m.Spec
	}
	return nil
}

func (m *DataSetConfig) GetExcludeSpec() []string {
	if m != nil {
		return m.ExcludeSpec
	}
	return nil
}

func (m *DataSetConfig) GetViewConfig() *DetailViewConfig {
	if m != nil {
		return m.ViewConfig
	}
	return nil
}

func (m *DataSetConfig) GetResultConfig() *ResultViewConfig {
	if m != nil {
		return m.ResultConfig
	}
	return nil
}

func (m *DataSetConfig) GetFilter() *DataSetFilter {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *DataSetConfig) GetOrgID() string {
	if m != nil {
		return m.OrgID
	}
	return ""
}

type DataSetFilter struct {
	QueryFilter          []*QueryFilter `protobuf:"bytes,1,rep,name=queryFilter,proto3" json:"queryFilter,omitempty"`
	Query                string         `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *DataSetFilter) Reset()         { *m = DataSetFilter{} }
func (m *DataSetFilter) String() string { return proto.CompactTextString(m) }
func (*DataSetFilter) ProtoMessage()    {}
func (*DataSetFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{1}
}

func (m *DataSetFilter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataSetFilter.Unmarshal(m, b)
}
func (m *DataSetFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataSetFilter.Marshal(b, m, deterministic)
}
func (m *DataSetFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataSetFilter.Merge(m, src)
}
func (m *DataSetFilter) XXX_Size() int {
	return xxx_messageInfo_DataSetFilter.Size(m)
}
func (m *DataSetFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_DataSetFilter.DiscardUnknown(m)
}

var xxx_messageInfo_DataSetFilter proto.InternalMessageInfo

func (m *DataSetFilter) GetQueryFilter() []*QueryFilter {
	if m != nil {
		return m.QueryFilter
	}
	return nil
}

func (m *DataSetFilter) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

type ResultViewConfig struct {
	ResultType           ResultType         `protobuf:"varint,1,opt,name=resultType,proto3,enum=fragments.ResultType" json:"resultType,omitempty"`
	Fields               *PresentationField `protobuf:"bytes,2,opt,name=fields,proto3" json:"fields,omitempty"`
	Inline               bool               `protobuf:"varint,3,opt,name=inline,proto3" json:"inline,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ResultViewConfig) Reset()         { *m = ResultViewConfig{} }
func (m *ResultViewConfig) String() string { return proto.CompactTextString(m) }
func (*ResultViewConfig) ProtoMessage()    {}
func (*ResultViewConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{2}
}

func (m *ResultViewConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResultViewConfig.Unmarshal(m, b)
}
func (m *ResultViewConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResultViewConfig.Marshal(b, m, deterministic)
}
func (m *ResultViewConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResultViewConfig.Merge(m, src)
}
func (m *ResultViewConfig) XXX_Size() int {
	return xxx_messageInfo_ResultViewConfig.Size(m)
}
func (m *ResultViewConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ResultViewConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ResultViewConfig proto.InternalMessageInfo

func (m *ResultViewConfig) GetResultType() ResultType {
	if m != nil {
		return m.ResultType
	}
	return ResultType_GRID
}

func (m *ResultViewConfig) GetFields() *PresentationField {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *ResultViewConfig) GetInline() bool {
	if m != nil {
		return m.Inline
	}
	return false
}

type PresentationField struct {
	I18NLabel            *I18NLabel `protobuf:"bytes,1,opt,name=i18nLabel,proto3" json:"i18nLabel,omitempty"`
	Clickable            bool       `protobuf:"varint,2,opt,name=clickable,proto3" json:"clickable,omitempty"`
	Searchable           bool       `protobuf:"varint,7,opt,name=searchable,proto3" json:"searchable,omitempty"`
	Predicate            string     `protobuf:"bytes,3,opt,name=predicate,proto3" json:"predicate,omitempty"`
	Single               bool       `protobuf:"varint,4,opt,name=single,proto3" json:"single,omitempty"`
	Order                int32      `protobuf:"varint,5,opt,name=order,proto3" json:"order,omitempty"`
	FieldType            FieldType  `protobuf:"varint,6,opt,name=fieldType,proto3,enum=fragments.FieldType" json:"fieldType,omitempty"`
	InlineType           InlineType `protobuf:"varint,8,opt,name=inlineType,proto3,enum=fragments.InlineType" json:"inlineType,omitempty"`
	InlineCSS            string     `protobuf:"bytes,9,opt,name=inlineCSS,proto3" json:"inlineCSS,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *PresentationField) Reset()         { *m = PresentationField{} }
func (m *PresentationField) String() string { return proto.CompactTextString(m) }
func (*PresentationField) ProtoMessage()    {}
func (*PresentationField) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{3}
}

func (m *PresentationField) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PresentationField.Unmarshal(m, b)
}
func (m *PresentationField) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PresentationField.Marshal(b, m, deterministic)
}
func (m *PresentationField) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PresentationField.Merge(m, src)
}
func (m *PresentationField) XXX_Size() int {
	return xxx_messageInfo_PresentationField.Size(m)
}
func (m *PresentationField) XXX_DiscardUnknown() {
	xxx_messageInfo_PresentationField.DiscardUnknown(m)
}

var xxx_messageInfo_PresentationField proto.InternalMessageInfo

func (m *PresentationField) GetI18NLabel() *I18NLabel {
	if m != nil {
		return m.I18NLabel
	}
	return nil
}

func (m *PresentationField) GetClickable() bool {
	if m != nil {
		return m.Clickable
	}
	return false
}

func (m *PresentationField) GetSearchable() bool {
	if m != nil {
		return m.Searchable
	}
	return false
}

func (m *PresentationField) GetPredicate() string {
	if m != nil {
		return m.Predicate
	}
	return ""
}

func (m *PresentationField) GetSingle() bool {
	if m != nil {
		return m.Single
	}
	return false
}

func (m *PresentationField) GetOrder() int32 {
	if m != nil {
		return m.Order
	}
	return 0
}

func (m *PresentationField) GetFieldType() FieldType {
	if m != nil {
		return m.FieldType
	}
	return FieldType_LITERAL
}

func (m *PresentationField) GetInlineType() InlineType {
	if m != nil {
		return m.InlineType
	}
	return InlineType_NONE
}

func (m *PresentationField) GetInlineCSS() string {
	if m != nil {
		return m.InlineCSS
	}
	return ""
}

type DetailViewConfig struct {
	Blocks               []*DetailBlock `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
	EntryType            string         `protobuf:"bytes,2,opt,name=entryType,proto3" json:"entryType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *DetailViewConfig) Reset()         { *m = DetailViewConfig{} }
func (m *DetailViewConfig) String() string { return proto.CompactTextString(m) }
func (*DetailViewConfig) ProtoMessage()    {}
func (*DetailViewConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{4}
}

func (m *DetailViewConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DetailViewConfig.Unmarshal(m, b)
}
func (m *DetailViewConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DetailViewConfig.Marshal(b, m, deterministic)
}
func (m *DetailViewConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DetailViewConfig.Merge(m, src)
}
func (m *DetailViewConfig) XXX_Size() int {
	return xxx_messageInfo_DetailViewConfig.Size(m)
}
func (m *DetailViewConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DetailViewConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DetailViewConfig proto.InternalMessageInfo

func (m *DetailViewConfig) GetBlocks() []*DetailBlock {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func (m *DetailViewConfig) GetEntryType() string {
	if m != nil {
		return m.EntryType
	}
	return ""
}

type I18NLabel struct {
	Lang                 string   `protobuf:"bytes,1,opt,name=lang,proto3" json:"lang,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *I18NLabel) Reset()         { *m = I18NLabel{} }
func (m *I18NLabel) String() string { return proto.CompactTextString(m) }
func (*I18NLabel) ProtoMessage()    {}
func (*I18NLabel) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{5}
}

func (m *I18NLabel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_I18NLabel.Unmarshal(m, b)
}
func (m *I18NLabel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_I18NLabel.Marshal(b, m, deterministic)
}
func (m *I18NLabel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_I18NLabel.Merge(m, src)
}
func (m *I18NLabel) XXX_Size() int {
	return xxx_messageInfo_I18NLabel.Size(m)
}
func (m *I18NLabel) XXX_DiscardUnknown() {
	xxx_messageInfo_I18NLabel.DiscardUnknown(m)
}

var xxx_messageInfo_I18NLabel proto.InternalMessageInfo

func (m *I18NLabel) GetLang() string {
	if m != nil {
		return m.Lang
	}
	return ""
}

func (m *I18NLabel) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type DetailBlock struct {
	I18NLabel            *I18NLabel           `protobuf:"bytes,1,opt,name=i18nLabel,proto3" json:"i18nLabel,omitempty"`
	ResourceLabel        string               `protobuf:"bytes,2,opt,name=resourceLabel,proto3" json:"resourceLabel,omitempty"`
	Order                int32                `protobuf:"varint,3,opt,name=order,proto3" json:"order,omitempty"`
	ResourceType         string               `protobuf:"bytes,4,opt,name=resourceType,proto3" json:"resourceType,omitempty"`
	Fields               []*PresentationField `protobuf:"bytes,5,rep,name=fields,proto3" json:"fields,omitempty"`
	InlineCSS            string               `protobuf:"bytes,6,opt,name=inlineCSS,proto3" json:"inlineCSS,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *DetailBlock) Reset()         { *m = DetailBlock{} }
func (m *DetailBlock) String() string { return proto.CompactTextString(m) }
func (*DetailBlock) ProtoMessage()    {}
func (*DetailBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a99035a5d046e2, []int{6}
}

func (m *DetailBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DetailBlock.Unmarshal(m, b)
}
func (m *DetailBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DetailBlock.Marshal(b, m, deterministic)
}
func (m *DetailBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DetailBlock.Merge(m, src)
}
func (m *DetailBlock) XXX_Size() int {
	return xxx_messageInfo_DetailBlock.Size(m)
}
func (m *DetailBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_DetailBlock.DiscardUnknown(m)
}

var xxx_messageInfo_DetailBlock proto.InternalMessageInfo

func (m *DetailBlock) GetI18NLabel() *I18NLabel {
	if m != nil {
		return m.I18NLabel
	}
	return nil
}

func (m *DetailBlock) GetResourceLabel() string {
	if m != nil {
		return m.ResourceLabel
	}
	return ""
}

func (m *DetailBlock) GetOrder() int32 {
	if m != nil {
		return m.Order
	}
	return 0
}

func (m *DetailBlock) GetResourceType() string {
	if m != nil {
		return m.ResourceType
	}
	return ""
}

func (m *DetailBlock) GetFields() []*PresentationField {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *DetailBlock) GetInlineCSS() string {
	if m != nil {
		return m.InlineCSS
	}
	return ""
}

func init() {
	proto.RegisterEnum("fragments.DataSetType", DataSetType_name, DataSetType_value)
	proto.RegisterEnum("fragments.ResultType", ResultType_name, ResultType_value)
	proto.RegisterEnum("fragments.FieldType", FieldType_name, FieldType_value)
	proto.RegisterEnum("fragments.InlineType", InlineType_name, InlineType_value)
	proto.RegisterType((*DataSetConfig)(nil), "fragments.DataSetConfig")
	proto.RegisterType((*DataSetFilter)(nil), "fragments.DataSetFilter")
	proto.RegisterType((*ResultViewConfig)(nil), "fragments.ResultViewConfig")
	proto.RegisterType((*PresentationField)(nil), "fragments.PresentationField")
	proto.RegisterType((*DetailViewConfig)(nil), "fragments.DetailViewConfig")
	proto.RegisterType((*I18NLabel)(nil), "fragments.I18NLabel")
	proto.RegisterType((*DetailBlock)(nil), "fragments.DetailBlock")
}

func init() { proto.RegisterFile("hub3/fragments/viewconfig.proto", fileDescriptor_15a99035a5d046e2) }

var fileDescriptor_15a99035a5d046e2 = []byte{
	// 860 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x5f, 0x6f, 0xdb, 0x54,
	0x14, 0x9f, 0xe3, 0xc4, 0x8d, 0x8f, 0xb7, 0xca, 0x5c, 0xb6, 0x61, 0x8d, 0x09, 0xa2, 0x88, 0x87,
	0xa8, 0x12, 0x1d, 0x4b, 0x87, 0xa8, 0xc4, 0x03, 0x72, 0x62, 0xb7, 0x5c, 0x70, 0x9d, 0xee, 0xc6,
	0x9d, 0xd4, 0x07, 0x14, 0x1c, 0xe7, 0x36, 0xb3, 0xea, 0x3a, 0xc1, 0x76, 0x18, 0xfd, 0x24, 0x7c,
	0x13, 0x3e, 0x0c, 0xdf, 0x83, 0x77, 0x74, 0xcf, 0x75, 0xe3, 0xdb, 0x86, 0x07, 0xb4, 0x37, 0x9f,
	0xf3, 0xfb, 0xfd, 0x7c, 0xfe, 0xfc, 0x8e, 0x13, 0xf8, 0xf2, 0xfd, 0x66, 0x7e, 0xf4, 0xea, 0xaa,
	0x88, 0x97, 0x37, 0x3c, 0xaf, 0xca, 0x57, 0xbf, 0xa7, 0xfc, 0x43, 0xb2, 0xca, 0xaf, 0xd2, 0xe5,
	0xe1, 0xba, 0x58, 0x55, 0x2b, 0x62, 0x6e, 0xb1, 0x17, 0xce, 0x03, 0x6e, 0xbc, 0x4e, 0x25, 0xa9,
	0xff, 0x97, 0x0e, 0x4f, 0xbc, 0xb8, 0x8a, 0xa7, 0xbc, 0x1a, 0xa3, 0x98, 0xec, 0x43, 0x8b, 0x7a,
	0x8e, 0xd5, 0xd3, 0x06, 0x26, 0x6b, 0x51, 0x8f, 0x3c, 0x85, 0x4e, 0x95, 0x56, 0x19, 0x77, 0x34,
	0x4c, 0xc9, 0x80, 0xbc, 0x80, 0xee, 0x4d, 0x5c, 0x5c, 0x2f, 0x56, 0x1f, 0x72, 0xa7, 0x85, 0xc0,
	0x36, 0x26, 0xc7, 0x60, 0x2d, 0xe4, 0x2b, 0xa3, 0xdb, 0x35, 0x77, 0xf4, 0x9e, 0x36, 0xd8, 0x1f,
	0x3e, 0x3f, 0xdc, 0x96, 0x3f, 0xf4, 0x1a, 0x94, 0xa9, 0x54, 0xf2, 0x35, 0x18, 0x57, 0x71, 0xc2,
	0xab, 0xd2, 0x69, 0xf7, 0xf4, 0x81, 0x35, 0x7c, 0xa6, 0x88, 0x4e, 0x04, 0x70, 0x92, 0xf2, 0x6c,
	0xc1, 0x6a, 0x12, 0x21, 0xd0, 0x2e, 0xd7, 0x3c, 0x71, 0x3a, 0x3d, 0x7d, 0x60, 0x32, 0x7c, 0x26,
	0x3d, 0xb0, 0xf8, 0x1f, 0x49, 0xb6, 0x59, 0xf0, 0xa9, 0x80, 0x00, 0x21, 0x35, 0x45, 0xbe, 0x07,
	0x10, 0xbb, 0x92, 0xe3, 0x3a, 0x46, 0x4f, 0x1b, 0x58, 0xc3, 0xcf, 0xd5, 0xee, 0x78, 0x15, 0xa7,
	0xd9, 0xbb, 0x2d, 0x85, 0x29, 0x74, 0xf2, 0x03, 0x3c, 0x2e, 0x78, 0xb9, 0xc9, 0xea, 0x6d, 0x39,
	0x7b, 0x3b, 0x72, 0x86, 0xb0, 0x22, 0xbf, 0x27, 0x20, 0xdf, 0x80, 0x71, 0x95, 0x66, 0x15, 0x2f,
	0x9c, 0x2e, 0x4a, 0x9d, 0xdd, 0xbd, 0x9c, 0x20, 0xce, 0x6a, 0x9e, 0x30, 0x60, 0x55, 0x2c, 0xa9,
	0xe7, 0x98, 0xd2, 0x00, 0x0c, 0xfa, 0xb3, 0xad, 0x6f, 0x92, 0x2e, 0xb6, 0xfe, 0xdb, 0x86, 0x17,
	0xb7, 0x32, 0x74, 0x34, 0x5c, 0xa0, 0xba, 0xf5, 0xb7, 0x0d, 0xca, 0x54, 0xaa, 0x28, 0x80, 0x61,
	0x6d, 0xa4, 0x0c, 0xfa, 0x7f, 0x6a, 0x60, 0x3f, 0x9c, 0x85, 0x7c, 0x0b, 0x20, 0xa7, 0x41, 0x67,
	0x35, 0x74, 0xf6, 0xd9, 0xce, 0xf0, 0x68, 0xac, 0x42, 0x24, 0x6f, 0xc4, 0xd0, 0x3c, 0x5b, 0x94,
	0x58, 0xc2, 0x1a, 0xbe, 0x54, 0x24, 0xe7, 0x05, 0x2f, 0x79, 0x5e, 0xc5, 0x55, 0xba, 0xca, 0xef,
	0xec, 0x45, 0x2e, 0x79, 0x0e, 0x46, 0x9a, 0x67, 0x69, 0x2e, 0x4f, 0xa8, 0xcb, 0xea, 0xa8, 0xff,
	0x77, 0x0b, 0x3e, 0xd9, 0x51, 0x91, 0x21, 0x98, 0xe9, 0xeb, 0xe3, 0x3c, 0x88, 0xe7, 0x3c, 0xc3,
	0xce, 0xac, 0xe1, 0x53, 0xa5, 0x0c, 0x7d, 0x7d, 0x1c, 0x22, 0xc6, 0x1a, 0x1a, 0x79, 0x09, 0x66,
	0x92, 0xa5, 0xc9, 0x75, 0x3c, 0xcf, 0x38, 0xb6, 0xd6, 0x65, 0x4d, 0x82, 0x7c, 0x01, 0x50, 0xf2,
	0xb8, 0x48, 0xde, 0x23, 0xbc, 0x87, 0xb0, 0x92, 0x11, 0xea, 0x75, 0xc1, 0x17, 0x69, 0x12, 0x57,
	0xb2, 0x45, 0x93, 0x35, 0x09, 0xd1, 0x7d, 0x99, 0xe6, 0xcb, 0x8c, 0x3b, 0x6d, 0xd9, 0xbd, 0x8c,
	0xa4, 0x9d, 0x0b, 0x5e, 0x38, 0x9d, 0x9e, 0x36, 0xe8, 0x30, 0x19, 0x88, 0xee, 0x71, 0x6a, 0xdc,
	0xab, 0x81, 0x7b, 0x55, 0xbb, 0x3f, 0xb9, 0xc3, 0x58, 0x43, 0x13, 0x66, 0xc8, 0x8d, 0xa0, 0xa8,
	0xbb, 0x63, 0x06, 0xdd, 0x82, 0x4c, 0x21, 0x8a, 0xb6, 0x65, 0x34, 0x9e, 0x4e, 0xeb, 0x9b, 0x6a,
	0x12, 0xfd, 0x5f, 0xc1, 0x7e, 0xf8, 0x01, 0x90, 0x43, 0x30, 0xe6, 0xd9, 0x2a, 0xb9, 0x2e, 0xff,
	0xe3, 0xaa, 0x24, 0x79, 0x24, 0x60, 0x56, 0xb3, 0x44, 0x05, 0x9e, 0x57, 0xc5, 0x2d, 0xf6, 0x25,
	0x8f, 0xaa, 0x49, 0xf4, 0x8f, 0xc0, 0xdc, 0x9a, 0x21, 0x3e, 0xe1, 0x2c, 0xce, 0x97, 0xf5, 0x8f,
	0x0b, 0x3e, 0x8b, 0x5c, 0x1e, 0xdf, 0xdc, 0x29, 0xf1, 0xb9, 0xff, 0x8f, 0x06, 0x96, 0x52, 0xea,
	0xa3, 0xdc, 0xfe, 0x0a, 0x9e, 0x14, 0xbc, 0x5c, 0x6d, 0x8a, 0x84, 0x4b, 0x9d, 0x2c, 0x70, 0x3f,
	0xd9, 0xf8, 0xa3, 0xab, 0xfe, 0xf4, 0xf1, 0xbb, 0x47, 0x1a, 0x4e, 0xd5, 0x46, 0xe9, 0xbd, 0x9c,
	0x72, 0xe5, 0x1d, 0x5c, 0xd3, 0xff, 0xbb, 0xf2, 0x7b, 0x76, 0x18, 0x0f, 0xec, 0x38, 0x78, 0x03,
	0x96, 0xf2, 0x6b, 0x49, 0x00, 0x8c, 0x29, 0x0d, 0x4f, 0x03, 0xdf, 0x7e, 0x44, 0x1e, 0x43, 0xf7,
	0xec, 0x22, 0x88, 0xe8, 0x79, 0xe0, 0xdb, 0x1a, 0xb1, 0x60, 0x6f, 0x74, 0xf9, 0xf6, 0xc2, 0x67,
	0x97, 0x76, 0xeb, 0xe0, 0x3b, 0x80, 0xe6, 0x4b, 0x24, 0x5d, 0x68, 0x9f, 0x32, 0xea, 0xd9, 0x8f,
	0x88, 0x09, 0x9d, 0xc8, 0x1d, 0x21, 0x7f, 0x0f, 0xf4, 0x33, 0xf7, 0xdc, 0x6e, 0x09, 0xa1, 0xcb,
	0xc6, 0x3f, 0xd2, 0x77, 0xbe, 0xad, 0x1f, 0xfc, 0x02, 0xe6, 0xf6, 0xd4, 0x04, 0x12, 0xd0, 0xc8,
	0x67, 0x6e, 0x20, 0xab, 0x31, 0x7f, 0x3a, 0xb9, 0x60, 0x63, 0xa1, 0xee, 0x42, 0xdb, 0x73, 0x23,
	0xdf, 0x6e, 0x89, 0x57, 0x9e, 0x4f, 0x68, 0x18, 0xd9, 0x3a, 0x21, 0xb0, 0xef, 0xd1, 0x53, 0x1a,
	0xb9, 0xc1, 0x6c, 0x32, 0xfa, 0xc9, 0x1f, 0x47, 0x76, 0x1b, 0x9b, 0x74, 0x43, 0x7a, 0xe2, 0x4f,
	0x23, 0xbb, 0x73, 0x30, 0x05, 0x68, 0x8e, 0x52, 0xbc, 0x24, 0x9c, 0x84, 0xf5, 0x28, 0x17, 0x8c,
	0xce, 0x26, 0x61, 0x70, 0x69, 0xb7, 0xc5, 0x2b, 0x03, 0x77, 0xe4, 0x07, 0xb6, 0x26, 0x1e, 0xcf,
	0x26, 0x9e, 0x1b, 0xd8, 0x2d, 0xf2, 0x19, 0x7c, 0x4a, 0xc3, 0x80, 0x86, 0xfe, 0xcc, 0xf3, 0x23,
	0x97, 0x06, 0xb3, 0x51, 0x30, 0x19, 0xff, 0x6c, 0xeb, 0x73, 0x03, 0xff, 0xc9, 0x8e, 0xfe, 0x0d,
	0x00, 0x00, 0xff, 0xff, 0xf7, 0xcb, 0x0a, 0xe1, 0x11, 0x07, 0x00, 0x00,
}
