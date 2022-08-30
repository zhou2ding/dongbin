package rpcpackage

import (
	"blog/pkg/rand"
	"bytes"
	"encoding/binary"
)

var PacketTail = []byte{'Z', 'D', 'B'}

var RpcVersion uint16 = 0x1010

type rpcMsgSplitter struct {
}

var gRPCMsgSplitter rpcMsgSplitter

func GetRPCMsgSplitterInstance() *rpcMsgSplitter {
	return &gRPCMsgSplitter
}

func (s *rpcMsgSplitter) splitPacketHelper(data []byte, isJson bool, totalPktLength int, totalSubPktCount int, subPktIdx int, seq uint16) [][]byte {
	var packets [][]byte
	dataLen := len(data)
	if dataLen < PacketPayloadLength { //one packet
		var header []byte
		if isJson {
			header = s.makeMsgHeader(uint32(totalPktLength), uint32(dataLen), uint32(totalSubPktCount), uint32(subPktIdx), uint32(dataLen), 0, seq)
		} else {
			header = s.makeMsgHeader(uint32(totalPktLength), uint32(dataLen), uint32(totalSubPktCount), uint32(subPktIdx), 0, uint32(dataLen), seq)
		}

		packets = append(packets, header)
		packets = append(packets, data)
		packets = append(packets, PacketTail)
	} else { //several packets
		pktCount := dataLen / (PacketPayloadLength)
		if dataLen%(PacketPayloadLength) != 0 {
			pktCount += 1
		}

		nextPos := 0
		for idx := 0; idx < pktCount; idx++ {
			if idx != pktCount-1 { //not the last packet
				if isJson {
					packets = append(packets, s.makeMsgHeader(uint32(totalPktLength), PacketPayloadLength, uint32(totalSubPktCount), uint32(subPktIdx),
						PacketPayloadLength, 0, seq))
				} else {
					packets = append(packets, s.makeMsgHeader(uint32(totalPktLength), PacketPayloadLength, uint32(totalSubPktCount), uint32(subPktIdx),
						0, PacketPayloadLength, seq))
				}
				packets = append(packets, data[nextPos:nextPos+PacketPayloadLength])
				packets = append(packets, PacketTail)
				nextPos += PacketPayloadLength
			} else { //the last packet
				if isJson {
					packets = append(packets, s.makeMsgHeader(uint32(totalPktLength), uint32(dataLen-nextPos), uint32(totalSubPktCount), uint32(subPktIdx),
						uint32(dataLen-nextPos), 0, seq))
				} else {
					packets = append(packets, s.makeMsgHeader(uint32(totalPktLength), uint32(dataLen-nextPos), uint32(totalSubPktCount), uint32(subPktIdx),
						0, uint32(dataLen-nextPos), seq))
				}
				packets = append(packets, data[nextPos:])
				packets = append(packets, PacketTail)
				nextPos += len(data[nextPos:])
			}
			subPktIdx++
		}
	}

	return packets
}

func (s *rpcMsgSplitter) SplitPacket(jsonData []byte, binaryData []byte) [][]byte {
	jsonLength := len(jsonData)
	binaryLength := len(binaryData)
	var seq uint16 = 0
	var allPackets [][]byte
	if jsonLength != 0 && binaryLength == 0 { //
		jsonPktCount := jsonLength / PacketPayloadLength
		if jsonLength%PacketPayloadLength != 0 {
			jsonPktCount += 1
		}
		if jsonPktCount > 1 {
			seq = uint16(rand.GetIdGeneratorInstance().GetId() % 0x10000)
		}
		allPackets = s.splitPacketHelper(jsonData, true, jsonLength, jsonPktCount, 0, seq)
	} else if jsonLength == 0 && binaryLength != 0 {
		binaryCount := binaryLength / PacketPayloadLength
		if binaryLength%PacketPayloadLength != 0 {
			binaryCount += 1
		}
		if binaryCount > 1 {
			seq = uint16(rand.GetIdGeneratorInstance().GetId() % 0x10000)
		}
		allPackets = s.splitPacketHelper(binaryData, false, binaryLength, binaryCount, 0, seq)
	} else if jsonLength != 0 && binaryLength != 0 {
		jsonPktCount := jsonLength / PacketPayloadLength
		if jsonLength%PacketPayloadLength != 0 {
			jsonPktCount += 1
		}

		binaryCount := binaryLength / PacketPayloadLength
		if binaryLength%PacketPayloadLength != 0 {
			binaryCount += 1
		}

		totalLength := jsonLength + binaryLength
		totalPktCount := jsonPktCount + binaryCount

		seq = uint16(rand.GetIdGeneratorInstance().GetId() % 0x10000)

		allPackets = s.splitPacketHelper(jsonData, true, totalLength, totalPktCount, 0, seq)

		packets := s.splitPacketHelper(binaryData, false, totalLength, totalPktCount, jsonPktCount, seq)
		allPackets = append(allPackets, packets...)
	}

	return allPackets
}

func (s *rpcMsgSplitter) makeMsgHeader(pktLength uint32, subPktLength uint32, subPktPktCount uint32,
	subPktIndex uint32, pktJsonLength uint32, pktBinaryLength uint32, seq uint16) []byte {
	header := msgHeader{
		Prefix:          [4]byte{'I', 'S', 'V', 'H'},
		Version:         RpcVersion,
		Seq:             seq,
		PktLength:       pktLength,
		SubPktLength:    subPktLength,
		SubPktCount:     subPktPktCount,
		SubPktIndex:     subPktIndex,
		PktJsonLength:   pktJsonLength,
		PktBinaryLength: pktBinaryLength,
	}

	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.BigEndian, header)
	return buffer.Bytes()
}
