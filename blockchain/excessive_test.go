package blockchain

import (
	"testing"
)

func TestHeight(t *testing.T) {
	var last *mockNode

	for i := uint32(0); i < 10; i++ {
		if height(last) != i {
			t.Fail()
		}

		last = &mockNode{
			parent: last,
		}
	}
}

func TestExcessive(t *testing.T) {
	excessive := &mockNode{
		excessive: true,
	}

	last := excessive
	for i := uint32(0); i < 10; i++ {
		for acceptDepth := uint32(0); acceptDepth < 10; acceptDepth++ {
			ex := excessiveChain(last, acceptDepth)

			if i >= acceptDepth && ex {
				t.Fail()
			} else if i < acceptDepth && !ex {
				t.Fail()
			}
		}

		last = &mockNode{
			parent: last,
		}
	}
}

func TestMaxAcceptible(t *testing.T) {
	root := &mockNode{}
	excessive := &mockNode{
		parent:    root,
		excessive: true,
	}

	rounds := uint32(3)

	for acceptDepth := uint32(0); acceptDepth < rounds; acceptDepth++ {
		branch := excessive

		for j := uint32(0); j <= acceptDepth; j++ {
			last := excessive
			anotherExcessiveAdded := false

			for i := uint32(0); i <= acceptDepth; i++ {
				max := maxAcceptable(last, acceptDepth)

				var expected *mockNode
				if anotherExcessiveAdded {
					expected = branch
				} else {
					expected = last
				}

				if i >= acceptDepth && max != expected {
					t.Fatal("invalid state")
				} else if i < acceptDepth && max != root {
					t.Fatal("invalid state")
				}

				if i == j {
					anotherExcessiveAdded = true
					branch = last
					last = &mockNode{
						parent:    last,
						excessive: true,
					}
				} else {
					last = &mockNode{
						parent: last,
					}
				}
			}
		}
	}
}

func TestMainChain(t *testing.T) {
	chain := makeMockChain(1)

	// Add one block. Is it the head?
	a := &mockNode{}
	chain.addBlock(a)
	if chain.main != a {
		t.Fatalf("A")
	}

	// Add another block. Is it the new head?
	b := &mockNode{parent: a}
	chain.addBlock(b)
	if chain.main != b {
		t.Fatalf("B")
	}

	// Add an excessive block. It should not be the new head.
	c := &mockNode{parent: b, excessive: true}
	chain.addBlock(c)
	if chain.main != b {
		t.Fatalf("C")
	}

	// Add another block. Now it should be the new head.
	d := &mockNode{parent: c}
	chain.addBlock(d)
	if chain.main != d {
		t.Fatalf("D")
	}

	// Add another excessive block. It should not be the new head.
	e := &mockNode{parent: d, excessive: true}
	chain.addBlock(e)
	if chain.main != d {
		t.Fatalf("E")
	}

	// Add another excessive block. The previous excessive block
	// should be the new head.
	f := &mockNode{parent: e, excessive: true}
	chain.addBlock(f)
	if chain.main != e {
		t.Fatalf("F")
	}
}
