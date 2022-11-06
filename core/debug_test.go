package core_test

import (
	"log"

	"github.com/DanielRenne/GoCore/core"
)

// GetDump many values to a string (this returns nothing and prints nothing to stdout)
func ExampleGetDump() {
	type example struct {
		bytes []byte
		data  *string
		fun   func()
	}

	type deeper struct {
		Test string `json:"Test"`
	}
	type data struct {
		Id        string `json:"Id"`
		Name      string `json:"Name"`
		Ordering  int    `json:"Ordering"`
		InfoPopup string `json:"InfoPopup"`
		Color     string `json:"Color"`
		Slug      string `json:"Slug"`
		Deeper    deeper `json:"Deeper"`
	}
	var records []data
	records = append(records, data{Id: "1", Name: "test", Ordering: 1, InfoPopup: "test", Color: "test", Slug: "test", Deeper: deeper{Test: "test"}})
	records = append(records, data{Id: "2", Name: "test", Ordering: 1, InfoPopup: "test", Color: "test", Slug: "test"})
	tmp := example{}
	ch := make(chan string)
	floatNumber := 5.999
	str := "1234"
	test := make(map[string]example, 0)
	test["testing"] = example{
		data:  &str,
		bytes: []byte{1, 2, 3, 4},
		fun: func() {
			log.Println("test")
		},
	}
	stringWithNonPrintables := "\x00 \x04 \x08 \x0c \x10 \x14 \x18 \x1c \x20 \x24 \x28 \x2c \x30 \x34 \x38 \x3c \x40 \x44 \x48 \x4c \x50 \x54 \x58 \x5c \x60 \x64 \x68 \x6c \x70 \x74 \x78 \x7c \x80 \x84 \x88 \x8c \x90 \x94 \x98 \x9c \xa0 \xa4 \xa8 \xac \xb0 \xb4 \xb8 \xbc \xc0 \xc4 \xc8 \xcc \xd0 \xd4 \xd8 \xdc \xe0 \xe4 \xe8 \xec \xf0 \xf4 \xf8 \xfc \x100 \x104 \x108 \x10c \x110 \x114 \x118 \x11c \x120 \x124 \x128 \x12c \x130 \x134 \x138 \x13c \x140 \x144 \x148 \x14c \x150 \x154 \x158 \x15c \x160 \x164 \x168 \x16c \x170 \x174 \x178 \x17c \x180 \x184 \x188 \x18c \x190 \x194 \x198 \x19c \x1a0 \x1a4 \x1a8 \x1ac \x1b0 \x1b4 \x1b8 \x1bc \x1c0 \x1c4 \x1c8 \x1cc \x1d0 \x1d4 \x1d8 \x1dc \x1e0 \x1e4 \x1e8 \x1ec \x1f0 \x1f4 \x1f8 \x1fc \x200 \x204 \x208 \x20c \x210 \x214 \x218 \x21c \x220 \x224 \x228 \x22c \x230 \x234 \x238 \x23c \x240 \x244 \x248 \x24c"
	stringWithEmojiOnly := "😄 🎉"

	dumpedData1 := core.GetDump(true, ch, stringWithEmojiOnly, stringWithNonPrintables, floatNumber, test, "just a string", floatNumber, example{
		data: &str,
	}, example{}, tmp.data, []string{"GoCore", "Rocks"}, records)

	dumpedData2 := core.GetDump("A second dump will output a new timestamp", "And more logs of anything you paste as a parameter")

	// You would log this how you want.  This goes to /dev/null and is here just so this compiles
	core.Debug.Nop(dumpedData1, dumpedData2)
	//Output:
}

