// +build ignore
package main

import "testing"

func TestQType(t *testing.T) {
	if getQType("u") != "uint" {
		t.Fatal(` "u" != "uint" ` + getQType("u"))
	}
	if getQType("ah") != "QList<quint32 >" {
		t.Fatal(` "ah" != "QList<quint32 >" ` + getQType("ah"))
	}
	if getQType("au") != "QList<uint >" {
		t.Fatal(` "au" != "QList<uint >" ` + getQType("au"))
	}
	if getQType("ao") != "QStringList" {
		t.Fatal(` "ao" != "QStringList" ` + getQType("ao"))
	}
	if getQType("as") != "QStringList" {
		t.Fatal(` "as" != "QStringList" ` + getQType("as"))
	}
	if getQType("av") != "QVariantList" {
		t.Fatal(` "av" != "QVariantList" ` + getQType("av"))
	}
	if getQType("a{ss}") != "QVariantMap" {
		t.Fatal(` "a{ss}" != "QVariantMap" ` + getQType("a{ss}"))
	}
}
