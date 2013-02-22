mjolnir
=======

_[mjolnir](https://github.com/daniel-fanjul-alcuten/mjolnir)_ is a launcher to simplify the test, compilation and execution of programs using libmjolnir.

libmjolnir
==========

_[libmjolnir](https://github.com/daniel-fanjul-alcuten/libmjolnir)_ is a go library to write build systems.

libmjolnir-demo
===============

_[libmjolnir-demo](https://github.com/daniel-fanjul-alcuten/libmjolnir-demo)_ is an example project using libmjolnir.

Workspaces
----------

A program using libmjolnir might use several go workspaces, and the setup of GOPATH could be not simple.

_mjolnir_ is a simple program that:

* reads its configuration from a a _mjolnir.json_ file,
* prepends the specified go workspaces to GOPATH,
* runs 'go test' on the specified go packages,
* runs 'go install' on the specified go packages,
* runs your main program with the same arguments as _mjolnir_ receives.

Example
-------

<pre>
{
  "GoDir": "mjolnir/src",
  "GoPaths": [
    "mjolnir"
  ],
  "TestPackages": [
    "libfoobar/foo",
    "libfoobar/bar",
  ],
  "InstallPackages": [
    "libfoobar/foo",
    "libfoobar/bar",
    "libfoobar/build"
  ],
  "Exec": "mjolnir/bin/build"
}
</pre>
