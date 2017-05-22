package utils

import "testing"

func TestPackageAndFunc(t *testing.T) {
    testCases := []struct{
        Offset   int
        Package  string
        Function string
    }{
        {
            Offset:  1,
            Package: "urlcheck/utils",
            Function: "packageAndFunc",
        },
        {
            Offset:  2,
            Package: "urlcheck/utils",
            Function: "TestPackageAndFunc",
        },
    }

    for _, c := range testCases {
        pkg, fun := packageAndFunc(c.Offset)
        if pkg != c.Package && fun != c.Function {
            t.Errorf("Package and function do not match at offset %d.  Expected package: %s and function: %s but got package: %s and function %s", c.Offset, c.Package, c.Function, pkg, fun)
        }
    }
}

func TestPackageName(t *testing.T) {
    pkg := PackageName()
    if pkg != "urlcheck/utils" {
        t.Errorf("PackageName() expected urlcheck/utils but returned %s", pkg)
    }
}

func TestFunctionName(t *testing.T) {
    pkg := FunctionName()
    if pkg != "TestFunctionName" {
        t.Errorf("FunctionName() expected TestFunctionName but returned %s", pkg)
    }
}

func TestGetFields(t *testing.T) {
    expected := map[string]interface{}{
        "key": "value",
        "package": "urlcheck/utils",
        "function": "TestGetFields",
    }

    fields := getFields(map[string]interface{}{ "key": "value"})

    match := (len(fields) == len(expected))

    // reflect.DeepEqual always shows these maps to be different.  Comparing manually instead.
    for key, _ := range fields {
        exp, ok := expected[key]; if !ok {
            match = false
            break
        }

        fld, ok := fields[key]; if !ok {
            match = false
            break
        }

        if fld != exp {
            match = false
            break
        }
    }

    if !match {
        t.Errorf("getFields() did not return the expected fields.")
        t.Errorf("%v\n%v", expected, fields)
    }
}
