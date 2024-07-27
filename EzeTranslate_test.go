package main

import (
	"fmt"
	"github.com/Ericwyn/EzeTranslate/conf"
	"github.com/Ericwyn/EzeTranslate/strutils"
	"testing"
)

func TestFormatAnnotation(t *testing.T) {
	conf.InitConfig()

	fmt.Println(strutils.FormatInputBoxText(`/**
 * The service class that manages LocationProviders and issues location
 * updates and alerts.
 */
`))
}

func TestFormatCamelCaseText(t *testing.T) {
	a := "hello friends"
	fmt.Println(a, "->", strutils.FormatCamelCaseText(a))

	b := "requestLocationUpdate"
	fmt.Println(b, "->", strutils.FormatCamelCaseText(b))

	c := "format_linux_function"
	fmt.Println(c, "->", strutils.FormatCamelCaseText(c))

	d := "CONFIG_KEY_TRANSLATE"
	fmt.Println(d, "->", strutils.FormatCamelCaseText(d))

	e := "Con_F_Key"
	fmt.Println(e, "->", strutils.FormatCamelCaseText(e))

}
