package main

import (
	"fmt"
	"strings"
)

func genEvents(i Interface) string {
	if len(i.Events) == 0 {
		return ""
	}
	var out strings.Builder

	for idx, ev := range i.Events {
		var argSig []string
		for _, a := range ev.Args {
			argSig = append(argSig, generateGoType(i.Name, a, true))
		}

		out.WriteString("\n")
		out.WriteString(fmt.Sprintf("func %s_%s(", i.Name, ev.Name))
		out.WriteString("s Sender, ")
		if ev.Since != nil && *ev.Since != "" {
			out.WriteString("boundVersion uint32, ")
		}
		out.WriteString(fmt.Sprintf("eventObjectID ObjectID[%s]", i.Name))
		if len(argSig) > 0 {
			out.WriteString(", ")
			out.WriteString(strings.Join(argSig, ", "))
		}
		out.WriteString(") {\n")

		if ev.Since != nil && *ev.Since != "" {
			out.WriteString(fmt.Sprintf("    if boundVersion < %s {\n", *ev.Since))
			out.WriteString("        // Event not available in this version; skip\n")
			out.WriteString("        return\n")
			out.WriteString("    }\n")
		}

		needUint32, needInt32 := false, false
		for _, a := range ev.Args {
			switch a.(type) {
			case *ArgFixed:
				needInt32 = true
			case *ArgInt:
				needInt32 = true
			case *ArgUint, *ArgNewID, *ArgObject:
				needUint32 = true
			case *ArgString, *ArgArray:
				needUint32 = true
			case *ArgFd:
			}
		}

		out.WriteString("    data := make([]byte, 0)\n")
		if needUint32 {
			out.WriteString("    putUint32 := func(v uint32) { data = append(data, byte(v), byte(v>>8), byte(v>>16), byte(v>>24)) }\n")
		}
		if needInt32 {
			if needUint32 {
				out.WriteString("    putInt32 := func(v int32) { putUint32(uint32(v)) }\n")
			} else {
				out.WriteString("    putInt32 := func(v int32) { uv := uint32(v); data = append(data, byte(uv), byte(uv>>8), byte(uv>>16), byte(uv>>24)) }\n")
			}
		}
		out.WriteString("    var fileDescriptor *FileDescriptor\n")

		for _, a := range ev.Args {
			name := sanitizedArgName(a)
			switch v := a.(type) {
			case *ArgFixed:
				out.WriteString(fmt.Sprintf("    // fixed: 24.8\n    putInt32(int32(%s * 256.0))\n", name))
			case *ArgInt:
				out.WriteString(fmt.Sprintf("    putInt32(int32(%s))\n", name))
			case *ArgUint:
				out.WriteString(fmt.Sprintf("    putUint32(uint32(%s))\n", name))
			case *ArgNewID:
				out.WriteString(fmt.Sprintf("    putUint32(uint32(%s))\n", name))
			case *ArgObject:
				if v.Interface != nil && v.AllowNull != nil && *v.AllowNull {
					tmp := "__tmp_" + name
					out.WriteString(fmt.Sprintf("    var %s uint32\n", tmp))
					out.WriteString(fmt.Sprintf("    if %s != nil { %s = uint32(*%s) }\n", name, tmp, name))
					out.WriteString(fmt.Sprintf("    putUint32(%s)\n", tmp))
				} else {
					out.WriteString(fmt.Sprintf("    putUint32(uint32(%s))\n", name))
				}
			case *ArgFd:
				out.WriteString(fmt.Sprintf("    fileDescriptor = &%s\n", name))
			case *ArgString:
				out.WriteString(fmt.Sprintf(
					"    {\n"+
						"        b := []byte(%s)\n"+
						"        total := len(b) + 1 // include null terminator\n"+
						"        putUint32(uint32(total))\n"+
						"        data = append(data, b...)\n"+
						"        data = append(data, 0)\n"+
						"        if pad := (4 - (total %% 4)) %% 4; pad != 0 {\n"+
						"            data = append(data, make([]byte, pad)...)\n"+
						"        }\n"+
						"    }\n", name))
			case *ArgArray:
				out.WriteString(fmt.Sprintf(
					"    {\n"+
						"        n := len(%s)\n"+
						"        putUint32(uint32(n))\n"+
						"        data = append(data, %s...)\n"+
						"        if pad := (4 - (n %% 4)) %% 4; pad != 0 {\n"+
						"            data = append(data, make([]byte, pad)...)\n"+
						"        }\n"+
						"    }\n", name, name))
			default:
				out.WriteString(fmt.Sprintf("    // TODO: unhandled arg kind: %T\n", a))
			}
		}
		out.WriteString("    obj := OutgoingEvent{\n")
		out.WriteString("        ObjectID:       AnyObjectID(eventObjectID),\n")
		out.WriteString(fmt.Sprintf("        Opcode:         %d,\n", idx))
		out.WriteString("        Data:           data,\n")
		out.WriteString("        FileDescriptor: fileDescriptor,\n")
		out.WriteString("    }\n")
		out.WriteString("    s.Send(obj)\n")

		out.WriteString("}\n")
	}

	return out.String()
}
