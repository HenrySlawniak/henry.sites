This is a single tool for serving multiple domains from one process. Supports
HTTP2, and 'vhosts'

Create a directory for individual sites in the `sites/` directory

For example, if you have a site `example.com`, put all of the HTML css, etc in
`sites/example.com`.

Then if you have `anotherexample.com`, put all of the HTML css, etc in
`sites/anotherexample.com`.

Any domain name that is pointed to this process will serve the assets in the
`client/` directory.

Configure your catch-all site in the `client/` directory. A simple static site
is already there
