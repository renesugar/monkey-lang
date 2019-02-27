#!./monkey-lang

fd := socket("tcp4")
bind(fd, "127.0.0.1:32535")
connect(fd, "127.0.0.1:8000")
write(fd, "Hello World")
print(read(fd))
close(fd)
