A simple proof of concept for parsing haproxy configurations and
transforming them into JSON.

Running:

    go run {parse,scanner}.go haproxy.cfg

This is far from feature complete, and approach is [based on this blog post](https://blog.gopheracademy.com/advent-2014/parsers-lexers/).