module github.com/tklauser/buildinfo/aprogram

go 1.18

require (
	github.com/tklauser/buildinfo/apackage v0.0.2-unpublished
)

replace github.com/tklauser/buildinfo/apackage v0.0.2-unpublished => ../replacement
