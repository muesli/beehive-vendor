# beehive-vendor

Vendor files for Beehive.

## Updating beehive dependencies

If you're adding new dependencies to beehive, you'll need to add them to this repository too.

Contributors are encouraged to use the [dep](https://github.com/golang/dep) tool for now.

The `dep` tool required manifest (`Gopkg.toml`) won't be commited here til that tool matures and recommends it:

> Note that the manifest and lock file formats are not finalized, and will likely change before the tool is released. We make no compatibility guarantees for the time being. Please don't commit any code or files created with the tool.

The recommended worflow for now is to:

1. Install `dep`: `go get -u github.com/golang/dep/cmd/dep`.
2. Get a fresh beehive checkout.
3. Delete the vendor folder.
4. Run `dep init`.
5. Run `dep ensure github.com/foo/bar@<version>` to add the new dependecie(s) you are interested in.
6. Get a fresh `beehive-vendor` checkout.
7. Replace the contents of `beehive-vendor` with the contents of the freshly generated `vendor` folder.
8. Open a pull request to `beehive-vendor`.

That'll add all the required dependencies and the new ones.
