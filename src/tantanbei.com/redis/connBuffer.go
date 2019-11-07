package redis

//Indexed ByteBuffer
//warning: Maybe the size not equal to the len(realbuffer)
type ConnBuffer struct {
	realBuffer []byte
	size       int
	capacity   int
	index      int
}

//	[/////////////////////////|*************************************|.............................]
//	|                         |                                     |                             |
//	0        pre-read        index        usable buffer           size         free buffer       cap
