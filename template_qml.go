package main

import "fmt"
import "os"
import "path"
import "strings"
import "text/template"

var __IFC_TEMPLATE_INIT_QML = `/*This file is auto generate by pkg.linuxdeepin.com/dbus-generator. Don't edit it*/
#include <QtDBus>
QVariant unmarsh(const QVariant&);
QVariant marsh(QDBusArgument target, const QVariant& arg, const QString& sig);
`

var __IFC_TEMPLATE_QML = `
#ifndef __{{ExportName}}_H__
#define __{{ExportName}}_H__

class {{ExportName}}Proxyer: public QDBusAbstractInterface
{
    Q_OBJECT
    void _hookSignals() {
        connection().connect(service(), path(), "org.freedesktop.DBus.Properties", "PropertiesChanged", parent(), SLOT(_propertiesChanged(QDBusMessage)));
    {{range .Signals}}
        connection().connect( service(), path(), interface(), "{{.Name}}", this, SLOT(_handleSignal(QDBusMessage)));{{end}}
    }

    Q_SLOT void _handleSignal( const QDBusMessage &msg ) {
	if ( msg.type() == QDBusMessage::SignalMessage ) {
		if ( msg.interface() != interface() )
			    return;

		if ( msg.path() != path() )
			    return;

		QVariantList args;
		foreach( QVariant arg, msg.arguments() ) {
			args << unmarsh(arg);
		}

		if (0) { {{range .Signals}}
		} else if ( QByteArray("{{.Name}}") == msg.member() ) {
                     QMetaObject::invokeMethod(parent(), "{{Lower .Name}}" {{range $i, $e := .Args}},
						Q_ARG({{getQType $e.Type}},args[{{$i}}].value<{{getQType $e.Type}} >()){{end}});{{end}}
                } else {
		    qDebug() << "Could not forward" << msg.member();
		}
	}

    }

public:
    {{ExportName}}Proxyer(const QString &path, QObject* parent)
          :QDBusAbstractInterface("{{DestName}}", path, "{{IfcName}}", QDBusConnection::{{BusType}}Bus(), parent)
    {
	    if (!isValid()) {
		    qDebug() << "Create {{ExportName}} remote object failed : " << lastError().message();
	    }
	    _hookSignals();
    }

    ~{{ExportName}}Proxyer() {
        connection().disconnect(service(), path(), "org.freedesktop.DBus.Properties", "PropertiesChanged",
	    				 parent(), SLOT(_propertiesChanged(QDBusMessage)));
    }

    QVariant fetchProperty(const char* name) const {
	QDBusMessage msg = QDBusMessage::createMethodCall(service(), path(),
		QLatin1String("org.freedesktop.DBus.Properties"),
		QLatin1String("Get"));
	msg << interface() << QString::fromUtf8(name);
	QDBusMessage reply = connection().call(msg, QDBus::Block, timeout());
	if (reply.type() != QDBusMessage::ReplyMessage) {
	    qDebug () << QDBusError(reply) << "at " << service() << path() << interface() << name;
	    return QVariant();
	}
	if (reply.signature() != QLatin1String("v")) {
	    QString errmsg = QLatin1String("Invalid signature org.freedesktop.DBus.Propertyies in return from call to ");
	    qDebug () << QDBusError(QDBusError::InvalidSignature, errmsg.arg(reply.signature()));
	    return QVariant();
	}
	QVariant value = unmarsh(reply.arguments().at(0));
	//qDebug() << Q_FUNC_INFO << name << value;
	return value;
    }

{{range .Properties}}
    QVariant __get_{{.Name}}__() const { return fetchProperty("{{.Name}}"); }
    {{if PropWritable .}}void __set_{{.Name}}__(const QVariant &v) { setProperty("{{.Name}}", v); }{{end}}
{{end}}

};

class {{ExportName}} : public QObject
{
    Q_OBJECT
private:
    QString m_path;
    Q_SLOT void _propertiesChanged(const QDBusMessage &msg) {
	    QList<QVariant> arguments = msg.arguments();
	    if (3 != arguments.count())
	    	return;
	    QString interfaceName = msg.arguments().at(0).toString();
	    if (interfaceName != "{{IfcName}}")
		    return;

	    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
	    foreach(const QString &prop, changedProps.keys()) {
	    //qDebug() << Q_FUNC_INFO << prop;
		    if (0) { {{range .Properties}}
		    } else if (prop == "{{.Name}}") {
			    Q_EMIT __{{Lower .Name}}Changed__();{{end}}
		    }
	    }
    }
    void _rebuild()
    {
	  delete m_ifc;
          m_ifc = new {{ExportName}}Proxyer(m_path, this);
    }

public:
    Q_PROPERTY(QString path READ path WRITE setPath NOTIFY pathChanged)
    const QString path() {
	    return m_path;
    }
    void setPath(const QString& path) {
	m_path = path;
	_rebuild();

    }
    Q_SIGNAL void pathChanged(QString);

    {{ExportName}}(QObject *parent=0) : QObject(parent), m_ifc(new {{ExportName}}Proxyer("{{Ifc2Obj IfcName}}", this))
    {

    }
    {{range .Properties}}
    Q_PROPERTY({{getQType .Type}} {{Lower .Name}} READ __get_{{.Name}}__ {{if PropWritable .}}WRITE __set_{{.Name}}__{{end}} NOTIFY __{{Lower .Name}}Changed__){{end}}

    //Property read methods{{range .Properties}}
    {{getQType .Type}} __get_{{.Name}}__() const { return m_ifc->__get_{{.Name}}__().value<{{getQType .Type}} >(); }{{end}}
    //Property set methods :TODO check access{{range .Properties}}{{if PropWritable .}}
    void __set_{{.Name}}__(const {{getQType .Type}} &v) {
            if ( v != __get_{{.Name}}__() ) {
	        QVariant marshedValue = marsh(QDBusArgument(), v, "{{.Type}}");
	        m_ifc->__set_{{.Name}}__(marshedValue);
	        Q_EMIT __{{Lower .Name}}Changed__();
	    }
    }{{end}}{{end}}

public Q_SLOTS:{{range .Methods}}
    QVariant {{.Name}}({{range $i, $e := GetOuts .Args}}{{if ne $i 0}}, {{end}}const QVariant &{{.Name}}{{end}}) {
	    QList<QVariant> argumentList;
	    {{range GetOuts .Args}}
	    argumentList << marsh(QDBusArgument(), {{.Name}}, "{{.Type}}");{{end}}

	    QDBusPendingReply<> call = m_ifc->asyncCallWithArgumentList(QLatin1String("{{.Name}}"), argumentList);
	    call.waitForFinished();
	    if (call.isValid()) {
		    QList<QVariant> args = call.reply().arguments();
		    switch (args.size()) {
			    case 0: return QVariant();
			    case 1: {
				    return unmarsh(args[0]);
			    }
		    default:
			    {
				    for (int i=0; i<args.size(); i++) {
					    args[i] = unmarsh(args[i]);
				    }
				    return args;
			    }
		    }
	    } else {
		    qDebug() << "Error:" << call.error().message();
		    return QVariant();
	    }
    }
{{end}}

Q_SIGNALS:
//Property changed notify signal{{range .Properties}}
    void __{{Lower .Name}}Changed__();{{end}}

//DBus Interface's signal{{range .Signals}}
    void {{Lower .Name}}({{range $i, $e := .Args}}{{if ne $i 0}},{{end}}{{getQType $e.Type}} {{$e.Name}}{{end}});{{end}}
private:
    {{ExportName}}Proxyer *m_ifc;
};

#endif
`

