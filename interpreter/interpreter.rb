require "sass"
require "socket"

class Interpreter
  SOCKET_NAME = "/tmp/train.interpreter.socket"

  def self.run
    server = listen

    loop {
      client = server.accept
      format, content = read_all(client)

      begin
        result = self.send("render_#{format}", content)
        client.write "success<<#{result}"
      rescue => e
        puts e
        client.write "error<<#{e}"
      end
      client.close
    }
  end

  private
  def self.listen
    begin
      `rm -f #{SOCKET_NAME}`
      server = UNIXServer.new(SOCKET_NAME)
      puts "<<ready"
      `touch #{SOCKET_NAME}`
      server
    rescue => e
      puts e
      exit 1
    end
  end


  def self.read_all client
    data = ""
    recv_length = 2000
    while tmp = client.recv(recv_length)
      data += tmp
      break if tmp.length < recv_length
    end
   data.split("<<")
  end

  def self.render_sass content
    engine = Sass::Engine.new(content, :load_paths => ["assets/stylesheets"])
    engine.render
  end
end

Interpreter.run
