* there is a bug where the client could crash the server with an index out of
bounds by passing in a command with more bytes than my read buffer has
provided

* points, locations, and ranges should really have an api
  - server should be 0 based
  - client can have a "to row point" where rows are one based and cols are
    zero based
  - once i take in input from the client then i'll need to also have location
    from cursor


* color rework probably is needed especially when it comes to location within app
