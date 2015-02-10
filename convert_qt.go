package main

import "strings"

var _sig2QType = map[byte]string{
	'y': "uchar",
	'b': "bool",
	'n': "short",
	'q': "ushort",
	'i': "int",
	'u': "uint",
	'h': "quint32",
	'x': "qlonglong",
	't': "qulonglong",
	'd': "double",
	's': "QString",
	'g': "QDBusSignature",
	'o': "QDBusObjectPath",
	'v': "QDBusVariant",
}

var _convertQDBus = map[string]string{
	"o": "QVariant::fromValue(QDBusObjectPath({{.Name}}.value<QString>()))",
}

func normalizeQDBus(v string) (r string) {
	return //TODO:
	if result, ok := _convertQDBus[v]; ok {
		r = result
		/*return "huhu" //result*/
	}
	return
}

func getQType(sig string) string {

        if sig[0] == 'o' {
		return "QString"
	}

        if sig == "ay" {
	       return "QString"
	}

	if qtype, ok := _sig2QType[sig[0]]; ok {
		return qtype
	}
	switch sig[0] {
	case 'a':
		if sig[1] == '{' {
			i := strings.LastIndex(sig, "}")
			r := "QMap<"
			r += getQType(string(sig[2])) + ", "
			r += getQType(sig[3:i])
			r += " >"
			if r == "QMap<QString, QVariant >" {
			    return "QVariantMap"
			} else if r == "QMap<QString, QVariantMap >" {
                             return "QVariantMap"
			} else if r == "QMap<QString, QDBusVariant >" {
			    return "QVariantMap"
			} else {
				r := "QMap<"
				r += getQType(string(sig[2])) + ", "
				r += getQType(sig[3:i])
				r += " >"
				return r
			}
		} else {
			r := "QList<"
			r += getQType(sig[1:])
			r += " >"
			if r == "QList<QString >" {
			    return "QStringList"
			} else if r == "QList<QVariant >" {
			    return "QVariantList"
			} else if r == "QList<QDBusVariant >" {
			    return "QVariantList"
			} else {
				return r
			}
		}
	case '(':
		return "QVariant"
	}
	panic("Unknow Type" + sig)
}
