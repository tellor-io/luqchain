// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: luqchain/bridge/inclusionproof.proto

package bridge

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type InclusionProofStuffFields struct {
	RootHash string `protobuf:"bytes,1,opt,name=RootHash,proto3" json:"RootHash,omitempty"`
	Version  int64  `protobuf:"varint,2,opt,name=Version,proto3" json:"Version,omitempty"`
	Key      string `protobuf:"bytes,3,opt,name=Key,proto3" json:"Key,omitempty"`
	DataHash string `protobuf:"bytes,4,opt,name=DataHash,proto3" json:"DataHash,omitempty"`
}

func (m *InclusionProofStuffFields) Reset()         { *m = InclusionProofStuffFields{} }
func (m *InclusionProofStuffFields) String() string { return proto.CompactTextString(m) }
func (*InclusionProofStuffFields) ProtoMessage()    {}
func (*InclusionProofStuffFields) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ae044566408f9b2, []int{0}
}
func (m *InclusionProofStuffFields) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InclusionProofStuffFields) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InclusionProofStuffFields.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InclusionProofStuffFields) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InclusionProofStuffFields.Merge(m, src)
}
func (m *InclusionProofStuffFields) XXX_Size() int {
	return m.Size()
}
func (m *InclusionProofStuffFields) XXX_DiscardUnknown() {
	xxx_messageInfo_InclusionProofStuffFields.DiscardUnknown(m)
}

var xxx_messageInfo_InclusionProofStuffFields proto.InternalMessageInfo

func (m *InclusionProofStuffFields) GetRootHash() string {
	if m != nil {
		return m.RootHash
	}
	return ""
}

func (m *InclusionProofStuffFields) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *InclusionProofStuffFields) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *InclusionProofStuffFields) GetDataHash() string {
	if m != nil {
		return m.DataHash
	}
	return ""
}

func init() {
	proto.RegisterType((*InclusionProofStuffFields)(nil), "luqchain.bridge.InclusionProofStuffFields")
}

func init() {
	proto.RegisterFile("luqchain/bridge/inclusionproof.proto", fileDescriptor_1ae044566408f9b2)
}

var fileDescriptor_1ae044566408f9b2 = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xc9, 0x29, 0x2d, 0x4c,
	0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x4f, 0x2a, 0xca, 0x4c, 0x49, 0x4f, 0xd5, 0xcf, 0xcc, 0x4b, 0xce,
	0x29, 0x2d, 0xce, 0xcc, 0xcf, 0x2b, 0x28, 0xca, 0xcf, 0x4f, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x87, 0xa9, 0xd2, 0x83, 0xa8, 0x52, 0xaa, 0xe7, 0x92, 0xf4, 0x84, 0x29, 0x0c, 0x00,
	0x29, 0x0c, 0x2e, 0x29, 0x4d, 0x4b, 0x73, 0xcb, 0x4c, 0xcd, 0x49, 0x29, 0x16, 0x92, 0xe2, 0xe2,
	0x08, 0xca, 0xcf, 0x2f, 0xf1, 0x48, 0x2c, 0xce, 0x90, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82,
	0xf3, 0x85, 0x24, 0xb8, 0xd8, 0xc3, 0x52, 0x8b, 0x40, 0xda, 0x24, 0x98, 0x14, 0x18, 0x35, 0x98,
	0x83, 0x60, 0x5c, 0x21, 0x01, 0x2e, 0x66, 0xef, 0xd4, 0x4a, 0x09, 0x66, 0xb0, 0x06, 0x10, 0x13,
	0x64, 0x8e, 0x4b, 0x62, 0x49, 0x22, 0xd8, 0x1c, 0x16, 0x88, 0x39, 0x30, 0xbe, 0x93, 0xe6, 0x89,
	0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3,
	0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0xf1, 0xa3, 0xf9, 0x28, 0x89, 0x0d, 0xec,
	0x07, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x37, 0x2c, 0xee, 0x50, 0xeb, 0x00, 0x00, 0x00,
}

func (m *InclusionProofStuffFields) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InclusionProofStuffFields) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InclusionProofStuffFields) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DataHash) > 0 {
		i -= len(m.DataHash)
		copy(dAtA[i:], m.DataHash)
		i = encodeVarintInclusionproof(dAtA, i, uint64(len(m.DataHash)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintInclusionproof(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Version != 0 {
		i = encodeVarintInclusionproof(dAtA, i, uint64(m.Version))
		i--
		dAtA[i] = 0x10
	}
	if len(m.RootHash) > 0 {
		i -= len(m.RootHash)
		copy(dAtA[i:], m.RootHash)
		i = encodeVarintInclusionproof(dAtA, i, uint64(len(m.RootHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintInclusionproof(dAtA []byte, offset int, v uint64) int {
	offset -= sovInclusionproof(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *InclusionProofStuffFields) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RootHash)
	if l > 0 {
		n += 1 + l + sovInclusionproof(uint64(l))
	}
	if m.Version != 0 {
		n += 1 + sovInclusionproof(uint64(m.Version))
	}
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovInclusionproof(uint64(l))
	}
	l = len(m.DataHash)
	if l > 0 {
		n += 1 + l + sovInclusionproof(uint64(l))
	}
	return n
}

func sovInclusionproof(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozInclusionproof(x uint64) (n int) {
	return sovInclusionproof(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *InclusionProofStuffFields) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInclusionproof
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: InclusionProofStuffFields: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InclusionProofStuffFields: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RootHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInclusionproof
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInclusionproof
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInclusionproof
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RootHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInclusionproof
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInclusionproof
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInclusionproof
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInclusionproof
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInclusionproof
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInclusionproof
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInclusionproof
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInclusionproof(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthInclusionproof
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipInclusionproof(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowInclusionproof
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInclusionproof
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInclusionproof
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthInclusionproof
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupInclusionproof
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthInclusionproof
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthInclusionproof        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowInclusionproof          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupInclusionproof = fmt.Errorf("proto: unexpected end of group")
)