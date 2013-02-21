class Interpreter
  def self.run socket_name
    server = listen(socket_name)

    loop {
      client = server.accept
      content = read_all(client)

      begin
        client.write render(content)
      rescue => e
        puts e
        client.puts "<<error"
        client.write e
      end
      client.close
    }
  end

  private
  def self.listen(socket_name)
    begin
      `rm #{socket_name}`
      server = UNIXServer.new(socket_name)
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
   data
  end
end
