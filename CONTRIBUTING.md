# Contributing to EMBD

This actually is really simple. A few simple guidelines and we can break for dinner:

* EMBD is designed with a lot of affection, with utmost importance given to the dev experience (read: the API feel and style.) So always think from that angle when creating the pull request
* [Documentation](https://godoc.org/github.com/kidoman/embd) helps drive adoption. No exceptions

When it comes to the code:

* Always [gofmt + goimports](https://michaelwhatcott.com/gosublime-goimports/) the code. We absolutely adore them. Sublime Text 3 with GoSublime + GoImports is a very potent combination in our opinion
* Often you will hear us mention idiomatic Go ([read](http://golang.org/doc/effective_go.html).) Writing Go "The Go Way™" helps keep the code readable and understandable by all Gophers.
* No blank lines where they don't belong

'commit'tee called:

* Commit messages should be all lower case, unless when absolutely required (proper nouns, acronyms, etc.)
* When possible, prefix the message with the general area the commit is regarding. Good examples would be ```doc```, ```gpio```, ```bbb```. You get the drift
* If the commit message is long, then follow this convention:

	```
	gpio: adding interrupts

	this is inspired by Dave Cheney's gpio library and his work on EPOLL
	```

* Individual lines must be wrapped at the 70-char limit. Yeah, old school	
* No trailing '.'

And:

* Real tabs for indentation
