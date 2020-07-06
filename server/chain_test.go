// Copyright 2017 Tim Oster. All rights reserved.
// Use of this source code is governed by the MIT license.
// More information can be found in the LICENSE file.

package server

import "testing"

func TestChainMaxLength(t *testing.T) {
	c := newChain(6)

	for i := 1; i <= 10; i++ {
		c.chain(&Packet{order: orderNumber(i)})
	}

	if n := c.start.packet.order; n != 5 {
		t.Errorf("Expected first order number to be 5 not %v", n)
	}

	if c.length != 6 {
		t.Errorf("Expected length to be 6 not %v", c.length)
	}
}

func TestChainPopConsecutive(t *testing.T) {
	c := newChain(10)

	c.chain(&Packet{order: 1})
	c.chain(&Packet{order: 2})
	c.chain(&Packet{order: 3})
	c.chain(&Packet{order: 5})
	c.chain(&Packet{order: 6})

	if c.popConsecutive() != nil {
		t.Error("Expected chain to contain no consecutive sequence starting from order number 0")
	}

	c.chain(&Packet{order: 0})
	p := c.popConsecutive()

	if p == nil {
		t.Error("Expected chain to contain consecutive sequence")
	} else {
		for e := p; e != nil; e = e.next {
			if e.next == nil && e.packet.order != 3 {
				t.Errorf("Expected end of consecutive sequence to be order number 3 not %v", e.packet.order)
			}
		}
	}

	if c.length != 2 {
		t.Errorf("Expected length to be 2 not %v after popping", c.length)
	}

	if c.start.packet.order != 5 {
		t.Errorf("Expected start link of chain to have order number 5 not %v", c.start.packet.order)
	}

	p = c.popConsecutive()

	if p != nil {
		t.Error("Expected chain to be missing order number 4 for a new consecutive sequence")
	}

	c.chain(&Packet{order: 4})
	p = c.popConsecutive()

	if p == nil && p.packet.order != 4 {
		t.Error("Expected new consecutive sequence starting with order number 4")
	}

	if c.length != 0 || c.start != nil {
		t.Error("Expected chain to be empty after popping twice")
	}

	if c.next != 7 {
		t.Error("Expected chain to be waiting for oder number 7 for a new sequence")
	}
}
