require "socket"

class Interpreter
  SOCKET_NAME = ARGV[0]
  MASTER_PID = ARGV[1].to_i
  ASSETS_PATH = ARGV[2]

  def self.serve
    prepare_for_automatic_termination

    server = listen
    loop do
      client = server.accept
      format, option, content = read_all(client)

      begin
        result = self.send("render_#{format}", content, option)
        client.write "success<<#{result}"
      rescue Exception => e
        puts e
        client.write "error<<#{e}"
      end
      client.close
    end
  end

  # Shut down the interpreter when master process does not exist.
  def self.prepare_for_automatic_termination
    interpreter_pid = Process.pid
    fork do
      $0 = "ruby [train] master process monitor (PID #{MASTER_PID})" # Set Process Name
      loop do
        sleep 2
        begin
          Process.getpgid(MASTER_PID)
        rescue #=> when master process cannot be found
          Interpreter.clean_up
          Process.kill 1, interpreter_pid
          exit
        end
      end
    end
  end

  def self.clean_up
    `rm -f #{SOCKET_NAME}`
  end

  private

  def self.listen
    begin
      clean_up

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

  def self.render_sass content, option
    _render_sass(content, :sass, option)
  end

  def self.render_scss content, option
    _render_sass(content, :scss, option)
  end

  def self._render_sass content, syntax, option
    require "sass"

    options = {
      :load_paths => ["#{ASSETS_PATH}/stylesheets"],
      :syntax => syntax
    }

    options[:debug_info] = true if option == "debug_info"
    options[:line_numbers] = true if option == "line_numbers"

    engine = Sass::Engine.new(content, options)
    engine.render
  end

  def self.render_coffee content, option
    require "coffee-script"
    CoffeeScript.compile content
  end
end

trap("INT"){} #=> make silent the "Interrupted" error
Interpreter.serve
