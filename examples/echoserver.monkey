#!./monkey-lang

fd := socket("tcp4")
bind(fd, "0.0.0.0:8000")
listen(fd, 1)

nfd := accept(fd)
msg := read(nfd)
write(nfd, msg)
close(nfd)
close(fd)
