package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

func ToPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == '_' })
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}

func enumName(interfaceName, enumNameWithDot string) string {

	newEnumName := enumNameWithDot
	if !strings.Contains(enumNameWithDot, ".") {
		newEnumName = interfaceName + newEnumName
	}

	parts := strings.FieldsFunc(newEnumName, func(r rune) bool { return r == '.' })

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	newEnumName = strings.Join(parts, "")
	return newEnumName + "_enum"
}

type Protocol struct {
	XMLName    xml.Name    `xml:"protocol"`
	Name       string      `xml:"name,attr"`
	Interfaces []Interface `xml:"interface"`
	Copyright  string      `xml:"copyright"`
}

type InterfaceAttr struct {
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
}

type Interface struct {
	InterfaceAttr
	Description any              `xml:"description"`
	Requests    []EventOrRequest `xml:"request"`
	Events      []EventOrRequest `xml:"event"`
	Enums       []Enum           `xml:"enum"`
}

type Description struct {
	Summary string `xml:"summary,attr"`
}

type Enum struct {
	Name    string  `xml:"name,attr"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Name    string `xml:"name,attr"`
	Value   string `xml:"value,attr"`
	Summary string `xml:"summary,attr"`
}

type EventOrRequestAttr struct {
	Name  string  `xml:"name,attr"`
	Since *string `xml:"since,attr,omitempty"`
}

type EventOrRequest struct {
	EventOrRequestAttr
	Args        []Arg        `xml:"arg"`
	Description *Description `xml:"description"`
}

type Arg interface {
	ArgKind() string
	Name() string
}

type ArgCommon struct {
	ArgName string
}

func (c ArgCommon) Name() string { return c.ArgName }

type ArgNewID struct {
	ArgCommon
	Interface *string
}

func (*ArgNewID) ArgKind() string { return "new_id" }

type ArgObject struct {
	ArgCommon
	Interface *string
	AllowNull *bool
}

func (*ArgObject) ArgKind() string { return "object" }

type ArgUint struct {
	ArgCommon
	Enum *string
}

func (*ArgUint) ArgKind() string { return "uint" }

type ArgString struct {
	ArgCommon
	AllowNull *bool
}

func (*ArgString) ArgKind() string { return "string" }

type ArgInt struct {
	ArgCommon
}

func (*ArgInt) ArgKind() string { return "int" }

type ArgFd struct {
	ArgCommon
}

func (*ArgFd) ArgKind() string { return "fd" }

type ArgFixed struct {
	ArgCommon
}

func (*ArgFixed) ArgKind() string { return "fixed" }

type ArgArray struct {
	ArgCommon
}

func (*ArgArray) ArgKind() string { return "array" }

func sanitizedArgName(arg Arg) string {
	switch arg.Name() {
	case "interface":
		return "interface_"
	case "class":
		return "class_"
	case "make":
		return "make_"
	default:
		return arg.Name()
	}
}

type argXML struct {
	Name      string  `xml:"name,attr"`
	Type      string  `xml:"type,attr"`
	Interface *string `xml:"interface,attr"`
	AllowNull *bool   `xml:"allow-null,attr"`
	Enum      *string `xml:"enum,attr"`
}

type eventOrRequestXML struct {
	Name        string       `xml:"name,attr"`
	Since       *string      `xml:"since,attr,omitempty"`
	Description *Description `xml:"description"`
	Args        []argXML     `xml:"arg"`
}

type interfaceXML struct {
	InterfaceAttr
	Description any                 `xml:"description"`
	Requests    []eventOrRequestXML `xml:"request"`
	Events      []eventOrRequestXML `xml:"event"`
	Enums       []Enum              `xml:"enum"`
}

type protocolXML struct {
	XMLName    xml.Name       `xml:"protocol"`
	Name       string         `xml:"name,attr"`
	Interfaces []interfaceXML `xml:"interface"`
	Copyright  string         `xml:"copyright"`
}

func convertEventOrRequestXML(x eventOrRequestXML) (EventOrRequest, error) {
	out := EventOrRequest{
		EventOrRequestAttr: EventOrRequestAttr{
			Name:  x.Name,
			Since: x.Since,
		},
		Description: x.Description,
		Args:        make([]Arg, 0, len(x.Args)),
	}

	for _, ax := range x.Args {
		var a Arg
		switch ax.Type {
		case "new_id":

			var iface *string
			if ax.Interface != nil {
				ifaceStr := ToPascalCase(*ax.Interface)
				iface = &ifaceStr
			}

			a = &ArgNewID{
				ArgCommon: ArgCommon{ArgName: ax.Name},
				Interface: iface,
			}
		case "object":
			var iface *string
			if ax.Interface != nil {
				ifaceStr := ToPascalCase(*ax.Interface)
				iface = &ifaceStr
			}
			a = &ArgObject{
				ArgCommon: ArgCommon{ArgName: ax.Name},
				Interface: iface,
				AllowNull: ax.AllowNull,
			}
		case "uint":
			var enum *string
			if ax.Enum != nil {
				enumStr := ToPascalCase(*ax.Enum)
				enum = &enumStr
			}
			a = &ArgUint{
				ArgCommon: ArgCommon{ArgName: ax.Name},
				Enum:      enum,
			}
		case "string":
			a = &ArgString{
				ArgCommon: ArgCommon{ArgName: ax.Name},
				AllowNull: ax.AllowNull,
			}
		case "int":
			a = &ArgInt{ArgCommon: ArgCommon{ArgName: ax.Name}}
		case "fd":
			a = &ArgFd{ArgCommon: ArgCommon{ArgName: ax.Name}}
		case "fixed":
			a = &ArgFixed{ArgCommon: ArgCommon{ArgName: ax.Name}}
		case "array":
			a = &ArgArray{ArgCommon: ArgCommon{ArgName: ax.Name}}
		default:
			return EventOrRequest{}, fmt.Errorf("unknown arg type: %q", ax.Type)
		}
		out.Args = append(out.Args, a)
	}

	return out, nil
}

func UnmarshalProtocolXML(data []byte) (*Protocol, error) {
	var tmp protocolXML
	if err := xml.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}

	out := &Protocol{
		XMLName:   tmp.XMLName,
		Name:      tmp.Name,
		Copyright: tmp.Copyright,
	}
	out.Interfaces = make([]Interface, len(tmp.Interfaces))

	for i, ix := range tmp.Interfaces {

		out_enums := make([]Enum, len(ix.Enums))
		for j, ex := range ix.Enums {
			out_enums[j] = Enum{
				Name:    ToPascalCase(ex.Name),
				Entries: ex.Entries,
			}
		}

		iface := Interface{
			InterfaceAttr: InterfaceAttr{
				Name:    ToPascalCase(ix.Name),
				Version: ix.Version,
			},
			Description: ix.Description,
			Enums:       out_enums,
			Requests:    make([]EventOrRequest, len(ix.Requests)),
			Events:      make([]EventOrRequest, len(ix.Events)),
		}

		for j, rx := range ix.Requests {
			ev, err := convertEventOrRequestXML(rx)
			if err != nil {
				return nil, fmt.Errorf("interface %q request %d (%s): %w", ix.Name, j, rx.Name, err)
			}
			iface.Requests[j] = ev
		}

		for j, ex := range ix.Events {
			ev, err := convertEventOrRequestXML(ex)
			if err != nil {
				return nil, fmt.Errorf("interface %q event %d (%s): %w", ix.Name, j, ex.Name, err)
			}
			iface.Events[j] = ev
		}

		out.Interfaces[i] = iface
	}

	return out, nil
}

func generateGoType(interfaceName string, a Arg, event bool) string {
	name := sanitizedArgName(a)

	switch v := a.(type) {
	case *ArgNewID:
		if v.Interface != nil {
			return fmt.Sprintf("%s ObjectID[%s]", name, *v.Interface)
		}
		return fmt.Sprintf("%sInterface string, %sVersion uint32, %sID AnyObjectID", name, name, name)

	case *ArgObject:
		if v.Interface != nil {
			if v.AllowNull != nil && *v.AllowNull {
				return fmt.Sprintf("%s *ObjectID[%s]", name, *v.Interface)
			}
			return fmt.Sprintf("%s ObjectID[%s]", name, *v.Interface)
		}
		return fmt.Sprintf("%s AnyObjectID", name)

	case *ArgUint:
		if v.Enum == nil {
			return fmt.Sprintf("%s uint32", name)
		}
		return fmt.Sprintf("%s %s", name, enumName(interfaceName, *v.Enum))

	case *ArgString:
		return fmt.Sprintf("%s string", name)

	case *ArgInt:
		return fmt.Sprintf("%s int32", name)

	case *ArgFd:
		if event {
			return fmt.Sprintf("%s FileDescriptor", name)
		}
		return fmt.Sprintf("%s *FileDescriptor", name)

	case *ArgFixed:
		if event {
			// TODO is this right?
			return fmt.Sprintf("%s float32", name)
		}
		return fmt.Sprintf("%s Fixed", name)

	case *ArgArray:
		return fmt.Sprintf("%s []byte", name)
	default:
		panic(fmt.Errorf("unknown arg kind: %T", a))
	}
}
