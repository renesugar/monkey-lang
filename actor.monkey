#!./monkey-lang

counter := actor(
  fn(sender, msg) {
    print("sender: " + sender)
    print("msg   : " + str(msg))
    if (type(msg) == "str" && msg == "started") {
      print("creating state")
      n := 0
    } else if (type(msg) == "int") {
      print("updating state")
      n = n + msg
    } else {
      print("unknown message " + type(msg) + " " + msg)
    }
  }
)
  
start(counter)