var __GLOBAL_TEMPLATE_QML = `
#ifndef __DBUS_H__
#define __DBUS_H__

{{range GetModules}}
#include "{{.}}.h"{{end}}
#include <QQmlExtensionPlugin>
#include <qqml.h>

struct TypeMapping {
    QString signature;
    int metatype;
};

static TypeMapping __types [] =
{
    { QLatin1String("iiii"), qMetaTypeId<QRect>() },
    {{ range $key, $value := GetQtSignaturesType }}
        { QLatin1String("{{$key}}"), qDBusRegisterMetaType<{{$value}} >()},{{end}}
   { QString(), -1 }
};

class DBusPlugin: public QQmlExtensionPlugin
{
    Q_OBJECT
	Q_PLUGIN_METADATA(IID "org.qt-project.Qt.QQmlExtensionInterface")

    public:

        void registerTypes(const char* uri) { {{range .Interfaces}}
             qmlRegisterType<{{.ObjectName}}>(uri, 1, 0, "{{.ObjectName}}");{{end}}


	}
};

` + _templateMarshUnMarsh + `
#endif
`

var __PROJECT_TEMPL_QML = `
TEMPLATE=lib
CONFIG += plugin
QT += qml dbus

OBJECTS_DIRS = tmp
MOC_DIR = tmp
DESTDIR = {{PkgPath}}

HEADERS += plugin.h {{range GetModules}}{{.}}.h {{end}}

IMPORT_VERSION = 1.0
TARGET = $$qtLibraryTarget({{PkgName}})
TARGETPATH = $$[QT_INSTALL_QML]/{{PkgPath}}
target.path = $$TARGETPATH

INSTALLS += target

QMAKE_EXTRA_TARGETS += check
check.commands = qmlscene -I . test.qml
check.depends = first

QMAKE_POST_LINK += qmlplugindump -noinstantiate {{ModuleName}} 1.0 . > plugins.qmltypes

typesout.files = plugins.qmltypes
typesout.CONFIG = no_check_exist
typesout.path = $$TARGETPATH
INSTALLS += typesout

qmldir.files = {{PkgPath}}/qmldir
qmldir.path = $$TARGETPATH
INSTALLS += qmldir

`