// Dump many values to stdout
func ExampleDump() {
	type example struct {
		bytes []byte
		data  *string
		fun   func()
	}

	type deeper struct {
		Test string `json:"Test"`
	}
	type data struct {
		Id        string `json:"Id"`
		Name      string `json:"Name"`
		Ordering  int    `json:"Ordering"`
		InfoPopup string `json:"InfoPopup"`
		Color     string `json:"Color"`
		Slug      string `json:"Slug"`
		Deeper    deeper `json:"Deeper"`
	}
	var records []data
	records = append(records, data{Id: "1", Name: "test", Ordering: 1, InfoPopup: "test", Color: "test", Slug: "test", Deeper: deeper{Test: "test"}})
	records = append(records, data{Id: "2", Name: "test", Ordering: 1, InfoPopup: "test", Color: "test", Slug: "test"})
	tmp := example{}
	ch := make(chan string)
	floatNumber := 5.999
	str := "1234"
	test := make(map[string]example, 0)
	test["testing"] = example{
		data:  &str,
		bytes: []byte{1, 2, 3, 4},
		fun: func() {
			log.Println("test")
		},
	}
	stringWithNonPrintables := "\x00 \x04 \x08 \x0c \x10 \x14 \x18 \x1c \x20 \x24 \x28 \x2c \x30 \x34 \x38 \x3c \x40 \x44 \x48 \x4c \x50 \x54 \x58 \x5c \x60 \x64 \x68 \x6c \x70 \x74 \x78 \x7c \x80 \x84 \x88 \x8c \x90 \x94 \x98 \x9c \xa0 \xa4 \xa8 \xac \xb0 \xb4 \xb8 \xbc \xc0 \xc4 \xc8 \xcc \xd0 \xd4 \xd8 \xdc \xe0 \xe4 \xe8 \xec \xf0 \xf4 \xf8 \xfc \x100 \x104 \x108 \x10c \x110 \x114 \x118 \x11c \x120 \x124 \x128 \x12c \x130 \x134 \x138 \x13c \x140 \x144 \x148 \x14c \x150 \x154 \x158 \x15c \x160 \x164 \x168 \x16c \x170 \x174 \x178 \x17c \x180 \x184 \x188 \x18c \x190 \x194 \x198 \x19c \x1a0 \x1a4 \x1a8 \x1ac \x1b0 \x1b4 \x1b8 \x1bc \x1c0 \x1c4 \x1c8 \x1cc \x1d0 \x1d4 \x1d8 \x1dc \x1e0 \x1e4 \x1e8 \x1ec \x1f0 \x1f4 \x1f8 \x1fc \x200 \x204 \x208 \x20c \x210 \x214 \x218 \x21c \x220 \x224 \x228 \x22c \x230 \x234 \x238 \x23c \x240 \x244 \x248 \x24c"
	stringWithEmojiOnly := "😄 🎉"

	core.Dump(true, ch, stringWithEmojiOnly, stringWithNonPrintables, floatNumber, test, "just a string", floatNumber, example{
		data: &str,
	}, example{}, tmp.data, []string{"GoCore", "Rocks"}, records)

	core.Dump("A second dump will output a new timestamp", "And more logs of anything you paste as a parameter")
	/*
		Output:
			!!!!!!!!!!!!! DEBUG 2022-10-04 13:24:34.336824!!!!!!!!!!!!!


			#### bool                                    ####
			true
			#### chan                                    ####
			0xc00007c0c0
			#### string                                  [len:9]####
			😄 🎉
			#### string (non printables -> dump hex)     [len:379]####
			00000000  00 20 04 20 08 20 0c 20  10 20 14 20 18 20 1c 20  |. . . . . . . . |
			00000010  20 20 24 20 28 20 2c 20  30 20 34 20 38 20 3c 20  |  $ ( , 0 4 8 < |
			00000020  40 20 44 20 48 20 4c 20  50 20 54 20 58 20 5c 20  |@ D H L P T X \ |
			00000030  60 20 64 20 68 20 6c 20  70 20 74 20 78 20 7c 20  |` d h l p t x | |
			00000040  80 20 84 20 88 20 8c 20  90 20 94 20 98 20 9c 20  |. . . . . . . . |
			00000050  a0 20 a4 20 a8 20 ac 20  b0 20 b4 20 b8 20 bc 20  |. . . . . . . . |
			00000060  c0 20 c4 20 c8 20 cc 20  d0 20 d4 20 d8 20 dc 20  |. . . . . . . . |
			00000070  e0 20 e4 20 e8 20 ec 20  f0 20 f4 20 f8 20 fc 20  |. . . . . . . . |
			00000080  10 30 20 10 34 20 10 38  20 10 63 20 11 30 20 11  |.0 .4 .8 .c .0 .|
			00000090  34 20 11 38 20 11 63 20  12 30 20 12 34 20 12 38  |4 .8 .c .0 .4 .8|
			000000a0  20 12 63 20 13 30 20 13  34 20 13 38 20 13 63 20  | .c .0 .4 .8 .c |
			000000b0  14 30 20 14 34 20 14 38  20 14 63 20 15 30 20 15  |.0 .4 .8 .c .0 .|
			000000c0  34 20 15 38 20 15 63 20  16 30 20 16 34 20 16 38  |4 .8 .c .0 .4 .8|
			000000d0  20 16 63 20 17 30 20 17  34 20 17 38 20 17 63 20  | .c .0 .4 .8 .c |
			000000e0  18 30 20 18 34 20 18 38  20 18 63 20 19 30 20 19  |.0 .4 .8 .c .0 .|
			000000f0  34 20 19 38 20 19 63 20  1a 30 20 1a 34 20 1a 38  |4 .8 .c .0 .4 .8|
			00000100  20 1a 63 20 1b 30 20 1b  34 20 1b 38 20 1b 63 20  | .c .0 .4 .8 .c |
			00000110  1c 30 20 1c 34 20 1c 38  20 1c 63 20 1d 30 20 1d  |.0 .4 .8 .c .0 .|
			00000120  34 20 1d 38 20 1d 63 20  1e 30 20 1e 34 20 1e 38  |4 .8 .c .0 .4 .8|
			00000130  20 1e 63 20 1f 30 20 1f  34 20 1f 38 20 1f 63 20  | .c .0 .4 .8 .c |
			00000140  20 30 20 20 34 20 20 38  20 20 63 20 21 30 20 21  | 0  4  8  c !0 !|
			00000150  34 20 21 38 20 21 63 20  22 30 20 22 34 20 22 38  |4 !8 !c "0 "4 "8|
			00000160  20 22 63 20 23 30 20 23  34 20 23 38 20 23 63 20  | "c #0 #4 #8 #c |
			00000170  24 30 20 24 34 20 24 38  20 24 63                 |$0 $4 $8 $c|
			#### float64                                 ####
			5.999
			#### map                                     ####
			{
				"testing": {}
			}
			#### string                                  [len:13]####
			just a string
			#### float64                                 ####
			5.999
			#### struct                                  ####
			{}
			#### main.example                            ####
			{bytes:[] data:<nil> fun:<nil>}
			#### *string                                 ####
			<nil>
			#### slice                                   [len:2]####
			[
				"GoCore",
				"Rocks"
			]
			#### slice                                   [len:2]####
			[
				{
					"Id": "1",
					"Name": "test",
					"Ordering": 1,
					"InfoPopup": "test",
					"Color": "test",
					"Slug": "test",
					"Deeper": {
						"Test": "test"
					}
				},
				{
					"Id": "2",
					"Name": "test",
					"Ordering": 1,
					"InfoPopup": "test",
					"Color": "test",
					"Slug": "test",
					"Deeper": {
						"Test": ""
					}
				}
			]

			!!!!!!!!!!!!! ENDDEBUG 2022-10-04 13:24:34.336824!!!!!!!!!!!!!
			!!!!!!!!!!!!! DEBUG 2022-10-04 13:24:34.337172!!!!!!!!!!!!!


			#### string                                  [len:41]####
			A second dump will output a new timestamp
			#### string                                  [len:50]####
			And more logs of anything you paste as a parameter

			!!!!!!!!!!!!! ENDDEBUG 2022-10-04 13:24:34.337172!!!!!!!!!!!!!
	*/
}
