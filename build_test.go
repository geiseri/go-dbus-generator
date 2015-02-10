package main

import "fmt"
import "os"
import "os/exec"
import C "launchpad.net/gocheck"

const (
	DBusInJson = "testdata/output/dbus.in.json"
)

func init() {
	_, err := os.Stat(DBusInJson)
	if err != nil {
		fmt.Println("There hasn't test data, try generating it")
		exec.Command("bash", "-c", fmt.Sprintf("cd testdata && go build && ./testdata -output output")).Run()
	}
}

func (*testWrap) TestGolang(c *C.C) {
	output := "tmp_golang__"
	infos, err := LoadInfos(DBusInJson, output, "Golang")
	c.Check(err, C.Equals, nil)

	geneateInit(infos)
	generateMain(infos)
	renderedEnd(infos)

	o, err := exec.Command("bash", "-c", fmt.Sprintf("cd %s && ls && go build", output)).CombinedOutput()
	if err != nil {
		c.Fatal(string(o))
	}
	c.Check(err, C.Equals, nil)

	err = exec.Command("rm", "-rf", "tmp_golang__").Run()
	c.Check(err, C.Equals, nil)
}

func (*testWrap) TestQML(c *C.C) {
	output := "tmp_qml__"
	infos, err := LoadInfos(DBusInJson, output, "Qml")
	c.Check(err, C.Equals, nil)

	geneateInit(infos)
	generateMain(infos)
	renderedEnd(infos)

	_, err = exec.Command("bash", "-c", fmt.Sprintf("cd %s && ls && qmake && make", output)).CombinedOutput()
	c.Check(err, C.Equals, nil)

	exec.Command("rm", "-rf", output).Run()
	c.Check(err, C.Equals, nil)
}
