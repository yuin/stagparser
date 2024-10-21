package stagparser_test

import (
	"testing"

	. "github.com/yuin/stagparser"
)

type StructA struct {
	f1 string `t1:"abc=1,def=ghi,jkl='mno',pkr=[1, -100.009, aaa, bbb, -56],stu(vwx=ccc, zzz=ddd), a1"` // nolint
	f2 string `t1:"abd='\\r\\n\\''"`                                                                    // nolint
}

func TestExampleSuccess(t *testing.T) {
	s := &StructA{}
	result, err := ParseStruct(s, "t1")
	if err != nil {
		t.Fatalf("parse failed: %s", err.Error())
	}
	f1, ok := result["f1"]
	if !ok {
		t.Fatalf("field f1 should be parsed")
	}
	if len(f1) != 6 {
		t.Fatalf("field f1 should be parsed into 6 definitions but %d", len(f1))
	}
	abc := f1[0]
	if abc.Name() != "abc" {
		t.Fatalf("1st f1 defition should be 'abc'")
	}
	if len(abc.Attributes()) != 1 {
		t.Fatalf("'abc' should have one attribute")
	}
	if v, ok := abc.Attribute("abc"); !ok || v.(int64) != 1 {
		t.Fatalf("abc attribute should be 1(int64) but got %v(%T)", v, v)
	}

	def := f1[1]
	if def.Name() != "def" {
		t.Fatalf("2nd f1 defition should be 'def'")
	}
	if len(def.Attributes()) != 1 {
		t.Fatalf("'def' should have one attribute")
	}
	if v, ok := def.Attribute("def"); !ok || v.(string) != "ghi" {
		t.Fatalf("def attribute should be \"ghi\" but got %v(%T)", v, v)
	}

	jkl := f1[2]
	if jkl.Name() != "jkl" {
		t.Fatalf("3rd f1 definition should be 'jkl'")
	}
	if len(jkl.Attributes()) != 1 {
		t.Fatalf("'jkl' should have one attribute")
	}
	if v, ok := jkl.Attribute("jkl"); !ok || v.(string) != "mno" {
		t.Fatalf("jkl attribute should be \"mno\" but got %v(%T)", v, v)
	}

	pkr := f1[3]
	if pkr.Name() != "pkr" {
		t.Fatalf("4th f1 defition should be 'pkr' but got %s", pkr.Name())
	}
	if len(pkr.Attributes()) != 1 {
		t.Fatalf("'pkr' should have one attribute")
	}
	if v, ok := pkr.Attribute("pkr"); !ok {
		t.Fatalf("pkr attribute should be exists")
	} else {
		v2, ok2 := v.([]interface{})
		if !ok2 {
			t.Fatalf("pkr attribute should be []interface{} but got %T", v)
		}
		if len(v2) != 5 {
			t.Fatalf("pkr attribute should be [5]interface{} but [%d]interface{}", len(v2))
		}
		if v2[0].(int64) != 1 {
			t.Fatalf("pkr attribute[0] should be 1(int64) but got %v(%T)", v, v)
		}
		if v2[1].(float64) != -100.009 {
			t.Fatalf("pkr attribute[1] should be 100.009(float64) but got %v(%T)", v, v)
		}
		if v2[2].(string) != "aaa" {
			t.Fatalf("pkr attribute[2] should be \"aaa\" but got %v(%T)", v, v)
		}
		if v2[3].(string) != "bbb" {
			t.Fatalf("pkr attribute[3] should be \"bbb\" but got %v(%T)", v, v)
		}
		if v2[4].(int64) != -56 {
			t.Fatalf("pkr attribute[4] should be 56(int64) but got %v(%T)", v, v)
		}
	}

	stu := f1[4]
	if stu.Name() != "stu" {
		t.Fatalf("5th f1 definition should be 'stu'")
	}
	if len(stu.Attributes()) != 2 {
		t.Fatalf("'stu' should have two attributes")
	}
	if v, ok := stu.Attribute("vwx"); !ok || v.(string) != "ccc" {
		t.Fatalf("vwx attribute should be \"ccc\" but got %v(%T)", v, v)
	}
	if v, ok := stu.Attribute("zzz"); !ok || v.(string) != "ddd" {
		t.Fatalf("zzz attribute should be \"ddd\" but got %v(%T)", v, v)
	}

	a1 := f1[5]
	if a1.Name() != "a1" {
		t.Fatalf("6th f1 definition should be 'a1'")
	}
	if len(a1.Attributes()) != 0 {
		t.Fatalf("'a1' should not have attributes")
	}

	f2, ok := result["f2"]
	if !ok {
		t.Fatalf("field f2 should be parsed")
	}
	if len(f2) != 1 {
		t.Fatalf("field f2 should be parsed into 1 definitions but %d", len(f2))
	}

	abd := f2[0]
	if abd.Name() != "abd" {
		t.Fatalf("1st f2 definition should be 'abd'")
	}
	if len(abd.Attributes()) != 1 {
		t.Fatalf("'abd' should have one attribute")
	}
	if v, ok := abd.Attribute("abd"); !ok || v.(string) != "\r\n'" {
		t.Fatalf("zzz attribute should be \"\\r\\n'\" but got %v(%T)", v, v)
	}
}
