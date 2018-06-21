package minimalrouter

import "testing"

func TestAddDuplicate(t *testing.T) {
	r := New()

	r.Add("GET", "/test1", true)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	r.Add("GET", "/test1", true)
}

func TestAddDuplicatePathOnly(t *testing.T) {
	r := New()

	r.Add("GET", "/test1", true)
	r.Add("POST", "/test1", true)
}

func TestAddDuplicateVariable(t *testing.T) {
	r := New()

	r.Add("GET", "/test/:id", true)
	r.Add("GET", "/test/:id/log", true)
}

func TestAddDirs(t *testing.T) {
	r := New()

	r.Add("GET", "/", true)
	r.Add("GET", "/test/", true)
}

func TestAddMismatchingVariable(t *testing.T) {
	r := New()

	r.Add("GET", "/test/:id", true)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	r.Add("GET", "/test/:other", true)
}

func TestAddMismatchingVariableLonger(t *testing.T) {
	r := New()

	r.Add("GET", "/test/:id", true)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	r.Add("GET", "/test/:other/id", true)
}

func TestMatchStatic(t *testing.T) {
	r := New()

	r.Add("GET", "/test1", "t1")
	r.Add("GET", "/test2", "t2")

	if m, p := r.Match("GET", "/test1"); m != "t1" {
		t.Errorf("The request did not match")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test2"); m != "t2" {
		t.Errorf("The request did not match")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("POST", "/test1"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test1/x"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}
}

func TestMatchVariable(t *testing.T) {
	r := New()

	r.Add("GET", "/test/:id", "s")
	r.Add("GET", "/test/:id/long", "l")

	if m, p := r.Match("GET", "/test"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test/1"); m != "s" {
		t.Errorf("The request did not match")
	} else if len(p) != 1 {
		t.Errorf("The request returned %d params", len(p))
	} else if p["id"] != "1" {
		t.Errorf("The request returned incorrect params: %v", p)
	}

	if m, p := r.Match("POST", "/test"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("POST", "/test/1"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test/1/longer"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}
}

func TestMatchDirs(t *testing.T) {
	r := New()

	r.Add("GET", "/", "t1")
	r.Add("GET", "/test/", "t2")

	if m, p := r.Match("GET", "/"); m != "t1" {
		t.Errorf("The request did not match")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test1"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test"); m != "t2" {
		t.Errorf("The request did not match")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test/"); m != "t2" {
		t.Errorf("The request did not match")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}

	if m, p := r.Match("GET", "/test/x"); m != nil {
		t.Errorf("The request matched erronously")
	} else if len(p) > 0 {
		t.Errorf("The request returned %d params", len(p))
	}
}