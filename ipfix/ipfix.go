package ipfix

const (
	VERSION = uint16(10) // 4 byte
	PADDING = uint8(0)   // 1 byte
)

type Message struct {
	Header             MessageHeader        `json:"header"`
	TemplateSet        []TemplateSet        `json:"templateSet"`
	OptionsTemplateSet []OptionsTemplateSet `json:"optionSet"`
	DataSet            []DataSet            `json:"dataSet"`
}

//IPFIX message header
//0                   1                   2                   3
//0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|       Version Number          |            Length             |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|                           Export Time                         |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|                       Sequence Number                         |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|                    Observation Domain ID                      |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
type MessageHeader struct {
	Version    uint16 `json:"version"`     // Version of IPFIX to which this Message conforms
	Length     uint16 `json:"length"`      // Total length of the IPFIX Message, measured in octets
	ExportTime uint32 `json:"export_time"` // Time at which the IPFIX Message Header leaves the Exporter
	SequenceNo uint32 `json:"sequence_no"` // Incremental sequence counter modulo 2^32
	DomainID   uint32 `json:"domain_id"`   // A 32-bit id that is locally unique to the Exporting Process
}

// template set header
//0                   1                   2                   3
//0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|          Set ID               |          Length               |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
type SetHeader struct {
	ID     uint16 `json:"id"`
	Length uint16 `json:"length"`
}

type TemplateSet struct {
	Header    SetHeader        `json:"header"`
	Templates []TemplateRecord `json:"templates"`
}

type OptionsTemplateSet struct {
	Header          SetHeader              `json:"header"`
	OptionTemplates []OptionTemplateRecord `json:"options"`
	padding         int                    //byte count
}

//template record header
//0                   1                   2                   3
//0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|      Template ID (> 255)      |         Field Count           |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
type TemplateRecord struct {
	ID         uint16           `json:"id"`
	FieldCount uint16           `json:"field_count"`
	Fields     []FieldSpecifier `json:"fields"`
}

//options template record head
//0                   1                   2                   3
//0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|         Template ID (> 255)   |         Field Count           |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|      Scope Field Count        |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
type OptionTemplateRecord struct {
	ID              uint16           `json:"id"`
	FieldCount      uint16           `json:"field_count"`
	ScopeFieldCount uint16           `json:"scope_field_count"`
	Fields          []FieldSpecifier `json:"fields"`
}

//Field Specifier
//0                   1                   2                   3
//0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|E|  Information Element ident. |        Field Length           |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|                      Enterprise Number                        |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//if EnterpriseNo is not zero,then ID contains E
type FieldSpecifier struct {
	ID           uint16 `json:"id"`
	Length       uint16 `json:"length"`
	EnterpriseNo uint32 `json:"enterprise_no"`
}

type DataField struct {
	FieldID uint16      `json:"field_id"`
	Value   interface{} `json:"value"`
}

//data set,containing data records
//0                   1                   2                   3
//0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|   Set ID = Template ID        |          Length               |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|   Record 1 - Field Value 1    |   Record 1 - Field Value 2    |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|   Record 1 - Field Value 3    |             ...               |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|   Record 2 - Field Value 1    |   Record 2 - Field Value 2    |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|   Record 2 - Field Value 3    |             ...               |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|   Record 3 - Field Value 1    |   Record 3 - Field Value 2    |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|   Record 3 - Field Value 3    |             ...               |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|              ...              |      Padding (optional)       |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
type DataSet struct {
	Header     SetHeader   `json:"header"`
	DataFields []DataField `json:"data_fields"`
	padding    int         //byte count
}
