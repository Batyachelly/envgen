package example

import "time"

type Config struct {
	// KEY1 Description
	KEY1 string `env:"KEY1,required"  envExample:"value_example_1"`
	// KEY2 Description
	// Additional info
	KEY2 string `env:"KEY2" envDefault:"value_2"`
	// KEY3 Description
	KEY3 []string `env:"KEY3,required" envSeparator:"," envExample:"value_3_1,value_3_2"`
	// KEY4 Description
	KEY4 []string `env:"KEY4" envSeparator:"," envDefault:"value_4_1,value_4_2"`

	InnerStruct struct {
		// KEY5 Description
		KEY5 string `env:"KEY5,required" envExample:"value_example_5" envDefault:"value_5"`
		// KEY6 Description
		KEY6 string `env:"KEY6" envDefault:"value_6"`
	}

	AnotherConfig1 AnotherConfig1

	AnotherConfig2 AnotherConfig2
} // file1.env

type AnotherConfig1 struct {
	// KEY7 Description
	KEY7 time.Duration `env:"KEY7" envDefault:"1m"`
} // file1.env

type AnotherConfig2 struct {
	// KEY8 Description
	KEY8 time.Duration `env:"KEY8" envDefault:"1m"`
} // file2.env
