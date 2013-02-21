require "sass"
require 'socket'

def read_all(client)
  data = ""
  recv_length = 2000
  while tmp = client.recv(recv_length)
      data += tmp
      break if tmp.length < recv_length
  end

 return data
end


SOCKET_NAME = "/tmp/train.sass.socket"

begin
  `rm #{SOCKET_NAME}`
  server = UNIXServer.new(SOCKET_NAME)
  puts "<<ready"
  `touch #{SOCKET_NAME}`
rescue => e
  puts e
  exit 1
end

loop {
  client = server.accept
  content = read_all(client)

  begin
    engine = Sass::Engine.new(content, :load_paths => ["assets/stylesheets"])
    client.write engine.render
  rescue => e
    puts e
    client.puts "<<error"
    client.write e
  end


  client.close
}


