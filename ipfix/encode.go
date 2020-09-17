package ipfix

import (
	"bytes"
	"encoding/binary"
)

//Encode a Message to a IPFIX packet byte array.
func Encode(msg Message, seqNo uint32) []byte {

	if msg.Header.Length == 0 {
		Filling(&msg)
	}

	buf := new(bytes.Buffer)
	//orginal flow header
	binary.Write(buf, binary.BigEndian, msg.Header.Version)
	binary.Write(buf, binary.BigEndian, msg.Header.Length)

	binary.Write(buf, binary.BigEndian, msg.Header.ExportTime)
	binary.Write(buf, binary.BigEndian, seqNo)
	binary.Write(buf, binary.BigEndian, msg.Header.DomainID)

	for _, template := range msg.TemplateSet {
		writeTemplateSet(buf, template)
	}
	for _, template := range msg.OptionsTemplateSet {
		writeOptionTemplateSet(buf, template)
	}
	writeDataSet(buf, msg.DataSet)

	result := buf.Bytes()
	return result
}

func writeTemplateSet(buf *bytes.Buffer, tplSet TemplateSet) {
	binary.Write(buf, binary.BigEndian, tplSet.Header.ID)
	binary.Write(buf, binary.BigEndian, tplSet.Header.Length)

	if len(tplSet.Templates) == 0 {
		return
	}
	for _, template := range tplSet.Templates {
		writeTemplate(buf, template)
	}
}

func writeTemplate(buf *bytes.Buffer, tplRecord TemplateRecord) {
	if tplRecord.FieldCount > 0 {
		binary.Write(buf, binary.BigEndian, tplRecord.ID)
		binary.Write(buf, binary.BigEndian, tplRecord.FieldCount)
		for _, field := range tplRecord.Fields {
			binary.Write(buf, binary.BigEndian, field.ID)
			binary.Write(buf, binary.BigEndian, field.Length)
			if field.ID >= 0x80 { // E == 1
				binary.Write(buf, binary.BigEndian, field.EnterpriseNo)
			}
		}
	}
}

func writeOptionTemplateSet(buf *bytes.Buffer, tplSet OptionsTemplateSet) {
	binary.Write(buf, binary.BigEndian, tplSet.Header.ID)
	binary.Write(buf, binary.BigEndian, tplSet.Header.Length)

	if len(tplSet.OptionTemplates) == 0 {
		return
	}
	for _, template := range tplSet.OptionTemplates {
		writeOptionsTemplate(buf, template)
	}
	for i := 0; i < tplSet.padding; i++ {
		binary.Write(buf, binary.BigEndian, PADDING)
	}
}

func writeOptionsTemplate(buf *bytes.Buffer, tplRecord OptionTemplateRecord) {
	if tplRecord.FieldCount > 0 {
		binary.Write(buf, binary.BigEndian, tplRecord.ID)
		binary.Write(buf, binary.BigEndian, tplRecord.FieldCount)
		binary.Write(buf, binary.BigEndian, tplRecord.ScopeFieldCount)
		for i := 0; i < int(tplRecord.ScopeFieldCount); i++ {
			binary.Write(buf, binary.BigEndian, (tplRecord.Fields[i]).ID)
			binary.Write(buf, binary.BigEndian, (tplRecord.Fields[i]).Length)
			if (tplRecord.Fields[i]).ID >= 0x80 { // E == 1
				binary.Write(buf, binary.BigEndian, (tplRecord.Fields[i]).EnterpriseNo)
			}
		}
		for i := int(tplRecord.ScopeFieldCount); i < int(tplRecord.FieldCount); i++ {
			binary.Write(buf, binary.BigEndian, (tplRecord.Fields[i]).ID)
			binary.Write(buf, binary.BigEndian, (tplRecord.Fields[i]).Length)
		}
	}
}

func writeDataSet(buf *bytes.Buffer, dataSet []DataSet) {
	for _, flowSet := range dataSet {
		binary.Write(buf, binary.BigEndian, flowSet.Header.ID)
		binary.Write(buf, binary.BigEndian, flowSet.Header.Length)
		for _, field := range flowSet.DataFields {
			//fmt.Printf("value:[id=%d,val=%+v,reflect.type=%s,reflect.size:%d,buildin.type=%s,buildin.size=%d]\n",
			//	field.FieldID, field.Value, reflect.TypeOf(field.Value).Name(), reflect.TypeOf(field.Value).Size(),
			//	InfoModel[ElementKey{0, field.FieldID}].Name, InfoModel[ElementKey{0, field.FieldID}].Type.minLen())
			binary.Write(buf, binary.BigEndian, field.Value)
		}
		for i := 0; i < flowSet.padding; i++ {
			binary.Write(buf, binary.BigEndian, PADDING)
		}
	}
}
