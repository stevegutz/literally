# literally

A confusingly named library for turning go values into their literal representations. This is meant to be used for generating tests and code in failure tolerant settings where speed is not a huge concern. Don't let me catch this in a vital code path.

There are some glaring holes in the library when it comes to certain types -- functions, unsafe pointers, various pointers to various types -- for which "nil" is always returned.

See [the godoc](http://godoc.org/github.com/stevegutz/literally)
