package ipfix

import (
	"bytes"
	"encoding/binary"
)

//无法控制输入的到底是什么啊，要求uint8也许会传int啊，这时候该怎么办啊?
//notice:map的读取顺序是不固定的，不固定的template顺序代表不同的template。
func buildMap(vals map[int]interface{}, templateID uint16) *Message {
	var fields []FieldSpecifier
	var dfs []DataField
	for i, v := range vals {
		fields = append(fields, FieldSpecifier{
			ID:           uint16(i),
			Length:       uint16(InfoModel[ElementKey{0, uint16(i)}].Type.minLen()),
			EnterpriseNo: 0,
		})
		dfs = append(dfs, DataField{
			FieldID: uint16(i),
			Value:   v,
		})
	}
	return &Message{
		Header: MessageHeader{
			Version:    VERSION,
			Length:     0,
			ExportTime: 0,
			SequenceNo: 0,
			DomainID:   0,
		},
		TemplateSet: []TemplateSet{
			{
				Header: SetHeader{
					ID:     2,
					Length: 0,
				},
				Templates: []TemplateRecord{{
					ID:         templateID,
					FieldCount: uint16(len(vals)),
					Fields:     fields,
				}},
			},
		},
		OptionsTemplateSet: nil,
		DataSet: []DataSet{
			{
				Header: SetHeader{
					ID:     templateID,
					Length: 0,
				},
				DataFields: dfs,
			},
		},
	}
}

func BuildArr(IDs []uint16, Vals []interface{}, templateID uint16) *Message {
	var fields []FieldSpecifier
	var dfs []DataField
	for i := 0; i < len(IDs); i++ {
		fields = append(fields, FieldSpecifier{
			ID:           IDs[i],
			Length:       uint16(InfoModel[ElementKey{0, IDs[i]}].Type.minLen()),
			EnterpriseNo: 0,
		})
		dfs = append(dfs, DataField{
			FieldID: IDs[i],
			Value:   Vals[i],
		})
	}
	return &Message{
		Header: MessageHeader{
			Version:    VERSION,
			Length:     0,
			ExportTime: 0,
			SequenceNo: 0,
			DomainID:   0,
		},
		TemplateSet: []TemplateSet{
			{
				Header: SetHeader{
					ID:     2,
					Length: 0,
				},
				Templates: []TemplateRecord{{
					ID:         templateID,
					FieldCount: uint16(len(IDs)),
					Fields:     fields,
				}},
			},
		},
		OptionsTemplateSet: nil,
		DataSet: []DataSet{
			{
				Header: SetHeader{
					ID:     templateID,
					Length: 0,
				},
				DataFields: dfs,
			},
		},
	}
}

func convert(id uint16, val interface{}) interface{} {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, val)
	bs := b.Bytes()
	return Interpret(&bs, InfoModel[ElementKey{0, id}].Type)
}