var __TEST_QML = `
import {{ModuleName}} 1.0
import QtQuick 2.0
import QtQuick.Controls 1.0

Item { {{range .Interfaces}}
    {{.ObjectName}} {
       id: "{{Lower .ObjectName}}ID"
       // path: "{{Ifc2Obj .Interface}}"
    } {{end}}
    width: 400; height: 400
    TabView {
	    anchors.fill  : parent
	    {{range .Interfaces}}
	    Tab {   {{$ifc := GetInterfaceInfo .}} {{$objName := Lower .ObjectName }}
		    title: "{{.ObjectName}}"
		    Column {
			    {{range $ifc.Properties}}
			    Row {
				    Label {
					    text: "{{.Name}}:"
				    }
				    Text {
					    text: JSON.stringify({{$objName}}ID.{{Lower .Name}})
				    }
			    }{{end}}
		    }
	    }
	    {{end}}
    }
}
`

func renderQMLProject() {

        moduleName := "DBus"
	modulePath := "DBus"
	for _, f := range strings.Split(INFOS.Config.DestName, ".") {
		moduleName += "." + upper(f)
		modulePath += "/" + upper(f)
	}

	writer, err := os.Create(path.Join(INFOS.Config.OutputDir, "tt.pro"))
	if err != nil {
		panic(err)
	}
	template.Must(template.New("main").Funcs(template.FuncMap{
		"BusType": func() string { return INFOS.Config.BusType },
		"PkgName": func() string { return INFOS.Config.PkgName },
		"PkgPath": func() string { return modulePath },
		"ModuleName": func() string { return moduleName },
		"GetModules": func() map[string]string {
			r := make(map[string]string)
			for _, ifc := range INFOS.Interfaces {
				r[ifc.OutFile] = ifc.OutFile
			}
			return r
		},
	}).Parse(__PROJECT_TEMPL_QML)).Execute(writer, INFOS)
	writer.Close()
}

