module github.com/ripx80/brewman

go 1.13

require (
	github.com/gdamore/tcell v1.3.0
	github.com/ghodss/yaml v1.0.0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/ripx80/brave v0.0.0-20200724063539-08b685d50e39
	github.com/ripx80/recipe v0.0.0-20210326064505-c8e4accab37d
	github.com/ripx80/signal v0.0.0-20200221114242-f8c0121af15b
	github.com/rivo/tview v0.0.0-20200414130344-8e06c826b3a5
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.5.1
	gopkg.in/validator.v2 v2.0.0-20191107172027-c3144fdedc21
	periph.io/x/periph v3.6.2+incompatible
)

//replace github.com/ripx80/recipe => ../recipe/
