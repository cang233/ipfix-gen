package ipfix

//fill every head in message,including length,padding.
func Filling(msg *Message) {
	length := uint16(16) //header 16 bytes

	//every sets length
	for i := range msg.TemplateSet {
		fillingTemplate(&(msg.TemplateSet[i]))
		length += msg.TemplateSet[i].Header.Length
	}
	for i := range msg.OptionsTemplateSet {
		fillingOptionTemplate(&(msg.OptionsTemplateSet[i]))
		length += msg.OptionsTemplateSet[i].Header.Length
	}
	for i := range msg.DataSet {
		fillingDataSet(&(msg.DataSet[i]))
		length += msg.DataSet[i].Header.Length
	}

	msg.Header.Length = length
}

func fillingTemplate(tplSet *TemplateSet) {
	length := uint16(4) //set head

	for _, tpl := range tplSet.Templates {
		length += 4 //t head
		for _, field := range tpl.Fields {
			length += 4 //field
			if field.ID > 0x80 {
				length += 4 //enterpriseNo
			}
		}
	}
	tplSet.Header.Length = length
}
func fillingOptionTemplate(tplSet *OptionsTemplateSet) {
	length := uint16(4) //set head
	for _, tpl := range tplSet.OptionTemplates {
		if tpl.FieldCount > 0 {
			length += 6 //options template head
			for i := 0; i < int(tpl.ScopeFieldCount); i++ {
				length += 4 // id + length
				if tpl.Fields[i].ID > 0x80 {
					length += 4
				}
			}
			length += (tpl.FieldCount - tpl.ScopeFieldCount) * 4
		}
	}
	//padding
	if len(tplSet.OptionTemplates)%2 == 1 {
		length += 2
		tplSet.padding = 2
	}
	tplSet.Header.Length = length
}

//can not cal,val is interface{},need cal from template
//or cal by user
func fillingDataSet(dataset *DataSet) {
	length := uint16(4) // set len
	for _, d := range dataset.DataFields {
		length += uint16(InfoModel[ElementKey{0, d.FieldID}].Type.minLen())
		//fmt.Printf("FieldID:%d,reflectLen:%d,InfoModelLen:%d\n", d.FieldID, reflect.TypeOf(d.Value).Size(), InfoModel[ElementKey{0, d.FieldID}].Type.minLen())
	}
	if length%4 != 0 {
		dataset.padding = int(4 - length%4)
		length += 4 - length%4
	}
	dataset.Header.Length = length
}