func testQML() {
	pkgName := INFOS.Config.PkgName
	if pkgName == "" {
		pkgName = getQMLPkgName("DBus." + INFOS.Config.DestName)
	}

	moduelPath := "DBus"
	moduleName := "DBus"
	for _, f := range strings.Split(INFOS.Config.DestName, ".") {
		moduleName += "." + upper(f)
		moduelPath += "/" + upper(f)
	}


	os.MkdirAll(INFOS.Config.OutputDir+"/" + moduelPath, 0755)

	qmldir, err := os.Create(path.Join(INFOS.Config.OutputDir, moduelPath, "qmldir"))
	if err != nil {
		panic(err)
	}


	qmldir.WriteString("module " + moduleName + "\n")
	qmldir.WriteString("plugin " + INFOS.Config.PkgName + "\n")
	qmldir.WriteString("typeinfo plugins.qmltypes\n" )
	qmldir.Close()

	writer, err := os.Create(path.Join(INFOS.Config.OutputDir, "test.qml"))
	if err != nil {
		panic(err)
	}
	template.Must(template.New("qmltest").Funcs(template.FuncMap{
		"Lower":            lower,
		"GetInterfaceInfo": GetInterfaceInfo,
		"BusType":          func() string { return INFOS.Config.BusType },
		"PkgName":          func() string { return pkgName },
		"Ifc2Obj":          ifc2obj,
		"ModuleName":       func() string { return moduleName },
		"GetModules": func() map[string]string {
			r := make(map[string]string)
			for _, ifc := range INFOS.Interfaces {
				r[ifc.OutFile] = ifc.OutFile
			}
			return r
		},
	}).Parse(__TEST_QML)).Execute(writer, INFOS)

}
func qtPropertyFilter(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "QMap") {
		return "QVariantMap"
	} else if strings.HasPrefix(s, "QList") {
		return "QVariantList"
	} else if strings.HasPrefix(s, "QValueList") {
		return "QVariantValueList"
	}
	return s
}

func getQtSignaturesType() (sigs map[string]string) {
	sigs = make(map[string]string)
	var store func(string)
	store = func(sig string) {
		if v, ok := _sig2QType[sig]; ok {
			sigs[sig] = v
		} else if sig == "as" {
			sigs[sig] = "QStringList"
		} else if sig == "ay" {
		        sigs[sig] = "QString"
		} else if sig == "so" {
			fmt.Println("Warning: `so` isn't supported")
			sigs[sig] = "QStringList"
		} else if sig[0] == 'a' {
			if sig[1] == '(' {
				sigs[sig] = "QVariantList"
			} else if sig[1] == '{' {
				if len(sig) < 5 { // a{xx} has at least five characters
					return
				}
				store(sig[2:3])
				store(sig[3:strings.LastIndex(sig, "}")])
				sigs[sig] = "QVariantMap"
			} else {
				store(sig[1:])
				sigs[sig] = "QList< " + sigs[sig[1:]] + " >"
			}
		} else if sig[0] == '(' {
			sigs[sig] = "QVariantList"
		} else {
			panic(fmt.Sprintf("parse signature failed:%q\n", sig))
		}
	}
	for _, ifc := range INFOS.Interfaces {
		info := GetInterfaceInfo(ifc)
		for _, m := range info.Methods {
			for _, a := range m.Args {
				store(a.Type)
			}
		}
		for _, p := range info.Properties {
			store(p.Type)
		}
		for _, s := range info.Signals {
			for _, ss := range s.Args {
				store(ss.Type)
			}
		}
	}
	return sigs
}

