package main

type Cursor struct {
	row int
	col int
}

func (c *Cursor) Init() {
	c.col = 0
	c.row = 0
}

func (c *Cursor) moveRow(num int) {
	c.row += num
}

func (c *Cursor) moveCol(num int) {
	c.col += num
}