var _templateMarshUnMarsh = `
inline

int getTypeId(const QString& sig) {
    TypeMapping *mapping = __types;
     while( mapping->signature != QString() ) {
         if ( mapping->signature == sig ) return mapping->metatype;
	 mapping++;
     }
     return -1;
}

inline
QVariant qstring2dbus(QString value, char sig) {
    switch (sig) {
        case 'y':
            return QVariant::fromValue(uchar(value[0].toLatin1()));
        case 'n':
            return QVariant::fromValue(value.toShort());
        case 'q':
            return QVariant::fromValue(value.toUShort());
        case 'i':
            return QVariant::fromValue(value.toInt());
        case 'u':
            return QVariant::fromValue(value.toUInt());
        case 'x':
            return QVariant::fromValue(value.toLongLong());
        case 't':
            return QVariant::fromValue(value.toULongLong());
        case 'd':
            return QVariant::fromValue(value.toDouble());
        case 's':
            return QVariant::fromValue(value);
        case 'o':
            return QVariant::fromValue(QDBusObjectPath(value));
        default:
            qDebug() << "Dict entry key should be an basic dbus type not an " << sig;
            return QVariant();
    }
}

QList<QString> splitStructureSignature(const QString& sig) {
    if (sig.size() < 3 || sig[0] != '(' || sig[sig.size()-1] != ')') {
        return QList<QString>();
    }

    QList<QString> sigs;

    QString tmp = sig.mid(1, sig.size()-2);
    while (tmp.size() != 0) {
        switch (tmp[0].toLatin1()) {
            case 'a':
                if (tmp.size() < 2) {
                    return QList<QString>();
                }
                if (tmp[1] == '{') {
                    int lastIndex = tmp.lastIndexOf('}') + 1;
                    if (lastIndex == 0) return QList<QString>();
                    sigs.append(tmp.mid(0, lastIndex));
                    tmp = tmp.mid(lastIndex);
                    break;
                } else if (tmp[1] == '(') {
                    int lastIndex = tmp.lastIndexOf(')') + 1;
                    if (lastIndex == 0) return QList<QString>();
                    sigs.append(tmp.mid(0, lastIndex));
                    tmp = tmp.mid(lastIndex);
                    break;
                } else {
                    sigs.append(tmp.mid(0, 2));
                    tmp = tmp.mid(2);
                    break;
                }
            case '(': {
                          int lastIndex = tmp.lastIndexOf(')') + 1;
                          if (lastIndex == 0) return QList<QString>();
                          sigs.append(tmp.mid(0, lastIndex));
                          tmp = tmp.mid(lastIndex);
                          break;
                      }
            case 'y': case 'b': case 'n': case 'q':
            case 'i': case 'u': case 'x': case 't':
            case 'd': case 's': case 'o': case 'g':
            case 'h': case 'v':
                sigs.append(QString(tmp[0]));
                tmp = tmp.mid(1, tmp.size() - 1);
                break;
            default:
                return QList<QString>();
        }
    }
    return sigs;
}

QVariant marsh(QDBusArgument target, const QVariant& arg, const QString& sig) {
    if (sig.size() == 0) {
        return QVariant::fromValue(target);
    }
    switch (sig[0].toLatin1()) {
        case 'y':
            target << qstring2dbus(arg.value<QString>(), 'y').value<uchar>();
            return QVariant::fromValue(target);
	    break;
        case 'b':
            target << arg.value<bool>();
            return QVariant::fromValue(target);
	    break;
        case 'n':
            target << arg.value<short>();
            return QVariant::fromValue(target);
	    break;
        case 'q':
            target << arg.value<ushort>();
            return QVariant::fromValue(target);
	    break;
        case 'i':
            target << arg.value<qint32>();
            return QVariant::fromValue(target);
	    break;
        case 'u':
            target << arg.value<quint32>();
            return QVariant::fromValue(target);
	    break;
        case 'x':
            target << arg.value<qlonglong>();
            return QVariant::fromValue(target);
	    break;
        case 't':
            target << arg.value<qulonglong>();
            return QVariant::fromValue(target);
	    break;
        case 'd':
            target << arg.value<double>();
            return QVariant::fromValue(target);
	    break;
        case 's':
            target << arg.value<QString>();
            return QVariant::fromValue(target);
	    break;
        case 'o':
            target << QDBusObjectPath(arg.value<QString>());
            return QVariant::fromValue(target);
	    break;
        case 'g':
            target << QDBusSignature(arg.value<QString>());
            return QVariant::fromValue(target);
	    break;
        case 'a':
            {
                if (sig.size() < 2) { return QVariant(); }
                char s = sig[1].toLatin1();
                if (s == '{') {
                    char key_sig = sig[2].toLatin1();
                    QString value_sig = sig.mid(3, sig.lastIndexOf('}') - 3);
                    target.beginMap(getTypeId(QString(key_sig)), getTypeId(value_sig));
                    //qDebug() << "BeginMap:" << key_sig << value_sig;
                    foreach(const QString& key, arg.value<QVariantMap>().keys()) {
                        //qDebug() << "KEY:" << key;
                        target.beginMapEntry();
                        //qDebug() <<"beginMapEntry";
                        marsh(target, qstring2dbus(key, key_sig), QString(key_sig));
                        marsh(target, arg.value<QVariantMap>()[key], value_sig);
                        //qDebug() <<"EndMapEntry";
                        target.endMapEntry();
                    }
                    //qDebug() << "EndMap";
                    target.endMap();
                    return QVariant::fromValue(target);
                } else {
                    QString next = sig.right(sig.size() - 1);
                    target.beginArray(getTypeId(next));
                    foreach(const QVariant& v, arg.value<QVariantList>()) {
                        marsh(target, v, next);
                    }
                    target.endArray();
                    return QVariant::fromValue(target);
                }
            }
	    break;
        case '(':
            {
                QList<QString> sigs = splitStructureSignature(sig);
                QVariantList values = arg.value<QVariantList>();
                if (values.size() != sigs.size()) {
                    qDebug() << "structure (" << arg << ") didn't match signature :" << sigs;
                    return QVariant();
                }
                target.beginStructure();
                for (int i=0; i < sigs.size(); i++) {
                    marsh(target, values[i], sigs[i]);
                }
                target.endStructure();
                return QVariant::fromValue(target);
            }
	    break;
        default:
            qDebug() << "Panic didn't support marsh" << sig;
	    break;
    }
    return QVariant::fromValue(target);
}

inline
QVariant unmarshDBus(const QDBusArgument &argument)
{
    switch (argument.currentType()) {
    case QDBusArgument::BasicType: {
        QVariant v = argument.asVariant();
        if (v.userType() == qMetaTypeId<QDBusObjectPath>())
            return v.value<QDBusObjectPath>().path();
        else if (v.userType() == qMetaTypeId<QDBusSignature>())
            return v.value<QDBusSignature>().signature();
        else
            return v;
    }
    case QDBusArgument::VariantType: {

        QVariant v = argument.asVariant().value<QDBusVariant>().variant();
        if (v.userType() == qMetaTypeId<QDBusArgument>())
            return unmarshDBus(v.value<QDBusArgument>());
	else if (v.userType() == qMetaTypeId<QDBusVariant>())
            return v.value<QDBusVariant>().variant();
        else
            return v;
    }
    case QDBusArgument::ArrayType: {
        QVariantList list;
        argument.beginArray();
        while (!argument.atEnd())
            list.append(unmarshDBus(argument));
        argument.endArray();
        return list;
    }
    case QDBusArgument::StructureType: {
        QVariantList list;
        argument.beginStructure();
        while (!argument.atEnd())
            list.append(unmarshDBus(argument));
        argument.endStructure();
        return QVariant::fromValue(list);
    }
    case QDBusArgument::MapType: {
        QVariantMap map;
        argument.beginMap();
        while (!argument.atEnd()) {
            argument.beginMapEntry();
            QVariant key = unmarshDBus(argument);
            QVariant value = unmarshDBus(argument);
            map.insert(key.toString(), value);
            argument.endMapEntry();
        }
        argument.endMap();
        return map;
    }
    default:
        return QVariant();
        break;
    }
}

QVariant unmarsh(const QVariant& v) {
    if (v.userType() == qMetaTypeId<QDBusObjectPath>()) {
        return QVariant::fromValue(v.value<QDBusObjectPath>().path());
    } else if (v.userType() == qMetaTypeId<QDBusArgument>()) {
        return unmarsh(unmarshDBus(v.value<QDBusArgument>()));
    } else if (v.userType() == qMetaTypeId<QDBusVariant>()) {
        return unmarsh( v.value<QDBusVariant>().variant());
    } else if (v.userType() == qMetaTypeId<QByteArray>()) {
        return QVariant::fromValue(QString::fromLatin1(v.value<QByteArray>()));
    }
    return v;
}
`
